package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"d/GITVIEW/PromeConfig/backend/internal/models"
	"d/GITVIEW/PromeConfig/backend/pkg/jwt"
	"d/GITVIEW/PromeConfig/backend/pkg/password"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	db *gorm.DB
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse 认证响应
type AuthResponse struct {
	Token string       `json:"token"`
	User  models.User `json:"user"`
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查邮箱是否已存在
	var existingUser models.User
	result := h.db.Where("email = ?", req.Email).First(&existingUser)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// 哈希密码
	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// 创建用户
	user := models.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// 生成JWT令牌
	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		Token: token,
		User:  user,
	})
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	log.Printf("🔐 收到登录请求 - IP: %s, UserAgent: %s", c.ClientIP(), c.GetHeader("User-Agent"))
	
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ 登录请求参数错误: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("📧 尝试登录用户: %s, 密码长度: %d", req.Email, len(req.Password))

	// 查找用户
	var user models.User
	result := h.db.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		log.Printf("❌ 数据库查询错误: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if result.RowsAffected == 0 {
		log.Printf("❌ 用户不存在: %s", req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	log.Printf("✅ 找到用户: ID=%s, Email=%s, 存储的密码哈希长度: %d", user.ID, user.Email, len(user.PasswordHash))

	// 验证密码
	log.Printf("🔍 开始验证密码...") 
	isValid := password.Verify(user.PasswordHash, req.Password)
	log.Printf("🔍 密码验证结果: %v", isValid)
	if !isValid {
		log.Printf("❌ 密码验证失败 - 用户: %s", req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// 生成JWT令牌
	log.Printf("🎫 开始生成JWT令牌...") 
	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		log.Printf("❌ JWT令牌生成失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	log.Printf("✅ 登录成功 - 用户: %s, 令牌长度: %d", user.Email, len(token))
	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User:  user,
	})
}

// Logout 用户登出
func (h *AuthHandler) Logout(c *gin.Context) {
	// 客户端应该删除令牌，服务器端不需要做特殊处理
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// GetUser 获取当前用户信息
func (h *AuthHandler) GetUser(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 查找用户
	var user models.User
	result := h.db.First(&user, "id = ?", userID)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}