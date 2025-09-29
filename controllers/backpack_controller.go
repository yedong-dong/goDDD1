package controllers

import (
	"goDDD1/services"
	"net/http"
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "参数错误",
			"message": "UID格式不正确",
		})
		return
	}

	backpackData, err := c.backpackService.GetBackpackByUID(uint(uid))
	if err != nil {
		// 根据错误类型返回不同的HTTP状态码
		switch err.Error() {
		case "用户不存在":
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    "50000",
				"error":   "用户不存在",
				"message": "用户不存在",
			})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "服务器内部错误",
				"message": err.Error(),
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": "20000",
		"data": gin.H{
			"backpack": backpackData,
		},
	})
}
