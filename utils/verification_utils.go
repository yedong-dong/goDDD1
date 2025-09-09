package utils

import (
	"fmt"
	"regexp"
)

// IsValidEmail 验证邮箱格式
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// IsValidVerificationCode 验证验证码格式（6位数字）
func IsValidVerificationCode(code string) bool {
	codeRegex := regexp.MustCompile(`^\d{6}$`)
	return codeRegex.MatchString(code)
}

// FormatVerificationCodeKey 格式化验证码在Redis中的key
func FormatVerificationCodeKey(email string) string {
	return fmt.Sprintf("verification_code:%s", email)
}
