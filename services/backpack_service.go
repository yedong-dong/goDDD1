package services

import (
	"errors"
	"fmt"
	"goDDD1/config"
	"goDDD1/models"
	"goDDD1/utils"
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

type BackpackService interface {
	// 基础查询方法
	GetBackpackByUID(uid uint) (map[string]interface{}, error)

	// 增删改查基础方法
	AddToBackpack(tx *gorm.DB, userID uint, storeID uint, quantity int64) error
	UpdateBackpackQuantity(tx *gorm.DB, backpackID uint, quantity int64) error
	GetBackpackItem(userID uint, storeID uint) (*models.Backpack, error)
	GetBackpackItemByID(backpackID uint) (*models.Backpack, error)
	ListBackpackItems(userID uint, page, pageSize int) ([]*models.Backpack, int64, error)
	DeleteBackpackItem(backpackID uint) error
	DeleteUserBackpack(userID uint) error

	// 新增实用方法
	GetUserBackpackItems(userID uint) ([]*models.Backpack, error)
	GetUserBackpackItemsByType(userID uint, itemType string) ([]*models.Backpack, error)
	ConsumeBackpackItem(tx *gorm.DB, userID uint, storeID uint, quantity int64) error
	BatchAddToBackpack(tx *gorm.DB, userID uint, items map[uint]int64) error
	SearchUserBackpackItems(userID uint, keyword string, page, pageSize int) ([]*models.Backpack, int64, error)
	GetUserBackpackSummary(userID uint) (map[string]interface{}, error)
	TransferBackpackItem(tx *gorm.DB, fromUserID, toUserID, storeID uint, quantity int64) error
}

// 修改结构体，添加RewardFlowService依赖

type backpackService struct {
	userService       UserService
	rewardFlowService RewardFlowService
}

func NewBackpackService() BackpackService {
	return &backpackService{
		userService:       NewUserService(),
		rewardFlowService: NewRewardFlowService(),
	}
}

// GetBackpackByUID 获取用户背包（已有方法，保持不变）
func (b *backpackService) GetBackpackByUID(uid uint) (map[string]interface{}, error) {

	// 尝试从redis获取 - 修复缓存键格式化问题
	cacheKey := fmt.Sprintf(models.CacheKeyUserBackpack, uid)
	var cacheResult map[string]interface{}
	err := utils.GetHashField(cacheKey, "data", &cacheResult)
	if err != nil {
		log.Printf("Error fetching backpack from cache: %v", err)
	}

	// 如果缓存存在，直接返回
	if cacheResult != nil {
		return cacheResult, nil
	}

	// 首先校验用户是否存在
	_, err = b.userService.GetUserByUID(uid)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	// 查询用户的所有背包物品
	var backpackItems []models.Backpack
	if err := config.Database.Where("user_id = ? AND quantity > 0", uid).Preload("Store").Find(&backpackItems).Error; err != nil {
		return nil, err
	}

	// 构建返回的map数据结构
	result := map[string]interface{}{
		"user_id":     uid,
		"total_items": len(backpackItems),
		"items":       make(map[uint]map[string]interface{}),
	}

	// 将背包物品转换为map格式
	itemsMap := make(map[uint]map[string]interface{})
	for _, item := range backpackItems {
		itemsMap[item.StoreID] = map[string]interface{}{
			"backpack_id": item.ID,
			"store_id":    item.StoreID,
			"quantity":    item.Quantity,
			"name":        item.Store.Name,
		}
	}

	result["items"] = itemsMap

	// 如果背包为空，返回空的items map
	if len(backpackItems) == 0 {
		result["message"] = "背包为空"
	}

	// 修复缓存键格式化问题，并使用固定的哈希字段名"data"
	utils.SetHashField(cacheKey, "data", result, time.Hour)

	return result, nil
}

func (b *backpackService) AddToBackpack(tx *gorm.DB, userID uint, storeID uint, quantity int64) error {
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

	// 查询用户是否存在
	var user models.User
	if err := tx.Raw("select * from users where uid = ? and is_deleted = ?", userID, 0).Scan(&user).Error; err != nil {
		if localTx != nil {
			localTx.Rollback()
		}
		return errors.New("用户不存在")
	}

	// 查询商品是否存在
	var store models.Store
	if err := tx.Raw("select * from stores where id = ? ", storeID).Scan(&store).Error; err != nil {
		if localTx != nil {
			localTx.Rollback()
		}
		return errors.New("商品不存在")
	}

	// 查询用户背包中是否已有该物品
	var backpack models.Backpack
	err := tx.Raw("select * from backpacks where user_id = ? and store_id = ? ", userID, storeID).Scan(&backpack).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			// 如果不存在，创建新记录
			backpack = models.Backpack{
				UserID:   userID,
				StoreID:  storeID,
				Quantity: quantity,
			}
			if err := tx.Create(&backpack).Error; err != nil {
				if localTx != nil {
					localTx.Rollback()
				}
				return err
			}
		} else {
			// 其他错误
			if localTx != nil {
				localTx.Rollback()
			}
			return err
		}
	} else {
		// 如果存在，更新数量
		backpack.Quantity += quantity
		if err := tx.Save(&backpack).Error; err != nil {
			if localTx != nil {
				localTx.Rollback()
			}
			return err
		}
	}

	// 添加奖励流水记录
	err = b.rewardFlowService.CreateRewardFlow(tx, userID, models.RewardTypeItem, storeID, quantity, "背包添加物品")
	if err != nil {
		if localTx != nil {
			localTx.Rollback()
		}
		return err
	}

	// 删除缓存记录
	cacheKey := fmt.Sprintf(models.CacheKeyUserBackpack, userID)
	err = utils.DelHashField(cacheKey, "data")
	if err == nil {
		log.Printf("successful delete cacheKey: %s backpack", cacheKey)
	}

	// 如果是本地事务，提交
	if localTx != nil {
		return localTx.Commit().Error
	}

	return nil
}

// UpdateBackpackQuantity 更新背包物品数量（已有方法，保持不变）
func (b *backpackService) UpdateBackpackQuantity(tx *gorm.DB, backpackID uint, quantity int64) error {
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

	// 查询背包物品
	var backpack models.Backpack
	if err := tx.First(&backpack, backpackID).Error; err != nil {
		if localTx != nil {
			localTx.Rollback()
		}
		return errors.New("背包物品不存在")
	}

	// 更新数量
	backpack.Quantity = quantity
	if err := tx.Save(&backpack).Error; err != nil {
		if localTx != nil {
			localTx.Rollback()
		}
		return err
	}

	// 如果是本地事务，提交
	if localTx != nil {
		return localTx.Commit().Error
	}

	return nil
}

// GetBackpackItem 获取用户特定物品（已有方法，保持不变）
func (b *backpackService) GetBackpackItem(userID uint, storeID uint) (*models.Backpack, error) {
	var backpack models.Backpack
	err := config.Database.Where("user_id = ? AND store_id = ?", userID, storeID).Preload("Store").First(&backpack).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New("背包中没有该物品")
		}
		return nil, err
	}
	return &backpack, nil
}

// GetBackpackItemByID 根据ID获取背包物品（已有方法，保持不变）
func (b *backpackService) GetBackpackItemByID(backpackID uint) (*models.Backpack, error) {
	var backpack models.Backpack
	err := config.Database.Preload("Store").First(&backpack, backpackID).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New("背包物品不存在")
		}
		return nil, err
	}
	return &backpack, nil
}

// ListBackpackItems 分页获取用户背包物品列表（已有方法，保持不变）
func (b *backpackService) ListBackpackItems(userID uint, page, pageSize int) ([]*models.Backpack, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	} else if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	var backpackItems []*models.Backpack
	var total int64

	if err := config.Database.Model(&models.Backpack{}).Where("user_id = ? AND quantity > 0", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := config.Database.Where("user_id = ? AND quantity > 0", userID).Preload("Store").Offset(offset).Limit(pageSize).Order("id desc").Find(&backpackItems).Error; err != nil {
		return nil, 0, err
	}

	return backpackItems, total, nil
}

// DeleteBackpackItem 删除背包物品（已有方法，保持不变）
func (b *backpackService) DeleteBackpackItem(backpackID uint) error {
	return config.Database.Delete(&models.Backpack{}, backpackID).Error
}

// DeleteUserBackpack 删除用户所有背包物品（已有方法，保持不变）
func (b *backpackService) DeleteUserBackpack(userID uint) error {
	return config.Database.Where("user_id = ?", userID).Delete(&models.Backpack{}).Error
}

// GetUserBackpackItems 获取用户所有背包物品（不分页）
func (b *backpackService) GetUserBackpackItems(userID uint) ([]*models.Backpack, error) {
	var backpackItems []*models.Backpack
	if err := config.Database.Where("user_id = ? AND quantity > 0", userID).Preload("Store").Find(&backpackItems).Error; err != nil {
		return nil, err
	}
	return backpackItems, nil
}

// GetUserBackpackItemsByType 根据商品类型获取用户背包物品
func (b *backpackService) GetUserBackpackItemsByType(userID uint, itemType string) ([]*models.Backpack, error) {
	var backpackItems []*models.Backpack
	if err := config.Database.Joins("JOIN stores ON backpacks.store_id = stores.id").
		Where("backpacks.user_id = ? AND backpacks.quantity > 0 AND stores.type = ?", userID, itemType).
		Preload("Store").
		Find(&backpackItems).Error; err != nil {
		return nil, err
	}
	return backpackItems, nil
}

// ConsumeBackpackItem 消耗背包物品（适用于使用道具等场景）
func (b *backpackService) ConsumeBackpackItem(tx *gorm.DB, userID uint, storeID uint, quantity int64) error {
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

	// 查询用户背包中是否有该物品
	var backpack models.Backpack
	err := tx.Where("user_id = ? AND store_id = ?", userID, storeID).First(&backpack).Error
	if err != nil {
		if localTx != nil {
			localTx.Rollback()
		}
		if gorm.IsRecordNotFoundError(err) {
			return errors.New("背包中没有该物品")
		}
		return err
	}

	// 检查数量是否足够
	if backpack.Quantity < quantity {
		if localTx != nil {
			localTx.Rollback()
		}
		return errors.New("物品数量不足")
	}

	// 更新数量
	backpack.Quantity -= quantity
	if err := tx.Save(&backpack).Error; err != nil {
		if localTx != nil {
			localTx.Rollback()
		}
		return err
	}

	// 如果是本地事务，提交
	if localTx != nil {
		return localTx.Commit().Error
	}

	return nil
}

// BatchAddToBackpack 批量添加物品到背包（适用于批量奖励等场景）
func (b *backpackService) BatchAddToBackpack(tx *gorm.DB, userID uint, items map[uint]int64) error {
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

	// 查询用户是否存在
	var user models.User
	if err := tx.First(&user, userID).Error; err != nil {
		if localTx != nil {
			localTx.Rollback()
		}
		return errors.New("用户不存在")
	}

	// 批量添加物品
	for storeID, quantity := range items {
		// 查询商品是否存在
		var store models.Store
		if err := tx.First(&store, storeID).Error; err != nil {
			if localTx != nil {
				localTx.Rollback()
			}
			return errors.New("商品ID " + string(storeID) + " 不存在")
		}

		// 查询用户背包中是否已有该物品
		var backpack models.Backpack
		err := tx.Where("user_id = ? AND store_id = ?", userID, storeID).First(&backpack).Error

		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				// 如果不存在，创建新记录
				backpack = models.Backpack{
					UserID:   userID,
					StoreID:  storeID,
					Quantity: quantity,
				}
				if err := tx.Create(&backpack).Error; err != nil {
					if localTx != nil {
						localTx.Rollback()
					}
					return err
				}
			} else {
				// 其他错误
				if localTx != nil {
					localTx.Rollback()
				}
				return err
			}
		} else {
			// 如果存在，更新数量
			backpack.Quantity += quantity
			if err := tx.Save(&backpack).Error; err != nil {
				if localTx != nil {
					localTx.Rollback()
				}
				return err
			}
		}
	}

	// 如果是本地事务，提交
	if localTx != nil {
		return localTx.Commit().Error
	}

	return nil
}

// SearchUserBackpackItems 搜索用户背包物品（适用于搜索功能）
func (b *backpackService) SearchUserBackpackItems(userID uint, keyword string, page, pageSize int) ([]*models.Backpack, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	} else if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	var backpackItems []*models.Backpack
	var total int64

	// 使用JOIN查询，根据商品名称搜索
	query := config.Database.Table("backpacks").
		Joins("JOIN stores ON backpacks.store_id = stores.id").
		Where("backpacks.user_id = ? AND backpacks.quantity > 0 AND stores.name LIKE ?", userID, "%"+keyword+"%")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Order("backpacks.id desc").
		Preload("Store").Find(&backpackItems).Error; err != nil {
		return nil, 0, err
	}

	return backpackItems, total, nil
}

// GetUserBackpackSummary 获取用户背包摘要信息（适用于统计分析）
func (b *backpackService) GetUserBackpackSummary(userID uint) (map[string]interface{}, error) {
	// 查询用户是否存在
	_, err := b.userService.GetUserByUID(userID)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	// 查询用户背包物品总数
	var totalItems int64
	if err := config.Database.Model(&models.Backpack{}).Where("user_id = ? AND quantity > 0", userID).Count(&totalItems).Error; err != nil {
		return nil, err
	}

	// 查询用户背包物品总价值
	type Result struct {
		TotalValue int64
	}
	var result Result
	if err := config.Database.Table("backpacks").
		Joins("JOIN stores ON backpacks.store_id = stores.id").
		Where("backpacks.user_id = ? AND backpacks.quantity > 0", userID).
		Select("SUM(backpacks.quantity * stores.price) as total_value").
		Scan(&result).Error; err != nil {
		return nil, err
	}

	// 查询用户背包物品类型分布
	type TypeDistribution struct {
		Type  string
		Count int64
	}
	var typeDistribution []TypeDistribution
	if err := config.Database.Table("backpacks").
		Joins("JOIN stores ON backpacks.store_id = stores.id").
		Where("backpacks.user_id = ? AND backpacks.quantity > 0", userID).
		Select("stores.type, COUNT(*) as count").
		Group("stores.type").
		Scan(&typeDistribution).Error; err != nil {
		return nil, err
	}

	// 构建返回结果
	typeMap := make(map[string]int64)
	for _, item := range typeDistribution {
		typeMap[item.Type] = item.Count
	}

	return map[string]interface{}{
		"user_id":           userID,
		"total_items":       totalItems,
		"total_value":       result.TotalValue,
		"type_distribution": typeMap,
	}, nil
}

// TransferBackpackItem 转移背包物品（适用于玩家之间交易）
func (b *backpackService) TransferBackpackItem(tx *gorm.DB, fromUserID, toUserID, storeID uint, quantity int64) error {
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

	// 检查两个用户是否存在
	var fromUser, toUser models.User
	if err := tx.First(&fromUser, fromUserID).Error; err != nil {
		if localTx != nil {
			localTx.Rollback()
		}
		return errors.New("转出用户不存在")
	}
	if err := tx.First(&toUser, toUserID).Error; err != nil {
		if localTx != nil {
			localTx.Rollback()
		}
		return errors.New("转入用户不存在")
	}

	// 检查商品是否存在
	var store models.Store
	if err := tx.First(&store, storeID).Error; err != nil {
		if localTx != nil {
			localTx.Rollback()
		}
		return errors.New("商品不存在")
	}

	// 检查转出用户是否有足够的物品
	var fromBackpack models.Backpack
	if err := tx.Where("user_id = ? AND store_id = ?", fromUserID, storeID).First(&fromBackpack).Error; err != nil {
		if localTx != nil {
			localTx.Rollback()
		}
		if gorm.IsRecordNotFoundError(err) {
			return errors.New("转出用户背包中没有该物品")
		}
		return err
	}

	if fromBackpack.Quantity < quantity {
		if localTx != nil {
			localTx.Rollback()
		}
		return errors.New("转出用户物品数量不足")
	}

	// 扣减转出用户物品数量
	fromBackpack.Quantity -= quantity
	if err := tx.Save(&fromBackpack).Error; err != nil {
		if localTx != nil {
			localTx.Rollback()
		}
		return err
	}

	// 增加转入用户物品数量
	var toBackpack models.Backpack
	err := tx.Where("user_id = ? AND store_id = ?", toUserID, storeID).First(&toBackpack).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			// 如果转入用户没有该物品，创建新记录
			toBackpack = models.Backpack{
				UserID:   toUserID,
				StoreID:  storeID,
				Quantity: quantity,
			}
			if err := tx.Create(&toBackpack).Error; err != nil {
				if localTx != nil {
					localTx.Rollback()
				}
				return err
			}
		} else {
			// 其他错误
			if localTx != nil {
				localTx.Rollback()
			}
			return err
		}
	} else {
		// 如果转入用户已有该物品，更新数量
		toBackpack.Quantity += quantity
		if err := tx.Save(&toBackpack).Error; err != nil {
			if localTx != nil {
				localTx.Rollback()
			}
			return err
		}
	}

	// 如果是本地事务，提交
	if localTx != nil {
		return localTx.Commit().Error
	}

	return nil
}
