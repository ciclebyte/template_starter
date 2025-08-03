package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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
	ID           int64  `json:"id"`
	TemplateId   int64  `json:"templateId"`
	Name         string `json:"name"`
	VariableType string `json:"variableType"`
	Description  string `json:"description"`
	DefaultValue string `json:"defaultValue"`
	IsRequired   int    `json:"isRequired"`
	Sort         int    `json:"sort"`
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

// TreeNode 树形文件结构
type TreeNode struct {
	ID          int64      `json:"id"`
	FilePath    string     `json:"filePath"`
	FileName    string     `json:"fileName"`
	FileContent string     `json:"fileContent"`
	FileSize    int64      `json:"fileSize"`
	IsDirectory int        `json:"isDirectory"`
	ParentID    int64      `json:"parentId"`
	Children    []TreeNode `json:"children"`
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

// GetTemplateVariables 获取模板变量
func (c *Client) GetTemplateVariables(templateID string) ([]TemplateVariable, error) {
	endpoint := fmt.Sprintf("/api/v1/templates/%s/variables", url.PathEscape(templateID))
	
	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("请求模板变量失败: %w", err)
	}
	
	if resp.Code != 0 {
		return nil, fmt.Errorf("获取模板变量失败: %s", resp.Message)
	}
	
	// 解析响应结构
	var variablesResponse struct {
		CustomVariables []TemplateVariable `json:"customVariables"`
	}
	if err := json.Unmarshal(resp.Data, &variablesResponse); err != nil {
		return nil, fmt.Errorf("解析模板变量失败: %w", err)
	}
	
	return variablesResponse.CustomVariables, nil
}

// RenderTemplate 渲染模板
func (c *Client) RenderTemplate(templateID string, variables map[string]interface{}) ([]RenderedFile, error) {
	endpoint := "/api/v1/templateFiles/renderFileTree"
	
	requestBody := map[string]interface{}{
		"templateId": templateID,
		"variables":  variables,
	}
	
	fmt.Printf("🔄 发送渲染请求: templateId=%s, variables=%+v\n", templateID, variables)
	
	resp, err := c.makeRequest("POST", endpoint, requestBody)
	if err != nil {
		return nil, fmt.Errorf("请求模板渲染失败: %w", err)
	}
	
	if resp.Code != 0 {
		return nil, fmt.Errorf("模板渲染失败: %s", resp.Message)
	}
	
	fmt.Printf("📥 收到渲染响应 (长度: %d): %s\n", len(resp.Data), string(resp.Data)[:min(500, len(resp.Data))])
	
	// 解析树形响应结构
	var renderResponse struct {
		TemplateID int64      `json:"templateId"`
		Tree       []TreeNode `json:"tree"`
	}
	if err := json.Unmarshal(resp.Data, &renderResponse); err != nil {
		return nil, fmt.Errorf("解析渲染结果失败: %w", err)
	}
	
	fmt.Printf("✅ 解析树形结构成功，根节点数量: %d\n", len(renderResponse.Tree))
	
	// 将树形结构转换为平铺的文件列表
	var renderedFiles []RenderedFile
	for _, node := range renderResponse.Tree {
		flattenTreeNode(node, &renderedFiles)
	}
	
	fmt.Printf("📄 转换后的文件数量: %d\n", len(renderedFiles))
	return renderedFiles, nil
}

// flattenTreeNode 将树形节点转换为平铺的文件列表
func flattenTreeNode(node TreeNode, files *[]RenderedFile) {
	// 使用node.FilePath作为完整路径，将反斜杠转换为正斜杠
	fullPath := strings.ReplaceAll(node.FilePath, "\\", "/")
	
	// 添加当前节点
	*files = append(*files, RenderedFile{
		Path:        fullPath,
		Content:     node.FileContent,
		IsDirectory: node.IsDirectory == 1,
	})
	
	// 递归处理子节点
	for _, child := range node.Children {
		flattenTreeNode(child, files)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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