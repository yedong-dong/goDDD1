package middleware

import (
	"net/http"
	"strings"

	"goDDD1/utils"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware JWT认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "缺少Authorization头",
				"code":  401,
			})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization头格式错误",
				"code":  401,
			})
			c.Abort()
			return
		}

		// 验证token
		tokenString := parts[1]
		if !utils.ValidateToken(tokenString) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "无效的token",
				"code":  401,
			})
			c.Abort()
			return
		}

		// 获取用户信息
		uid, email, err := utils.GetUserInfoFromToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "解析token失败",
				"code":  401,
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("uid", uid)
		c.Set("email", email)
		c.Set("token", tokenString)

		// 继续处理请求
		c.Next()
	}
}

// OptionalJWTAuthMiddleware 可选的JWT认证中间件（用于某些可选登录的接口）
func OptionalJWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString := parts[1]
				if utils.ValidateToken(tokenString) {
					uid, email, _ := utils.GetUserInfoFromToken(tokenString)
					c.Set("uid", uid)
					c.Set("email", email)
					c.Set("token", tokenString)
					c.Set("authenticated", true)
				} else {
					c.Set("authenticated", false)
				}
			} else {
				c.Set("authenticated", false)
			}
		} else {
			c.Set("authenticated", false)
		}
		c.Next()
	}
}

// AdminAuthMiddleware 管理员权限中间件
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先执行JWT验证
		JWTAuthMiddleware()(c)
		if c.IsAborted() {
			return
		}

		// 这里可以添加管理员权限检查逻辑
		// 例如：检查用户角色、权限等
		uid, exists := c.Get("uid")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "权限不足",
				"code":  403,
			})
			c.Abort()
			return
		}

		// TODO: 根据uid查询用户角色，验证是否为管理员
		_ = uid

		c.Next()
	}
}
