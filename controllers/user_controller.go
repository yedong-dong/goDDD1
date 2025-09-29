package controllers

import (
	"goDDD1/models"
	"goDDD1/services"
	"goDDD1/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	userService services.UserService
}

// NewUserController 创建用户控制器实例
func NewUserController() *UserController {
	return &UserController{
		userService: services.NewUserService(),
	}
}

// Register 注册用户
func (c *UserController) Register(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		utils.ResClientError(ctx, "JSON绑定失败: "+err.Error())
		return
	}

	// 调试日志：打印接收到的用户数据
	ctx.Header("Content-Type", "application/json")
	if user.Username == "" {
		utils.ResClientError(ctx, "用户名为空")
		return
	}
	if user.Email == "" {
		utils.ResClientError(ctx, "邮箱为空")
		return
	}
	if user.Password == "" {
		utils.ResClientError(ctx, "密码为空")
		return
	}

	if err := c.userService.CreateUser(&user); err != nil {
		utils.ResServerError(ctx, err)
		return
	}

	utils.ResSuccess(ctx, "用户注册成功", user)
}

// GetUserByID 根据ID获取用户
func (c *UserController) GetUserByUID(ctx *gin.Context) {
	// 优先使用uid参数，如果没有则使用id参数
	uidStr := ctx.Query("uid")
	if uidStr == "" {
		utils.ResClientError(ctx, "缺少uid参数")
		return
	}

	var err error

	uid, err := strconv.ParseUint(uidStr, 10, 32)
	if err != nil {
		utils.ResClientError(ctx, "无效的用户ID")
		return
	}

	user, err := c.userService.GetUserByUID(uint(uid))
	if err != nil {
		utils.ResServerError(ctx, err)
		return
	}

	utils.ResSuccess(ctx, "获取用户成功", user)
}

// GetUserByID 根据ID获取用户
func (c *UserController) GetUserByUIDDetail(ctx *gin.Context) {
	// 优先使用uid参数，如果没有则使用id参数
	uidStr := ctx.Query("uid")
	if uidStr == "" {
		utils.ResClientError(ctx, "缺少uid参数")
		return
	}

	var err error

	uid, err := strconv.ParseUint(uidStr, 10, 32)
	if err != nil {
		utils.ResServerError(ctx, err)
		return
	}

	user, err := c.userService.GetUserByUID(uint(uid))
	if err != nil {
		utils.ResServerError(ctx, err)
		return
	}

	utils.ResSuccess(ctx, "获取用户详情成功", user)
}

func (c *UserController) GetAllUsers(ctx *gin.Context) {
	users, err := c.userService.GetAllUsers()
	if err != nil {
		utils.ResClientError(ctx, "用户不存在")
		return
	}

	utils.ResSuccess(ctx, "获取所有用户成功", users)
}

// UpdateUser 更新用户信息
func (c *UserController) UpdateUser(ctx *gin.Context) {
	// 接收POST请求，通过请求体传递uid和要更新的用户信息
	// 请求格式: POST /api/users/update
	// 请求体: JSON格式，包含uid和用户信息

	// 1. 请求体解析：绑定JSON数据到结构体
	var requestData struct {
		UID       uint   `json:"uid" binding:"required"`
		Username  string `json:"username,omitempty"`
		Email     string `json:"email,omitempty"`
		Password  string `json:"password,omitempty"`
		IsDeleted string `json:"is_deleted,omitempty"`
	}

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		utils.ResClientError(ctx, "JSON数据格式错误: "+err.Error())
		return
	}

	// 2. 参数验证：检查UID是否有效
	if requestData.UID == 0 {
		utils.ResClientError(ctx, "无效的用户UID: UID必须是大于0的数字")
		return
	}

	// 3. 先从数据库中查询出用户
	user, err := c.userService.GetUserByUID(uint(requestData.UID))
	if err != nil {
		utils.ResClientError(ctx, "用户不存在: 未找到指定UID的用户")
		return
	}

	// 4. 从requestData中赋值给user（只更新非空字段）
	if requestData.Username != "" {
		user.Username = requestData.Username
	}
	if requestData.Email != "" {
		user.Email = requestData.Email
	}
	if requestData.Password != "" {
		user.Password = requestData.Password
	}
	if requestData.IsDeleted != "" {
		user.IsDeleted = requestData.IsDeleted
	} else {
		user.IsDeleted = "0"
	}

	// 5. 业务逻辑：调用服务层更新用户信息到数据库
	if err := c.userService.UpdateUser(user); err != nil {
		utils.ResServerError(ctx, err)
		return
	}

	// 6. 成功响应：返回更新后的用户信息
	utils.ResSuccess(ctx, "更新用户成功", user)
}
