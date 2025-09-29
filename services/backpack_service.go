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
	AddToBackpack(tx *gorm.DB, userID uint, storeID uint, quantity int64) error
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

func (b *backpackService) GetBackpackByUID(uid uint) (map[string]interface{}, error) {

	// 尝试从redis获取 - 修复缓存键格式化问题
	cacheKey := fmt.Sprintf(models.CacheKeyUserBackpack, uid)
	var cacheResult map[string]interface{}

	err := utils.GetHashField(cacheKey, "data", &cacheResult)
	if err == nil && cacheResult != nil {
		return cacheResult, nil
	}

	// 首先校验用户是否存在
	_, err = b.userService.GetUserByUID(uid)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New("用户不存在")
		}
		return nil, fmt.Errorf("查询用户背包失败: %w", err)
	}

	// 查询用户的所有背包物品
	var backpackItems []models.Backpack
	if err := config.Database.Where("user_id = ? AND quantity > 0", uid).Preload("Store").Find(&backpackItems).Error; err != nil {
		return nil, fmt.Errorf("查询用户背包物品失败: %w", err)
	}

	// 构建返回的map数据结构
	result := map[string]interface{}{
		"user_id":     uid,
		"total_items": len(backpackItems),
		"items":       buildItemsMap(backpackItems),
	}

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

// 辅助函数：构建物品映射
func buildItemsMap(items []models.Backpack) map[uint]map[string]interface{} {
	itemsMap := make(map[uint]map[string]interface{})
	for _, item := range items {
		itemsMap[item.StoreID] = map[string]interface{}{
			"backpack_id": item.ID,
			"store_id":    item.StoreID,
			"quantity":    item.Quantity,
			"name":        item.Store.Name,
		}
	}
	return itemsMap
}
