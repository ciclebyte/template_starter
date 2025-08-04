package libAI

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	internalModel "github.com/ciclebyte/template_starter/internal/model"
	"github.com/ciclebyte/template_starter/library/libConfig"
	"github.com/gogf/gf/v2/frame/g"
)

// AIClient AI客户端接口
type AIClient interface {
	TestConnection(ctx context.Context) error
	GenerateTemplate(ctx context.Context, req *TemplateGenerateRequest) (*TemplateGenerateResponse, error)
	SuggestVariables(ctx context.Context, req *VariableSuggestRequest) (*VariableSuggestResponse, error)
	Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
	ChatStream(ctx context.Context, req *ChatRequest, onChunk func(chunk string)) (*ChatResponse, error)
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
	ProjectStructure []FileInfo     `json:"projectStructure"` // 项目结构
	Variables        []VariableInfo `json:"variables"`        // 推荐变量
	Instructions     string         `json:"instructions"`     // 使用说明
	EstimatedTime    int            `json:"estimatedTime"`    // 预估完成时间(分钟)
}

// VariableSuggestRequest 变量建议请求
type VariableSuggestRequest struct {
	ProjectType string   `json:"projectType"`
	TechStack   []string `json:"techStack"`
	Description string   `json:"description"`
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
	Name         string   `json:"name"`
	Type         string   `json:"type"`         // string, number, boolean
	Description  string   `json:"description"`
	DefaultValue string   `json:"defaultValue"`
	Required     bool     `json:"required"`
	Options      []string `json:"options,omitempty"` // 选项值（如果是选择类型）
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Action      string                 `json:"action"`      // 操作类型
	Context     map[string]interface{} `json:"context"`     // 上下文信息
	UserInput   string                 `json:"userInput"`   // 用户输入
	Preferences map[string]interface{} `json:"preferences"` // 用户偏好设置
	ChatHistory []ChatMessage          `json:"chatHistory"` // 聊天历史
}

// ChatResponse 聊天响应
type ChatResponse struct {
	Content     string                 `json:"content"`     // AI回复的主要内容
	Suggestions []ChatSuggestion       `json:"suggestions"` // 建议项
	Metadata    map[string]interface{} `json:"metadata"`    // 元数据
}

// ChatMessage 聊天消息
type ChatMessage struct {
	Role      string `json:"role"`      // user, assistant
	Content   string `json:"content"`   // 消息内容
	Timestamp string `json:"timestamp"` // 时间戳
}

// ChatSuggestion 聊天建议
type ChatSuggestion struct {
	Type        string  `json:"type"`        // code, variable, file, action
	Name        string  `json:"name"`        // 建议名称
	Description string  `json:"description"` // 建议描述
	Code        string  `json:"code"`        // 建议的代码内容
	Confidence  float64 `json:"confidence"`  // 置信度
	Priority    string  `json:"priority"`    // high, medium, low
}

// SimpleAIClient 简化的AI客户端
type SimpleAIClient struct {
	config *internalModel.AIConfig
}

// NewAIClient 创建AI客户端
func NewAIClient(ctx context.Context) (AIClient, error) {
	// 优先尝试创建SimpleEinoClient，现在它也支持真实AI调用
	einoClient, err := NewSimpleEinoClient(ctx)
	if err != nil {
		g.Log().Warning(ctx, "SimpleEinoClient创建失败，使用SimpleAIClient:", err)
		
		// 回退到SimpleAIClient
		config, err := libConfig.GetAIConfig(ctx)
		if err != nil {
			return nil, err
		}

		if !config.Enabled {
			return nil, fmt.Errorf("AI功能未启用")
		}

		g.Log().Info(ctx, "==== 创建SimpleAIClient ====", "provider", config.Provider, "model", config.OpenAI.Model)
		return &SimpleAIClient{
			config: config,
		}, nil
	}

	g.Log().Info(ctx, "==== 使用SimpleEinoClient（支持真实AI） ====")
	return einoClient, nil
}

// TestConnection 测试连接
func (c *SimpleAIClient) TestConnection(ctx context.Context) error {
	g.Log().Info(ctx, "AI连接测试", "provider", c.config.Provider)
	
	// TODO: 实际的AI服务连接测试
	// 这里可以添加真实的HTTP请求来测试AI服务连接
	
	return nil
}

// GenerateTemplate 生成模板
func (c *SimpleAIClient) GenerateTemplate(ctx context.Context, req *TemplateGenerateRequest) (*TemplateGenerateResponse, error) {
	g.Log().Info(ctx, "AI模板生成", "request", req)
	
	// 根据不同项目类型生成不同的模板结构
	projectStructure := c.generateProjectStructure(req)
	variables := c.generateVariables(req)
	
	return &TemplateGenerateResponse{
		ProjectStructure: projectStructure,
		Variables:        variables,
		Instructions:     c.generateInstructions(req),
		EstimatedTime:    c.estimateTime(req),
	}, nil
}

// SuggestVariables 建议变量
func (c *SimpleAIClient) SuggestVariables(ctx context.Context, req *VariableSuggestRequest) (*VariableSuggestResponse, error) {
	g.Log().Info(ctx, "AI变量建议", "request", req)
	
	variables := c.generateVariablesByType(req)
	
	return &VariableSuggestResponse{
		Variables: variables,
	}, nil
}

// generateProjectStructure 生成项目结构
func (c *SimpleAIClient) generateProjectStructure(req *TemplateGenerateRequest) []FileInfo {
	var files []FileInfo
	
	// 基础文件
	files = append(files, FileInfo{
		Path:        "README.md",
		Content:     "# {{.ProjectName}}\n\n{{.Description}}\n\n## 安装\n\n```bash\nnpm install\n```\n\n## 使用\n\n```bash\nnpm run dev\n```",
		IsDirectory: false,
		Description: "项目说明文档",
	})
	
	files = append(files, FileInfo{
		Path:        ".gitignore",
		Content:     "node_modules/\n.env\n.DS_Store\ndist/\nbuild/",
		IsDirectory: false,
		Description: "Git忽略配置",
	})
	
	// 根据项目类型添加特定文件
	switch req.ProjectType {
	case "web", "frontend":
		files = append(files, c.generateWebFiles(req)...)
	case "backend", "api":
		files = append(files, c.generateBackendFiles(req)...)
	case "mobile":
		files = append(files, c.generateMobileFiles(req)...)
	default:
		files = append(files, c.generateGenericFiles(req)...)
	}
	
	return files
}

// generateWebFiles 生成Web项目文件
func (c *SimpleAIClient) generateWebFiles(req *TemplateGenerateRequest) []FileInfo {
	files := []FileInfo{
		{
			Path:        "src/",
			Content:     "",
			IsDirectory: true,
			Description: "源代码目录",
		},
		{
			Path:        "src/main.js",
			Content:     "import { createApp } from 'vue'\nimport App from './App.vue'\n\ncreateApp(App).mount('#app')",
			IsDirectory: false,
			Description: "应用入口文件",
		},
		{
			Path:        "src/App.vue",
			Content:     "<template>\n  <div id=\"app\">\n    <h1>{{.ProjectName}}</h1>\n    <p>{{.Description}}</p>\n  </div>\n</template>\n\n<script>\nexport default {\n  name: 'App'\n}\n</script>",
			IsDirectory: false,
			Description: "主应用组件",
		},
		{
			Path:        "package.json",
			Content:     "{\n  \"name\": \"{{.ProjectName}}\",\n  \"version\": \"{{.Version}}\",\n  \"description\": \"{{.Description}}\",\n  \"main\": \"src/main.js\",\n  \"scripts\": {\n    \"dev\": \"vite\",\n    \"build\": \"vite build\",\n    \"preview\": \"vite preview\"\n  },\n  \"dependencies\": {\n    \"vue\": \"^3.0.0\"\n  },\n  \"devDependencies\": {\n    \"vite\": \"^4.0.0\"\n  }\n}",
			IsDirectory: false,
			Description: "项目配置文件",
		},
		{
			Path:        "vite.config.js",
			Content:     "import { defineConfig } from 'vite'\nimport vue from '@vitejs/plugin-vue'\n\nexport default defineConfig({\n  plugins: [vue()],\n  server: {\n    port: {{.Port}}\n  }\n})",
			IsDirectory: false,
			Description: "Vite配置文件",
		},
	}
	
	return files
}

// generateBackendFiles 生成后端项目文件
func (c *SimpleAIClient) generateBackendFiles(req *TemplateGenerateRequest) []FileInfo {
	files := []FileInfo{
		{
			Path:        "main.go",
			Content:     "package main\n\nimport (\n\t\"fmt\"\n\t\"net/http\"\n)\n\nfunc main() {\n\thttp.HandleFunc(\"/\", func(w http.ResponseWriter, r *http.Request) {\n\t\tfmt.Fprintf(w, \"Hello from {{.ProjectName}}!\")\n\t})\n\n\tfmt.Println(\"Server starting on port {{.Port}}...\")\n\thttp.ListenAndServe(\":{{.Port}}\", nil)\n}",
			IsDirectory: false,
			Description: "应用入口文件",
		},
		{
			Path:        "go.mod",
			Content:     "module {{.ProjectName}}\n\ngo 1.21\n\nrequire (\n\t// Add your dependencies here\n)",
			IsDirectory: false,
			Description: "Go模块配置",
		},
		{
			Path:        "api/",
			Content:     "",
			IsDirectory: true,
			Description: "API接口目录",
		},
		{
			Path:        "internal/",
			Content:     "",
			IsDirectory: true,
			Description: "内部代码目录",
		},
	}
	
	return files
}

// generateMobileFiles 生成移动应用文件
func (c *SimpleAIClient) generateMobileFiles(req *TemplateGenerateRequest) []FileInfo {
	files := []FileInfo{
		{
			Path:        "App.js",
			Content:     "import React from 'react';\nimport { Text, View } from 'react-native';\n\nexport default function App() {\n  return (\n    <View style={{flex: 1, justifyContent: 'center', alignItems: 'center'}}>\n      <Text>{{.ProjectName}}</Text>\n      <Text>{{.Description}}</Text>\n    </View>\n  );\n}",
			IsDirectory: false,
			Description: "应用主组件",
		},
		{
			Path:        "package.json",
			Content:     "{\n  \"name\": \"{{.ProjectName}}\",\n  \"version\": \"{{.Version}}\",\n  \"main\": \"App.js\",\n  \"scripts\": {\n    \"start\": \"expo start\",\n    \"android\": \"expo start --android\",\n    \"ios\": \"expo start --ios\"\n  },\n  \"dependencies\": {\n    \"react\": \"^18.0.0\",\n    \"react-native\": \"^0.72.0\"\n  }\n}",
			IsDirectory: false,
			Description: "项目配置文件",
		},
	}
	
	return files
}

// generateGenericFiles 生成通用项目文件
func (c *SimpleAIClient) generateGenericFiles(req *TemplateGenerateRequest) []FileInfo {
	return []FileInfo{
		{
			Path:        "src/",
			Content:     "",
			IsDirectory: true,
			Description: "源代码目录",
		},
		{
			Path:        "docs/",
			Content:     "",
			IsDirectory: true,
			Description: "文档目录",
		},
	}
}

// generateVariables 生成推荐变量
func (c *SimpleAIClient) generateVariables(req *TemplateGenerateRequest) []VariableInfo {
	variables := []VariableInfo{
		{
			Name:         "ProjectName",
			Type:         "string",
			Description:  "项目名称",
			DefaultValue: "my-project",
			Required:     true,
		},
		{
			Name:         "Description",
			Type:         "string",
			Description:  "项目描述",
			DefaultValue: req.Description,
			Required:     false,
		},
		{
			Name:         "Author",
			Type:         "string",
			Description:  "作者姓名",
			DefaultValue: "开发者",
			Required:     false,
		},
		{
			Name:         "Version",
			Type:         "string",
			Description:  "版本号",
			DefaultValue: "1.0.0",
			Required:     false,
		},
	}
	
	// 根据项目类型添加特定变量
	switch req.ProjectType {
	case "web", "frontend":
		variables = append(variables, VariableInfo{
			Name:         "Port",
			Type:         "number",
			Description:  "开发服务器端口",
			DefaultValue: "3000",
			Required:     false,
		})
	case "backend", "api":
		variables = append(variables, VariableInfo{
			Name:         "Port",
			Type:         "number",
			Description:  "服务端口",
			DefaultValue: "8080",
			Required:     false,
		})
	}
	
	return variables
}

// generateVariablesByType 根据类型生成变量建议
func (c *SimpleAIClient) generateVariablesByType(req *VariableSuggestRequest) []VariableInfo {
	var variables []VariableInfo
	
	// 基础变量
	variables = append(variables, []VariableInfo{
		{
			Name:         "ProjectName",
			Type:         "string",
			Description:  "项目名称",
			DefaultValue: "my-app",
			Required:     true,
		},
		{
			Name:         "Description",
			Type:         "string",
			Description:  "项目描述",
			DefaultValue: req.Description,
			Required:     false,
		},
		{
			Name:         "Author",
			Type:         "string",
			Description:  "作者姓名",
			DefaultValue: "开发者",
			Required:     false,
		},
	}...)
	
	// 根据项目类型添加特定变量
	switch req.ProjectType {
	case "web", "frontend":
		variables = append(variables, []VariableInfo{
			{
				Name:         "Theme",
				Type:         "select",
				Description:  "UI主题",
				DefaultValue: "light",
				Required:     false,
				Options:      []string{"light", "dark"},
			},
			{
				Name:         "Language",
				Type:         "select",
				Description:  "开发语言",
				DefaultValue: "javascript",
				Required:     false,
				Options:      []string{"javascript", "typescript"},
			},
		}...)
		
	case "backend", "api":
		variables = append(variables, []VariableInfo{
			{
				Name:         "DatabaseType",
				Type:         "select",
				Description:  "数据库类型",
				DefaultValue: "mysql",
				Required:     false,
				Options:      []string{"mysql", "postgresql", "mongodb"},
			},
		}...)
	}
	
	return variables
}

// generateInstructions 生成使用说明
func (c *SimpleAIClient) generateInstructions(req *TemplateGenerateRequest) string {
	instructions := "## 使用说明\n\n1. 解压模板文件到目标目录\n2. 根据需要修改模板变量\n3. 运行初始化命令\n\n"
	
	switch req.ProjectType {
	case "web", "frontend":
		instructions += "## 开发步骤\n1. 安装依赖：npm install\n2. 启动开发服务器：npm run dev\n3. 构建生产版本：npm run build"
	case "backend", "api":
		instructions += "## 开发步骤\n1. 安装Go依赖：go mod tidy\n2. 运行应用：go run main.go\n3. 构建二进制：go build"
	case "mobile":
		instructions += "## 开发步骤\n1. 安装依赖：npm install\n2. 启动开发服务器：npm start\n3. 在模拟器中运行：npm run android 或 npm run ios"
	}
	
	return instructions
}

// estimateTime 估算完成时间
func (c *SimpleAIClient) estimateTime(req *TemplateGenerateRequest) int {
	baseTime := 10 // 基础时间10分钟
	
	// 根据技术栈复杂度调整时间
	techComplexity := len(req.TechStack) * 5
	featureComplexity := len(req.Features) * 3
	
	return baseTime + techComplexity + featureComplexity
}

// Chat 聊天接口
func (c *SimpleAIClient) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	g.Log().Info(ctx, "==== SimpleAIClient.Chat 调用 ====", "action", req.Action, "userInput", req.UserInput)
	
	// 检查配置
	if !c.config.Enabled {
		return nil, fmt.Errorf("AI功能未启用，请在系统配置中设置 ai.enabled = true")
	}
	
	if c.config.OpenAI.APIKey == "" {
		return nil, fmt.Errorf("AI API密钥未配置，请在系统配置中设置 ai.openai.api_key")
	}
	
	// 调用真实AI
	apiKeyPrefix := c.config.OpenAI.APIKey
	if len(apiKeyPrefix) > 10 {
		apiKeyPrefix = apiKeyPrefix[:10] + "..."
	}
	g.Log().Info(ctx, "调用真实AI服务", "provider", c.config.Provider, "api_key_prefix", apiKeyPrefix)
	
	// 尝试真实AI调用，失败时返回详细错误
	response, err := c.callRealAI(ctx, req)
	if err != nil {
		// 如果是网络或API错误，提供详细的错误信息和建议
		g.Log().Error(ctx, "真实AI调用失败", "error", err)
		return nil, fmt.Errorf("AI服务调用失败: %v\n\n建议检查：\n1. 网络连接是否正常\n2. API密钥是否正确\n3. API服务器地址是否正确: %s", err, c.config.OpenAI.BaseURL)
	}
	
	return response, nil
}

// callRealAI 调用真实的AI服务
func (c *SimpleAIClient) callRealAI(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	g.Log().Info(ctx, "调用真实AI服务", "provider", c.config.Provider, "model", c.config.OpenAI.Model)
	
	// 构建提示词
	prompt := c.buildPromptForAction(req)
	
	// 根据provider调用不同的AI服务
	switch c.config.Provider {
	case "openai", "moonshot":
		return c.callOpenAICompatible(ctx, prompt, req)
	case "claude":
		return c.callClaude(ctx, prompt, req)
	default:
		// 默认使用OpenAI兼容接口
		return c.callOpenAICompatible(ctx, prompt, req)
	}
}

// callOpenAICompatible 调用OpenAI兼容的API
func (c *SimpleAIClient) callOpenAICompatible(ctx context.Context, prompt string, req *ChatRequest) (*ChatResponse, error) {
	// 构建消息历史
	messages := []map[string]interface{}{
		{
			"role":    "system",
			"content": prompt,
		},
	}
	
	// 添加聊天历史
	for _, msg := range req.ChatHistory {
		messages = append(messages, map[string]interface{}{
			"role":    msg.Role,
			"content": msg.Content,
		})
	}
	
	// 添加用户当前输入
	messages = append(messages, map[string]interface{}{
		"role":    "user",
		"content": req.UserInput,
	})
	
	// 构建请求体
	requestBody := map[string]interface{}{
		"model":       c.config.OpenAI.Model,
		"messages":    messages,
		"max_tokens":  c.config.OpenAI.MaxTokens,
		"temperature": c.config.OpenAI.Temperature,
		"stream":      false,
	}
	
	// 创建HTTP客户端 - 使用链式调用方式
	client := g.Client()
	
	apiURL := strings.TrimRight(c.config.OpenAI.BaseURL, "/") + "/v1/chat/completions"
	
	g.Log().Info(ctx, "发送AI请求", "url", apiURL, "model", c.config.OpenAI.Model)
	
	// 序列化请求体
	requestJson, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %v", err)
	}
	
	g.Log().Debug(ctx, "请求体", "json", string(requestJson))
	
	// 使用链式调用方式发送请求
	responseBody := client.Header(map[string]string{
		"Authorization": "Bearer " + c.config.OpenAI.APIKey,
		"Content-Type":  "application/json",
		"User-Agent":    "GoFrame-AI-Client/1.0",
	}).Timeout(60*time.Second).PostContent(ctx, apiURL, string(requestJson))
	
	g.Log().Info(ctx, "收到AI响应", "body_length", len(responseBody))
	
	// 解析响应
	var apiResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			TotalTokens int `json:"total_tokens"`
		} `json:"usage"`
		Error struct {
			Message string `json:"message"`
			Type    string `json:"type"`
		} `json:"error"`
	}
	
	// 检查响应是否为空或错误
	if len(responseBody) == 0 {
		return nil, fmt.Errorf("AI API返回空响应")
	}
	
	if err := json.Unmarshal([]byte(responseBody), &apiResponse); err != nil {
		g.Log().Error(ctx, "AI响应解析失败:", err, "response:", responseBody)
		return nil, fmt.Errorf("AI API响应格式错误: %v\n原始响应: %s", err, responseBody)
	}
	
	// 检查是否有错误
	if apiResponse.Error.Message != "" {
		g.Log().Error(ctx, "AI API错误:", apiResponse.Error.Message, apiResponse.Error.Type)
		return nil, fmt.Errorf("AI API调用失败: %s (类型: %s)", apiResponse.Error.Message, apiResponse.Error.Type)
	}
	
	// 检查是否有有效响应
	if len(apiResponse.Choices) == 0 {
		g.Log().Warning(ctx, "AI API返回空响应")
		return nil, fmt.Errorf("AI API返回空响应，请检查API配置或稍后重试")
	}
	
	content := apiResponse.Choices[0].Message.Content
	
	// 生成建议（可以根据内容解析或使用固定建议）
	suggestions := c.generateSuggestionsForAction(req)
	
	// 生成元数据
	metadata := map[string]interface{}{
		"model":         c.config.OpenAI.Model,
		"provider":      c.config.Provider,
		"tokens_used":   apiResponse.Usage.TotalTokens,
		"response_time": 2.0, // 可以记录实际响应时间
		"prompt_version": "v1.0",
		"real_ai":       true,
	}
	
	g.Log().Info(ctx, "AI真实响应成功", "tokens", apiResponse.Usage.TotalTokens, "length", len(content))
	
	return &ChatResponse{
		Content:     content,
		Suggestions: suggestions,
		Metadata:    metadata,
	}, nil
}

// callClaude 调用Claude API
func (c *SimpleAIClient) callClaude(ctx context.Context, prompt string, req *ChatRequest) (*ChatResponse, error) {
	return nil, fmt.Errorf("Claude API集成尚未实现，请使用 openai 或 moonshot 作为 AI provider")
}

// buildPromptForAction 根据操作类型构建提示词
func (c *SimpleAIClient) buildPromptForAction(req *ChatRequest) string {
	var basePrompt string
	
	switch req.Action {
	case "optimize_code":
		basePrompt = `你是一个专业的代码优化专家。请分析用户的代码并提供具体的优化建议。

要求：
1. 指出性能瓶颈和改进点
2. 提供优化后的代码示例
3. 解释优化的理由和预期效果
4. 注意代码的可读性和维护性

请用中文回复，使用Markdown格式。`

	case "explain_code":
		basePrompt = `你是一个代码分析专家。请详细解释用户提供的代码的功能和逻辑。

要求：
1. 解释代码的整体功能和目的
2. 分析关键部分的实现逻辑
3. 指出重要的设计模式或算法
4. 说明代码的优缺点

请用中文回复，使用Markdown格式。`

	case "suggest_variables":
		basePrompt = `你是一个模板系统专家。请为用户的项目推荐合适的模板变量。

要求：
1. 分析项目类型和技术栈
2. 推荐必要的模板变量
3. 为每个变量提供合理的默认值
4. 说明变量的用途和重要性

请用中文回复，使用Markdown格式。`

	case "generate_template":
		basePrompt = `你是一个项目模板生成专家。请根据用户需求生成完整的项目模板。

要求：
1. 创建合理的项目目录结构
2. 生成必要的配置文件
3. 提供基础的代码模板
4. 包含使用说明和最佳实践

请用中文回复，使用Markdown格式。`

	case "refactor_code":
		basePrompt = `你是一个代码重构专家。请提供具体的代码重构建议。

要求：
1. 识别代码中的问题和改进点
2. 提供重构后的代码示例
3. 解释重构的理由和好处
4. 保持功能不变的前提下改进结构

请用中文回复，使用Markdown格式。`

	case "add_comments":
		basePrompt = `你是一个代码文档专家。请为用户的代码添加适当的注释。

要求：
1. 为函数和类添加说明注释
2. 为复杂逻辑添加解释注释
3. 遵循代码注释的最佳实践
4. 注释要简洁明了，不冗余

请用中文回复，提供带注释的代码。`

	default:
		basePrompt = `你是一个AI编程助手，可以帮助用户解决各种编程相关问题。

请根据用户的问题提供有帮助的建议和解决方案。用中文回复，使用清晰的格式。`
	}
	
	// 添加上下文信息
	if req.Context != nil {
		if fileName, ok := req.Context["fileName"].(string); ok && fileName != "" {
			basePrompt += fmt.Sprintf("\n\n文件名：%s", fileName)
		}
		
		if selectedText, ok := req.Context["selectedText"].(string); ok && selectedText != "" {
			basePrompt += fmt.Sprintf("\n\n用户选中的代码：\n```\n%s\n```", selectedText)
		} else if fileContent, ok := req.Context["fileContent"].(string); ok && fileContent != "" {
			// 如果没有选中文本但有文件内容，截取前1000个字符
			content := fileContent
			if len(content) > 1000 {
				content = content[:1000] + "..."
			}
			basePrompt += fmt.Sprintf("\n\n完整文件内容：\n```\n%s\n```", content)
		}
		
		if variables, ok := req.Context["variables"].([]interface{}); ok && len(variables) > 0 {
			basePrompt += fmt.Sprintf("\n\n模板变量：%v", variables)
		}
	}
	
	return basePrompt
}

// generateSuggestionsForAction 根据操作类型生成建议
func (c *SimpleAIClient) generateSuggestionsForAction(req *ChatRequest) []ChatSuggestion {
	switch req.Action {
	case "optimize_code":
		return []ChatSuggestion{
			{
				Type:        "code",
				Name:        "应用优化建议",
				Description: "将AI建议的优化应用到代码中",
				Confidence:  0.9,
				Priority:    "high",
			},
		}
	case "suggest_variables":
		return []ChatSuggestion{
			{
				Type:        "variable",
				Name:        "添加推荐变量",
				Description: "将AI推荐的变量添加到模板中",
				Confidence:  0.95,
				Priority:    "high",
			},
		}
	default:
		return []ChatSuggestion{}
	}
}

// generateMockResponse 生成模拟响应（回退方案）
func (c *SimpleAIClient) generateMockResponse(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	g.Log().Warning(ctx, "使用模拟AI响应")
	
	var content string
	var suggestions []ChatSuggestion
	
	// 根据不同的操作类型生成不同的响应
	switch req.Action {
	case "optimize_code":
		content = c.generateOptimizeCodeResponse(req)
		suggestions = c.generateOptimizeCodeSuggestions(req)
	case "explain_code":
		content = c.generateExplainCodeResponse(req)
		suggestions = c.generateExplainCodeSuggestions(req)
	case "suggest_variables":
		content = c.generateSuggestVariablesResponse(req)
		suggestions = c.generateVariableSuggestions(req)
	case "generate_template":
		content = c.generateTemplateResponse(req)
		suggestions = c.generateTemplateSuggestions(req)
	case "refactor_code":
		content = c.generateRefactorCodeResponse(req)
		suggestions = c.generateRefactorCodeSuggestions(req)
	case "add_comments":
		content = c.generateAddCommentsResponse(req)
		suggestions = c.generateAddCommentsSuggestions(req)
	case "general_chat":
		content = c.generateGeneralChatResponse(req)
		suggestions = []ChatSuggestion{}
	default:
		content = "我是AI助手，可以帮您优化代码、解释代码、建议变量、生成模板等。请告诉我您需要什么帮助？"
		suggestions = []ChatSuggestion{}
	}
	
	// 生成元数据
	metadata := map[string]interface{}{
		"model":        c.config.OpenAI.Model,
		"provider":     c.config.Provider,
		"tokens_used":  len(content) / 4, // 简单估算
		"response_time": 1.2,
		"prompt_version": "v1.0",
		"real_ai":       false,
	}
	
	return &ChatResponse{
		Content:     content,
		Suggestions: suggestions,
		Metadata:    metadata,
	}, nil
}

// generateOptimizeCodeResponse 生成代码优化响应
func (c *SimpleAIClient) generateOptimizeCodeResponse(req *ChatRequest) string {
	fileName := ""
	selectedText := ""
	hasSelection := false
	
	if context, ok := req.Context["fileName"].(string); ok {
		fileName = context
	}
	if selection, ok := req.Context["selectedText"].(string); ok {
		selectedText = selection
		hasSelection = len(selectedText) > 0
	}
	
	if hasSelection {
		return fmt.Sprintf("## 代码优化建议\n\n我已分析您选中的代码片段，以下是优化建议：\n\n### 优化点\n1. **性能优化**：可以通过缓存和减少不必要的计算来提升性能\n2. **代码结构**：建议使用更清晰的函数命名和模块化设计\n3. **错误处理**：添加适当的错误处理机制\n\n### 建议的优化代码\n```%s\n// 优化后的代码\n%s\n```\n\n**注意事项**：请在应用优化前进行充分测试。", getFileExtension(fileName), selectedText)
	}
	
	return fmt.Sprintf("## 文件优化分析\n\n我已分析您的文件 `%s`，以下是优化建议：\n\n### 主要优化点\n1. **代码组织**：建议重新组织代码结构\n2. **性能提升**：识别出几个性能瓶颈点\n3. **可读性**：改进变量命名和添加注释\n\n### 下一步\n点击下方的具体建议来查看详细的优化方案。", fileName)
}

// generateExplainCodeResponse 生成代码解释响应
func (c *SimpleAIClient) generateExplainCodeResponse(req *ChatRequest) string {
	fileName := ""
	selectedText := ""
	hasSelection := false
	
	if context, ok := req.Context["fileName"].(string); ok {
		fileName = context
	}
	if selection, ok := req.Context["selectedText"].(string); ok {
		selectedText = selection
		hasSelection = len(selectedText) > 0
	}
	
	if hasSelection {
		return fmt.Sprintf("## 代码解释\n\n### 选中代码功能\n这段代码的主要功能是：\n\n1. **核心逻辑**：执行特定的业务逻辑处理\n2. **输入处理**：接收和验证输入参数\n3. **输出结果**：返回处理结果\n\n### 代码流程\n```\n输入 → 验证 → 处理 → 输出\n```\n\n### 关键点\n- 该代码片段属于 `%s` 文件\n- 主要使用了编程语言的核心特性\n- 建议关注错误处理和边界条件", fileName)
	}
	
	return fmt.Sprintf("## 文件代码解释\n\n### `%s` 文件分析\n\n这个文件的主要职责是：\n\n1. **模块功能**：实现特定的业务功能模块\n2. **接口定义**：定义了相关的接口和数据结构\n3. **逻辑实现**：包含核心的业务逻辑实现\n\n### 架构说明\n该文件在整个项目中扮演重要角色，与其他模块协同工作。", fileName)
}

// generateSuggestVariablesResponse 生成变量建议响应
func (c *SimpleAIClient) generateSuggestVariablesResponse(req *ChatRequest) string {
	return "## 模板变量建议\n\n基于您的项目类型和技术栈，我推荐以下模板变量：\n\n### 基础变量\n- **ProjectName**: 项目名称\n- **Description**: 项目描述\n- **Author**: 作者信息\n- **Version**: 版本号\n\n### 技术相关变量\n根据您的技术栈，还建议添加框架特定的变量。\n\n点击下方建议查看具体的变量配置。"
}

// generateTemplateResponse 生成模板响应
func (c *SimpleAIClient) generateTemplateResponse(req *ChatRequest) string {
	return "## 模板生成方案\n\n我已根据您的需求生成了项目模板：\n\n### 项目结构\n- 📁 **src/**: 源代码目录\n- 📁 **docs/**: 文档目录\n- 📄 **README.md**: 项目说明\n- 📄 **package.json**: 项目配置\n\n### 特性包含\n✅ 基础项目结构\n✅ 配置文件模板\n✅ 开发工具配置\n✅ 文档模板\n\n模板已准备就绪，您可以开始使用了！"
}

// generateRefactorCodeResponse 生成重构响应
func (c *SimpleAIClient) generateRefactorCodeResponse(req *ChatRequest) string {
	return "## 代码重构建议\n\n### 重构目标\n1. **提高可维护性**：优化代码结构\n2. **增强可读性**：改进命名和注释\n3. **减少复杂度**：简化逻辑流程\n\n### 重构方案\n- 提取公共方法\n- 优化数据结构\n- 改进错误处理\n- 添加适当注释\n\n### 重构后的效果\n代码将更加清晰、易于维护和扩展。"
}

// generateAddCommentsResponse 生成添加注释响应
func (c *SimpleAIClient) generateAddCommentsResponse(req *ChatRequest) string {
	return "## 注释添加建议\n\n### 注释策略\n1. **函数注释**：说明函数的用途、参数和返回值\n2. **复杂逻辑注释**：解释复杂的算法和业务逻辑\n3. **TODO注释**：标记待优化和待完成的部分\n\n### 注释规范\n- 使用清晰简洁的语言\n- 保持注释与代码同步\n- 避免显而易见的注释\n\n点击建议查看具体的注释方案。"
}

// generateGeneralChatResponse 生成通用聊天响应
func (c *SimpleAIClient) generateGeneralChatResponse(req *ChatRequest) string {
	userInput := req.UserInput
	
	// 简单的关键词匹配回复
	if contains(userInput, []string{"你好", "hello", "hi"}) {
		return "您好！我是AI编程助手，可以帮您优化代码、解释代码逻辑、生成模板、建议变量等。请告诉我您需要什么帮助？"
	}
	
	if contains(userInput, []string{"帮助", "help", "功能"}) {
		return "## 我的功能\n\n🔧 **代码优化**: 分析并优化您的代码\n💡 **代码解释**: 解释代码的功能和逻辑\n📝 **模板生成**: 根据需求生成项目模板\n🏷️ **变量建议**: 为模板推荐合适的变量\n✂️ **代码重构**: 重构代码结构\n📄 **添加注释**: 为代码添加合适的注释\n\n您可以选择文件或代码片段后使用这些功能。"
	}
	
	return fmt.Sprintf("我理解您想了解：%s\n\n作为AI编程助手，我可以帮您分析代码、优化性能、解释逻辑等。您可以：\n\n1. 选择代码片段让我分析\n2. 询问具体的编程问题\n3. 让我帮您生成模板或建议变量\n\n请告诉我具体需要什么帮助？", userInput)
}

// 生成各种建议的辅助方法
func (c *SimpleAIClient) generateOptimizeCodeSuggestions(req *ChatRequest) []ChatSuggestion {
	return []ChatSuggestion{
		{
			Type:        "code",
			Name:        "性能优化",
			Description: "优化循环和算法效率",
			Code:        "// 优化后的代码示例",
			Confidence:  0.9,
			Priority:    "high",
		},
		{
			Type:        "code",
			Name:        "内存优化",
			Description: "减少内存占用",
			Code:        "// 内存优化代码",
			Confidence:  0.8,
			Priority:    "medium",
		},
	}
}

func (c *SimpleAIClient) generateExplainCodeSuggestions(req *ChatRequest) []ChatSuggestion {
	return []ChatSuggestion{
		{
			Type:        "action",
			Name:        "查看相关函数",
			Description: "查看调用的其他函数",
			Confidence:  0.9,
			Priority:    "medium",
		},
	}
}

func (c *SimpleAIClient) generateVariableSuggestions(req *ChatRequest) []ChatSuggestion {
	return []ChatSuggestion{
		{
			Type:        "variable",
			Name:        "ProjectName",
			Description: "项目名称变量",
			Code:        "{{.ProjectName}}",
			Confidence:  0.95,
			Priority:    "high",
		},
		{
			Type:        "variable",
			Name:        "Version",
			Description: "版本号变量",
			Code:        "{{.Version}}",
			Confidence:  0.9,
			Priority:    "medium",
		},
	}
}

func (c *SimpleAIClient) generateTemplateSuggestions(req *ChatRequest) []ChatSuggestion {
	return []ChatSuggestion{
		{
			Type:        "file",
			Name:        "添加配置文件",
			Description: "添加项目配置文件模板",
			Confidence:  0.9,
			Priority:    "high",
		},
	}
}

func (c *SimpleAIClient) generateRefactorCodeSuggestions(req *ChatRequest) []ChatSuggestion {
	return []ChatSuggestion{
		{
			Type:        "code",
			Name:        "提取方法",
			Description: "将重复代码提取为独立方法",
			Confidence:  0.9,
			Priority:    "high",
		},
	}
}

func (c *SimpleAIClient) generateAddCommentsSuggestions(req *ChatRequest) []ChatSuggestion {
	return []ChatSuggestion{
		{
			Type:        "code",
			Name:        "函数注释",
			Description: "添加函数说明注释",
			Code:        "// 函数功能描述\n// @param {type} param 参数说明\n// @return {type} 返回值说明",
			Confidence:  0.95,
			Priority:    "high",
		},
	}
}

// 辅助方法
func getFileExtension(fileName string) string {
	if fileName == "" {
		return ""
	}
	
	ext := ""
	if idx := strings.LastIndex(fileName, "."); idx != -1 {
		ext = fileName[idx+1:]
	}
	
	// 映射到语言标识
	langMap := map[string]string{
		"js":   "javascript",
		"ts":   "typescript",
		"vue":  "vue",
		"go":   "go",
		"py":   "python",
		"java": "java",
		"cpp":  "cpp",
		"c":    "c",
	}
	
	if lang, exists := langMap[ext]; exists {
		return lang
	}
	
	return ext
}

// ChatStream 流式聊天接口 - SimpleAIClient的实现
func (c *SimpleAIClient) ChatStream(ctx context.Context, req *ChatRequest, onChunk func(chunk string)) (*ChatResponse, error) {
	g.Log().Info(ctx, "==== SimpleAIClient.ChatStream 调用 ====", "action", req.Action, "userInput", req.UserInput)
	
	// SimpleAIClient暂时不支持真正的流式，先调用普通Chat然后模拟流式输出
	response, err := c.Chat(ctx, req)
	if err != nil {
		return nil, err
	}
	
	// 模拟流式输出
	if onChunk != nil {
		content := response.Content
		chunkSize := 20 // 每次发送的字符数
		
		for i := 0; i < len(content); i += chunkSize {
			end := i + chunkSize
			if end > len(content) {
				end = len(content)
			}
			
			chunk := content[i:end]
			onChunk(chunk)
			
			// 模拟打字延迟
			time.Sleep(50 * time.Millisecond)
		}
	}
	
	return response, nil
}

func contains(text string, keywords []string) bool {
	textLower := strings.ToLower(text)
	for _, keyword := range keywords {
		if strings.Contains(textLower, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}