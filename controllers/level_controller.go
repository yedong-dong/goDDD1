package controllers

import (
	"goDDD1/services"
	"goDDD1/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// LevelController 等级控制器
type LevelController struct {
	levelService services.LevelService
}

// NewLevelController 创建等级控制器
func NewLevelController() *LevelController {
	return &LevelController{
		levelService: services.NewLevelService(),
	}
}

// GetUserLevel 获取用户当前等级信息
func (c *LevelController) GetUserLevel(ctx *gin.Context) {
	userIDStr := ctx.Query("user_id")
	if userIDStr == "" {
		utils.ResClientError(ctx, "缺少user_id参数")
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		utils.ResClientError(ctx, "user_id参数格式错误")
		return
	}

	user, err := c.levelService.GetUserLevel(uint(userID))
	if err != nil {
		utils.ResServerError(ctx, err)
		return
	}

	// 获取当前等级配置
	currentLevelConfig, err := c.levelService.GetLevelConfig(user.Level)
	if err != nil {
		utils.ResServerError(ctx, err)
		return
	}

	// 获取下一等级配置
	nextLevelConfig, err := c.levelService.GetLevelConfig(user.Level + 1)
	if err != nil {
		// 如果没有下一级，使用当前等级配置
		nextLevelConfig = currentLevelConfig
	}

	utils.ResSuccess(ctx, "获取用户等级信息成功", gin.H{
		"user_id":           user.UID,
		"username":          user.Username,
		"level":             user.Level,
		"experience":        user.Experience,
		"total_spent":       user.TotalSpent,
		"current_level_exp": currentLevelConfig.RequiredExp,
		"next_level_exp":    nextLevelConfig.RequiredExp,
		"discount_percent":  currentLevelConfig.DiscountPercent,
		"description":       currentLevelConfig.Description,
	})
}

// GetLevelHistory 获取用户等级历史记录
func (c *LevelController) GetLevelHistory(ctx *gin.Context) {
	userIDStr := ctx.Query("user_id")
	if userIDStr == "" {
		utils.ResClientError(ctx, "缺少user_id参数")
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		utils.ResClientError(ctx, "user_id参数格式错误")
		return
	}

	histories, err := c.levelService.GetLevelHistory(uint(userID))
	if err != nil {
		utils.ResServerError(ctx, err)
		return
	}

	utils.ResSuccess(ctx, "获取用户等级历史记录成功", gin.H{
		"histories": histories,
		"total":     len(histories),
	})
}

// GetAllLevelConfigs 获取所有等级配置
func (c *LevelController) GetAllLevelConfigs(ctx *gin.Context) {
	configs, err := c.levelService.GetAllLevelConfigs()
	if err != nil {
		utils.ResServerError(ctx, err)
		return
	}

	utils.ResSuccess(ctx, "获取等级配置成功", gin.H{
		"configs": configs,
		"total":   len(configs),
	})
}
