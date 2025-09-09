package controllers

import (
	"crypto/md5"
	"fmt"
	"goDDD1/models"
	"goDDD1/services"
	"goDDD1/utils"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

type AuthorizationController struct {
	userService services.UserService
}

func NewAuthorizationController() *AuthorizationController {
	return &AuthorizationController{
		userService: services.NewUserService(),
	}
}

// RegisterRequest 注册请求结构体
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest 登录请求结构体
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserResponse 用户响应结构体（不包含密码）
type UserResponse struct {
	ID       uint   `json:"id"`
	UID      uint   `json:"uid"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Register 用户注册
func (c *AuthorizationController) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求参数错误",
			"message": err.Error(),
		})
		return
	}

	// 验证用户名格式（只允许字母、数字、下划线，长度3-20）
	if !isValidUsername(req.Username) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "用户名格式错误",
			"message": "用户名只能包含字母、数字、下划线，长度3-20位",
		})
		return
	}

	// 检查用户名是否已存在
	existingUser, _ := c.userService.GetUserByUsername(req.Username)
	if existingUser != nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"error":   "用户名已存在",
			"message": "该用户名已被注册，请选择其他用户名",
		})
		return
	}

	// 创建用户对象
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password, // 加密密码
	}

	// 创建用户
	if err := c.userService.CreateUser(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "注册失败",
			"message": err.Error(),
		})
		return
	}

	// 返回成功响应（不包含密码）
	userResp := UserResponse{
		ID:       user.ID,
		UID:      user.UID,
		Username: user.Username,
		Email:    user.Email,
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "注册成功",
		"user":    userResp,
	})
}

// Login 用户登录
func (c *AuthorizationController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求参数错误",
			"message": err.Error(),
		})
		return
	}

	// 根据用户名查找用户
	user, err := c.userService.GetUserByEmail(req.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   "登录失败",
			"message": "用户名或密码错误",
		})
		return
	}

	// 验证密码
	if !verifyPassword(req.Password, user.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   "登录失败",
			"message": "用户名或密码错误",
		})
		return
	}

	// 检查用户是否被删除
	if user.IsDeleted == "1" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   "登录失败",
			"message": "账户已被禁用",
		})
		return
	}

	// 返回成功响应（不包含密码）
	userResp := UserResponse{
		ID:       user.ID,
		UID:      user.UID,
		Username: user.Username,
		Email:    user.Email,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"user":    userResp,
		"token":   "Bearer " + generateToken(user), // 这里可以生成JWT token
		// "token": "1", // 这里可以生成JWT token
	})
}

// 辅助函数：验证用户名格式
func isValidUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	// 只允许字母、数字、下划线
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
	return matched
}

// 辅助函数：密码加密（使用MD5，生产环境建议使用bcrypt）
func hashPassword(password string) string {
	hash := md5.Sum([]byte(password))
	return fmt.Sprintf("%x", hash)
}

// 辅助函数：验证密码
func verifyPassword(password, hashedPassword string) bool {
	// return hashPassword(password) == hashedPassword
	return password == hashedPassword
}

// 辅助函数：生成token（简单实现，生产环境建议使用JWT）
func generateToken(user *models.User) string {
	token, err := utils.GenerateToken(user.UID, user.Email)
	if err != nil {
		return ""
	}
	return token
}
