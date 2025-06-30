package main

import (
	"log"
	"net/http"
	"nav-panel-backend/internal/database"
	"nav-panel-backend/internal/handler"
	"nav-panel-backend/internal/repository"
	"nav-panel-backend/internal/service"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	db, err := database.InitDatabase()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// 运行数据库迁移
	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// 初始化仓库层
	websiteRepo := repository.NewWebsiteRepository(db)
	searchEngineRepo := repository.NewSearchEngineRepository(db)
	settingsRepo := repository.NewSettingsRepository(db)

	// 初始化服务层
	websiteService := service.NewWebsiteService(websiteRepo)
	searchEngineService := service.NewSearchEngineService(searchEngineRepo)
	settingsService := service.NewSettingsService(settingsRepo)
	
	// 为设置服务设置其他仓库依赖（用于导入导出功能）
	settingsService.SetRepositories(websiteRepo, searchEngineRepo)

	// 初始化处理器层
	websiteHandler := handler.NewWebsiteHandler(websiteService)
	searchEngineHandler := handler.NewSearchEngineHandler(searchEngineService)
	settingsHandler := handler.NewSettingsHandler(settingsService)

	// 初始化 Gin 路由
	r := gin.Default()

	// 配置 CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:5173"} // 前端地址
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// 注册路由
	api := r.Group("/api")
	{
		// 健康检查接口
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":    "ok",
				"message":   "service is healthy",
				"timestamp": time.Now().Unix(),
			})
		})

		// 网站路由
		websites := api.Group("/websites")
		{
			websites.GET("", websiteHandler.GetAll)
			websites.POST("", websiteHandler.Create)
			websites.PUT("/reorder", websiteHandler.Reorder)
			websites.PUT("/:id", websiteHandler.Update)
			websites.DELETE("/:id", websiteHandler.Delete)
		}

		// 搜索引擎路由
		searchEngines := api.Group("/search-engines")
		{
			searchEngines.GET("", searchEngineHandler.GetAll)
			searchEngines.POST("", searchEngineHandler.Create)
			searchEngines.POST("/with-icon", searchEngineHandler.CreateWithIcon)
			searchEngines.PUT("/:id", searchEngineHandler.Update)
			searchEngines.PUT("/:id/with-icon", searchEngineHandler.UpdateWithIcon)
			searchEngines.DELETE("/:id", searchEngineHandler.Delete)
		}

		// 设置路由
		settings := api.Group("/settings")
		{
			settings.GET("", settingsHandler.Get)
			settings.PUT("", settingsHandler.Update)
		}

		// 数据导入导出
		api.GET("/export", settingsHandler.Export)
		api.POST("/import", settingsHandler.Import)
	}

	// 启动服务器
	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
} 