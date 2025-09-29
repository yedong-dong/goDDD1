package controllers

import (
	"goDDD1/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserCurrencyFlowController struct {
	userCurrencyFlowService services.UserCurrencyFlowServiceInterface
}

// NewUserCurrencyFlowController 创建用户货币流控制器实例
func NewUserCurrencyFlowController() *UserCurrencyFlowController {
	return &UserCurrencyFlowController{
		userCurrencyFlowService: services.NewUserCurrencyFlowService(),
	}
}

func (c *UserCurrencyFlowController) GetUserCurrencyFlow(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Query("user_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	userCurrencyFlow, err := c.userCurrencyFlowService.GetUserCurrencyFlow(strconv.Itoa(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户货币流失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user_currency_flow": userCurrencyFlow})
}

func (c *UserCurrencyFlowController) GetAllUserCurrencyFlow(ctx *gin.Context) {
	userCurrencyFlow, err := c.userCurrencyFlowService.GetAllUserCurrencyFlow()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取所有用户货币流失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": "20000",
		"data": gin.H{
			"user_currency_flow": userCurrencyFlow,
			"message":            "获取所有用户货币流成功",
		},
	})
}
