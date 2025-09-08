package config

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"goDDD1/models"
	"log"
	"os"
)

// DBConfig 数据库配置结构体
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// Database 全局数据库连接实例
var Database *gorm.DB

// InitDB 初始化数据库连接
func InitDB() *gorm.DB {
	// 加载.env文件中的环境变量
	err := godotenv.Load()
	if err != nil {
		log.Println("未找到.env文件，将使用默认配置")
	}

	// 从环境变量中获取数据库配置
	dbConfig := DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "3306"),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "goDDD1"),
	}

	// 构建数据库连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName)

	// 连接数据库
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 设置连接池
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	// 启用日志
	db.LogMode(true)

	// 自动迁移数据库表结构
	db.AutoMigrate(&models.User{}, &models.UserWallet{})

	// 保存全局数据库连接实例
	Database = db

	return db
}

// CloseDB 关闭数据库连接
func CloseDB() {
	if Database != nil {
		Database.Close()
	}
}

// 从环境变量获取值，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}