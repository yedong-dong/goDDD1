package controllers

import (
	"goDDD1/models"
	"goDDD1/services"
	"goDDD1/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserWalletController 用户钱包控制器
type UserWalletController struct {
	walletService services.UserWalletService
}

// NewUserWalletController 创建用户钱包控制器实例
func NewUserWalletController() *UserWalletController {
	return &UserWalletController{
		walletService: services.NewUserWalletService(),
	}
}

// GetUserWallets 获取用户所有钱包
func (c *UserWalletController) GetUserWallets(ctx *gin.Context) {
	userIDStr := ctx.Query("user_id")
	if userIDStr == "" {
		utils.ResClientError(ctx, "缺少user_id参数")
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.ResClientError(ctx, "无效的用户ID")
		return
	}

	wallets, err := c.walletService.GetUserWallets(uint(userID))
	if err != nil {
		utils.ResServerError(ctx, err)
		return
	}

	utils.ResSuccess(ctx, "获取钱包信息成功", wallets)
}

// GetWalletByType 根据类型获取用户钱包
func (c *UserWalletController) GetWalletByType(ctx *gin.Context) {
	userIDStr := ctx.Query("user_id")
	if userIDStr == "" {
		utils.ResClientError(ctx, "缺少user_id参数")
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.ResClientError(ctx, "无效的用户ID")
		return
	}

	walletTypeStr := ctx.Query("type")
	if walletTypeStr == "" {
		utils.ResClientError(ctx, "缺少type参数")
		return
	}

	walletType := models.WalletType(walletTypeStr)
	if walletType != models.Coin && walletType != models.Diamond {
		utils.ResClientError(ctx, "无效的钱包类型")
		return
	}

	wallet, err := c.walletService.GetWalletByUserIDAndType(uint(userID), walletType)
	if err != nil {
		utils.ResClientError(ctx, "钱包不存在")
		return
	}

	utils.ResSuccess(ctx, "获取钱包成功", wallet)
}

// UpdateWalletBalance 更新钱包余额
func (c *UserWalletController) UpdateWalletBalance(ctx *gin.Context) {
	var request struct {
		UserID     uint              `json:"user_id" binding:"required"`
		WalletType models.WalletType `json:"type" binding:"required"`
		Amount     int64             `json:"amount" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		utils.ResClientError(ctx, err.Error())
		return
	}

	if request.WalletType != models.Coin && request.WalletType != models.Diamond {
		utils.ResClientError(ctx, "无效的钱包类型")
		return
	}

	if err := c.walletService.UpdateWalletBalance(request.UserID, request.WalletType, request.Amount); err != nil {
		utils.ResServerError(ctx, err)
		return
	}

	wallets, err := NewUserWalletController().walletService.GetWalletByUserIDAndType(request.UserID, request.WalletType)
	if err != nil {
		utils.ResServerError(ctx, err)
		return
	}

	utils.ResSuccess(ctx, "钱包余额更新成功", wallets)
}
