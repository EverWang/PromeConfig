package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"d/GITVIEW/PromeConfig/backend/internal/models"
)

// TargetHandler 监控目标处理器
type TargetHandler struct {
	db *gorm.DB
}

// NewTargetHandler 创建监控目标处理器
func NewTargetHandler(db *gorm.DB) *TargetHandler {
	return &TargetHandler{db: db}
}

// GetTargets 获取所有监控目标
func (h *TargetHandler) GetTargets(c *gin.Context) {
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

	// 查询用户的所有监控目标
	var targets []models.Target
	if err := h.db.Where("user_id = ?", userStr).Order("created_at DESC").Find(&targets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch targets"})
		return
	}

	c.JSON(http.StatusOK, targets)
}

// CreateTarget 创建监控目标
func (h *TargetHandler) CreateTarget(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 解析请求体
	var target models.Target
	if err := c.ShouldBindJSON(&target); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// 设置用户ID
	target.UserID = userStr

	// 创建监控目标
	if err := h.db.Create(&target).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create target"})
		return
	}

	c.JSON(http.StatusCreated, target)
}

// UpdateTarget 更新监控目标
func (h *TargetHandler) UpdateTarget(c *gin.Context) {
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

	// 获取目标ID
	id := c.Param("id")

	// 检查目标是否存在且属于当前用户
	var existingTarget models.Target
	result := h.db.Where("id = ? AND user_id = ?", id, userStr).First(&existingTarget)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Target not found or not authorized"})
		return
	}

	// 解析请求体
	var updatedTarget models.Target
	if err := c.ShouldBindJSON(&updatedTarget); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新允许修改的字段
	updates := map[string]interface{}{
		"job_name":                updatedTarget.JobName,
		"targets":                 updatedTarget.Targets,
		"scrape_interval":         updatedTarget.ScrapeInterval,
		"metrics_path":            updatedTarget.MetricsPath,
		"relabel_configs":         updatedTarget.RelabelConfigs,
		"metric_relabel_configs":  updatedTarget.MetricRelabelConfigs,
	}

	// 更新监控目标
	if err := h.db.Model(&existingTarget).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update target"})
		return
	}

	// 获取更新后的目标
	h.db.First(&existingTarget, id)

	c.JSON(http.StatusOK, existingTarget)
}

// DeleteTarget 删除监控目标
func (h *TargetHandler) DeleteTarget(c *gin.Context) {
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

	// 获取目标ID
	id := c.Param("id")

	// 检查目标是否存在且属于当前用户
	var existingTarget models.Target
	result := h.db.Where("id = ? AND user_id = ?", id, userStr).First(&existingTarget)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Target not found or not authorized"})
		return
	}

	// 删除监控目标
	if err := h.db.Delete(&existingTarget).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete target"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Target deleted successfully"})
}