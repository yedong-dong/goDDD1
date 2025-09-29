package controllers

import (
	"goDDD1/models"
	"goDDD1/services"
	"goDDD1/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StoreController struct {
	storeService      services.StoreService
	backpackService   services.BackpackService
	userWalletService services.UserWalletService
}

func NewStoreController() *StoreController {
	return &StoreController{
		storeService:      services.NewStoreService(),
		backpackService:   services.NewBackpackService(),
		userWalletService: services.NewUserWalletService(),
	}
}

func (c *StoreController) CreateStore(ctx *gin.Context) {
	var store models.Store
	if err := ctx.ShouldBindJSON(&store); err != nil {
		utils.ResClientError(ctx, "JSON数据格式错误")
		return
	}
	if store.CostType != models.CostTypeCoin && store.CostType != models.CostTypeDiamond {
		utils.ResClientError(ctx, "CostTpye类型错误")
		return
	}

	if store.StoreType != models.StoreTypeGood && store.StoreType != models.StoreTypeGift {
		utils.ResClientError(ctx, "StoreType类型错误")
		return
	}
	if store.Tag != models.TagNormal && store.Tag != models.TagClothes && store.Tag != models.TagWeapon && store.Tag != models.TagArtifact && store.Tag != models.TagConsumable {
		utils.ResClientError(ctx, "Tag类型错误")
		return
	}

	if err := c.storeService.CreateStore(&store); err != nil {
		utils.ResServerError(ctx, err)
		return
	}
	utils.ResSuccess(ctx, "创建成功", gin.H{
		"store": store,
	})
}

func (c *StoreController) UpdateStore(ctx *gin.Context) {
	// 定义更新请求结构体
	type UpdateRequest struct {
		ID       *uint            `json:"id" binding:"required"`
		Name     *string          `json:"name,omitempty"`
		Price    *int64           `json:"price,omitempty"`
		Stock    *int64           `json:"stock,omitempty"`
		Status   *int             `json:"status,omitempty"`
		CostType *models.CostType `json:"cost_type,omitempty"`
	}

	var requestData UpdateRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		utils.ResClientError(ctx, "json数据格式错误")
		return
	}

	// 验证 ID
	if requestData.ID != nil && *requestData.ID <= 0 {
		utils.ResClientError(ctx, "id不能小于等于0")
		return
	}

	// 验证 Price
	if requestData.Price != nil && *requestData.Price < 0 {
		utils.ResClientError(ctx, "price不能小于0")
		return
	}

	// 验证 Stock
	if requestData.Stock != nil && *requestData.Stock < 0 {
		utils.ResClientError(ctx, "stock不能小于0")
		return
	}

	// 验证 CostType
	if requestData.CostType != nil &&
		*requestData.CostType != models.CostTypeCoin &&
		*requestData.CostType != models.CostTypeDiamond {
		utils.ResClientError(ctx, "cost_type必须是coin或diamond")
		return
	}

	// 验证 Status
	if requestData.Status != nil &&
		*requestData.Status != 0 &&
		*requestData.Status != 1 {
		utils.ResClientError(ctx, "status必须是0或1")
		return
	}

	// 首先查询现有的store
	storeID := strconv.Itoa(int(*requestData.ID))
	existingStore, err := c.storeService.GetStoreByID(storeID)
	if err != nil {
		utils.ResClientError(ctx, "指定的商店ID不存在")
		return
	}

	// 将UpdateRequest转换为models.Store
	// 只更新提供的字段，保留其他字段的原值
	if requestData.Name != nil {
		existingStore.Name = *requestData.Name
	}
	if requestData.Price != nil {
		existingStore.Price = *requestData.Price
	}
	if requestData.Stock != nil {
		existingStore.Stock = *requestData.Stock
	}
	if requestData.Status != nil {
		existingStore.Status = *requestData.Status
	}
	if requestData.CostType != nil {
		existingStore.CostType = *requestData.CostType
	}

	// 调用服务层更新
	updatedStore, err := c.storeService.UpdateStore(existingStore)
	if err != nil {
		utils.ResServerError(ctx, err)
		return
	}

	updateStoreDTO := updatedStore.ToStoreDTO()
	// 返回成功响应
	utils.ResSuccess(ctx, "商店更新成功", gin.H{
		"data": updateStoreDTO,
	})
}

func (c *StoreController) GetStoreByID(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		utils.ResClientError(ctx, "id不能为空")
		return
	}
	store, err := c.storeService.GetStoreByID(id)
	if err != nil {
		utils.ResServerError(ctx, err)
		return
	}

	storeDTO := store.ToStoreDTO()

	utils.ResSuccess(ctx, "获取成功", gin.H{
		"store": storeDTO,
	})
}

func (c *StoreController) GetStoreByTag(ctx *gin.Context) {
	tag := ctx.Query("tag")
	if tag == "" {
		utils.ResClientError(ctx, "tag不能为空")
		return
	}
	stores, err := c.storeService.GetStoreByTag(models.Tag(tag))
	if err != nil {
		utils.ResServerError(ctx, err)
		return
	}
	utils.ResSuccess(ctx, "获取成功", gin.H{
		"stores": stores,
	})
}

func (c *StoreController) GetStoreByTagPage(ctx *gin.Context) {
	tag := ctx.Query("tag")
	if tag == "" {
		utils.ResClientError(ctx, "tag不能为空")
		return
	}
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		utils.ResClientError(ctx, "page必须是整数")
		return
	}
	pageSize, err := strconv.Atoi(ctx.Query("page_size"))
	if err != nil {
		utils.ResClientError(ctx, "page_size必须是整数")
		return
	}
	stores, total, err := c.storeService.GetStoreByTagPage(models.Tag(tag), page, pageSize)
	if err != nil {
		utils.ResServerError(ctx, err)
		return
	}
	utils.ResSuccess(ctx, "获取成功", gin.H{
		"stores": stores,
		"total":  total,
	})
}

func (c *StoreController) GetAllStores(ctx *gin.Context) {
	stores, err := c.storeService.GetAllStores()
	if err != nil {
		utils.ResServerError(ctx, err)
		return
	}

	utils.ResSuccess(ctx, "获取商店列表成功", gin.H{
		"stores": stores,
		"total":  len(stores),
	})
}

func (c *StoreController) BuyGoods(ctx *gin.Context) {
	// 定义购买请求结构体
	type BuyRequest struct {
		UserID  uint `json:"user_id" binding:"required"`
		StoreID uint `json:"store_id" binding:"required"`
		Num     uint `json:"num" binding:"required"`
	}

	var requestData BuyRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		utils.ResClientError(ctx, "json数据格式错误")
		return
	}

	if requestData.Num <= 0 {
		utils.ResClientError(ctx, "num不能小于等于0")
		return
	}

	err := c.storeService.BuyGoods(requestData.UserID, requestData.StoreID, requestData.Num)
	if err != nil {
		utils.ResServerError(ctx, err)
		return
	}

	utils.ResSuccess(ctx, "购买成功", nil)
}
