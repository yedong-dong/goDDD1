package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

// RedisConfig Redis配置结构体
type RedisConfig struct {
	Host         string
	Port         string
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	MaxRetries   int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// RedisClient 全局Redis客户端实例
var RedisClient *redis.Client

// InitRedis 初始化Redis连接
func InitRedis() *redis.Client {
	// 加载.env文件中的环境变量
	err := godotenv.Load()
	if err != nil {
		log.Println("未找到.env文件，将使用默认Redis配置")
	}

	// 从环境变量中获取Redis配置
	redisConfig := RedisConfig{
		Host:         getEnv("REDIS_HOST", "localhost"),
		Port:         getEnv("REDIS_PORT", "6378"),
		Password:     getEnv("REDIS_PASSWORD", ""),
		DB:           getEnvAsInt("REDIS_DB", 0),
		PoolSize:     getEnvAsInt("REDIS_POOL_SIZE", 10),
		MinIdleConns: getEnvAsInt("REDIS_MIN_IDLE_CONNS", 5),
		MaxRetries:   getEnvAsInt("REDIS_MAX_RETRIES", 3),
		DialTimeout:  time.Duration(getEnvAsInt("REDIS_DIAL_TIMEOUT", 5)) * time.Second,
		ReadTimeout:  time.Duration(getEnvAsInt("REDIS_READ_TIMEOUT", 3)) * time.Second,
		WriteTimeout: time.Duration(getEnvAsInt("REDIS_WRITE_TIMEOUT", 3)) * time.Second,
	}

	// 创建Redis客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", redisConfig.Host, redisConfig.Port),
		Password:     redisConfig.Password,
		DB:           redisConfig.DB,
		PoolSize:     redisConfig.PoolSize,
		MinIdleConns: redisConfig.MinIdleConns,
		MaxRetries:   redisConfig.MaxRetries,
		DialTimeout:  redisConfig.DialTimeout,
		ReadTimeout:  redisConfig.ReadTimeout,
		WriteTimeout: redisConfig.WriteTimeout,
	})

	// 测试连接
	ctx := context.Background()
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis连接失败: %v", err)
	}

	log.Println("Redis连接成功")

	// 保存全局Redis客户端实例
	RedisClient = rdb

	return rdb
}

// CloseRedis 关闭Redis连接
func CloseRedis() {
	if RedisClient != nil {
		err := RedisClient.Close()
		if err != nil {
			log.Printf("关闭Redis连接时出错: %v", err)
		}
	}
}

// GetRedisClient 获取Redis客户端实例
func GetRedisClient() *redis.Client {
	return RedisClient
}

// 从环境变量获取整数值，如果不存在或转换失败则返回默认值
func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("环境变量 %s 转换为整数失败，使用默认值 %d", key, defaultValue)
		return defaultValue
	}

	return intValue
}
