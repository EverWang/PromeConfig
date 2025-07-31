package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"d/GITVIEW/PromeConfig/backend/internal/models"
	"d/GITVIEW/PromeConfig/backend/pkg/prometheus"
)

// PrometheusHandler Prometheus API处理器
type PrometheusHandler struct {
	db     *gorm.DB
	client *prometheus.Client
}

// NewPrometheusHandler 创建Prometheus API处理器
func NewPrometheusHandler(db *gorm.DB, client *prometheus.Client) *PrometheusHandler {
	return &PrometheusHandler{db: db, client: client}
}

// ReloadConfig 重新加载Prometheus配置
func (h *PrometheusHandler) ReloadConfig(c *gin.Context) {
	// 调用Prometheus API重新加载配置
	err := h.client.ReloadConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reload Prometheus configuration: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Prometheus configuration reloaded successfully"})
}

// GetConfig 获取Prometheus配置
func (h *PrometheusHandler) GetConfig(c *gin.Context) {
	// 调用Prometheus API获取配置
	config, err := h.client.GetConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Prometheus configuration: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

// UploadConfig 上传Prometheus配置文件
func (h *PrometheusHandler) UploadConfig(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("config")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// 确保配置目录存在
	configDir := "./configs"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create config directory"})
		return
	}

	// 保存文件
	destPath := filepath.Join(configDir, "prometheus.yml")
	if err := c.SaveUploadedFile(file, destPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save config file"})
		return
	}

	// 尝试重新加载配置
	if err := h.client.ReloadConfig(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "Config file uploaded successfully, but failed to reload Prometheus: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Config file uploaded and Prometheus reloaded successfully"})
}

// GenerateConfig 生成Prometheus配置文件
func (h *PrometheusHandler) GenerateConfig(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 查询用户的所有监控目标
	var targets []models.Target
	if err := h.db.Where("user_id = ?", userID).Find(&targets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch targets"})
		return
	}

	// 查询用户的所有告警规则
	var alertRules []models.AlertRule
	if err := h.db.Where("user_id = ?", userID).Find(&alertRules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch alert rules"})
		return
	}

	// 生成配置文件内容
	config := generatePrometheusConfig(targets, alertRules)

	// 确保配置目录存在
	configDir := "./configs"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create config directory"})
		return
	}

	// 保存配置文件
	configPath := filepath.Join(configDir, "prometheus.yml")
	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save config file"})
		return
	}

	// 尝试重新加载配置
	if err := h.client.ReloadConfig(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "Config file generated successfully, but failed to reload Prometheus: " + err.Error(),
			"config":   config,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Config file generated and Prometheus reloaded successfully",
		"config":  config,
	})
}

// GetAlerts 获取Prometheus告警
func (h *PrometheusHandler) GetAlerts(c *gin.Context) {
	// 调用Prometheus API获取告警
	alerts, err := h.client.GetAlerts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get alerts: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, alerts)
}

// Query 执行PromQL查询
func (h *PrometheusHandler) Query(c *gin.Context) {
	// 获取查询参数
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter is required"})
		return
	}

	// 调用Prometheus API执行查询
	result, err := h.client.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// QueryRange 执行PromQL范围查询
func (h *PrometheusHandler) QueryRange(c *gin.Context) {
	// 获取查询参数
	query := c.Query("query")
	start := c.Query("start")
	end := c.Query("end")
	step := c.Query("step")

	if query == "" || start == "" || end == "" || step == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query, start, end, and step parameters are required"})
		return
	}

	// 调用Prometheus API执行范围查询
	result, err := h.client.QueryRange(query, start, end, step)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute range query: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GenerateAlertRule 使用AI生成告警规则
func (h *PrometheusHandler) GenerateAlertRule(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 解析请求体
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// 解析请求参数
	var request struct {
		Description string `json:"description"`
	}
	if err := json.Unmarshal(body, &request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if request.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Description is required"})
		return
	}

	// 查询用户的AI设置
	var aiSettings models.AISettings
	result := h.db.Where("user_id = ?", userID).First(&aiSettings)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "AI settings not configured"})
		return
	}

	// TODO: 实现AI生成告警规则的逻辑
	// 这里需要根据用户的AI设置调用相应的AI服务生成告警规则
	// 为了示例，我们返回一个模拟的告警规则

	// 创建示例标签和注释
	labelsJSON, _ := json.Marshal(map[string]string{
		"severity": "critical",
	})
	annotationsJSON, _ := json.Marshal(map[string]string{
		"summary": "Generated alert rule",
		"description": request.Description,
	})

	alertRule := models.AlertRule{
		AlertName:   "Generated Alert",
		Expr:        "up == 0",
		ForDuration: "5m",
		Severity:    "warning",
		Labels:      labelsJSON,
		Annotations: annotationsJSON,
	}

	c.JSON(http.StatusOK, alertRule)
}

// 生成Prometheus配置文件内容
func generatePrometheusConfig(targets []models.Target, alertRules []models.AlertRule) string {
	// Generate alerts.yml file
	generateAlertsFile(alertRules)
	
	config := `global:
  scrape_interval: 15s
  evaluation_interval: 15s

alerting:
  alertmanagers:
  - static_configs:
    - targets:
      - alertmanager:9093

rule_files:
  - "alerts.yml"

scrape_configs:
`

	// 添加默认的Prometheus job
	config += `  - job_name: "prometheus"
    static_configs:
    - targets: ["localhost:9090"]
`

	// 添加用户配置的targets
	for _, target := range targets {
		config += `
  - job_name: "` + target.JobName + `"
`
		if target.ScrapeInterval != "" {
			config += `    scrape_interval: ` + target.ScrapeInterval + `
`
		}
		if target.MetricsPath != "" && target.MetricsPath != "/metrics" {
			config += `    metrics_path: ` + target.MetricsPath + `
`
		}

		// 添加relabel_configs
		if target.RelabelConfigs != nil {
			var relabelConfigs []models.RelabelConfig
			if err := json.Unmarshal(target.RelabelConfigs, &relabelConfigs); err == nil && len(relabelConfigs) > 0 {
				config += `    relabel_configs:
`
				for _, relabel := range relabelConfigs {
					config += `      - `
					if len(relabel.SourceLabels) > 0 {
						config += `source_labels: [`
						for i, label := range relabel.SourceLabels {
							if i > 0 {
								config += `, `
							}
							config += `'` + label + `'`
						}
						config += `]\n        `
					}
					if relabel.Separator != "" {
						config += `separator: '` + relabel.Separator + `'\n        `
					}
					if relabel.TargetLabel != "" {
						config += `target_label: ` + relabel.TargetLabel + `\n        `
					}
					if relabel.Regex != "" {
						config += `regex: '` + relabel.Regex + `'\n        `
					}
					if relabel.Replacement != "" {
						config += `replacement: '` + relabel.Replacement + `'\n        `
					}
					if relabel.Action != "" {
						config += `action: ` + relabel.Action + `\n`
					}
				}
			}
		}

		// 添加static_configs
		config += `    static_configs:
    - targets: [`
		var targets []string
		if err := json.Unmarshal(target.Targets, &targets); err == nil {
			for i, targetAddr := range targets {
				if i > 0 {
					config += `, `
				}
				config += `"` + targetAddr + `"`
			}
		}
		config += `]\n`

		// 添加metric_relabel_configs
		if target.MetricRelabelConfigs != nil {
			var metricRelabelConfigs []models.RelabelConfig
			if err := json.Unmarshal(target.MetricRelabelConfigs, &metricRelabelConfigs); err == nil && len(metricRelabelConfigs) > 0 {
				config += `    metric_relabel_configs:
`
				for _, relabel := range metricRelabelConfigs {
					config += `      - `
					if len(relabel.SourceLabels) > 0 {
						config += `source_labels: [`
						for i, label := range relabel.SourceLabels {
							if i > 0 {
								config += `, `
							}
							config += `'` + label + `'`
						}
						config += `]\n        `
					}
					if relabel.Separator != "" {
						config += `separator: '` + relabel.Separator + `'\n        `
					}
					if relabel.TargetLabel != "" {
						config += `target_label: ` + relabel.TargetLabel + `\n        `
					}
					if relabel.Regex != "" {
						config += `regex: '` + relabel.Regex + `'\n        `
					}
					if relabel.Replacement != "" {
						config += `replacement: '` + relabel.Replacement + `'\n        `
					}
					if relabel.Action != "" {
						config += `action: ` + relabel.Action + `\n`
				}
			}
		}
	}
	}

	return config
}

// generateAlertsFile 生成告警规则文件
func generateAlertsFile(alertRules []models.AlertRule) {
	if len(alertRules) == 0 {
		return
	}

	alertsConfig := `groups:
  - name: user_alerts
    rules:
`

	for _, rule := range alertRules {
		alertsConfig += `      - alert: ` + rule.AlertName + `
`
		alertsConfig += `        expr: ` + rule.Expr + `
`
		if rule.ForDuration != "" {
			alertsConfig += `        for: ` + rule.ForDuration + `
`
		}

		// 添加labels
		if len(rule.Labels) > 0 {
			var labels map[string]string
			if err := json.Unmarshal(rule.Labels, &labels); err == nil && len(labels) > 0 {
				alertsConfig += `        labels:
`
				for key, value := range labels {
					alertsConfig += `          ` + key + `: "` + value + `"
`
				}
			}
		}

		// 添加annotations
		if len(rule.Annotations) > 0 {
			var annotations map[string]string
			if err := json.Unmarshal(rule.Annotations, &annotations); err == nil && len(annotations) > 0 {
				alertsConfig += `        annotations:
`
				for key, value := range annotations {
					alertsConfig += `          ` + key + `: "` + value + `"
`
				}
			}
		}
	}

	// 确保配置目录存在
	configDir := "./configs"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return
	}

	// 保存告警规则文件
	alertsPath := filepath.Join(configDir, "alerts.yml")
	if err := os.WriteFile(alertsPath, []byte(alertsConfig), 0644); err != nil {
		return
	}
}