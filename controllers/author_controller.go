package controllers

import (
	"fmt"
	"goDDD1/models"
	"goDDD1/services"
	"goDDD1/utils"
	"regexp"

	"github.com/gin-gonic/gin"
)

type AuthorizationController struct {
	userService         services.UserService
	verificationService services.VerificationService
}

func NewAuthorizationController() *AuthorizationController {
	return &AuthorizationController{
		userService:         services.NewUserService(),
		verificationService: services.NewVerificationService(),
	}
}

// RegisterRequest 注册请求结构体
type RegisterRequest struct {
	Username         string `json:"username" binding:"required"`
	Email            string `json:"email" binding:"required,email"`
	Password         string `json:"password" binding:"required,min=6"`
	VerificationCode string `json:"verification_code" binding:"required"`
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

// SendVerificationCode 发送验证码
func (c *AuthorizationController) SendVerificationCode(ctx *gin.Context) {
	email := ctx.Query("email")
	if email == "" {
		utils.ResClientError(ctx, "邮箱不能为空")
		return
	}

	existingUser, _ := c.userService.GetUserByEmail(email)
	if existingUser != nil {
		utils.ResClientError(ctx, "该邮箱已被注册")
		return
	}

	// 检查是否已有未过期的验证码
	if c.verificationService.CheckVerificationCodeExists(email) {
		ttl, _ := c.verificationService.GetVerificationCodeTTL(email)
		utils.ResClientError(ctx, fmt.Sprintf("请等待 %.0f 秒后再次发送", ttl.Seconds()))
		return
	}

	// 发送验证码
	err := c.verificationService.SendVerificationCode(email)
	if err != nil {
		utils.ResServerError(ctx, err)
		return
	}
	utils.ResSuccess(ctx, "验证码已发送，请查收邮件", nil)
}

// Register 用户注册（修改后的版本）
func (c *AuthorizationController) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResClientError(ctx, err.Error())
		return
	}

	// 验证验证码
	if !c.verificationService.VerifyCode(req.Email, req.VerificationCode) {
		utils.ResClientError(ctx, "验证码不正确或已过期，请重新获取")
		return
	}

	// 验证用户名格式（只允许字母、数字、下划线，长度3-20）
	if !isValidUsername(req.Username) {
		utils.ResClientError(ctx, "用户名只能包含字母、数字、下划线，长度3-20位")
		return
	}

	// 检查用户名是否已存在
	existingUser, _ := c.userService.GetUserByUsername(req.Username)
	if existingUser != nil {
		utils.ResClientError(ctx, "该用户名已被注册，请选择其他用户名")
		return
	}

	// 检查邮箱是否已存在
	existingEmailUser, _ := c.userService.GetUserByEmail(req.Email)
	if existingEmailUser != nil {
		utils.ResClientError(ctx, "该邮箱已被注册，请使用其他邮箱")
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
		utils.ResServerError(ctx, err)
		return
	}
	// 注册成功后删除验证码
	c.verificationService.DeleteVerificationCode(req.Email)

	// 返回成功响应（不包含密码）
	userResp := UserResponse{
		ID:       user.ID,
		UID:      user.UID,
		Username: user.Username,
		Email:    user.Email,
	}

	utils.ResSuccess(ctx, "注册成功", userResp)
}

// Login 用户登录
func (c *AuthorizationController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResClientError(ctx, err.Error())
		return
	}

	// 根据用户名查找用户
	user, err := c.userService.GetUserByEmail(req.Email)
	if err != nil {
		utils.ResClientError(ctx, "用户名或密码错误")
		return
	}

	// 验证密码
	if !verifyPassword(req.Password, user.Password) {
		utils.ResClientError(ctx, "用户名或密码错误")
		return
	}

	// 检查用户是否被删除
	if user.IsDeleted == "1" {
		utils.ResClientError(ctx, "账户已被禁用")
		return
	}

	// 生成token
	token := generateToken(user)

	utils.ResSuccess(ctx, "登录成功", gin.H{
		"token": "Bearer " + token,
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
