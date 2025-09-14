package services

import (
	"errors" // 添加这行
	"goDDD1/config"
	"goDDD1/models"
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
}

func NewStoreService() StoreService {
	return &storeService{}
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
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	//2、检查是否有该用户
	var user models.User
	if err := tx.Where("uid = ?", userID).First(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	//3、检查库存是否充足
	var store models.Store
	if err := tx.Where("id = ? and status = 1", storeID).First(&store).Error; err != nil {
		tx.Rollback()
		return err
	}
	if store.Stock < int64(num) {
		tx.Rollback()
		return errors.New("库存不足")
	}

	//4、检查wallet是否充足
	var wallet models.UserWallet
	if err := tx.Where("user_id = ? and type = ?", userID, store.CostType).First(&wallet).Error; err != nil {
		tx.Rollback()
		return err
	}
	if wallet.Num < store.Price*int64(num) {
		tx.Rollback()
		return errors.New("钱包不足")
	}

	//5、扣减余额
	wallet.Num -= store.Price * int64(num)
	if err := tx.Save(&wallet).Error; err != nil {
		tx.Rollback()
		return err
	}

	//6、扣减库存
	store.Stock -= int64(num)
	if err := tx.Save(&store).Error; err != nil {
		tx.Rollback()
		return err
	}

	//7、增加用户背包
	var bag models.Backpack
	if err := tx.Where("user_id = ? and store_id = ?", userID, storeID).FirstOrCreate(&bag, models.Backpack{
		UserID:   userID,
		StoreID:  storeID,
		Quantity: 0,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	bag.Quantity += int64(num)
	if err := tx.Save(&bag).Error; err != nil {
		tx.Rollback()
		return err
	}

	//8、添加交易流水
	var userCurrencyFlowService = NewUserCurrencyFlowService()
	if err := userCurrencyFlowService.CreateUserCurrencyFlow(&models.UserCurrencyFlow{
		UserID:   userID,
		StoreID:  storeID,
		CostType: string(store.CostType),
		Price:    -store.Price * int64(num),
	}); err != nil {
		return err
	}

	//9、提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil

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
