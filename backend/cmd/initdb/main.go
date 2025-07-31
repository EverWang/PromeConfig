package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	"d/GITVIEW/PromeConfig/backend/config"
	_ "github.com/lib/pq"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 构建数据库连接字符串
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
	)
	
	fmt.Printf("Connecting with DSN: %s\n", dsn)

	// 连接数据库
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("Connected to PostgreSQL database successfully!")

	// 读取初始化SQL文件
	sqlContent, err := ioutil.ReadFile("init.sql")
	if err != nil {
		log.Fatal("Failed to read init.sql:", err)
	}

	// 执行初始化SQL
	_, err = db.Exec(string(sqlContent))
	if err != nil {
		log.Fatal("Failed to execute init.sql:", err)
	}

	fmt.Println("Database initialized successfully!")
	fmt.Println("Admin user created: admin@example.com / password123")
}