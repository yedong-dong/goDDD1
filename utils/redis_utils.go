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

// SetHashField 将指定 key-field 对应的 value 以 JSON 序列化后写入 Redis 哈希表，
// 并支持为整个哈希 key 设置过期时间。
// 参数说明：
//
//	key       - Redis 哈希表的主键
//	field     - 哈希表中的字段名
//	value     - 任意可 JSON 序列化的 Go 值
//	expiration - key 的过期时间；传入 0 表示不设置过期
//
// 返回值：
//
//	成功返回 nil；序列化失败、HSet 失败或设置过期失败时返回对应的 error
func SetHashField(key string, field string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	rdb := config.GetRedisClient()

	// 将 value 序列化为 JSON 字节数组，便于统一存储
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// 将序列化后的值写入哈希字段
	if err := rdb.HSet(ctx, key, field, jsonValue).Err(); err != nil {
		return err
	}

	// 若指定了过期时间，则给整个哈希 key 设置过期
	if expiration > 0 {
		return rdb.Expire(ctx, key, expiration).Err()
	}

	return nil

}

// GetHashField 从 Redis 哈希表中获取指定 key-field 对应的值，
// 并将其反序列化到 dest 指针指向的结构体或映射中。
// 参数说明：
//
//	key   - Redis 哈希表的主键
//	field - 哈希表中的字段名
//	dest  - 指向目标结构体或映射的指针，用于存储反序列化后的值
//
// 返回值：
//
//	成功返回 nil；获取失败或反序列化失败时返回对应的 error
func GetHashField(key string, field string, dest interface{}) error {
	ctx := context.Background()
	rdb := config.GetRedisClient()

	// 从哈希表中获取字段值
	jsonValue, err := rdb.HGet(ctx, key, field).Result()
	if err != nil {
		return err
	}

	// 将 JSON 字节数组反序列化到 dest 指向的结构体或映射
	return json.Unmarshal([]byte(jsonValue), dest)

}

func IncrBy(key string, value int64) (int64, error) {
	ctx := context.Background()
	rdb := config.GetRedisClient()
	return rdb.IncrBy(ctx, key, value).Result()
}

func DelHashField(key string, field string) error {
	ctx := context.Background()
	rdb := config.GetRedisClient()
	return rdb.HDel(ctx, key, field).Err()
}
