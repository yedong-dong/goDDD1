package utils

import (
	"context"
	"encoding/json"
	"goDDD1/config"
	"time"
)

// SetCache 设置缓存
func SetCache(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	rdb := config.GetRedisClient()

	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return rdb.Set(ctx, key, jsonValue, expiration).Err()
}

// GetCache 获取缓存
func GetCache(key string, dest interface{}) error {
	ctx := context.Background()
	rdb := config.GetRedisClient()

	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

// DeleteCache 删除缓存
func DeleteCache(key string) error {
	ctx := context.Background()
	rdb := config.GetRedisClient()

	return rdb.Del(ctx, key).Err()
}

// ExistsCache 检查缓存是否存在
func ExistsCache(key string) (bool, error) {
	ctx := context.Background()
	rdb := config.GetRedisClient()

	result, err := rdb.Exists(ctx, key).Result()
	return result > 0, err
}
