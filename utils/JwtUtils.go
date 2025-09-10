package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims JWT 声明结构体
type JWTClaims struct {
	UID   uint   `json:"uid"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// JWT 配置
var (
	JWTSecret           = []byte("your-secret-key-change-in-production") // 在生产环境中应该从环境变量读取
	TokenExpireDuration = time.Hour * 2
)

// GenerateToken 生成 JWT token
func GenerateToken(uid uint, email string) (string, error) {
	// 创建声明
	claims := JWTClaims{
		UID:   uid,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "goDDD1",
			Subject:   email,
		},
	}

	// 创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名并获取完整的编码后的字符串 token
	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析 JWT token
func ParseToken(tokenString string) (*JWTClaims, error) {
	// 解析 token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// 检查 token 是否有效
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ValidateToken 验证 token 是否有效
func ValidateToken(tokenString string) bool {
	_, err := ParseToken(tokenString)
	return err == nil
}

// RefreshToken 刷新 token（在即将过期时）
func RefreshToken(tokenString string) (string, error) {
	// 解析旧 token
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	// 检查 token 是否即将过期（在过期前1小时内可以刷新）
	if time.Until(claims.ExpiresAt.Time) > time.Hour {
		return "", errors.New("token is not close to expiration")
	}

	// 生成新 token
	return GenerateToken(claims.UID, claims.Email)
}

// GetUserInfoFromToken 从 token 中获取用户信息
func GetUserInfoFromToken(tokenString string) (uint, string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return 0, "", err
	}

	return claims.UID, claims.Email, nil
}

// SetJWTSecret 设置 JWT 密钥（用于配置）
func SetJWTSecret(secret string) {
	JWTSecret = []byte(secret)
}

// SetTokenExpireDuration 设置 token 过期时间
func SetTokenExpireDuration(duration time.Duration) {
	TokenExpireDuration = duration
}

// GetTokenExpireDuration 获取 token 过期时间
func GetTokenExpireDuration() time.Duration {
	return TokenExpireDuration
}

// IsTokenExpired 检查 token 是否已过期
func IsTokenExpired(tokenString string) bool {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return true
	}

	return time.Now().After(claims.ExpiresAt.Time)
}

// GetTokenRemainingTime 获取 token 剩余有效时间
func GetTokenRemainingTime(tokenString string) (time.Duration, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return 0, err
	}

	remainingTime := time.Until(claims.ExpiresAt.Time)
	if remainingTime < 0 {
		return 0, errors.New("token has expired")
	}

	return remainingTime, nil
}
