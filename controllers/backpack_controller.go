package controllers

import (
	"goDDD1/services"
	"goDDD1/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BackpackController struct {
	backpackService services.BackpackService
	storeService    services.StoreService
	userService     services.UserWalletService
}

func NewBackpackController() *BackpackController {
	return &BackpackController{
		backpackService: services.NewBackpackService(),
		storeService:    services.NewStoreService(),
		userService:     services.NewUserWalletService(),
	}
}

func (c *BackpackController) GetBackpack(ctx *gin.Context) {
	uid, err := strconv.Atoi(ctx.Query("uid"))
	if err != nil {
		utils.ResClientError(ctx, "UID格式不正确")
		return
	}

	backpackData, err := c.backpackService.GetBackpackByUID(uint(uid))
	if err != nil {
		// 根据错误类型返回不同的响应
		if err.Error() == "用户不存在" {
			utils.ResClientError(ctx, "用户不存在")
		} else {
			utils.ResServerError(ctx, err)
		}
		return
	}

	utils.ResSuccess(ctx, "获取背包成功", gin.H{
		"backpack": backpackData,
	})
}
