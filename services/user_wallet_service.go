package services

import (
	"fmt"
	"goDDD1/config"
	"goDDD1/models"
	"goDDD1/utils"
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

// UserWalletService 用户钱包服务接口

type UserWalletService interface {
	InitializeWallet(userID uint) error
	InitializeWalletWithTx(tx *gorm.DB, userID uint) error
	GetWalletByUserIDAndType(userID uint, walletType models.WalletType) (*models.UserWallet, error)
	GetWalletByUserIDAndTypeWithTx(tx *gorm.DB, userID uint, walletType models.WalletType) (*models.UserWallet, error)
	UpdateWalletBalance(userID uint, walletType models.WalletType, amount int64) error
	UpdateWalletBalance2(userID uint, walletType models.WalletType, amount int64) error
	UpdateWalletBalanceWithTx(tx *gorm.DB, userID uint, walletType models.WalletType, amount int64, description string) error
	GetUserWallets(userID uint) ([]models.UserWallet, error)
}

// userWalletService 用户钱包服务实现

type userWalletService struct {
	UserCurrencyFlow  UserCurrencyFlowServiceInterface
	rewardFlowService RewardFlowService
}

// NewUserWalletService 创建用户钱包服务实例
func NewUserWalletService() UserWalletService {
	return &userWalletService{
		UserCurrencyFlow:  NewUserCurrencyFlowService(),
		rewardFlowService: NewRewardFlowService(),
	}
}

// InitializeWallet 初始化用户钱包（创建用户时调用）
func (s *userWalletService) InitializeWallet(userID uint) error {
	tx := config.Database.Begin()
	if err := s.InitializeWalletWithTx(tx, userID); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// InitializeWalletWithTx 使用事务初始化用户钱包
func (s *userWalletService) InitializeWalletWithTx(tx *gorm.DB, userID uint) error {
	// 创建coin钱包，初始化1000个coin
	coinWallet := models.UserWallet{
		UserID: userID,
		Num:    1000,
		Type:   models.Coin,
	}
	if err := tx.Create(&coinWallet).Error; err != nil {
		return err
	}

	// 创建diamond钱包，初始化200个diamond
	diamondWallet := models.UserWallet{
		UserID: userID,
		Num:    200,
		Type:   models.Diamond,
	}
	return tx.Create(&diamondWallet).Error
}

// GetWalletByUserIDAndType 根据用户ID和钱包类型获取钱包
func (s *userWalletService) GetWalletByUserIDAndType(userID uint, walletType models.WalletType) (*models.UserWallet, error) {
	var wallet models.UserWallet
	result := config.Database.Where("user_id = ? AND type = ?", userID, walletType).First(&wallet)
	if result.Error != nil {
		return nil, result.Error
	}
	return &wallet, nil
}

// GetWalletByUserIDAndTypeWithTx 使用事务根据用户ID和钱包类型获取钱包
func (s *userWalletService) GetWalletByUserIDAndTypeWithTx(tx *gorm.DB, userID uint, walletType models.WalletType) (*models.UserWallet, error) {
	var wallet models.UserWallet
	result := tx.Where("user_id = ? AND type = ?", userID, walletType).First(&wallet)
	if result.Error != nil {
		return nil, result.Error
	}
	return &wallet, nil
}

// UpdateWalletBalance 更新钱包余额
func (s *userWalletService) UpdateWalletBalance(userID uint, walletType models.WalletType, amount int64) error {
	tx := config.Database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新余额
	if err := tx.Where("user_id = ? AND type = ?", userID, walletType).Update("num", gorm.Expr("num + ?", amount)).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 添加奖励流水记录
	rewardType := models.RewardTypeCoin
	if walletType == models.Diamond {
		rewardType = models.RewardTypeDiamond
	}

	err := s.rewardFlowService.CreateRewardFlow(tx, userID, rewardType, 0, amount, "钱包余额更新")
	if err != nil {
		tx.Rollback()
		return err
	}

	// 删除缓存记录
	cacheKey := fmt.Sprintf(models.CacheKeyUserBackpack, userID)
	err = utils.DelHashField(cacheKey, "wallets")
	if err == nil {
		log.Printf("successful delete cacheKey: %s wallets", cacheKey)
	}

	return tx.Commit().Error
}

// UpdateWalletBalance2 更新钱包余额
func (s *userWalletService) UpdateWalletBalance2(userID uint, walletType models.WalletType, amount int64) error {
	tx := config.Database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新余额
	if err := tx.Where("user_id = ? AND type = ?", userID, walletType).Update("num", gorm.Expr("num + ?", amount)).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 添加流水代码
	if err := s.UserCurrencyFlow.CreateUserCurrencyFlow(&models.UserCurrencyFlow{
		UserID:   userID,
		CostType: string(walletType),
		Price:    amount,
	}); err != nil {
		tx.Rollback()
		return err
	}

	// 添加奖励流水记录
	rewardType := models.RewardTypeCoin
	if walletType == models.Diamond {
		rewardType = models.RewardTypeDiamond
	}

	err := s.rewardFlowService.CreateRewardFlow(tx, userID, rewardType, 0, amount, "钱包余额更新")
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// UpdateWalletBalanceWithTx 使用事务更新钱包余额
func (s *userWalletService) UpdateWalletBalanceWithTx(tx *gorm.DB, userID uint, walletType models.WalletType, amount int64, description string) error {
	var wallet models.UserWallet

	tx = tx.Set("gorm:query_option", "FOR UPDATE")

	if err := tx.Where("user_id = ? AND type = ?", userID, walletType).
		First(&wallet).Error; err != nil {
		return err
	}

	// 更新余额
	wallet.Num += amount

	// 添加流水代码
	if err := s.UserCurrencyFlow.CreateUserCurrencyFlow(&models.UserCurrencyFlow{
		UserID:      userID,
		CostType:    string(walletType),
		Price:       amount,
		Description: description,
	}); err != nil {
		return err
	}

	// 添加奖励流水记录
	rewardType := models.RewardTypeCoin
	if walletType == models.Diamond {
		rewardType = models.RewardTypeDiamond
	}

	var itemID uint
	if walletType == models.Diamond {
		itemID = 0
	}
	if walletType == models.Coin {
		itemID = 1
	}

	err := s.rewardFlowService.CreateRewardFlow(tx, userID, rewardType, itemID, amount, description)
	if err != nil {
		return err
	}

	// 保存更新
	return tx.Save(&wallet).Error
}

// GetUserWallets 获取用户所有钱包
func (s *userWalletService) GetUserWallets(userID uint) ([]models.UserWallet, error) {
	cacheKey := fmt.Sprintf(models.CacheKeyUserBackpack, userID)
	var wallets []models.UserWallet
	err := utils.GetHashField(cacheKey, "wallets", &wallets)
	if err == nil && len(wallets) > 0 {
		return wallets, nil
	}

	result := config.Database.Where("user_id = ?", userID).Find(&wallets)
	if result.Error != nil {
		return nil, result.Error
	}

	if len(wallets) > 0 {
		utils.SetHashField(cacheKey, "wallets", wallets, time.Hour)
	}

	return wallets, nil

}
