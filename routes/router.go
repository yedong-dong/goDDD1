package routes

import (
	"goDDD1/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置路由
func SetupRouter() *gin.Engine {
	// 创建默认的gin路由引擎
	r := gin.Default()

	// 创建控制器实例
	backpackController := controllers.NewBackpackController()
	userController := controllers.NewUserController()
	userWalletController := controllers.NewUserWalletController()
	storeController := controllers.NewStoreController()
	userCurrencyFlowController := controllers.NewUserCurrencyFlowController()

	// API路由组
	api := r.Group("/api")
	{
		// 用户相关路由
		users := api.Group("/users")
		{
			users.POST("/register", userController.Register) // 注册用户
			users.GET("/", userController.GetUserByUID)      // 获取用户信息 ?uid=1
			users.GET("/all", userController.GetAllUsers)    // 获取用户信息 ?uid=1
			users.POST("/update", userController.UpdateUser) // 更新用户信息
		}

		// 用户钱包相关路由
		wallets := api.Group("/wallets")
		{
			wallets.GET("/user", userWalletController.GetUserWallets)
			wallets.GET("/user/type", userWalletController.GetWalletByType)        // 获取指定类型钱包 ?user_id=1&type=coin
			wallets.POST("/user/update", userWalletController.UpdateWalletBalance) // 更新钱包余额
		}

		store := api.Group("/store")
		{
			store.POST("/create", storeController.CreateStore)
			store.GET("/get/:id", storeController.GetStoreByID)
			store.POST("/update", storeController.UpdateStore)
			store.POST("/buy", storeController.BuyGoods)
		}

		backpack := api.Group("/backpack")
		{
			backpack.GET("/get", backpackController.GetBackpack)
		}

		userCurrencyFlow := api.Group("/userCurrencyFlow")
		{
			userCurrencyFlow.GET("/get", userCurrencyFlowController.GetUserCurrencyFlow)
		}

	}

	return r
}
