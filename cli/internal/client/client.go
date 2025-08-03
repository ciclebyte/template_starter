package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client 模板服务客户端
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	APIKey     string
}

// NewClient 创建新的客户端实例
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		APIKey: apiKey,
	}
}

// Template 模板信息结构
type Template struct {
	ID           int64                  `json:"id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Introduction string                 `json:"introduction"`
	CategoryId   int                    `json:"categoryId"`
	IsFeatured   int                    `json:"isFeatured"`
	Logo         string                 `json:"logo"`
	Icon         string                 `json:"icon"`
	Variables    []TemplateVariable     `json:"variables"`
	Files        []TemplateFile         `json:"files"`
	Languages    []TemplateLanguage     `json:"languages"`
}

// TemplateLanguage 模板语言结构
type TemplateLanguage struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// TemplateVariable 模板变量结构
type TemplateVariable struct {
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	Description  string      `json:"description"`
	DefaultValue interface{} `json:"defaultValue"`
	Required     bool        `json:"required"`
	Options      []string    `json:"options,omitempty"`
}

// TemplateFile 模板文件结构
type TemplateFile struct {
	Path        string `json:"path"`
	Content     string `json:"content"`
	IsDirectory bool   `json:"isDirectory"`
	Condition   string `json:"condition,omitempty"`
}

// RenderedFile 渲染后的文件结构
type RenderedFile struct {
	Path        string `json:"path"`
	Content     string `json:"content"`
	IsDirectory bool   `json:"isDirectory"`
}

// ListTemplates 获取模板列表
func (c *Client) ListTemplates(category string) ([]Template, error) {
	endpoint := "/api/v1/templates/list"
	params := url.Values{}
	if category != "" {
		params.Add("categoryId", category)
	}
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}
	
	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("请求模板列表失败: %w", err)
	}
	
	if resp.Code != 0 {
		return nil, fmt.Errorf("获取模板列表失败: %s", resp.Message)
	}
	
	// 解析嵌套的响应结构
	var listResponse struct {
		TemplatesList []Template `json:"templatesList"`
		Total         int        `json:"total"`
	}
	if err := json.Unmarshal(resp.Data, &listResponse); err != nil {
		return nil, fmt.Errorf("解析模板列表失败: %w", err)
	}
	
	return listResponse.TemplatesList, nil
}

// GetTemplate 获取指定模板详情
func (c *Client) GetTemplate(templateID string) (*Template, error) {
	endpoint := "/api/v1/templates/detail?id=" + url.QueryEscape(templateID)
	
	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("请求模板详情失败: %w", err)
	}
	
	if resp.Code != 0 {
		return nil, fmt.Errorf("获取模板详情失败: %s", resp.Message)
	}
	
	// 解析嵌套的响应结构
	var detailResponse struct {
		TemplatesInfo *Template `json:"data"`
	}
	if err := json.Unmarshal(resp.Data, &detailResponse); err != nil {
		return nil, fmt.Errorf("解析模板详情失败: %w", err)
	}
	
	return detailResponse.TemplatesInfo, nil
}

// GetTemplateInfo 获取模板基本信息（不包含文件内容）
func (c *Client) GetTemplateInfo(templateID string) (*Template, error) {
	endpoint := "/api/v1/templates/detail?id=" + url.QueryEscape(templateID)
	
	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("请求模板信息失败: %w", err)
	}
	
	if resp.Code != 0 {
		return nil, fmt.Errorf("获取模板信息失败: %s", resp.Message)
	}
	
	// 解析嵌套的响应结构
	var detailResponse struct {
		TemplatesInfo *Template `json:"data"`
	}
	if err := json.Unmarshal(resp.Data, &detailResponse); err != nil {
		return nil, fmt.Errorf("解析模板信息失败: %w", err)
	}
	
	return detailResponse.TemplatesInfo, nil
}

// RenderTemplate 渲染模板
func (c *Client) RenderTemplate(templateID string, variables map[string]interface{}) ([]RenderedFile, error) {
	endpoint := "/api/v1/templateFiles/renderFileTree"
	
	requestBody := map[string]interface{}{
		"templateId": templateID,
		"variables":  variables,
	}
	
	resp, err := c.makeRequest("POST", endpoint, requestBody)
	if err != nil {
		return nil, fmt.Errorf("请求模板渲染失败: %w", err)
	}
	
	if resp.Code != 0 {
		return nil, fmt.Errorf("模板渲染失败: %s", resp.Message)
	}
	
	var renderedFiles []RenderedFile
	if err := json.Unmarshal(resp.Data, &renderedFiles); err != nil {
		return nil, fmt.Errorf("解析渲染结果失败: %w", err)
	}
	
	return renderedFiles, nil
}

// SearchTemplates 搜索模板
func (c *Client) SearchTemplates(keyword, category string) ([]Template, error) {
	endpoint := "/api/v1/templates/search"
	params := url.Values{}
	if keyword != "" {
		params.Add("keyword", keyword)
	}
	if category != "" {
		params.Add("category", category)
	}
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}
	
	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("请求模板搜索失败: %w", err)
	}
	
	if resp.Code != 0 {
		return nil, fmt.Errorf("搜索模板失败: %s", resp.Message)
	}
	
	var templates []Template
	if err := json.Unmarshal(resp.Data, &templates); err != nil {
		return nil, fmt.Errorf("解析搜索结果失败: %w", err)
	}
	
	return templates, nil
}

// Response 通用响应结构
type Response struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// makeRequest 发送HTTP请求的通用方法
func (c *Client) makeRequest(method, endpoint string, body interface{}) (*Response, error) {
	// 构建完整URL
	fullURL := c.BaseURL + endpoint
	
	// 准备请求体
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("序列化请求体失败: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}
	
	// 创建HTTP请求
	req, err := http.NewRequest(method, fullURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}
	
	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}
	
	// 发送请求
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("无法连接到服务器 %s: %w\n请检查:\n1. 服务器是否运行\n2. 地址是否正确\n3. 网络连接是否正常", c.BaseURL, err)
	}
	defer resp.Body.Close()
	
	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}
	
	// 检查HTTP状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP请求失败: %d %s", resp.StatusCode, string(respBody))
	}
	
	// 解析响应
	var response Response
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}
	
	return &response, nil
}