package routes

import (
	"goDDD1/controllers"
	"goDDD1/middleware"

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
	authorController := controllers.NewAuthorizationController()
	vueController := controllers.NewVueController()

	public := r.Group("/api")
	{
		test := public.Group("/test")
		{
			test.GET("/list", vueController.Table)
		}

		author := public.Group("/author")
		{
			author.POST("/register", authorController.Register)              // 注册用户
			author.POST("/login", authorController.Login)                    // 登录用户
			author.POST("/send_code", authorController.SendVerificationCode) // 发送验证码
			author.GET("/info", vueController.Info)
		}
	}

	// API路由组
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		// 用户相关路由
		users := protected.Group("/users")
		{
			users.POST("/register", userController.Register) // 注册用户
			users.GET("/", userController.GetUserByUID)      // 获取用户信息 ?uid=1
			users.GET("/all", userController.GetAllUsers)    // 获取用户信息 ?uid=1
			users.POST("/update", userController.UpdateUser) // 更新用户信息
		}

		// 用户钱包相关路由
		wallets := protected.Group("/wallets")
		{
			wallets.GET("/user", userWalletController.GetUserWallets)
			wallets.GET("/user/type", userWalletController.GetWalletByType)        // 获取指定类型钱包 ?user_id=1&type=coin
			wallets.POST("/user/update", userWalletController.UpdateWalletBalance) // 更新钱包余额
		}

		store := protected.Group("/store")
		{
			store.POST("/create", storeController.CreateStore)
			store.GET("/get", storeController.GetStoreByID)
			store.POST("/update", storeController.UpdateStore)
			store.POST("/buy", storeController.BuyGoods)
			store.GET("/tag", storeController.GetStoreByTag)
			store.GET("/tag/page", storeController.GetStoreByTagPage)
			store.GET("/all", storeController.GetAllStores)
		}

		backpack := protected.Group("/backpack")
		{
			backpack.GET("/get", backpackController.GetBackpack)
		}

		userCurrencyFlow := protected.Group("/userCurrencyFlow")
		{
			userCurrencyFlow.GET("/get", userCurrencyFlowController.GetUserCurrencyFlow)
		}

		// 在 SetupRouter 函数中添加
		// 创建控制器实例
		levelController := controllers.NewLevelController()

		// 在 protected 路由组中添加
		// 用户等级相关路由
		level := protected.Group("/level")
		{
			level.GET("/user", levelController.GetUserLevel)          // 获取用户等级信息
			level.GET("/history", levelController.GetLevelHistory)    // 获取用户等级历史记录
			level.GET("/configs", levelController.GetAllLevelConfigs) // 获取所有等级配置
		}
	}

	return r
}
