package utils

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"goDDD1/config"
	"time"
)

type CachedUserInfo struct {
	UID   uint   `json:"uid"`
	Email string `json:"email"`
}

func GetUserInfoFromCacheRedis(token string) (*CachedUserInfo, error) {
	cacheKey := fmt.Sprintf("jwt:token:%s", hashToken(token))

	result := config.RedisClient.Get(context.Background(), cacheKey)
	if result.Err() != nil {
		return nil, result.Err()
	}

	var userInfo CachedUserInfo
	err := json.Unmarshal([]byte(result.Val()), &userInfo)
	if err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func CacheUserInfo(token string, uid uint, email string) {
	userInfo := CachedUserInfo{
		UID:   uid,
		Email: email,
	}

	data, _ := json.Marshal(userInfo)
	cacheKey := fmt.Sprintf("jwt:token:%s", hashToken(token))

	// 缓存时间设置为token剩余有效期的一半，避免缓存过期问题
	// remainingTime, _ := utils.GetTokenRemainingTime(token)
	// cacheExpiry := remainingTime / 2
	// if cacheExpiry > time.Hour {
	// 	cacheExpiry = time.Hour // 最大缓存1小时
	// }
	// cacheExpiry := time.Hour // 最大缓存1小时

	config.RedisClient.Set(context.Background(), cacheKey, data, time.Minute)
}

// 对token进行哈希，避免Redis key过长
func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
