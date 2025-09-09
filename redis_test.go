package main

import (
	"context"
	"fmt"
	"goDDD1/config"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestRedisConnection 测试Redis连接
func TestRedisConnection(t *testing.T) {
	// 设置测试环境变量
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("REDIS_PORT", "6378")
	os.Setenv("REDIS_PASSWORD", "")
	os.Setenv("REDIS_DB", "0")
	os.Setenv("REDIS_POOL_SIZE", "10")

	// 初始化Redis连接
	client := config.InitRedis()
	assert.NotNil(t, client, "Redis客户端不应该为nil")

	// 测试Ping命令
	ctx := context.Background()
	pong, err := client.Ping(ctx).Result()
	assert.NoError(t, err, "Redis连接应该成功")
	assert.Equal(t, "PONG", pong, "Ping命令应该返回PONG")

	// 测试基本的SET/GET操作
	testKey := "test:simple:key"
	testValue := "Hello Redis!"

	// SET操作
	err = client.Set(ctx, testKey, testValue, time.Minute).Err()
	assert.NoError(t, err, "SET操作应该成功")

	// GET操作
	val, err := client.Get(ctx, testKey).Result()
	assert.NoError(t, err, "GET操作应该成功")
	assert.Equal(t, testValue, val, "获取的值应该与设置的值相同")

	// 清理测试数据
	client.Del(ctx, testKey)

	// 关闭连接
	config.CloseRedis()

	fmt.Println("✅ Redis连接测试通过！")
}

// TestRedisBasicOperations 测试Redis基本操作
func TestRedisBasicOperations(t *testing.T) {
	// 设置测试环境
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("REDIS_PORT", "6378")
	os.Setenv("REDIS_DB", "0")

	client := config.InitRedis()
	defer config.CloseRedis()

	ctx := context.Background()
	testKey := "test:operations"

	// 测试字符串操作
	err := client.Set(ctx, testKey, "测试数据", time.Minute).Err()
	assert.NoError(t, err)

	// 测试EXISTS
	exists, err := client.Exists(ctx, testKey).Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), exists)

	// 测试TTL
	ttl, err := client.TTL(ctx, testKey).Result()
	assert.NoError(t, err)
	assert.Greater(t, ttl.Seconds(), float64(0))

	// 测试DEL
	deleted, err := client.Del(ctx, testKey).Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), deleted)

	fmt.Println("✅ Redis基本操作测试通过！")
}
