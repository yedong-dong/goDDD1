package controllers

import (
	"goDDD1/services"
	"net/http"
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "缺少user_id参数",
		})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "user_id参数格式错误",
			"message": err.Error(),
		})
		return
	}

	user, err := c.levelService.GetUserLevel(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "获取用户等级信息失败",
			"message": err.Error(),
		})
		return
	}

	// 获取当前等级配置
	currentLevelConfig, err := c.levelService.GetLevelConfig(user.Level)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "获取等级配置失败",
			"message": err.Error(),
		})
		return
	}

	// 获取下一等级配置
	nextLevelConfig, err := c.levelService.GetLevelConfig(user.Level + 1)
	if err != nil {
		// 如果没有下一级，使用当前等级配置
		nextLevelConfig = currentLevelConfig
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 20000,
		"data": gin.H{
			"user_id":           user.UID,
			"username":          user.Username,
			"level":             user.Level,
			"experience":        user.Experience,
			"total_spent":       user.TotalSpent,
			"current_level_exp": currentLevelConfig.RequiredExp,
			"next_level_exp":    nextLevelConfig.RequiredExp,
			"discount_percent":  currentLevelConfig.DiscountPercent,
			"description":       currentLevelConfig.Description,
		},
		"message": "获取用户等级信息成功",
	})
}

// GetLevelHistory 获取用户等级历史记录
func (c *LevelController) GetLevelHistory(ctx *gin.Context) {
	userIDStr := ctx.Query("user_id")
	if userIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "缺少user_id参数",
		})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "user_id参数格式错误",
			"message": err.Error(),
		})
		return
	}

	histories, err := c.levelService.GetLevelHistory(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "获取用户等级历史记录失败",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 20000,
		"data": gin.H{
			"histories": histories,
			"total":     len(histories),
		},
		"message": "获取用户等级历史记录成功",
	})
}

// GetAllLevelConfigs 获取所有等级配置
func (c *LevelController) GetAllLevelConfigs(ctx *gin.Context) {
	configs, err := c.levelService.GetAllLevelConfigs()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "获取等级配置失败",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 20000,
		"data": gin.H{
			"configs": configs,
			"total":   len(configs),
		},
		"message": "获取等级配置成功",
	})
}
