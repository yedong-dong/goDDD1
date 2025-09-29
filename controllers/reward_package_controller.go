package controllers

import (
	"goDDD1/models"
	"goDDD1/services"
	"goDDD1/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RewardPackageController struct {
	rewardPackageService services.RewardPackageService
}

func NewRewardPackageController() *RewardPackageController {
	return &RewardPackageController{
		rewardPackageService: services.NewRewardPackageService(),
	}
}

// CreateRewardPackage 创建奖励包
func (c *RewardPackageController) CreateRewardPackage(ctx *gin.Context) {
	// 定义请求结构体
	type CreateRequest struct {
		Package models.RewardPackage        `json:"package" binding:"required"`
		Items   []*models.RewardPackageItem `json:"items" binding:"required"`
	}

	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResClientError(ctx, "JSON数据格式错误")
		return
	}

	// 验证奖励包内容
	if len(req.Items) == 0 {
		utils.ResClientError(ctx, "奖励包内容不能为空")
		return
	}

	// 创建奖励包
	if err := c.rewardPackageService.CreateRewardPackage(&req.Package, req.Items); err != nil {
		utils.ResServerError(ctx, err)
		return
	}

	utils.ResSuccess(ctx, "创建成功", gin.H{
		"package_id": req.Package.ID,
		"name":       req.Package.Name,
	})
}

// UpdateRewardPackage 更新奖励包
func (c *RewardPackageController) UpdateRewardPackage(ctx *gin.Context) {
	// 定义请求结构体
	type UpdateRequest struct {
		Package models.RewardPackage        `json:"package" binding:"required"`
		Items   []*models.RewardPackageItem `json:"items" binding:"required"`
	}

	var req UpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResClientError(ctx, "JSON数据格式错误")
		return
	}

	// 验证奖励包ID
	if req.Package.ID == 0 {
		utils.ResClientError(ctx, "奖励包ID不能为空")
		return
	}

	// 验证奖励包内容
	if len(req.Items) == 0 {
		utils.ResClientError(ctx, "奖励包内容不能为空")
		return
	}

	// 更新奖励包基本信息
	if err := c.rewardPackageService.UpdateRewardPackage(&req.Package); err != nil {
		utils.ResServerError(ctx, err)
		return
	}

	// 更新奖励包内容
	if err := c.rewardPackageService.UpdateRewardPackageItems(req.Package.ID, req.Items); err != nil {
		utils.ResServerError(ctx, err)
		return
	}

	utils.ResSuccess(ctx, "更新成功", gin.H{
		"package_id": req.Package.ID,
		"name":       req.Package.Name,
	})
}

// GetRewardPackage 获取奖励包详情
func (c *RewardPackageController) GetRewardPackage(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResClientError(ctx, "无效的奖励包ID")
		return
	}

	// 获取奖励包基本信息
	pkg, err := c.rewardPackageService.GetRewardPackageByID(uint(id))
	if err != nil {
		utils.ResClientError(ctx, "奖励包不存在")
		return
	}

	// 获取奖励包内容
	items, err := c.rewardPackageService.GetRewardPackageItems(uint(id))
	if err != nil {
		utils.ResServerError(ctx, err)
		return
	}

	utils.ResSuccess(ctx, "查询成功", gin.H{
		"package": pkg,
		"items":   items,
	})
}

// ListRewardPackages 获取奖励包列表
func (c *RewardPackageController) ListRewardPackages(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	packages, total, err := c.rewardPackageService.ListRewardPackages(page, pageSize)
	if err != nil {
		utils.ResServerError(ctx, err)
		return
	}

	utils.ResSuccess(ctx, "查询成功", gin.H{
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
		"packages": packages,
	})
}

// DeleteRewardPackage 删除奖励包
func (c *RewardPackageController) DeleteRewardPackage(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResClientError(ctx, "无效的奖励包ID")
		return
	}

	if err := c.rewardPackageService.DeleteRewardPackage(uint(id)); err != nil {
		utils.ResServerError(ctx, err)
		return
	}

	utils.ResSuccess(ctx, "删除成功", nil)
}

// GetUserRewardRecords 获取用户奖励记录
func (c *RewardPackageController) GetUserRewardRecords(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.ResClientError(ctx, "无效的用户ID")
		return
	}

	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	records, total, err := c.rewardPackageService.GetRewardRecordsByUserID(uint(userID), page, pageSize)
	if err != nil {
		utils.ResServerError(ctx, err)
		return
	}

	utils.ResSuccess(ctx, "查询成功", gin.H{
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
		"records":  records,
	})
}

// GrantReward 手动发放奖励
func (c *RewardPackageController) GrantReward(ctx *gin.Context) {
	type GrantRequest struct {
		UserID    uint   `json:"user_id" binding:"required"`
		PackageID uint   `json:"package_id" binding:"required"`
		Source    string `json:"source" binding:"required"`
	}

	var req GrantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResClientError(ctx, "JSON数据格式错误")
		return
	}

	record, err := c.rewardPackageService.GrantReward(nil, req.UserID, req.PackageID, req.Source)
	if err != nil {
		utils.ResServerError(ctx, err)
		return
	}

	utils.ResSuccess(ctx, "发放奖励成功", gin.H{
		"record_id":  record.ID,
		"user_id":    record.UserID,
		"package_id": record.PackageID,
	})
}
