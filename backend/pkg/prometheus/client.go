package prometheus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client Prometheus客户端
type Client struct {
	BaseURL  string
	Username string
	Password string
	Client   *http.Client
}

// NewClient 创建新的Prometheus客户端
func NewClient(baseURL, username, password string) *Client {
	return &Client{
		BaseURL:  baseURL,
		Username: username,
		Password: password,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Reload 重新加载Prometheus配置
func (c *Client) Reload() error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/-/reload", c.BaseURL), nil)
	if err != nil {
		return err
	}

	// 添加基本认证
	if c.Username != "" && c.Password != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to reload Prometheus: %s, status code: %d", string(body), resp.StatusCode)
	}

	return nil
}

// GetConfig 获取Prometheus配置
func (c *Client) GetConfig() (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/status/config", c.BaseURL), nil)
	if err != nil {
		return "", err
	}

	// 添加基本认证
	if c.Username != "" && c.Password != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get Prometheus config: %s, status code: %d", string(body), resp.StatusCode)
	}

	var result struct {
		Status string `json:"status"`
		Data   struct {
			YAML string `json:"yaml"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Data.YAML, nil
}

// GetAlerts 获取当前告警
func (c *Client) GetAlerts() ([]map[string]interface{}, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/alerts", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}

	// 添加基本认证
	if c.Username != "" && c.Password != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get alerts: %s, status code: %d", string(body), resp.StatusCode)
	}

	var result struct {
		Status string `json:"status"`
		Data   struct {
			Alerts []map[string]interface{} `json:"alerts"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data.Alerts, nil
}

// Query 执行PromQL查询
func (c *Client) Query(query string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/query?query=%s", c.BaseURL, query), nil)
	if err != nil {
		return nil, err
	}

	// 添加基本认证
	if c.Username != "" && c.Password != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to execute query: %s, status code: %d", string(body), resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// ReloadConfig 重新加载Prometheus配置（别名）
func (c *Client) ReloadConfig() error {
	return c.Reload()
}

// QueryRange 执行PromQL范围查询
func (c *Client) QueryRange(query, start, end, step string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/query_range?query=%s&start=%s&end=%s&step=%s", 
		c.BaseURL, query, start, end, step)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 添加基本认证
	if c.Username != "" && c.Password != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to execute range query: %s, status code: %d", string(body), resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// UploadConfig 上传Prometheus配置
func (c *Client) UploadConfig(config string, path string) error {
	// 这个函数需要根据实际Prometheus部署方式来实现
	// 这里只是一个示例，实际上可能需要通过SSH或其他方式将配置文件上传到Prometheus服务器
	return fmt.Errorf("not implemented")
}
