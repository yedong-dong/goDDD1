package controllers

import (
	"goDDD1/services"
	"goDDD1/utils"
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
		utils.ResClientError(ctx, "无效的用户ID")
		return
	}

	userCurrencyFlow, err := c.userCurrencyFlowService.GetUserCurrencyFlow(strconv.Itoa(userID))
	if err != nil {
		utils.ResServerError(ctx, err)
		return
	}

	utils.ResSuccess(ctx, "获取用户货币流成功", userCurrencyFlow)
}

func (c *UserCurrencyFlowController) GetAllUserCurrencyFlow(ctx *gin.Context) {
	userCurrencyFlow, err := c.userCurrencyFlowService.GetAllUserCurrencyFlow()
	if err != nil {
		utils.ResServerError(ctx, err)
		return
	}

	utils.ResSuccess(ctx, "获取所有用户货币流成功", userCurrencyFlow)
}
