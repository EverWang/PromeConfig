package models

import (
	"fmt"

	"d/GITVIEW/PromeConfig/backend/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB 初始化数据库连接
func InitDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// MigrateDB 自动迁移数据库模型
func MigrateDB(db *gorm.DB) error {
	// 逐个迁移模型以便定位问题
	if err := db.AutoMigrate(&User{}); err != nil {
		return fmt.Errorf("failed to migrate User: %w", err)
	}
	
	if err := db.AutoMigrate(&Target{}); err != nil {
		return fmt.Errorf("failed to migrate Target: %w", err)
	}
	
	if err := db.AutoMigrate(&AlertRule{}); err != nil {
		return fmt.Errorf("failed to migrate AlertRule: %w", err)
	}
	
	if err := db.AutoMigrate(&AISettings{}); err != nil {
		return fmt.Errorf("failed to migrate AISettings: %w", err)
	}
	
	return nil
}