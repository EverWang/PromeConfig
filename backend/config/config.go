package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config 应用程序配置结构
type Config struct {
	Environment string
	Server      ServerConfig
	Database    DatabaseConfig
	JWT         JWTConfig
	Prometheus  PrometheusConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret    string
	ExpiresIn time.Duration
}

// PrometheusConfig Prometheus配置
type PrometheusConfig struct {
	URL      string
	Username string
	Password string
}

// Load 加载配置
func Load() (*Config, error) {
	// 加载.env文件
	godotenv.Load()
	
	// 设置默认值
	viper.SetDefault("environment", "development")
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.read_timeout", "10s")
	viper.SetDefault("server.write_timeout", "10s")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("jwt.expires_in", "24h")

	// 从环境变量加载配置
	viper.AutomaticEnv()

	// 从配置文件加载
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../../config")

	if err := viper.ReadInConfig(); err != nil {
		// 配置文件不存在时不返回错误
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	// 优先使用环境变量
	dbHost := getEnv("DB_HOST", viper.GetString("database.host"))
	dbPort := getEnv("DB_PORT", viper.GetString("database.port"))
	dbUser := getEnv("DB_USER", viper.GetString("database.user"))
	dbPassword := getEnv("DB_PASSWORD", viper.GetString("database.password"))
	dbName := getEnv("DB_NAME", viper.GetString("database.dbname"))
	dbSSLMode := getEnv("DB_SSLMODE", viper.GetString("database.sslmode"))

	jwtSecret := getEnv("JWT_SECRET", viper.GetString("jwt.secret"))
	jwtExpiresIn := getEnv("JWT_EXPIRES_IN", viper.GetString("jwt.expires_in"))

	promURL := getEnv("PROMETHEUS_URL", viper.GetString("prometheus.url"))
	promUsername := getEnv("PROMETHEUS_USERNAME", viper.GetString("prometheus.username"))
	promPassword := getEnv("PROMETHEUS_PASSWORD", viper.GetString("prometheus.password"))

	// 解析时间
	readTimeout, err := time.ParseDuration(getEnv("SERVER_READ_TIMEOUT", viper.GetString("server.read_timeout")))
	if err != nil {
		readTimeout = 10 * time.Second
	}

	writeTimeout, err := time.ParseDuration(getEnv("SERVER_WRITE_TIMEOUT", viper.GetString("server.write_timeout")))
	if err != nil {
		writeTimeout = 10 * time.Second
	}

	jwtDuration, err := time.ParseDuration(jwtExpiresIn)
	if err != nil {
		jwtDuration = 24 * time.Hour
	}

	return &Config{
		Environment: getEnv("ENVIRONMENT", viper.GetString("environment")),
		Server: ServerConfig{
			Port:         getEnv("PORT", viper.GetString("server.port")),
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
		Database: DatabaseConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			DBName:   dbName,
			SSLMode:  dbSSLMode,
		},
		JWT: JWTConfig{
			Secret:    jwtSecret,
			ExpiresIn: jwtDuration,
		},
		Prometheus: PrometheusConfig{
			URL:      promURL,
			Username: promUsername,
			Password: promPassword,
		},
	}, nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvBool 获取布尔类型的环境变量
func getEnvBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}