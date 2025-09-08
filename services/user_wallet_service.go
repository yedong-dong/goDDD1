package services

import (
	"goDDD1/config"
	"goDDD1/models"

	"github.com/jinzhu/gorm"
)

// UserWalletService 用户钱包服务接口
type UserWalletService interface {
	InitializeWallet(userID uint) error
	InitializeWalletWithTx(tx *gorm.DB, userID uint) error
	GetWalletByUserIDAndType(userID uint, walletType models.WalletType) (*models.UserWallet, error)
	UpdateWalletBalance(userID uint, walletType models.WalletType, amount int64) error
	GetUserWallets(userID uint) ([]models.UserWallet, error)
}

// userWalletService 用户钱包服务实现
type userWalletService struct{}

// NewUserWalletService 创建用户钱包服务实例
func NewUserWalletService() UserWalletService {
	return &userWalletService{}
}

// InitializeWallet 初始化用户钱包（创建用户时调用）
func (s *userWalletService) InitializeWallet(userID uint) error {
	return s.InitializeWalletWithTx(config.Database, userID)
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

// UpdateWalletBalance 更新钱包余额
func (s *userWalletService) UpdateWalletBalance(userID uint, walletType models.WalletType, amount int64) error {
	return config.Database.Model(&models.UserWallet{}).
		Where("user_id = ? AND type = ?", userID, walletType).
		Update("num", gorm.Expr("num + ?", amount)).Error
}

// GetUserWallets 获取用户所有钱包
func (s *userWalletService) GetUserWallets(userID uint) ([]models.UserWallet, error) {
	var wallets []models.UserWallet
	result := config.Database.Where("user_id = ?", userID).Find(&wallets)
	if result.Error != nil {
		return nil, result.Error
	}
	return wallets, nil
}
