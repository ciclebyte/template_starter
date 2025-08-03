package libAI

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ciclebyte/template_starter/internal/model"
	"github.com/ciclebyte/template_starter/library/libConfig"
	"github.com/gogf/gf/v2/frame/g"
)

// AIClient AI客户端接口
type AIClient interface {
	TestConnection(ctx context.Context) error
	GenerateTemplate(ctx context.Context, req *TemplateGenerateRequest) (*TemplateGenerateResponse, error)
	SuggestVariables(ctx context.Context, req *VariableSuggestRequest) (*VariableSuggestResponse, error)
}

// TemplateGenerateRequest 模板生成请求
type TemplateGenerateRequest struct {
	Description   string            `json:"description"`   // 项目描述
	ProjectType   string            `json:"projectType"`   // 项目类型
	TechStack     []string          `json:"techStack"`     // 技术栈
	Variables     map[string]string `json:"variables"`     // 变量
	Features      []string          `json:"features"`      // 功能特性
}

// TemplateGenerateResponse 模板生成响应
type TemplateGenerateResponse struct {
	ProjectStructure []FileInfo            `json:"projectStructure"` // 项目结构
	Variables        []VariableInfo        `json:"variables"`        // 推荐变量
	Instructions     string                `json:"instructions"`     // 使用说明
	EstimatedTime    int                   `json:"estimatedTime"`    // 预估完成时间(分钟)
}

// VariableSuggestRequest 变量建议请求
type VariableSuggestRequest struct {
	ProjectType   string   `json:"projectType"`
	TechStack     []string `json:"techStack"`
	Description   string   `json:"description"`
}

// VariableSuggestResponse 变量建议响应
type VariableSuggestResponse struct {
	Variables []VariableInfo `json:"variables"`
}

// FileInfo 文件信息
type FileInfo struct {
	Path        string `json:"path"`
	Content     string `json:"content"`
	IsDirectory bool   `json:"isDirectory"`
	Description string `json:"description"`
}

// VariableInfo 变量信息
type VariableInfo struct {
	Name         string `json:"name"`
	Type         string `json:"type"`         // string, number, boolean
	Description  string `json:"description"`
	DefaultValue string `json:"defaultValue"`
	Required     bool   `json:"required"`
	Options      []string `json:"options,omitempty"` // 选项值（如果是选择类型）
}

// OpenAIClient OpenAI客户端
type OpenAIClient struct {
	APIKey  string
	BaseURL string
	Model   string
	client  *http.Client
}

// ClaudeClient Claude客户端
type ClaudeClient struct {
	APIKey  string
	BaseURL string
	Model   string
	client  *http.Client
}

// NewAIClient 创建AI客户端
func NewAIClient(ctx context.Context) (AIClient, error) {
	config, err := libConfig.GetAIConfig(ctx)
	if err != nil {
		return nil, err
	}

	if !config.Enabled {
		return nil, fmt.Errorf("AI功能未启用")
	}

	switch config.Provider {
	case "openai":
		return NewOpenAIClient(config.OpenAI.APIKey, config.OpenAI.BaseURL, config.OpenAI.Model), nil
	case "claude":
		return NewClaudeClient(config.Claude.APIKey, config.Claude.BaseURL, config.Claude.Model), nil
	default:
		return nil, fmt.Errorf("不支持的AI服务商: %s", config.Provider)
	}
}

// NewOpenAIClient 创建OpenAI客户端
func NewOpenAIClient(apiKey, baseURL, model string) *OpenAIClient {
	return &OpenAIClient{
		APIKey:  apiKey,
		BaseURL: baseURL,
		Model:   model,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewClaudeClient 创建Claude客户端
func NewClaudeClient(apiKey, baseURL, model string) *ClaudeClient {
	return &ClaudeClient{
		APIKey:  apiKey,
		BaseURL: baseURL,
		Model:   model,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// TestConnection 测试OpenAI连接
func (c *OpenAIClient) TestConnection(ctx context.Context) error {
	reqBody := map[string]interface{}{
		"model": c.Model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": "Hello, this is a connection test.",
			},
		},
		"max_tokens": 10,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("连接失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API调用失败 (状态码: %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// TestConnection 测试Claude连接
func (c *ClaudeClient) TestConnection(ctx context.Context) error {
	reqBody := map[string]interface{}{
		"model":      c.Model,
		"max_tokens": 10,
		"messages": []map[string]string{
			{
				"role":    "user", 
				"content": "Hello, this is a connection test.",
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("连接失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API调用失败 (状态码: %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// GenerateTemplate 生成模板 (OpenAI实现)
func (c *OpenAIClient) GenerateTemplate(ctx context.Context, req *TemplateGenerateRequest) (*TemplateGenerateResponse, error) {
	// 构建prompt
	prompt := buildTemplateGeneratePrompt(req)
	
	reqBody := map[string]interface{}{
		"model": c.Model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"max_tokens": 4000,
		"temperature": 0.7,
	}

	// 发送请求并解析响应
	response, err := c.sendRequest(ctx, "/v1/chat/completions", reqBody)
	if err != nil {
		return nil, err
	}

	// 解析AI响应为模板结构
	return parseTemplateResponse(response)
}

// SuggestVariables 建议变量 (OpenAI实现)
func (c *OpenAIClient) SuggestVariables(ctx context.Context, req *VariableSuggestRequest) (*VariableSuggestResponse, error) {
	prompt := buildVariableSuggestPrompt(req)
	
	reqBody := map[string]interface{}{
		"model": c.Model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"max_tokens": 1000,
		"temperature": 0.5,
	}

	response, err := c.sendRequest(ctx, "/v1/chat/completions", reqBody)
	if err != nil {
		return nil, err
	}

	return parseVariableResponse(response)
}

// GenerateTemplate Claude实现
func (c *ClaudeClient) GenerateTemplate(ctx context.Context, req *TemplateGenerateRequest) (*TemplateGenerateResponse, error) {
	// Claude实现类似，但API格式略有不同
	return nil, fmt.Errorf("Claude模板生成功能开发中")
}

// SuggestVariables Claude实现
func (c *ClaudeClient) SuggestVariables(ctx context.Context, req *VariableSuggestRequest) (*VariableSuggestResponse, error) {
	return nil, fmt.Errorf("Claude变量建议功能开发中")
}

// 辅助方法
func (c *OpenAIClient) sendRequest(ctx context.Context, endpoint string, reqBody interface{}) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API调用失败 (状态码: %d): %s", resp.StatusCode, string(body))
	}

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	return response, err
}

func buildTemplateGeneratePrompt(req *TemplateGenerateRequest) string {
	return fmt.Sprintf(`作为一个专业的软件架构师和代码生成专家，请根据以下需求生成一个完整的项目模板：

项目描述：%s
项目类型：%s
技术栈：%v
功能特性：%v

请返回JSON格式的响应，包含：
1. projectStructure: 项目文件结构和内容
2. variables: 推荐的模板变量
3. instructions: 使用说明
4. estimatedTime: 预估完成时间

要求：
- 文件结构要完整实用
- 代码要有注释和最佳实践
- 变量要有合理的默认值
- 说明要清晰易懂`, 
		req.Description, req.ProjectType, req.TechStack, req.Features)
}

func buildVariableSuggestPrompt(req *VariableSuggestRequest) string {
	return fmt.Sprintf(`请为以下项目推荐合适的模板变量：

项目类型：%s
技术栈：%v
项目描述：%s

请返回JSON格式的变量列表，每个变量包含：
- name: 变量名
- type: 类型(string/number/boolean)
- description: 描述
- defaultValue: 默认值
- required: 是否必填
- options: 选项(如果是选择类型)

要求变量实用、常见，有助于模板的灵活性。`,
		req.ProjectType, req.TechStack, req.Description)
}

func parseTemplateResponse(response map[string]interface{}) (*TemplateGenerateResponse, error) {
	// 简化实现，实际应该解析AI返回的JSON
	choices, ok := response["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return nil, fmt.Errorf("无效的AI响应")
	}

	choice := choices[0].(map[string]interface{})
	message := choice["message"].(map[string]interface{})
	content := message["content"].(string)

	g.Log().Debug(context.Background(), "AI响应内容:", content)

	// 这里应该解析AI返回的JSON结构
	// 为了演示，返回一个示例结构
	return &TemplateGenerateResponse{
		ProjectStructure: []FileInfo{
			{
				Path:        "README.md",
				Content:     "# {{.ProjectName}}\n\n{{.Description}}",
				IsDirectory: false,
				Description: "项目说明文档",
			},
		},
		Variables: []VariableInfo{
			{
				Name:         "ProjectName",
				Type:         "string",
				Description:  "项目名称",
				DefaultValue: "my-project",
				Required:     true,
			},
		},
		Instructions: "AI生成的模板使用说明",
		EstimatedTime: 15,
	}, nil
}

func parseVariableResponse(response map[string]interface{}) (*VariableSuggestResponse, error) {
	// 简化实现
	return &VariableSuggestResponse{
		Variables: []VariableInfo{
			{
				Name:         "ProjectName",
				Type:         "string",
				Description:  "项目名称",
				DefaultValue: "my-app",
				Required:     true,
			},
			{
				Name:         "Author",
				Type:         "string", 
				Description:  "作者姓名",
				DefaultValue: "开发者",
				Required:     false,
			},
		},
	}, nil
}