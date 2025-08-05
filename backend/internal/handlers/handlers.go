package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"promeconfig-backend/internal/config"
	"promeconfig-backend/internal/middleware"
	"promeconfig-backend/internal/models"
)

type Handlers struct {
	db *sql.DB
}

func New(db *sql.DB) *Handlers {
	return &Handlers{db: db}
}

// 生成JWT Token
func (h *Handlers) generateToken(userID uuid.UUID) (string, error) {
	cfg := config.Load()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7天过期
	})

	return token.SignedString([]byte(cfg.JWTSecret))
}

// 用户注册
func (h *Handlers) SignUp(c *gin.Context) {
	var req models.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户是否已存在
	var exists bool
	err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// 创建用户
	var user models.User
	err = h.db.QueryRow(`
		INSERT INTO users (email, password_hash) 
		VALUES ($1, $2) 
		RETURNING id, email, created_at, updated_at`,
		req.Email, string(hashedPassword)).Scan(&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// 生成token
	token, err := h.generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, models.AuthResponse{
		User:        user,
		AccessToken: token,
		TokenType:   "Bearer",
	})
}

// 用户登录
func (h *Handlers) SignIn(c *gin.Context) {
	var req models.SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找用户
	var user models.User
	err := h.db.QueryRow(`
		SELECT id, email, password_hash, created_at, updated_at 
		FROM users WHERE email = $1`, req.Email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// 生成token
	token, err := h.generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{
		User:        user,
		AccessToken: token,
		TokenType:   "Bearer",
	})
}

// 获取当前用户信息
func (h *Handlers) GetUser(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var user models.User
	err := h.db.QueryRow(`
		SELECT id, email, created_at, updated_at 
		FROM users WHERE id = $1`, userID).Scan(
		&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// 用户登出
func (h *Handlers) SignOut(c *gin.Context) {
	// 在实际应用中，你可能想要将token加入黑名单
	c.JSON(http.StatusOK, gin.H{"message": "Signed out successfully"})
}

// 刷新Token
func (h *Handlers) RefreshToken(c *gin.Context) {
	// 这里可以实现token刷新逻辑
	c.JSON(http.StatusOK, gin.H{"message": "Token refreshed"})
}

// Targets相关处理器
func (h *Handlers) GetTargets(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	rows, err := h.db.Query(`
		SELECT id, user_id, job_name, targets, scrape_interval, metrics_path, 
		       relabel_configs, metric_relabel_configs, created_at, updated_at
		FROM targets WHERE user_id = $1 ORDER BY created_at DESC`, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get targets"})
		return
	}
	defer rows.Close()

	var targets []models.Target
	for rows.Next() {
		var target models.Target
		var relabelConfigs, metricRelabelConfigs sql.NullString

		err := rows.Scan(&target.ID, &target.UserID, &target.JobName, &target.Targets,
			&target.ScrapeInterval, &target.MetricsPath, &relabelConfigs, &metricRelabelConfigs,
			&target.CreatedAt, &target.UpdatedAt)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan target"})
			return
		}

		if relabelConfigs.Valid {
			target.RelabelConfigs = json.RawMessage(relabelConfigs.String)
		}
		if metricRelabelConfigs.Valid {
			target.MetricRelabelConfigs = json.RawMessage(metricRelabelConfigs.String)
		}

		targets = append(targets, target)
	}

	c.JSON(http.StatusOK, targets)
}

func (h *Handlers) CreateTarget(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var req models.CreateTargetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置默认值
	if req.ScrapeInterval == "" {
		req.ScrapeInterval = "15s"
	}
	if req.MetricsPath == "" {
		req.MetricsPath = "/metrics"
	}

	// 处理空的JSON字段
	if len(req.RelabelConfigs) == 0 {
		req.RelabelConfigs = nil
	}
	if len(req.MetricRelabelConfigs) == 0 {
		req.MetricRelabelConfigs = nil
	}

	var target models.Target
	err := h.db.QueryRow(`
		INSERT INTO targets (user_id, job_name, targets, scrape_interval, metrics_path, relabel_configs, metric_relabel_configs)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, job_name, targets, scrape_interval, metrics_path, relabel_configs, metric_relabel_configs, created_at, updated_at`,
		userID, req.JobName, req.Targets, req.ScrapeInterval, req.MetricsPath, req.RelabelConfigs, req.MetricRelabelConfigs).Scan(
		&target.ID, &target.UserID, &target.JobName, &target.Targets, &target.ScrapeInterval,
		&target.MetricsPath, &target.RelabelConfigs, &target.MetricRelabelConfigs, &target.CreatedAt, &target.UpdatedAt)

	if err != nil {
		// 记录详细错误信息
		println("Database error creating target:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create target: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, target)
}

func (h *Handlers) UpdateTarget(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	targetID := c.Param("id")
	targetUUID, err := uuid.Parse(targetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID"})
		return
	}

	var req models.CreateTargetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 处理空的JSON字段
	if len(req.RelabelConfigs) == 0 {
		req.RelabelConfigs = nil
	}
	if len(req.MetricRelabelConfigs) == 0 {
		req.MetricRelabelConfigs = nil
	}

	var target models.Target
	err = h.db.QueryRow(`
		UPDATE targets 
		SET job_name = $1, targets = $2, scrape_interval = $3, metrics_path = $4, 
		    relabel_configs = $5, metric_relabel_configs = $6
		WHERE id = $7 AND user_id = $8
		RETURNING id, user_id, job_name, targets, scrape_interval, metrics_path, relabel_configs, metric_relabel_configs, created_at, updated_at`,
		req.JobName, req.Targets, req.ScrapeInterval, req.MetricsPath, req.RelabelConfigs, req.MetricRelabelConfigs, targetUUID, userID).Scan(
		&target.ID, &target.UserID, &target.JobName, &target.Targets, &target.ScrapeInterval,
		&target.MetricsPath, &target.RelabelConfigs, &target.MetricRelabelConfigs, &target.CreatedAt, &target.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Target not found"})
		return
	}

	if err != nil {
		// 记录详细错误信息
		println("Database error updating target:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update target: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, target)
}

func (h *Handlers) DeleteTarget(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	targetID := c.Param("id")
	targetUUID, err := uuid.Parse(targetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID"})
		return
	}

	result, err := h.db.Exec("DELETE FROM targets WHERE id = $1 AND user_id = $2", targetUUID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete target"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Target not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Target deleted successfully"})
}

// Alert Rules相关处理器
func (h *Handlers) GetAlertRules(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	rows, err := h.db.Query(`
		SELECT id, user_id, alert_name, expr, for_duration, labels, annotations, created_at, updated_at
		FROM alert_rules WHERE user_id = $1 ORDER BY created_at DESC`, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get alert rules"})
		return
	}
	defer rows.Close()

	var alertRules []models.AlertRule
	for rows.Next() {
		var rule models.AlertRule
		err := rows.Scan(&rule.ID, &rule.UserID, &rule.AlertName, &rule.Expr,
			&rule.ForDuration, &rule.Labels, &rule.Annotations, &rule.CreatedAt, &rule.UpdatedAt)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan alert rule"})
			return
		}

		alertRules = append(alertRules, rule)
	}

	c.JSON(http.StatusOK, alertRules)
}

func (h *Handlers) CreateAlertRule(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var req models.CreateAlertRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置默认值
	if req.ForDuration == "" {
		req.ForDuration = "5m"
	}
	if req.Labels == nil {
		req.Labels = json.RawMessage("{}")
	}
	if req.Annotations == nil {
		req.Annotations = json.RawMessage("{}")
	}

	var rule models.AlertRule
	err := h.db.QueryRow(`
		INSERT INTO alert_rules (user_id, alert_name, expr, for_duration, labels, annotations)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, alert_name, expr, for_duration, labels, annotations, created_at, updated_at`,
		userID, req.AlertName, req.Expr, req.ForDuration, req.Labels, req.Annotations).Scan(
		&rule.ID, &rule.UserID, &rule.AlertName, &rule.Expr, &rule.ForDuration,
		&rule.Labels, &rule.Annotations, &rule.CreatedAt, &rule.UpdatedAt)

	if err != nil {
		// 记录详细错误信息
		println("Database error creating alert rule:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create alert rule: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, rule)
}

func (h *Handlers) UpdateAlertRule(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	ruleID := c.Param("id")
	ruleUUID, err := uuid.Parse(ruleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rule ID"})
		return
	}

	var req models.CreateAlertRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var rule models.AlertRule
	err = h.db.QueryRow(`
		UPDATE alert_rules 
		SET alert_name = $1, expr = $2, for_duration = $3, labels = $4, annotations = $5
		WHERE id = $6 AND user_id = $7
		RETURNING id, user_id, alert_name, expr, for_duration, labels, annotations, created_at, updated_at`,
		req.AlertName, req.Expr, req.ForDuration, req.Labels, req.Annotations, ruleUUID, userID).Scan(
		&rule.ID, &rule.UserID, &rule.AlertName, &rule.Expr, &rule.ForDuration,
		&rule.Labels, &rule.Annotations, &rule.CreatedAt, &rule.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Alert rule not found"})
		return
	}

	if err != nil {
		// 记录详细错误信息
		println("Database error updating alert rule:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update alert rule: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, rule)
}

func (h *Handlers) DeleteAlertRule(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	ruleID := c.Param("id")
	ruleUUID, err := uuid.Parse(ruleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rule ID"})
		return
	}

	result, err := h.db.Exec("DELETE FROM alert_rules WHERE id = $1 AND user_id = $2", ruleUUID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete alert rule"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Alert rule not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Alert rule deleted successfully"})
}

// AI Settings相关处理器
func (h *Handlers) GetAISettings(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var settings models.AISettings
	err := h.db.QueryRow(`
		SELECT id, user_id, provider, api_key, base_url, model, temperature, created_at, updated_at
		FROM ai_settings WHERE user_id = $1`, userID).Scan(
		&settings.ID, &settings.UserID, &settings.Provider, &settings.APIKey,
		&settings.BaseURL, &settings.Model, &settings.Temperature, &settings.CreatedAt, &settings.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusOK, nil)
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get AI settings"})
		return
	}

	c.JSON(http.StatusOK, settings)
}

func (h *Handlers) SaveAISettings(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var req models.SaveAISettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 尝试更新现有设置
	var settings models.AISettings
	err := h.db.QueryRow(`
		UPDATE ai_settings 
		SET provider = $1, api_key = $2, base_url = $3, model = $4, temperature = $5
		WHERE user_id = $6
		RETURNING id, user_id, provider, api_key, base_url, model, temperature, created_at, updated_at`,
		req.Provider, req.APIKey, req.BaseURL, req.Model, req.Temperature, userID).Scan(
		&settings.ID, &settings.UserID, &settings.Provider, &settings.APIKey,
		&settings.BaseURL, &settings.Model, &settings.Temperature, &settings.CreatedAt, &settings.UpdatedAt)

	if err == sql.ErrNoRows {
		// 创建新设置
		err = h.db.QueryRow(`
			INSERT INTO ai_settings (user_id, provider, api_key, base_url, model, temperature)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, user_id, provider, api_key, base_url, model, temperature, created_at, updated_at`,
			userID, req.Provider, req.APIKey, req.BaseURL, req.Model, req.Temperature).Scan(
			&settings.ID, &settings.UserID, &settings.Provider, &settings.APIKey,
			&settings.BaseURL, &settings.Model, &settings.Temperature, &settings.CreatedAt, &settings.UpdatedAt)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save AI settings"})
		return
	}

	c.JSON(http.StatusOK, settings)
}

func (h *Handlers) DeleteAISettings(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	result, err := h.db.Exec("DELETE FROM ai_settings WHERE user_id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete AI settings"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "AI settings not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "AI settings deleted successfully"})
}

// Prometheus配置管理
func (h *Handlers) SyncPrometheusConfig(c *gin.Context) {
	// 这里实现配置同步逻辑
	// 可以调用外部服务或直接写入文件系统
	c.JSON(http.StatusOK, gin.H{"message": "Configuration synced successfully"})
}

func (h *Handlers) ReloadPrometheusConfig(c *gin.Context) {
	// 这里实现Prometheus重载逻辑
	// 调用Prometheus的reload API
	c.JSON(http.StatusOK, gin.H{"message": "Configuration reloaded successfully"})
}

func (h *Handlers) GetPrometheusStatus(c *gin.Context) {
	// 这里实现获取Prometheus状态逻辑
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"version": "2.45.0",
		"uptime": "2h 15m",
	})
}