package services

import (
	"database/sql"
	"errors" // 添加这行
	"fmt"
	"goDDD1/config"
	"goDDD1/models"
	"goDDD1/utils"
	"log"

	"github.com/jinzhu/gorm"
)

type StoreService interface {
	CreateStore(store *models.Store) error
	GetStoreByID(id string) (*models.Store, error)
	UpdateStore(store *models.Store) (*models.Store, error) // 修改方法签名
	BuyGoods(userID uint, storeID uint, num uint) error
	GetStoreByTag(tag models.Tag) ([]*models.StoreDTO, error)
	GetStoreByTagPage(tag models.Tag, page, pageSize int) ([]*models.StoreDTO, int64, error)
	GetAllStores() ([]*models.StoreDTO, error)
}

type storeService struct {
	levelService LevelService
}

func NewStoreService() StoreService {
	return &storeService{
		levelService: NewLevelService(),
	}
}

func (s *storeService) GetStoreByTag(tag models.Tag) ([]*models.StoreDTO, error) {
	var stores []*models.Store
	result := config.Database.Where("tag = ?", tag).Find(&stores)
	if result.Error != nil {
		return nil, result.Error
	}

	// 转换为 DTO
	storeDTOs := make([]*models.StoreDTO, len(stores))
	for i, store := range stores {
		storeDTOs[i] = store.ToStoreDTO()
	}

	return storeDTOs, nil
}

func (s *storeService) CreateStore(store *models.Store) error {
	tx := config.Database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(store).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
func (s *storeService) GetStoreByID(id string) (*models.Store, error) {
	var store models.Store
	result := config.Database.Raw("Select * from stores where id = ?", id).Scan(&store)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("store not found")
	}

	return &store, nil
}

// 简化的UpdateStore方法
func (s *storeService) UpdateStore(store *models.Store) (*models.Store, error) {
	// 开始事务
	tx := config.Database.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// 确保事务回滚
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 直接保存store
	if err := tx.Save(store).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return store, nil
}
func (s *storeService) BuyGoods(userID uint, storeID uint, num uint) error {
	//1、开始事务
	//2、检查是否有该用户
	//3、检查库存是否充足
	//4、检查wallet是否充足
	//5、扣减余额
	//6、扣减库存
	//7、增加用户背包
	//8、提交事务
	tx := config.Database.Begin()
	defer func() {
		// 使用recover确保在panic时事务被回滚
		if r := recover(); r != nil {
			// 使用SafeRollback避免重复回滚错误
			SafeRollback(tx)
		}
	}()

	//2、检查是否有该用户
	var user models.User
	if err := tx.Where("uid = ?", userID).First(&user).Error; err != nil {
		SafeRollback(tx)
		return err
	}

	//3、检查库存是否充足
	var store models.Store
	if err := tx.Where("id = ? and status = 1", storeID).First(&store).Error; err != nil {
		SafeRollback(tx)
		return err
	}
	if store.Stock < int64(num) {
		SafeRollback(tx)
		return errors.New("库存不足")
	}

	//4、检查wallet是否充足
	var wallet models.UserWallet
	if err := tx.Where("user_id = ? and type = ?", userID, store.CostType).First(&wallet).Error; err != nil {
		SafeRollback(tx)
		return err
	}

	originalPrice := store.Price * int64(num)
	discountPrice, err := s.levelService.CalculateDiscountPrice(user.UID, uint(originalPrice))
	if err != nil {
		// 如果计算折扣失败，使用原价
		discountPrice = uint(originalPrice)
	}

	if wallet.Num < int64(discountPrice) {
		SafeRollback(tx)
		return errors.New("钱包不足")
	}

	//5、扣减余额
	wallet.Num -= store.Price * int64(num)
	if err := tx.Save(&wallet).Error; err != nil {
		SafeRollback(tx)
		return err
	}

	//8、添加交易流水 - 修改为使用事务
	var userCurrencyFlow = models.UserCurrencyFlow{
		UserID:   userID,
		StoreID:  storeID,
		CostType: string(store.CostType),
		Price:    -store.Price * int64(num),
	}
	if err := tx.Create(&userCurrencyFlow).Error; err != nil {
		SafeRollback(tx)
		return err
	}

	//6、扣减库存
	store.Stock -= int64(num)
	if err := tx.Save(&store).Error; err != nil {
		SafeRollback(tx)
		return err
	}

	//7、增加用户背包
	var bag models.Backpack
	if err := tx.Where("user_id = ? and store_id = ?", userID, storeID).FirstOrCreate(&bag, models.Backpack{
		UserID:   userID,
		StoreID:  storeID,
		Quantity: 0,
	}).Error; err != nil {
		SafeRollback(tx)
		return err
	}

	// 删除缓存记录
	cacheKey := fmt.Sprintf(models.CacheKeyUserBackpack, userID)
	err = utils.DelHashField(cacheKey, "data")
	if err == nil {
		log.Printf("successful delete cacheKey: %s backpack", cacheKey)
	}

	// 修改BuyGoods方法中的经验值增加部分
	//7.1、 增加经验值
	switch store.CostType {
	case "coin":
		//增加经验值 - 使用levelService处理
		expToAdd := uint(store.Price*int64(num)) / 2
		description := fmt.Sprintf("购买%s商品:%s, 价格为:%d", store.CostType, store.Name, store.Price*int64(num))

		// 使用levelService处理经验值增加和可能的升级
		if _, err := s.levelService.AddExpeirence(tx, userID, expToAdd, description); err != nil {
			SafeRollback(tx)
			return err
		}

	case "diamond":
		//增加经验值 - 使用levelService处理
		expToAdd := uint(store.Price * int64(num))
		description := fmt.Sprintf("购买%s商品:%s, 价格为:%d", store.CostType, store.Name, store.Price*int64(num))

		// 使用levelService处理经验值增加和可能的升级
		if _, err := s.levelService.AddExpeirence(tx, userID, expToAdd, description); err != nil {
			SafeRollback(tx)
			return err
		}
	}

	bag.Quantity += int64(num)
	if err := tx.Save(&bag).Error; err != nil {
		SafeRollback(tx)
		return err
	}

	//9、提交事务
	return tx.Commit().Error
}

// SafeRollback 安全回滚事务，忽略"已回滚"错误
func SafeRollback(tx *gorm.DB) {
	err := tx.Rollback().Error
	if err != nil && err != sql.ErrTxDone {
		// 可以记录日志，但不要返回错误
		fmt.Println("事务回滚失败:", err)
	}
}

func (s *storeService) GetStoreByTagPage(tag models.Tag, page, pageSize int) ([]*models.StoreDTO, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	} else if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	//查询总数
	var total int64
	if err := config.Database.Model(&models.Store{}).Where("tag = ?", tag).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var stores []*models.Store
	result := config.Database.Where("tag = ?", tag).Offset(offset).Limit(pageSize).Find(&stores)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	// 转换为 DTO
	storeDTOs := make([]*models.StoreDTO, len(stores))
	for i, store := range stores {
		storeDTOs[i] = store.ToStoreDTO()
	}

	return storeDTOs, total, nil

}

func (s *storeService) GetAllStores() ([]*models.StoreDTO, error) {
	var stores []*models.Store
	result := config.Database.Find(&stores)
	if result.Error != nil {
		return nil, result.Error
	}

	// 转换为 DTO
	storeDTOs := make([]*models.StoreDTO, len(stores))
	for i, store := range stores {
		storeDTOs[i] = store.ToStoreDTO()
	}

	return storeDTOs, nil
}
