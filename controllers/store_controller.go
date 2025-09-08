package controllers

import (
	"goDDD1/models"
	"goDDD1/services"
	"net/http"
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "JSON数据格式错误",
			"message": err.Error(),
		})
		return
	}
	if store.CostType != models.CostTypeCoin && store.CostType != models.CostTypeDiamond {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "JSON数据格式错误",
			"message": "CostTpye类型错误",
		})
		return
	}

	if store.StoreType != models.StoreTypeGood && store.StoreType != models.StoreTypeGift {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "JSON数据格式错误",
			"message": "StoreType类型错误",
		})
		return
	}

	if err := c.storeService.CreateStore(&store); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "创建失败",
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "创建成功",
		"store":   store,
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "json数据格式错误",
			"msg":   err.Error(),
		})
		return
	}

	// 验证 ID
	if requestData.ID != nil && *requestData.ID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "无效id",
			"msg":   "id不能小于等于0",
		})
		return
	}

	// 验证 Price
	if requestData.Price != nil && *requestData.Price < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "无效price",
			"msg":   "price不能小于0",
		})
		return
	}

	// 验证 Stock
	if requestData.Stock != nil && *requestData.Stock < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "无效stock",
			"msg":   "stock不能小于0",
		})
		return
	}

	// 验证 CostType
	if requestData.CostType != nil &&
		*requestData.CostType != models.CostTypeCoin &&
		*requestData.CostType != models.CostTypeDiamond {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "无效cost_type",
			"msg":   "cost_type必须是coin或diamond",
		})
		return
	}

	// 验证 Status
	if requestData.Status != nil &&
		*requestData.Status != 0 &&
		*requestData.Status != 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "无效status",
			"msg":   "status必须是0或1",
		})
		return
	}

	// 首先查询现有的store
	storeID := strconv.Itoa(int(*requestData.ID))
	existingStore, err := c.storeService.GetStoreByID(storeID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "商店不存在",
			"msg":   "指定的商店ID不存在",
		})
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
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新失败",
			"msg":   err.Error(),
		})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message": "商店更新成功",
		"data":    updatedStore,
	})
}

func (c *StoreController) GetStoreByID(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "无效id",
			"msg":   "id不能为空",
		})
		return
	}
	store, err := c.storeService.GetStoreByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取失败",
			"msg":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"store":   store,
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "json数据格式错误",
			"msg":   err.Error(),
		})
		return
	}

	if requestData.Num <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "无效num",
			"msg":   "num不能小于等于0",
		})
		return
	}

	err := c.storeService.BuyGoods(requestData.UserID, requestData.StoreID, requestData.Num)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "购买失败",
			"msg":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "购买成功",
	})

}
