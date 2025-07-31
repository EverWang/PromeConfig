package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"d/GITVIEW/PromeConfig/backend/internal/models"
)

// AlertRuleHandler 告警规则处理器
type AlertRuleHandler struct {
	db *gorm.DB
}

// NewAlertRuleHandler 创建告警规则处理器
func NewAlertRuleHandler(db *gorm.DB) *AlertRuleHandler {
	return &AlertRuleHandler{db: db}
}

// GetAlertRules 获取所有告警规则
func (h *AlertRuleHandler) GetAlertRules(c *gin.Context) {
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

	// 查询用户的所有告警规则
	var alertRules []models.AlertRule
	if err := h.db.Where("user_id = ?", userStr).Order("created_at DESC").Find(&alertRules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch alert rules"})
		return
	}

	c.JSON(http.StatusOK, alertRules)
}

// CreateAlertRule 创建告警规则
func (h *AlertRuleHandler) CreateAlertRule(c *gin.Context) {
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
	var alertRule models.AlertRule
	if err := c.ShouldBindJSON(&alertRule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置用户ID
	alertRule.UserID = userStr

	// 创建告警规则
	if err := h.db.Create(&alertRule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create alert rule"})
		return
	}

	c.JSON(http.StatusCreated, alertRule)
}

// UpdateAlertRule 更新告警规则
func (h *AlertRuleHandler) UpdateAlertRule(c *gin.Context) {
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

	// 获取规则ID
	id := c.Param("id")

	// 检查规则是否存在且属于当前用户
	var existingRule models.AlertRule
	result := h.db.Where("id = ? AND user_id = ?", id, userStr).First(&existingRule)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Alert rule not found or not authorized"})
		return
	}

	// 解析请求体
	var updatedRule models.AlertRule
	if err := c.ShouldBindJSON(&updatedRule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新允许修改的字段
	updates := map[string]interface{}{
		"alert_name":   updatedRule.AlertName,
		"expr":         updatedRule.Expr,
		"for_duration": updatedRule.ForDuration,
		"severity":     updatedRule.Severity,
		"labels":       updatedRule.Labels,
		"annotations":  updatedRule.Annotations,
	}

	// 更新告警规则
	if err := h.db.Model(&existingRule).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update alert rule"})
		return
	}

	// 获取更新后的规则
	h.db.First(&existingRule, id)

	c.JSON(http.StatusOK, existingRule)
}

// DeleteAlertRule 删除告警规则
func (h *AlertRuleHandler) DeleteAlertRule(c *gin.Context) {
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

	// 获取规则ID
	id := c.Param("id")

	// 检查规则是否存在且属于当前用户
	var existingRule models.AlertRule
	result := h.db.Where("id = ? AND user_id = ?", id, userStr).First(&existingRule)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Alert rule not found or not authorized"})
		return
	}

	// 删除告警规则
	if err := h.db.Delete(&existingRule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete alert rule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Alert rule deleted successfully"})
}