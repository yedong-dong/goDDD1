package services

import (
	"errors"
	"goDDD1/config"
	"goDDD1/models"
	"log"
)

// UserService 用户服务接口
type UserService interface {
	CreateUser(user *models.User) error
	GetUserByUID(id uint) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error
	GetAllUsers(isDeleted string) ([]*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}

// userService 用户服务实现
type userService struct{}

// NewUserService 创建用户服务实例
func NewUserService() UserService {
	return &userService{}
}

// GetUserByEmail 根据邮箱获取用户
func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := config.Database.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// CreateUser 创建用户
func (s *userService) CreateUser(user *models.User) error {
	// 开始事务
	tx := config.Database.Begin()
	// 使用defer确保在函数退出时处理panic情况
	// 如果发生panic，自动回滚事务以保证数据一致性
	defer func() {
		if r := recover(); r != nil {
			// 捕获panic并回滚事务，防止数据不一致
			tx.Rollback()
		}
	}()

	// 创建用户
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 初始化用户钱包
	walletService := NewUserWalletService()
	if err := walletService.InitializeWalletWithTx(tx, user.UID); err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	// 当所有数据库操作都成功完成后，提交事务以确保数据的一致性和持久化
	// 此时会将事务中的所有更改永久保存到数据库中
	// 包括：1. 用户基本信息的创建 2. 用户coin钱包的初始化(1000个coin) 3. 用户diamond钱包的初始化(200个diamond)
	// 如果提交失败，GORM会自动处理错误，确保数据库状态的完整性
	return tx.Commit().Error
}

// GetUserByID 根据ID获取用户
func (s *userService) GetUserByUID(uid uint) (*models.User, error) {
	var user models.User
	result := config.Database.Where("uid = ?", uid).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (s *userService) GetAllUsers(isDeleted string) ([]*models.User, error) {
	var users []*models.User
	result := config.Database.Where("is_deleted = ?", isDeleted).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// GetUserByUsername 根据用户名获取用户
func (s *userService) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	result := config.Database.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func (s *userService) UpdateUser(user *models.User) error {
	if user.ID == 0 {
		return errors.New("用户ID不能为空")
	}
	// 修改：将 log.Fatal 改为 log.Println，避免程序退出
	log.Println("更新用户信息", user)
	return config.Database.Save(user).Error
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(id uint) error {
	return config.Database.Delete(&models.User{}, id).Error
}
