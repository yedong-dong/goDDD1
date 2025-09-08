package services

import (
	"errors"
	"goDDD1/config"
	"goDDD1/models"

	"github.com/jinzhu/gorm"
)

type BackpackService interface {
	GetBackpackByUID(uid uint) (map[string]interface{}, error)
}

type backpackService struct {
	userService UserService
}

func NewBackpackService() BackpackService {
	return &backpackService{
		userService: NewUserService(),
	}
}

func (b *backpackService) GetBackpackByUID(uid uint) (map[string]interface{}, error) {
	// 首先校验用户是否存在
	_, err := b.userService.GetUserByUID(uid)
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
			// "store_info": map[string]interface{}{
			// 	"name":      item.Store.Name,
			// 	"price":     item.Store.Price,
			// 	"cost_type": item.Store.CostType,
			// 	"status":    item.Store.Status,
			// },
		}
	}

	result["items"] = itemsMap

	// 如果背包为空，返回空的items map
	if len(backpackItems) == 0 {
		result["message"] = "背包为空"
	}

	return result, nil

}
