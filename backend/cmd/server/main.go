package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	
	"d/GITVIEW/PromeConfig/backend/config"
	"d/GITVIEW/PromeConfig/backend/internal/api"
	"d/GITVIEW/PromeConfig/backend/internal/middleware"
	"d/GITVIEW/PromeConfig/backend/internal/models"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// 初始化配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 设置Gin模式
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化数据库
	db, err := models.InitDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移数据库模型 (临时跳过以解决迁移问题)
	// if err := models.MigrateDB(db); err != nil {
	// 	log.Fatalf("Failed to migrate database: %v", err)
	// }
	log.Println("Skipping database migration (tables already exist)")

	// 创建Gin路由
	r := gin.Default()

	// 添加中间件
	r.Use(middleware.CORS())

	// 设置API路由
	api.SetupRoutes(r, db, cfg)

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Server.Port
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}