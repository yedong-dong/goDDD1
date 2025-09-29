package services

import (
	"fmt"
	"goDDD1/config"
	"goDDD1/models"

	"github.com/jinzhu/gorm"
)

type LevelService interface {
	// 获取用户当前等级信息
	GetUserLevel(userID uint) (*models.User, error)
	// 增加用户经验值，检查是否升级
	AddExpeirence(tx *gorm.DB, userID uint, exp uint, description string) (*models.LevelHistory, error)
	// 获取用户等级历史记录
	GetLevelHistory(userID uint) ([]*models.LevelHistory, error)
	// 获取等级配置
	GetLevelConfig(level uint) (*models.LevelConfig, error)
	// 获取所有等级配置
	GetAllLevelConfigs() ([]*models.LevelConfig, error)
	// 计算用户购买商品的折扣价格
	CalculateDiscountPrice(userID uint, originalPrice uint) (uint, error)
}

type levelService struct {
	userWalletService UserWalletService
}

// NewLevelService 创建等级服务
func NewLevelService() LevelService {
	return &levelService{
		userWalletService: NewUserWalletService(),
	}
}

func (s *levelService) GetUserLevel(userID uint) (*models.User, error) {
	var user models.User
	if err := config.Database.Where("uid = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// 修改AddExpeirence方法，不在内部提交事务
func (s *levelService) AddExpeirence(tx *gorm.DB, userID uint, exp uint, description string) (*models.LevelHistory, error) {
	var user models.User
	if err := tx.Where("uid=?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	oldLevel := user.Level
	user.Experience += exp

	var allLevelConfigs []models.LevelConfig
	if err := tx.Order("level asc").Find(&allLevelConfigs).Error; err != nil {
		return nil, err
	}

	levelConfigMap := make(map[uint]*models.LevelConfig)
	var maxLevel uint = 0
	for _, config := range allLevelConfigs {
		levelConfigMap[config.Level] = &config
		if config.Level > maxLevel {
			maxLevel = config.Level
		}
	}

	// 记录总奖励
	var totalCoinReward uint = 0
	var totalDiamondReward uint = 0
	var newLevel = user.Level

	// 检查是否升级，处理可能的多级跳升
	for level := user.Level + 1; level <= maxLevel; level++ {
		config, exists := levelConfigMap[level]
		if !exists {
			continue
		}

		// 如果经验值达到要求，则升级并累加奖励
		if user.Experience >= config.RequiredExp {
			newLevel = level
			totalCoinReward += config.CoinReward
			totalDiamondReward += config.DiamondReward
		} else {
			// 经验不足以升到下一级，终止检查
			break
		}
	}

	// 如果有升级，更新用户等级并发放奖励
	if newLevel > user.Level {
		user.Level = newLevel

		// 发放金币奖励
		if totalCoinReward > 0 {
			if err := s.userWalletService.UpdateWalletBalanceWithTx(tx, userID, "coin", int64(totalCoinReward), fmt.Sprintf("升级奖励%s", user.Level)); err != nil {
				return nil, err
			}
		}

		// 发放钻石奖励
		if totalDiamondReward > 0 {
			if err := s.userWalletService.UpdateWalletBalanceWithTx(tx, userID, "diamond", int64(totalDiamondReward), fmt.Sprintf("升级奖励%s", user.Level)); err != nil {
				return nil, err
			}
		}
	}

	// 保存用户信息
	if err := tx.Save(&user).Error; err != nil {
		return nil, err
	}

	// 无论是否升级，都记录经验值变化历史
	history := &models.LevelHistory{
		UserID:          userID,
		OldLevel:        oldLevel,
		NewLevel:        user.Level,
		ExpGained:       exp,
		Experience:      user.Experience,
		CoinRewarded:    totalCoinReward,
		DiamondRewarded: totalDiamondReward,
		Description:     description,
	}
	if err := tx.Create(history).Error; err != nil {
		return nil, err
	}

	// 移除事务提交，由调用方负责提交
	return history, nil
}

// GetLevelConfig 获取等级配置
func (s *levelService) GetLevelConfig(level uint) (*models.LevelConfig, error) {
	var levelConfig models.LevelConfig
	if err := config.Database.Where("level = ?", level).First(&levelConfig).Error; err != nil {

		return nil, err
	}
	return &levelConfig, nil
}

// GetLevelHistory 获取用户等级历史记录
func (s *levelService) GetLevelHistory(userID uint) ([]*models.LevelHistory, error) {
	var histories []*models.LevelHistory
	if err := config.Database.Where("user_id = ?", userID).Order("id desc").Find(&histories).Error; err != nil {
		return nil, err
	}
	return histories, nil
}

// GetAllLevelConfigs 获取所有等级配置
func (s *levelService) GetAllLevelConfigs() ([]*models.LevelConfig, error) {
	var configs []*models.LevelConfig
	if err := config.Database.Order("level asc").Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

// CalculateDiscountPrice 计算用户购买商品的折扣价格
func (s *levelService) CalculateDiscountPrice(userID uint, originalPrice uint) (uint, error) {
	// 获取用户信息
	user, err := s.GetUserLevel(userID)
	if err != nil {
		return originalPrice, err
	}

	// 获取用户等级配置
	levelConfig, err := s.GetLevelConfig(user.Level)
	if err != nil {
		return originalPrice, err
	}

	// 计算折扣价格
	discountPrice := originalPrice * levelConfig.DiscountPercent / 100
	return discountPrice, nil
}
