package controllers

import (
	"goDDD1/models"
	"goDDD1/services"
	"net/http"
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "JSON绑定失败: " + err.Error()})
		return
	}

	// 调试日志：打印接收到的用户数据
	ctx.Header("Content-Type", "application/json")
	if user.Username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "用户名为空", "username": user.Username})
		return
	}
	if user.Email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "邮箱为空", "email": user.Email})
		return
	}
	if user.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "密码为空", "password": "[隐藏]"})
		return
	}

	if err := c.userService.CreateUser(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "用户注册成功", "user": user})
}

// GetUserByID 根据ID获取用户
func (c *UserController) GetUserByUID(ctx *gin.Context) {
	// 优先使用uid参数，如果没有则使用id参数
	uidStr := ctx.Query("uid")
	if uidStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "缺少uid参数"})
		return
	}

	var err error

	uid, err := strconv.ParseUint(uidStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	user, err := c.userService.GetUserByUID(uint(uid))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    "50000",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

// GetUserByID 根据ID获取用户
func (c *UserController) GetUserByUIDDetail(ctx *gin.Context) {
	// 优先使用uid参数，如果没有则使用id参数
	uidStr := ctx.Query("uid")
	if uidStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "缺少uid参数"})
		return
	}

	var err error

	uid, err := strconv.ParseUint(uidStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    "50000",
			"message": err.Error(),
		})
		return
	}

	user, err := c.userService.GetUserByUID(uint(uid))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    "50000",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func (c *UserController) GetAllUsers(ctx *gin.Context) {
	users, err := c.userService.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 20000,
		"data": gin.H{
			"users":   users,
			"message": "success",
		},
	})
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "JSON数据格式错误",
			"message": err.Error(),
		})
		return
	}

	// 2. 参数验证：检查UID是否有效
	if requestData.UID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "无效的用户UID",
			"message": "UID必须是大于0的数字",
		})
		return
	}

	// 3. 先从数据库中查询出用户
	user, err := c.userService.GetUserByUID(uint(requestData.UID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "用户不存在",
			"message": "未找到指定UID的用户",
		})
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
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "更新用户失败",
			"message": "服务器内部错误，请稍后重试",
		})
		return
	}

	// 6. 成功响应：返回更新后的用户信息
	ctx.JSON(http.StatusOK, gin.H{
		"code": "20000",
		"data": gin.H{
			"message": "success",
			"user":    user,
		},
	})
}
