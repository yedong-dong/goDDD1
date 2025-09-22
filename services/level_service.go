package services

import (
	"goDDD1/config"
	"goDDD1/models"
)

type LevelService interface {
	// 获取用户当前等级信息
	GetUserLevel(userID uint) (*models.User, error)
	// 增加用户经验值，检查是否升级
	AddExpeirence(userID uint, exp uint, description string) (*models.LevelHistory, error)
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

func (s *levelService) AddExpeirence(userID uint, exp uint, description string) (*models.LevelHistory, error) {
	tx := config.Database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var user models.User
	if err := tx.Where("uid=?", userID).First(&user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	oldLevel := user.Level
	user.Experience += exp

	// 检查是否升级
	levelConfig, err := s.GetLevelConfig(user.Level + 1)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if user.Experience >= levelConfig.RequiredExp {
		// 升级
		user.Level++
	}
	// 保存用户信息
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 记录等级变更历史
	history := &models.LevelHistory{
		UserID:          userID,
		OldLevel:        oldLevel,
		NewLevel:        user.Level,
		ExpGained:       exp,
		Experience:      user.Experience,
		CoinRewarded:    levelConfig.CoinReward,
		DiamondRewarded: levelConfig.DiamondReward,
		Description:     description,
	}
	if err := tx.Create(history).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

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
