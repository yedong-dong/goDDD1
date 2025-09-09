# Go Redis 使用文档

## 目录
1. [简介](#简介)
2. [安装与配置](#安装与配置)
3. [基础连接](#基础连接)
4. [常用API](#常用api)
5. [高级功能](#高级功能)
6. [最佳实践](#最佳实践)
7. [错误处理](#错误处理)
8. [性能优化](#性能优化)

## 简介

Redis是一个开源的内存数据结构存储系统，可以用作数据库、缓存和消息代理。在Go中，我们主要使用`go-redis`库来操作Redis。

## 安装与配置

### 1. 安装依赖

```bash
go mod init your-project
go get github.com/redis/go-redis/v9
```

### 2. 环境配置

创建`.env`文件：

```env
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_POOL_SIZE=10
REDIS_MIN_IDLE_CONNS=5
```

### 3. 配置结构体

```go
// config/redis.go
package config

import (
    "context"
    "fmt"
    "log"
    "os"
    "strconv"
    "time"

    "github.com/redis/go-redis/v9"
)

type RedisConfig struct {
    Host         string
    Port         string
    Password     string
    DB           int
    PoolSize     int
    MinIdleConns int
}

var RedisClient *redis.Client

// InitRedis 初始化Redis连接
func InitRedis() {
    config := &RedisConfig{
        Host:         getEnv("REDIS_HOST", "localhost"),
        Port:         getEnv("REDIS_PORT", "6379"),
        Password:     getEnv("REDIS_PASSWORD", ""),
        DB:           getEnvAsInt("REDIS_DB", 0),
        PoolSize:     getEnvAsInt("REDIS_POOL_SIZE", 10),
        MinIdleConns: getEnvAsInt("REDIS_MIN_IDLE_CONNS", 5),
    }

    RedisClient = redis.NewClient(&redis.Options{
        Addr:         fmt.Sprintf("%s:%s", config.Host, config.Port),
        Password:     config.Password,
        DB:           config.DB,
        PoolSize:     config.PoolSize,
        MinIdleConns: config.MinIdleConns,
        DialTimeout:  5 * time.Second,
        ReadTimeout:  3 * time.Second,
        WriteTimeout: 3 * time.Second,
        PoolTimeout:  4 * time.Second,
    })

    // 测试连接
    ctx := context.Background()
    _, err := RedisClient.Ping(ctx).Result()
    if err != nil {
        log.Fatalf("Redis连接失败: %v", err)
    }
    log.Println("Redis连接成功")
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

// getEnvAsInt 获取环境变量并转换为int
func getEnvAsInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}

// CloseRedis 关闭Redis连接
func CloseRedis() {
    if RedisClient != nil {
        RedisClient.Close()
    }
}
```

## 基础连接

### 在main.go中初始化

```go
package main

import (
    "log"
    "your-project/config"
)

func main() {
    // 初始化Redis
    config.InitRedis()
    defer config.CloseRedis()

    // 你的应用逻辑
    log.Println("应用启动成功")
}
```

## 常用API

### 1. 字符串操作

```go
// utils/redis_string.go
package utils

import (
    "context"
    "time"
    "your-project/config"
)

// SetString 设置字符串值
func SetString(key, value string, expiration time.Duration) error {
    ctx := context.Background()
    return config.RedisClient.Set(ctx, key, value, expiration).Err()
}

// GetString 获取字符串值
func GetString(key string) (string, error) {
    ctx := context.Background()
    return config.RedisClient.Get(ctx, key).Result()
}

// SetNX 仅当key不存在时设置
func SetNX(key, value string, expiration time.Duration) (bool, error) {
    ctx := context.Background()
    return config.RedisClient.SetNX(ctx, key, value, expiration).Result()
}

// Incr 递增
func Incr(key string) (int64, error) {
    ctx := context.Background()
    return config.RedisClient.Incr(ctx, key).Result()
}

// IncrBy 按指定值递增
func IncrBy(key string, value int64) (int64, error) {
    ctx := context.Background()
    return config.RedisClient.IncrBy(ctx, key, value).Result()
}

// Decr 递减
func Decr(key string) (int64, error) {
    ctx := context.Background()
    return config.RedisClient.Decr(ctx, key).Result()
}
```

### 2. 哈希操作

```go
// utils/redis_hash.go
package utils

import (
    "context"
    "your-project/config"
)

// HSet 设置哈希字段
func HSet(key, field, value string) error {
    ctx := context.Background()
    return config.RedisClient.HSet(ctx, key, field, value).Err()
}

// HGet 获取哈希字段值
func HGet(key, field string) (string, error) {
    ctx := context.Background()
    return config.RedisClient.HGet(ctx, key, field).Result()
}

// HMSet 批量设置哈希字段
func HMSet(key string, fields map[string]interface{}) error {
    ctx := context.Background()
    return config.RedisClient.HMSet(ctx, key, fields).Err()
}

// HGetAll 获取所有哈希字段
func HGetAll(key string) (map[string]string, error) {
    ctx := context.Background()
    return config.RedisClient.HGetAll(ctx, key).Result()
}

// HDel 删除哈希字段
func HDel(key string, fields ...string) (int64, error) {
    ctx := context.Background()
    return config.RedisClient.HDel(ctx, key, fields...).Result()
}

// HExists 检查哈希字段是否存在
func HExists(key, field string) (bool, error) {
    ctx := context.Background()
    return config.RedisClient.HExists(ctx, key, field).Result()
}
```

### 3. 列表操作

```go
// utils/redis_list.go
package utils

import (
    "context"
    "time"
    "your-project/config"
)

// LPush 从左侧推入元素
func LPush(key string, values ...interface{}) (int64, error) {
    ctx := context.Background()
    return config.RedisClient.LPush(ctx, key, values...).Result()
}

// RPush 从右侧推入元素
func RPush(key string, values ...interface{}) (int64, error) {
    ctx := context.Background()
    return config.RedisClient.RPush(ctx, key, values...).Result()
}

// LPop 从左侧弹出元素
func LPop(key string) (string, error) {
    ctx := context.Background()
    return config.RedisClient.LPop(ctx, key).Result()
}

// RPop 从右侧弹出元素
func RPop(key string) (string, error) {
    ctx := context.Background()
    return config.RedisClient.RPop(ctx, key).Result()
}

// LRange 获取列表范围内的元素
func LRange(key string, start, stop int64) ([]string, error) {
    ctx := context.Background()
    return config.RedisClient.LRange(ctx, key, start, stop).Result()
}

// LLen 获取列表长度
func LLen(key string) (int64, error) {
    ctx := context.Background()
    return config.RedisClient.LLen(ctx, key).Result()
}

// BLPop 阻塞式左侧弹出
func BLPop(timeout time.Duration, keys ...string) ([]string, error) {
    ctx := context.Background()
    return config.RedisClient.BLPop(ctx, timeout, keys...).Result()
}
```

### 4. 集合操作

```go
// utils/redis_set.go
package utils

import (
    "context"
    "your-project/config"
)

// SAdd 添加集合成员
func SAdd(key string, members ...interface{}) (int64, error) {
    ctx := context.Background()
    return config.RedisClient.SAdd(ctx, key, members...).Result()
}

// SMembers 获取集合所有成员
func SMembers(key string) ([]string, error) {
    ctx := context.Background()
    return config.RedisClient.SMembers(ctx, key).Result()
}

// SIsMember 检查是否为集合成员
func SIsMember(key string, member interface{}) (bool, error) {
    ctx := context.Background()
    return config.RedisClient.SIsMember(ctx, key, member).Result()
}

// SRem 移除集合成员
func SRem(key string, members ...interface{}) (int64, error) {
    ctx := context.Background()
    return config.RedisClient.SRem(ctx, key, members...).Result()
}

// SCard 获取集合成员数量
func SCard(key string) (int64, error) {
    ctx := context.Background()
    return config.RedisClient.SCard(ctx, key).Result()
}

// SUnion 求集合并集
func SUnion(keys ...string) ([]string, error) {
    ctx := context.Background()
    return config.RedisClient.SUnion(ctx, keys...).Result()
}

// SInter 求集合交集
func SInter(keys ...string) ([]string, error) {
    ctx := context.Background()
    return config.RedisClient.SInter(ctx, keys...).Result()
}
```

### 5. 有序集合操作

```go
// utils/redis_zset.go
package utils

import (
    "context"
    "github.com/redis/go-redis/v9"
    "your-project/config"
)

// ZAdd 添加有序集合成员
func ZAdd(key string, members ...redis.Z) (int64, error) {
    ctx := context.Background()
    return config.RedisClient.ZAdd(ctx, key, members...).Result()
}

// ZRange 按索引范围获取成员
func ZRange(key string, start, stop int64) ([]string, error) {
    ctx := context.Background()
    return config.RedisClient.ZRange(ctx, key, start, stop).Result()
}

// ZRangeByScore 按分数范围获取成员
func ZRangeByScore(key string, min, max string) ([]string, error) {
    ctx := context.Background()
    opt := &redis.ZRangeBy{
        Min: min,
        Max: max,
    }
    return config.RedisClient.ZRangeByScore(ctx, key, opt).Result()
}

// ZRank 获取成员排名
func ZRank(key, member string) (int64, error) {
    ctx := context.Background()
    return config.RedisClient.ZRank(ctx, key, member).Result()
}

// ZScore 获取成员分数
func ZScore(key, member string) (float64, error) {
    ctx := context.Background()
    return config.RedisClient.ZScore(ctx, key, member).Result()
}

// ZRem 移除成员
func ZRem(key string, members ...interface{}) (int64, error) {
    ctx := context.Background()
    return config.RedisClient.ZRem(ctx, key, members...).Result()
}
```

### 6. 通用操作

```go
// utils/redis_common.go
package utils

import (
    "context"
    "time"
    "your-project/config"
)

// Exists 检查key是否存在
func Exists(keys ...string) (int64, error) {
    ctx := context.Background()
    return config.RedisClient.Exists(ctx, keys...).Result()
}

// Del 删除key
func Del(keys ...string) (int64, error) {
    ctx := context.Background()
    return config.RedisClient.Del(ctx, keys...).Result()
}

// Expire 设置过期时间
func Expire(key string, expiration time.Duration) (bool, error) {
    ctx := context.Background()
    return config.RedisClient.Expire(ctx, key, expiration).Result()
}

// TTL 获取剩余过期时间
func TTL(key string) (time.Duration, error) {
    ctx := context.Background()
    return config.RedisClient.TTL(ctx, key).Result()
}

// Keys 查找匹配模式的key
func Keys(pattern string) ([]string, error) {
    ctx := context.Background()
    return config.RedisClient.Keys(ctx, pattern).Result()
}

// Type 获取key的类型
func Type(key string) (string, error) {
    ctx := context.Background()
    return config.RedisClient.Type(ctx, key).Result()
}

// FlushDB 清空当前数据库
func FlushDB() error {
    ctx := context.Background()
    return config.RedisClient.FlushDB(ctx).Err()
}
```

## 高级功能

### 1. 发布订阅

```go
// utils/redis_pubsub.go
package utils

import (
    "context"
    "log"
    "your-project/config"
)

// Publish 发布消息
func Publish(channel, message string) error {
    ctx := context.Background()
    return config.RedisClient.Publish(ctx, channel, message).Err()
}

// Subscribe 订阅频道
func Subscribe(channels ...string) {
    ctx := context.Background()
    pubsub := config.RedisClient.Subscribe(ctx, channels...)
    defer pubsub.Close()

    // 等待订阅确认
    _, err := pubsub.Receive(ctx)
    if err != nil {
        log.Printf("订阅失败: %v", err)
        return
    }

    // 监听消息
    ch := pubsub.Channel()
    for msg := range ch {
        log.Printf("收到消息 - 频道: %s, 内容: %s", msg.Channel, msg.Payload)
        // 处理消息逻辑
        handleMessage(msg.Channel, msg.Payload)
    }
}

// handleMessage 处理接收到的消息
func handleMessage(channel, payload string) {
    // 根据频道处理不同的消息
    switch channel {
    case "user_notifications":
        // 处理用户通知
        log.Printf("处理用户通知: %s", payload)
    case "system_events":
        // 处理系统事件
        log.Printf("处理系统事件: %s", payload)
    default:
        log.Printf("未知频道消息: %s - %s", channel, payload)
    }
}
```

### 2. 管道操作

```go
// utils/redis_pipeline.go
package utils

import (
    "context"
    "time"
    "your-project/config"
)

// PipelineExample 管道操作示例
func PipelineExample() error {
    ctx := context.Background()
    pipe := config.RedisClient.Pipeline()

    // 批量操作
    incr := pipe.Incr(ctx, "pipeline_counter")
    pipe.Expire(ctx, "pipeline_counter", time.Hour)
    pipe.Set(ctx, "pipeline_key", "pipeline_value", time.Hour)

    // 执行管道
    _, err := pipe.Exec(ctx)
    if err != nil {
        return err
    }

    // 获取结果
    val, err := incr.Result()
    if err != nil {
        return err
    }

    log.Printf("计数器值: %d", val)
    return nil
}

// TxPipelineExample 事务管道示例
func TxPipelineExample() error {
    ctx := context.Background()
    
    // 使用事务管道
    _, err := config.RedisClient.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
        pipe.Set(ctx, "tx_key1", "value1", time.Hour)
        pipe.Set(ctx, "tx_key2", "value2", time.Hour)
        pipe.Incr(ctx, "tx_counter")
        return nil
    })
    
    return err
}
```

### 3. Lua脚本

```go
// utils/redis_lua.go
package utils

import (
    "context"
    "github.com/redis/go-redis/v9"
    "your-project/config"
)

// 预定义Lua脚本
var (
    // 原子性增加并设置过期时间
    incrWithExpireScript = redis.NewScript(`
        local key = KEYS[1]
        local expire = ARGV[1]
        local current = redis.call('INCR', key)
        if current == 1 then
            redis.call('EXPIRE', key, expire)
        end
        return current
    `)

    // 限流脚本
    rateLimitScript = redis.NewScript(`
        local key = KEYS[1]
        local limit = tonumber(ARGV[1])
        local window = tonumber(ARGV[2])
        
        local current = redis.call('GET', key)
        if current == false then
            redis.call('SET', key, 1)
            redis.call('EXPIRE', key, window)
            return {1, limit}
        end
        
        current = tonumber(current)
        if current < limit then
            local new_val = redis.call('INCR', key)
            return {new_val, limit}
        else
            return {current, limit}
        end
    `)
)

// IncrWithExpire 原子性递增并设置过期时间
func IncrWithExpire(key string, expire int) (int64, error) {
    ctx := context.Background()
    return incrWithExpireScript.Run(ctx, config.RedisClient, []string{key}, expire).Int64()
}

// RateLimit 限流检查
func RateLimit(key string, limit, window int) (bool, error) {
    ctx := context.Background()
    result, err := rateLimitScript.Run(ctx, config.RedisClient, []string{key}, limit, window).Result()
    if err != nil {
        return false, err
    }
    
    values := result.([]interface{})
    current := values[0].(int64)
    maxLimit := values[1].(int64)
    
    return current <= maxLimit, nil
}
```

## 最佳实践

### 1. 缓存服务

```go
// services/cache_service.go
package services

import (
    "encoding/json"
    "fmt"
    "time"
    "your-project/utils"
)

type CacheService struct{}

func NewCacheService() *CacheService {
    return &CacheService{}
}

// SetCache 设置缓存
func (c *CacheService) SetCache(key string, data interface{}, expiration time.Duration) error {
    jsonData, err := json.Marshal(data)
    if err != nil {
        return fmt.Errorf("序列化失败: %v", err)
    }
    
    return utils.SetString(key, string(jsonData), expiration)
}

// GetCache 获取缓存
func (c *CacheService) GetCache(key string, dest interface{}) error {
    data, err := utils.GetString(key)
    if err != nil {
        return fmt.Errorf("获取缓存失败: %v", err)
    }
    
    return json.Unmarshal([]byte(data), dest)
}

// DeleteCache 删除缓存
func (c *CacheService) DeleteCache(keys ...string) error {
    _, err := utils.Del(keys...)
    return err
}

// GetOrSet 获取缓存，如果不存在则设置
func (c *CacheService) GetOrSet(key string, dest interface{}, fetcher func() (interface{}, error), expiration time.Duration) error {
    // 尝试从缓存获取
    err := c.GetCache(key, dest)
    if err == nil {
        return nil // 缓存命中
    }
    
    // 缓存未命中，调用fetcher获取数据
    data, err := fetcher()
    if err != nil {
        return fmt.Errorf("获取数据失败: %v", err)
    }
    
    // 设置缓存
    if err := c.SetCache(key, data, expiration); err != nil {
        return fmt.Errorf("设置缓存失败: %v", err)
    }
    
    // 将数据复制到dest
    jsonData, _ := json.Marshal(data)
    return json.Unmarshal(jsonData, dest)
}
```

### 2. 分布式锁

```go
// utils/redis_lock.go
package utils

import (
    "context"
    "fmt"
    "time"
    "github.com/redis/go-redis/v9"
    "your-project/config"
)

type DistributedLock struct {
    key        string
    value      string
    expiration time.Duration
}

// NewDistributedLock 创建分布式锁
func NewDistributedLock(key, value string, expiration time.Duration) *DistributedLock {
    return &DistributedLock{
        key:        key,
        value:      value,
        expiration: expiration,
    }
}

// Lock 获取锁
func (l *DistributedLock) Lock() (bool, error) {
    ctx := context.Background()
    return config.RedisClient.SetNX(ctx, l.key, l.value, l.expiration).Result()
}

// Unlock 释放锁
func (l *DistributedLock) Unlock() error {
    ctx := context.Background()
    
    // Lua脚本确保原子性释放
    script := `
        if redis.call("GET", KEYS[1]) == ARGV[1] then
            return redis.call("DEL", KEYS[1])
        else
            return 0
        end
    `
    
    _, err := config.RedisClient.Eval(ctx, script, []string{l.key}, l.value).Result()
    return err
}

// TryLockWithRetry 重试获取锁
func (l *DistributedLock) TryLockWithRetry(maxRetries int, retryInterval time.Duration) (bool, error) {
    for i := 0; i < maxRetries; i++ {
        locked, err := l.Lock()
        if err != nil {
            return false, err
        }
        if locked {
            return true, nil
        }
        time.Sleep(retryInterval)
    }
    return false, fmt.Errorf("获取锁失败，已重试%d次", maxRetries)
}
```

### 3. 会话管理

```go
// services/session_service.go
package services

import (
    "encoding/json"
    "fmt"
    "time"
    "your-project/utils"
)

type SessionData struct {
    UserID    int64     `json:"user_id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    LoginTime time.Time `json:"login_time"`
    LastSeen  time.Time `json:"last_seen"`
}

type SessionService struct {
    prefix     string
    expiration time.Duration
}

func NewSessionService() *SessionService {
    return &SessionService{
        prefix:     "session:",
        expiration: 24 * time.Hour, // 24小时过期
    }
}

// CreateSession 创建会话
func (s *SessionService) CreateSession(sessionID string, data *SessionData) error {
    key := s.prefix + sessionID
    data.LoginTime = time.Now()
    data.LastSeen = time.Now()
    
    jsonData, err := json.Marshal(data)
    if err != nil {
        return fmt.Errorf("序列化会话数据失败: %v", err)
    }
    
    return utils.SetString(key, string(jsonData), s.expiration)
}

// GetSession 获取会话
func (s *SessionService) GetSession(sessionID string) (*SessionData, error) {
    key := s.prefix + sessionID
    data, err := utils.GetString(key)
    if err != nil {
        return nil, fmt.Errorf("获取会话失败: %v", err)
    }
    
    var sessionData SessionData
    if err := json.Unmarshal([]byte(data), &sessionData); err != nil {
        return nil, fmt.Errorf("反序列化会话数据失败: %v", err)
    }
    
    return &sessionData, nil
}

// UpdateSession 更新会话
func (s *SessionService) UpdateSession(sessionID string, data *SessionData) error {
    key := s.prefix + sessionID
    data.LastSeen = time.Now()
    
    jsonData, err := json.Marshal(data)
    if err != nil {
        return fmt.Errorf("序列化会话数据失败: %v", err)
    }
    
    return utils.SetString(key, string(jsonData), s.expiration)
}

// DeleteSession 删除会话
func (s *SessionService) DeleteSession(sessionID string) error {
    key := s.prefix + sessionID
    _, err := utils.Del(key)
    return err
}

// RefreshSession 刷新会话过期时间
func (s *SessionService) RefreshSession(sessionID string) error {
    key := s.prefix + sessionID
    _, err := utils.Expire(key, s.expiration)
    return err
}
```

## 错误处理

```go
// utils/redis_errors.go
package utils

import (
    "errors"
    "fmt"
    "log"
    "github.com/redis/go-redis/v9"
)

// IsRedisNil 检查是否为Redis nil错误
func IsRedisNil(err error) bool {
    return errors.Is(err, redis.Nil)
}

// HandleRedisError 统一处理Redis错误
func HandleRedisError(err error) error {
    if err == nil {
        return nil
    }
    
    if IsRedisNil(err) {
        return errors.New("key不存在")
    }
    
    // 其他Redis错误
    return fmt.Errorf("Redis操作失败: %v", err)
}

// SafeGet 安全获取，不存在时返回默认值
func SafeGet(key, defaultValue string) string {
    value, err := GetString(key)
    if IsRedisNil(err) {
        return defaultValue
    }
    if err != nil {
        log.Printf("获取key %s 失败: %v", key, err)
        return defaultValue
    }
    return value
}
```

## 性能优化

### 1. 连接池配置

```go
// 在config/redis.go中优化连接池配置
RedisClient = redis.NewClient(&redis.Options{
    Addr:         fmt.Sprintf("%s:%s", config.Host, config.Port),
    Password:     config.Password,
    DB:           config.DB,
    
    // 连接池配置
    PoolSize:     20,              // 连接池大小
    MinIdleConns: 5,               // 最小空闲连接数
    MaxConnAge:   time.Hour,       // 连接最大存活时间
    PoolTimeout:  30 * time.Second, // 获取连接超时时间
    IdleTimeout:  5 * time.Minute,  // 空闲连接超时时间
    
    // 网络配置
    DialTimeout:  5 * time.Second,
    ReadTimeout:  3 * time.Second,
    WriteTimeout: 3 * time.Second,
    
    // 重试配置
    MaxRetries:      3,
    MinRetryBackoff: 8 * time.Millisecond,
    MaxRetryBackoff: 512 * time.Millisecond,
})
```

### 2. 批量操作

```go
// utils/redis_batch.go
package utils

import (
    "context"
    "time"
    "your-project/config"
    "github.com/redis/go-redis/v9"
)

// BatchSet 批量设置
func BatchSet(data map[string]interface{}, expiration time.Duration) error {
    ctx := context.Background()
    pipe := config.RedisClient.Pipeline()
    
    for key, value := range data {
        pipe.Set(ctx, key, value, expiration)
    }
    
    _, err := pipe.Exec(ctx)
    return err
}

// BatchGet 批量获取
func BatchGet(keys []string) (map[string]string, error) {
    ctx := context.Background()
    pipe := config.RedisClient.Pipeline()
    
    // 创建命令映射
    cmds := make(map[string]*redis.StringCmd)
    for _, key := range keys {
        cmds[key] = pipe.Get(ctx, key)
    }
    
    // 执行管道
    _, err := pipe.Exec(ctx)
    if err != nil {
        return nil, err
    }
    
    // 收集结果
    result := make(map[string]string)
    for key, cmd := range cmds {
        val, err := cmd.Result()
        if err == nil {
            result[key] = val
        }
    }
    
    return result, nil
}
```

## 使用示例

### 完整示例：用户缓存系统

```go
// examples/user_cache_example.go
package main

import (
    "fmt"
    "log"
    "time"
    "your-project/config"
    "your-project/services"
    "your-project/utils"
)

type User struct {
    ID       int64  `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Avatar   string `json:"avatar"`
}

func main() {
    // 初始化Redis
    config.InitRedis()
    defer config.CloseRedis()
    
    // 创建缓存服务
    cacheService := services.NewCacheService()
    
    // 示例用户数据
    user := &User{
        ID:       1001,
        Username: "john_doe",
        Email:    "john@example.com",
        Avatar:   "avatar.jpg",
    }
    
    // 1. 设置用户缓存
    userKey := fmt.Sprintf("user:%d", user.ID)
    if err := cacheService.SetCache(userKey, user, time.Hour); err != nil {
        log.Printf("设置用户缓存失败: %v", err)
        return
    }
    log.Println("用户缓存设置成功")
    
    // 2. 获取用户缓存
    var cachedUser User
    if err := cacheService.GetCache(userKey, &cachedUser); err != nil {
        log.Printf("获取用户缓存失败: %v", err)
        return
    }
    log.Printf("从缓存获取用户: %+v", cachedUser)
    
    // 3. 使用GetOrSet模式
    var user2 User
    user2Key := "user:1002"
    err := cacheService.GetOrSet(user2Key, &user2, func() (interface{}, error) {
        // 模拟从数据库获取用户
        log.Println("从数据库获取用户数据...")
        return &User{
            ID:       1002,
            Username: "jane_doe",
            Email:    "jane@example.com",
            Avatar:   "jane_avatar.jpg",
        }, nil
    }, time.Hour)
    
    if err != nil {
        log.Printf("GetOrSet失败: %v", err)
        return
    }
    log.Printf("GetOrSet获取用户: %+v", user2)
    
    // 4. 分布式锁示例
    lockKey := "user_update_lock:1001"
    lock := utils.NewDistributedLock(lockKey, "unique_value", 30*time.Second)
    
    locked, err := lock.Lock()
    if err != nil {
        log.Printf("获取锁失败: %v", err)
        return
    }
    
    if locked {
        log.Println("获取锁成功，执行业务逻辑...")
        // 执行需要锁保护的业务逻辑
        time.Sleep(2 * time.Second)
        
        // 释放锁
        if err := lock.Unlock(); err != nil {
            log.Printf("释放锁失败: %v", err)
        } else {
            log.Println("锁释放成功")
        }
    } else {
        log.Println("获取锁失败，其他进程正在处理")
    }
    
    // 5. 限流示例
    rateLimitKey := "api_rate_limit:user:1001"
    allowed, err := utils.RateLimit(rateLimitKey, 10, 60) // 每分钟最多10次
    if err != nil {
        log.Printf("限流检查失败: %v", err)
        return
    }
    
    if allowed {
        log.Println("请求通过限流检查")
    } else {
        log.Println("请求被限流")
    }
    
    // 6. 会话管理示例
    sessionService := services.NewSessionService()
    sessionID := "session_123456"
    
    sessionData := &services.SessionData{
        UserID:   user.ID,
        Username: user.Username,
        Email:    user.Email,
    }
    
    // 创建会话
    if err := sessionService.CreateSession(sessionID, sessionData); err != nil {
        log.Printf("创建会话失败: %v", err)
        return
    }
    log.Println("会话创建成功")
    
    // 获取会话
    retrievedSession, err := sessionService.GetSession(sessionID)
    if err != nil {
        log.Printf("获取会话失败: %v", err)
        return
    }
    log.Printf("获取会话: %+v", retrievedSession)
}
```

## 监控和调试

### 1. Redis监控

```go
// utils/redis_monitor.go
package utils

import (
    "context"
    "log"
    "time"
    "your-project/config"
)

// MonitorRedis 监控Redis状态
func MonitorRedis() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            checkRedisHealth()
        }
    }
}

// checkRedisHealth 检查Redis健康状态
func checkRedisHealth() {
    ctx := context.Background()
    
    // 检查连接
    _, err := config.RedisClient.Ping(ctx).Result()
    if err != nil {
        log.Printf("Redis连接异常: %v", err)
        return
    }
    
    // 获取连接池状态
    stats := config.RedisClient.PoolStats()
    log.Printf("Redis连接池状态 - 总连接数: %d, 空闲连接数: %d, 过期连接数: %d",
        stats.TotalConns, stats.IdleConns, stats.StaleConns)
    
    // 获取内存使用情况
    info, err := config.RedisClient.Info(ctx, "memory").Result()
    if err == nil {
        log.Printf("Redis内存信息: %s", info)
    }
}
```

### 2. 性能测试

```go
// tests/redis_benchmark_test.go
package tests

import (
    "fmt"
    "testing"
    "time"
    "your-project/config"
    "your-project/utils"
)

func BenchmarkRedisSet(b *testing.B) {
    config.InitRedis()
    defer config.CloseRedis()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        key := fmt.Sprintf("benchmark_key_%d", i)
        utils.SetString(key, "benchmark_value", time.Hour)
    }
}

func BenchmarkRedisGet(b *testing.B) {
    config.InitRedis()
    defer config.CloseRedis()
    
    // 预设数据
    for i := 0; i < 1000; i++ {
        key := fmt.Sprintf("benchmark_key_%d", i)
        utils.SetString(key, "benchmark_value", time.Hour)
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        key := fmt.Sprintf("benchmark_key_%d", i%1000)
        utils.GetString(key)
    }
}
```

## 总结

这个Redis使用文档涵盖了：

1. **基础配置**：环境变量、连接池配置
2. **常用API**：字符串、哈希、列表、集合、有序集合操作
3. **高级功能**：发布订阅、管道、Lua脚本
4. **最佳实践**：缓存服务、分布式锁、会话管理
5. **错误处理**：统一错误处理机制
6. **性能优化**：连接池优化、批量操作
7. **监控调试**：健康检查、性能测试

通过这个文档，您可以在Go项目中高效地使用Redis，实现缓存、会话管理、分布式锁等功能。记住要根据实际业务需求调整配置参数和实现细节。