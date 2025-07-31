package api

import (
	"d/GITVIEW/PromeConfig/backend/config"
	"d/GITVIEW/PromeConfig/backend/internal/middleware"
	"d/GITVIEW/PromeConfig/backend/pkg/jwt"
	"d/GITVIEW/PromeConfig/backend/pkg/prometheus"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes 设置API路由
func SetupRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// 初始化JWT
	jwt.InitJWT(cfg.JWT.Secret, cfg.JWT.ExpiresIn)

	// 创建Prometheus客户端
	promClient := prometheus.NewClient(
		cfg.Prometheus.URL,
		cfg.Prometheus.Username,
		cfg.Prometheus.Password,
	)

	// 创建API处理器
	authHandler := NewAuthHandler(db)
	targetHandler := NewTargetHandler(db)
	alertRuleHandler := NewAlertRuleHandler(db)
	aiSettingsHandler := NewAISettingsHandler(db)
	prometheusHandler := NewPrometheusHandler(db, promClient)

	// API路由组
	api := r.Group("/api")

	// 认证路由
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", middleware.Auth(), authHandler.Logout)
		auth.GET("/user", middleware.Auth(), authHandler.GetUser)
	}

	// 需要认证的路由
	authorized := api.Group("/")
	authorized.Use(middleware.Auth())
	{
		// 监控目标路由
		targets := authorized.Group("/targets")
		{
			targets.GET("", targetHandler.GetTargets)
			targets.POST("", targetHandler.CreateTarget)
			targets.PUT("/:id", targetHandler.UpdateTarget)
			targets.DELETE("/:id", targetHandler.DeleteTarget)
		}

		// 告警规则路由
		alertRules := authorized.Group("/alertrules")
		{
			alertRules.GET("", alertRuleHandler.GetAlertRules)
			alertRules.POST("", alertRuleHandler.CreateAlertRule)
			alertRules.PUT("/:id", alertRuleHandler.UpdateAlertRule)
			alertRules.DELETE("/:id", alertRuleHandler.DeleteAlertRule)
		}

		// AI设置路由
		aiSettings := authorized.Group("/aisettings")
		{
			aiSettings.GET("", aiSettingsHandler.GetAISettings)
			aiSettings.POST("", aiSettingsHandler.SaveAISettings)
			aiSettings.DELETE("", aiSettingsHandler.DeleteAISettings)
		}

		// Prometheus API路由
		prometheus := authorized.Group("/prometheus")
		{
			prometheus.POST("/reload", prometheusHandler.ReloadConfig)
			prometheus.GET("/config", prometheusHandler.GetConfig)
			prometheus.POST("/config", prometheusHandler.UploadConfig)
			prometheus.GET("/alerts", prometheusHandler.GetAlerts)
			prometheus.GET("/query", prometheusHandler.Query)
		}
	}
}