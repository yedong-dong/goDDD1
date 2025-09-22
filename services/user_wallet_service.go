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
	GetWalletByUserIDAndTypeWithTx(tx *gorm.DB, userID uint, walletType models.WalletType) (*models.UserWallet, error)
	UpdateWalletBalance(userID uint, walletType models.WalletType, amount int64) error
	UpdateWalletBalanceWithTx(tx *gorm.DB, userID uint, walletType models.WalletType, amount int64) error
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
	// 创建事务
	tx := config.Database.Begin()

	// 使用事务更新钱包余额
	if err := s.UpdateWalletBalanceWithTx(tx, userID, walletType, amount); err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}

// UpdateWalletBalanceWithTx 使用事务更新钱包余额
func (s *userWalletService) UpdateWalletBalanceWithTx(tx *gorm.DB, userID uint, walletType models.WalletType, amount int64) error {
	// 设置超时时间（通过数据库层面的超时控制）
	// 注意：GORM v1.9.16不支持WithContext，所以我们使用其他方式处理超时

	// 使用FOR UPDATE锁定行，避免并发更新冲突
	var wallet models.UserWallet

	// 设置语句超时（如果数据库支持）
	tx = tx.Set("gorm:query_option", "FOR UPDATE")

	if err := tx.Where("user_id = ? AND type = ?", userID, walletType).
		First(&wallet).Error; err != nil {
		return err
	}

	// 更新余额
	wallet.Num += amount

	// 保存更新
	return tx.Save(&wallet).Error
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
