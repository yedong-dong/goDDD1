package services

import (
	"errors"
	"fmt"
	"goDDD1/config"
	"goDDD1/models"

	"github.com/jinzhu/gorm"
)

// RewardPackageService 奖励包服务接口
type RewardPackageService interface {
	// 奖励包管理
	CreateRewardPackage(pkg *models.RewardPackage, items []*models.RewardPackageItem) error
	UpdateRewardPackage(pkg *models.RewardPackage) error
	UpdateRewardPackageItems(packageID uint, items []*models.RewardPackageItem) error
	GetRewardPackageByID(id uint) (*models.RewardPackage, error)
	GetRewardPackageItems(packageID uint) ([]*models.RewardPackageItem, error)
	ListRewardPackages(page, pageSize int) ([]*models.RewardPackage, int64, error)
	DeleteRewardPackage(id uint) error

	// 奖励记录管理
	CreateRewardRecord(record *models.RewardRecord) error
	GetRewardRecordsByUserID(userID uint, page, pageSize int) ([]*models.RewardRecord, int64, error)
	GetRewardRecordByID(id uint) (*models.RewardRecord, error)

	// 奖励发放
	GrantReward(tx *gorm.DB, userID uint, packageID uint, source string) (*models.RewardRecord, error)
}

// rewardPackageService 奖励包服务实现
type rewardPackageService struct {
	userWalletService UserWalletService
	backpackService   BackpackService
}

// NewRewardPackageService 创建奖励包服务实例
func NewRewardPackageService() RewardPackageService {
	return &rewardPackageService{
		userWalletService: NewUserWalletService(),
		backpackService:   NewBackpackService(),
	}
}

// CreateRewardPackage 创建奖励包
func (s *rewardPackageService) CreateRewardPackage(pkg *models.RewardPackage, items []*models.RewardPackageItem) error {
	tx := config.Database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建奖励包
	if err := tx.Create(pkg).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 创建奖励包内容
	for _, item := range items {
		item.PackageID = pkg.ID
		if err := tx.Create(item.ID).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// UpdateRewardPackage 更新奖励包基本信息
func (s *rewardPackageService) UpdateRewardPackage(pkg *models.RewardPackage) error {
	return config.Database.Save(pkg).Error
}

// UpdateRewardPackageItems 更新奖励包内容
func (s *rewardPackageService) UpdateRewardPackageItems(packageID uint, items []*models.RewardPackageItem) error {
	tx := config.Database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除原有奖励包内容
	if err := tx.Where("package_id = ?", packageID).Delete(models.RewardPackageItem{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 创建新的奖励包内容
	for _, item := range items {
		item.PackageID = packageID
		if err := tx.Create(item).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// GetRewardPackageByID 根据ID获取奖励包
func (s *rewardPackageService) GetRewardPackageByID(id uint) (*models.RewardPackage, error) {
	var pkg models.RewardPackage
	if err := config.Database.First(&pkg, id).Error; err != nil {
		return nil, err
	}
	return &pkg, nil
}

// GetRewardPackageItems 获取奖励包内容
func (s *rewardPackageService) GetRewardPackageItems(packageID uint) ([]*models.RewardPackageItem, error) {
	var items []*models.RewardPackageItem
	if err := config.Database.Where("package_id = ?", packageID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// ListRewardPackages 分页获取奖励包列表
func (s *rewardPackageService) ListRewardPackages(page, pageSize int) ([]*models.RewardPackage, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	} else if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	var packages []*models.RewardPackage
	var total int64

	if err := config.Database.Model(&models.RewardPackage{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := config.Database.Offset(offset).Limit(pageSize).Order("id desc").Find(&packages).Error; err != nil {
		return nil, 0, err
	}

	return packages, total, nil
}

// DeleteRewardPackage 删除奖励包
func (s *rewardPackageService) DeleteRewardPackage(id uint) error {
	tx := config.Database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除奖励包内容
	if err := tx.Where("package_id = ?", id).Delete(models.RewardPackageItem{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除奖励包
	if err := tx.Delete(&models.RewardPackage{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// CreateRewardRecord 创建奖励记录
func (s *rewardPackageService) CreateRewardRecord(record *models.RewardRecord) error {
	return config.Database.Create(record).Error
}

// GetRewardRecordsByUserID 获取用户的奖励记录
func (s *rewardPackageService) GetRewardRecordsByUserID(userID uint, page, pageSize int) ([]*models.RewardRecord, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	} else if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	var records []*models.RewardRecord
	var total int64

	if err := config.Database.Model(&models.RewardRecord{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := config.Database.Where("user_id = ?", userID).Offset(offset).Limit(pageSize).Order("id desc").Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// GetRewardRecordByID 根据ID获取奖励记录
func (s *rewardPackageService) GetRewardRecordByID(id uint) (*models.RewardRecord, error) {
	var record models.RewardRecord
	if err := config.Database.First(&record, id).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// GrantReward 发放奖励
func (s *rewardPackageService) GrantReward(tx *gorm.DB, userID uint, packageID uint, source string) (*models.RewardRecord, error) {
	// 判断是否需要开启事务
	var localTx *gorm.DB
	if tx == nil {
		localTx = config.Database.Begin()
		tx = localTx
	}

	defer func() {
		if r := recover(); r != nil && localTx != nil {
			localTx.Rollback()
		}
	}()

	// 查询奖励包
	var pkg models.RewardPackage
	if err := tx.First(&pkg, packageID).Error; err != nil {
		if localTx != nil {
			localTx.Rollback()
		}
		return nil, err
	}

	// 查询奖励包内容
	var items []*models.RewardPackageItem
	if err := tx.Where("package_id = ?", packageID).Find(&items).Error; err != nil {
		if localTx != nil {
			localTx.Rollback()
		}
		return nil, err
	}

	if len(items) == 0 {
		if localTx != nil {
			localTx.Rollback()
		}
		return nil, errors.New("奖励包内容为空")
	}

	// 创建奖励记录
	record := &models.RewardRecord{
		UserID:    userID,
		PackageID: packageID,
		Source:    source,
	}

	if err := tx.Create(record).Error; err != nil {
		if localTx != nil {
			localTx.Rollback()
		}
		return nil, err
	}

	// 发放奖励
	for _, item := range items {
		// 修复奖励包发放时的参数类型和fmt.Sprintf拼接问题
		switch item.ItemType {
		case models.ItemTypeGoods: // 商品
			// 添加到背包
			if err := s.backpackService.AddToBackpack(tx, userID, item.ItemID, int64(item.Num)); err != nil {
				if localTx != nil {
					localTx.Rollback()
				}
				return nil, err
			}
		case models.ItemTypeCurrency: // 货币
			if item.ItemID == 0 {
				// 更新钱包 - 使用正确的models.WalletType类型参数
				if err := s.userWalletService.UpdateWalletBalanceWithTx(tx, userID, models.Diamond, int64(item.Num), fmt.Sprintf("奖励包发放，奖励包ID：%d", item.PackageID)); err != nil {
					if localTx != nil {
						localTx.Rollback()
					}
					return nil, err
				}
			}
			if item.ItemID == 1 {
				// 更新钱包 - 使用正确的models.WalletType类型参数
				if err := s.userWalletService.UpdateWalletBalanceWithTx(tx, userID, models.Coin, int64(item.Num), fmt.Sprintf("奖励包发放，奖励包ID：%d", item.PackageID)); err != nil {
					if localTx != nil {
						localTx.Rollback()
					}
					return nil, err
				}
			}
		default:
			// 未知类型，记录日志但不中断流程
			// 这里可以添加日志记录
		}
	}

	// 如果是本地事务，提交
	if localTx != nil {
		if err := localTx.Commit().Error; err != nil {
			return nil, err
		}
	}

	return record, nil
}
