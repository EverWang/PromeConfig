package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"promeconfig-backend/internal/config"
	"promeconfig-backend/internal/database"
	"promeconfig-backend/internal/handlers"
	"promeconfig-backend/internal/middleware"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// 初始化配置
	cfg := config.Load()

	// 初始化数据库
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// 运行数据库迁移
	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// 初始化Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 初始化处理器
	h := handlers.New(db)

	// 公共路由
	public := r.Group("/api")
	{
		public.POST("/auth/signup", h.SignUp)
		public.POST("/auth/signin", h.SignIn)
		public.POST("/auth/refresh", h.RefreshToken)
	}

	// 需要认证的路由
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// 用户相关
		protected.GET("/user", h.GetUser)
		protected.POST("/auth/signout", h.SignOut)

		// Targets管理
		protected.GET("/targets", h.GetTargets)
		protected.POST("/targets", h.CreateTarget)
		protected.PUT("/targets/:id", h.UpdateTarget)
		protected.DELETE("/targets/:id", h.DeleteTarget)

		// Alert Rules管理
		protected.GET("/alert-rules", h.GetAlertRules)
		protected.POST("/alert-rules", h.CreateAlertRule)
		protected.PUT("/alert-rules/:id", h.UpdateAlertRule)
		protected.DELETE("/alert-rules/:id", h.DeleteAlertRule)

		// AI Settings管理
		protected.GET("/ai-settings", h.GetAISettings)
		protected.POST("/ai-settings", h.SaveAISettings)
		protected.DELETE("/ai-settings", h.DeleteAISettings)

		// Prometheus配置管理
		protected.POST("/prometheus/sync", h.SyncPrometheusConfig)
		protected.POST("/prometheus/reload", h.ReloadPrometheusConfig)
		protected.GET("/prometheus/status", h.GetPrometheusStatus)
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(r.Run(":" + port))
}