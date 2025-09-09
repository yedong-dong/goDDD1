package services

import (
	"context"
	"fmt"
	"goDDD1/config"
	"goDDD1/utils"
	"math/rand"
	"time"
)

// VerificationService 验证码服务接口
type VerificationService interface {
	SendVerificationCode(email string) error
	VerifyCode(email, code string) bool
	CheckVerificationCodeExists(email string) bool
	GetVerificationCodeTTL(email string) (time.Duration, error)
	DeleteVerificationCode(email string) error
}

type verificationService struct{}

// NewVerificationService 创建验证码服务实例
func NewVerificationService() VerificationService {
	return &verificationService{}
}

// SendVerificationCode 发送验证码
func (s *verificationService) SendVerificationCode(email string) error {
	// 生成6位数字验证码
	code := generateVerificationCode()

	// 验证码在Redis中的key
	key := fmt.Sprintf("verification_code:%s", email)

	// 将验证码存储到Redis，有效期5分钟
	err := utils.SetCache(key, code, 5*time.Minute)
	if err != nil {
		return fmt.Errorf("存储验证码失败: %v", err)
	}

	// 这里应该调用邮件服务发送验证码
	// 为了演示，我们只是打印验证码
	fmt.Printf("发送验证码到 %s: %s\n", email, code)

	// TODO: 集成真实的邮件服务
	// err = emailService.SendVerificationCode(email, code)
	// if err != nil {
	//     return fmt.Errorf("发送邮件失败: %v", err)
	// }

	return nil
}

// VerifyCode 验证验证码
func (s *verificationService) VerifyCode(email, code string) bool {
	key := fmt.Sprintf("verification_code:%s", email)

	// 从Redis获取存储的验证码
	var storedCode string
	err := utils.GetCache(key, &storedCode)
	if err != nil {
		// 验证码不存在或已过期
		return false
	}

	// 比较验证码
	return storedCode == code
}

// DeleteVerificationCode 删除验证码
func (s *verificationService) DeleteVerificationCode(email string) error {
	key := fmt.Sprintf("verification_code:%s", email)
	return utils.DeleteCache(key)
}

// generateVerificationCode 生成6位数字验证码
func generateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(900000) + 100000 // 生成100000-999999之间的数字
	return fmt.Sprintf("%06d", code)
}

// CheckVerificationCodeExists 检查验证码是否存在
func (s *verificationService) CheckVerificationCodeExists(email string) bool {
	key := fmt.Sprintf("verification_code:%s", email)
	exists, err := utils.ExistsCache(key)
	return err == nil && exists
}

// GetVerificationCodeTTL 获取验证码剩余有效时间
func (s *verificationService) GetVerificationCodeTTL(email string) (time.Duration, error) {
	key := fmt.Sprintf("verification_code:%s", email)
	rdb := config.GetRedisClient()
	ctx := context.Background()

	ttl, err := rdb.TTL(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	return ttl, nil
}
