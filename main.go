package main

import (
	"fmt"
	"goDDD1/config"
	"goDDD1/models"
	"goDDD1/routes"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	err := godotenv.Load()
	if err != nil {
		log.Println("未找到.env文件，将使用默认配置")
	}

	// 初始化数据库
	db := config.InitDB()
	defer config.CloseDB()

	// 初始化Redis
	config.InitRedis()
	defer config.CloseRedis()

	// 自动迁移数据库表结构
	db.AutoMigrate(&models.User{},
		&models.Store{},
		&models.Backpack{},
		&models.UserCurrencyFlow{},
		&models.LevelConfig{},
		&models.LevelHistory{},
	)

	// 设置服务器端口
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	// 设置路由
	router := routes.SetupRouter()

	// 启动服务器
	fmt.Printf("服务器已启动，监听端口: %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}

}
