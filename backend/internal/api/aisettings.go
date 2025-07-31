package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"d/GITVIEW/PromeConfig/backend/internal/models"
)

// AISettingsHandler AI设置处理器
type AISettingsHandler struct {
	db *gorm.DB
}

// NewAISettingsHandler 创建AI设置处理器
func NewAISettingsHandler(db *gorm.DB) *AISettingsHandler {
	return &AISettingsHandler{db: db}
}

// GetAISettings 获取用户的AI设置
func (h *AISettingsHandler) GetAISettings(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// 查询用户的AI设置
	var aiSettings models.AISettings
	result := h.db.Where("user_id = ?", userStr).First(&aiSettings)

	// 如果没有找到设置，返回空对象
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	// 如果查询出错（非记录不存在错误）
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch AI settings"})
		return
	}

	c.JSON(http.StatusOK, aiSettings)
}

// SaveAISettings 保存用户的AI设置
func (h *AISettingsHandler) SaveAISettings(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// 解析请求体
	var aiSettings models.AISettings
	if err := c.ShouldBindJSON(&aiSettings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置用户ID
	aiSettings.UserID = userStr

	// 检查是否已存在设置
	var existingSettings models.AISettings
	result := h.db.Where("user_id = ?", userStr).First(&existingSettings)

	// 如果已存在，则更新
	if result.RowsAffected > 0 {
		// 更新允许修改的字段
		updates := map[string]interface{}{
			"provider":  aiSettings.Provider,
			"api_key":   aiSettings.APIKey,
			"base_url":  aiSettings.BaseURL,
			"model":     aiSettings.Model,
			"temperature": aiSettings.Temperature,
		}

		// 更新AI设置
		if err := h.db.Model(&existingSettings).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update AI settings"})
			return
		}

		// 获取更新后的设置
		h.db.Where("user_id = ?", userID).First(&existingSettings)
		c.JSON(http.StatusOK, existingSettings)
		return
	}

	// 如果不存在，则创建
	if err := h.db.Create(&aiSettings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create AI settings"})
		return
	}

	c.JSON(http.StatusCreated, aiSettings)
}

// DeleteAISettings 删除用户的AI设置
func (h *AISettingsHandler) DeleteAISettings(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// 查询用户的AI设置
	var aiSettings models.AISettings
	result := h.db.Where("user_id = ?", userStr).First(&aiSettings)

	// 如果没有找到设置
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "AI settings not found"})
		return
	}

	// 删除AI设置
	if err := h.db.Delete(&aiSettings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete AI settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "AI settings deleted successfully"})
}