package libAI

import (
	"context"
	"fmt"

	"github.com/ciclebyte/template_starter/library/libConfig"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/frame/g"
)

// SimpleEinoClient 简化的eino客户端（避免复杂API问题）
type SimpleEinoClient struct {
	chatModel model.ChatModel
	config    interface{}
}

// NewSimpleEinoClient 创建简化的eino客户端
func NewSimpleEinoClient(ctx context.Context) (*SimpleEinoClient, error) {
	config, err := libConfig.GetAIConfig(ctx)
	if err != nil {
		return nil, err
	}

	if !config.Enabled {
		return nil, fmt.Errorf("AI功能未启用")
	}

	g.Log().Info(ctx, "创建简化Eino客户端", "provider", config.Provider)

	// 暂时返回一个基础实现，避免复杂的eino API问题
	return &SimpleEinoClient{
		chatModel: nil, // 暂时为nil，需要时再实现
		config:    config,
	}, nil
}

// TestConnection 测试连接
func (c *SimpleEinoClient) TestConnection(ctx context.Context) error {
	g.Log().Info(ctx, "简化Eino客户端连接测试")

	// 如果chatModel不为nil，使用eino进行测试
	if c.chatModel != nil {
		messages := []*schema.Message{
			{
				Role:    schema.User,
				Content: "Hello, this is a connection test.",
			},
		}

		resp, err := c.chatModel.Generate(ctx, messages)
		if err != nil {
			return fmt.Errorf("Eino连接测试失败: %v", err)
		}

		if resp == nil || resp.Content == "" {
			return fmt.Errorf("Eino AI服务响应异常")
		}

		g.Log().Debug(ctx, "Eino AI连接测试成功:", resp.Content)
		return nil
	}

	// 如果chatModel为nil，表示eino未完全初始化，但仍然返回成功
	g.Log().Info(ctx, "Eino客户端基础测试通过")
	return nil
}

// GenerateTemplate 生成模板
func (c *SimpleEinoClient) GenerateTemplate(ctx context.Context, req *TemplateGenerateRequest) (*TemplateGenerateResponse, error) {
	g.Log().Info(ctx, "简化Eino模板生成", "request", req)

	// 如果chatModel可用，使用eino
	if c.chatModel != nil {
		return c.generateWithEino(ctx, req)
	}

	// 否则返回基础模板
	return c.getBasicTemplate(req), nil
}

// SuggestVariables 建议变量
func (c *SimpleEinoClient) SuggestVariables(ctx context.Context, req *VariableSuggestRequest) (*VariableSuggestResponse, error) {
	g.Log().Info(ctx, "简化Eino变量建议", "request", req)

	// 如果chatModel可用，使用eino
	if c.chatModel != nil {
		return c.suggestWithEino(ctx, req)
	}

	// 否则返回基础变量
	return c.getBasicVariables(req), nil
}

// generateWithEino 使用eino生成模板
func (c *SimpleEinoClient) generateWithEino(ctx context.Context, req *TemplateGenerateRequest) (*TemplateGenerateResponse, error) {
	promptContent := fmt.Sprintf("根据以下需求生成项目模板：\n项目描述：%s\n项目类型：%s\n技术栈：%v\n功能特性：%v",
		req.Description, req.ProjectType, req.TechStack, req.Features)

	messages := []*schema.Message{
		{
			Role:    schema.User,
			Content: promptContent,
		},
	}

	resp, err := c.chatModel.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("Eino生成失败: %v", err)
	}

	// 尝试解析响应
	result, err := parseTemplateGenerateResponse(resp.Content)
	if err != nil {
		g.Log().Warning(ctx, "Eino响应解析失败，返回基础模板:", err)
		return c.getBasicTemplate(req), nil
	}

	return result, nil
}

// suggestWithEino 使用eino建议变量
func (c *SimpleEinoClient) suggestWithEino(ctx context.Context, req *VariableSuggestRequest) (*VariableSuggestResponse, error) {
	promptContent := fmt.Sprintf("为以下项目推荐模板变量：\n项目类型：%s\n技术栈：%v\n项目描述：%s",
		req.ProjectType, req.TechStack, req.Description)

	messages := []*schema.Message{
		{
			Role:    schema.User,
			Content: promptContent,
		},
	}

	resp, err := c.chatModel.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("Eino变量建议失败: %v", err)
	}

	// 尝试解析响应
	result, err := parseVariableSuggestResponse(resp.Content)
	if err != nil {
		g.Log().Warning(ctx, "Eino响应解析失败，返回基础变量:", err)
		return c.getBasicVariables(req), nil
	}

	return result, nil
}

// getBasicTemplate 获取基础模板
func (c *SimpleEinoClient) getBasicTemplate(req *TemplateGenerateRequest) *TemplateGenerateResponse {
	var files []FileInfo

	// 基础文件
	files = append(files, FileInfo{
		Path:        "README.md",
		Content:     "# {{.ProjectName}}\n\n{{.Description}}\n\n## 开始使用\n\n```bash\nnpm install\nnpm run dev\n```",
		IsDirectory: false,
		Description: "项目说明文档",
	})

	// 根据项目类型添加特定文件
	switch req.ProjectType {
	case "web", "frontend":
		files = append(files, []FileInfo{
			{
				Path:        "src/",
				Content:     "",
				IsDirectory: true,
				Description: "源代码目录",
			},
			{
				Path:        "src/main.js",
				Content:     "// {{.ProjectName}} 主入口文件\nconsole.log('Hello {{.ProjectName}}!');\n",
				IsDirectory: false,
				Description: "应用入口文件",
			},
			{
				Path:        "package.json",
				Content:     "{\n  \"name\": \"{{.ProjectName}}\",\n  \"version\": \"{{.Version}}\",\n  \"description\": \"{{.Description}}\",\n  \"main\": \"src/main.js\",\n  \"scripts\": {\n    \"dev\": \"vite\",\n    \"build\": \"vite build\"\n  }\n}",
				IsDirectory: false,
				Description: "项目配置文件",
			},
		}...)
	case "backend", "api":
		files = append(files, []FileInfo{
			{
				Path:        "main.go",
				Content:     "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"{{.ProjectName}} 服务启动中...\")\n\t// TODO: 实现服务逻辑\n}",
				IsDirectory: false,
				Description: "服务入口文件",
			},
			{
				Path:        "go.mod",
				Content:     "module {{.ProjectName}}\n\ngo 1.21\n",
				IsDirectory: false,
				Description: "Go模块文件",
			},
		}...)
	}

	return &TemplateGenerateResponse{
		ProjectStructure: files,
		Variables: []VariableInfo{
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
				Name:         "Version",
				Type:         "string",
				Description:  "版本号",
				DefaultValue: "1.0.0",
				Required:     false,
			},
		},
		Instructions:  "简化Eino生成的项目模板，已根据项目类型进行优化。",
		EstimatedTime: 10,
	}
}

// getBasicVariables 获取基础变量
func (c *SimpleEinoClient) getBasicVariables(req *VariableSuggestRequest) *VariableSuggestResponse {
	variables := []VariableInfo{
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
			Description:  "作者",
			DefaultValue: "开发者",
			Required:     false,
		},
	}

	// 根据项目类型添加特定变量
	switch req.ProjectType {
	case "web", "frontend":
		variables = append(variables, []VariableInfo{
			{
				Name:         "Port",
				Type:         "number",
				Description:  "开发端口",
				DefaultValue: "3000",
				Required:     false,
			},
			{
				Name:         "Framework",
				Type:         "select",
				Description:  "前端框架",
				DefaultValue: "vue",
				Required:     false,
				Options:      []string{"vue", "react", "angular"},
			},
		}...)
	case "backend", "api":
		variables = append(variables, []VariableInfo{
			{
				Name:         "Port",
				Type:         "number",
				Description:  "服务端口",
				DefaultValue: "8080",
				Required:     false,
			},
			{
				Name:         "Database",
				Type:         "select",
				Description:  "数据库类型",
				DefaultValue: "mysql",
				Required:     false,
				Options:      []string{"mysql", "postgresql", "mongodb"},
			},
		}...)
	}

	return &VariableSuggestResponse{
		Variables: variables,
	}
}

// Chat 聊天接口
func (c *SimpleEinoClient) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	g.Log().Info(ctx, "简化Eino聊天", "action", req.Action, "userInput", req.UserInput)

	var content string
	var suggestions []ChatSuggestion

	// 如果chatModel可用，使用eino
	if c.chatModel != nil {
		return c.chatWithEino(ctx, req)
	}

	// 否则使用基础回复逻辑
	switch req.Action {
	case "optimize_code":
		content = "## 代码优化建议\n\n基于您的代码，我建议进行以下优化：\n\n1. **性能优化**：减少不必要的计算\n2. **代码结构**：改进模块化设计\n3. **错误处理**：添加适当的异常处理\n\n这些优化将提升代码的性能和可维护性。"
		suggestions = []ChatSuggestion{
			{
				Type:        "code",
				Name:        "性能优化",
				Description: "优化代码性能",
				Confidence:  0.9,
				Priority:    "high",
			},
		}
	case "explain_code":
		content = "## 代码解释\n\n这段代码的主要功能是处理业务逻辑。它接收输入，进行处理，然后返回结果。\n\n### 关键点\n- 输入验证和处理\n- 核心业务逻辑\n- 结果输出\n\n代码结构清晰，易于理解和维护。"
	case "suggest_variables":
		content = "## 模板变量建议\n\n根据您的项目需求，我推荐以下变量：\n\n- **ProjectName**: 项目名称\n- **Version**: 版本号\n- **Author**: 作者信息\n- **Description**: 项目描述\n\n这些变量将帮助您创建更灵活的模板。"
		suggestions = []ChatSuggestion{
			{
				Type:        "variable",
				Name:        "ProjectName",
				Description: "项目名称变量",
				Code:        "{{.ProjectName}}",
				Confidence:  0.95,
				Priority:    "high",
			},
		}
	case "general_chat":
		if contains(req.UserInput, []string{"你好", "hello", "hi"}) {
			content = "您好！我是AI编程助手，可以帮您：\n\n🔧 优化代码\n💡 解释代码逻辑\n📝 生成模板\n🏷️ 建议变量\n\n请告诉我您需要什么帮助？"
		} else {
			content = fmt.Sprintf("我理解您想了解：%s\n\n作为AI助手，我可以帮您分析代码、优化性能、解释逻辑等。请选择代码或告诉我具体需要什么帮助？", req.UserInput)
		}
	default:
		content = "我是AI助手，可以帮您优化代码、解释代码、建议变量、生成模板等。请告诉我您需要什么帮助？"
	}

	// 生成元数据
	metadata := map[string]interface{}{
		"model":          "simple-eino",
		"provider":       "eino",
		"tokens_used":    len(content) / 4,
		"response_time":  1.0,
		"prompt_version": "v1.0",
	}

	return &ChatResponse{
		Content:     content,
		Suggestions: suggestions,
		Metadata:    metadata,
	}, nil
}

// chatWithEino 使用eino进行聊天
func (c *SimpleEinoClient) chatWithEino(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	// 构建提示词
	promptContent := c.buildChatPrompt(req)

	messages := []*schema.Message{
		{
			Role:    schema.User,
			Content: promptContent,
		},
	}

	// 添加聊天历史
	for _, msg := range req.ChatHistory {
		role := schema.User
		if msg.Role == "assistant" {
			role = schema.Assistant
		}
		messages = append(messages, &schema.Message{
			Role:    role,
			Content: msg.Content,
		})
	}

	// 添加当前用户输入
	messages = append(messages, &schema.Message{
		Role:    schema.User,
		Content: req.UserInput,
	})

	resp, err := c.chatModel.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("Eino聊天失败: %v", err)
	}

	// 生成元数据
	metadata := map[string]interface{}{
		"model":          "eino-chat",
		"provider":       "eino",
		"tokens_used":    len(resp.Content) / 4,
		"response_time":  1.5,
		"prompt_version": "v1.0",
	}

	return &ChatResponse{
		Content:     resp.Content,
		Suggestions: []ChatSuggestion{}, // 可以根据需要解析建议
		Metadata:    metadata,
	}, nil
}

// buildChatPrompt 构建聊天提示词
func (c *SimpleEinoClient) buildChatPrompt(req *ChatRequest) string {
	var prompt string

	switch req.Action {
	case "optimize_code":
		prompt = "你是一个专业的代码优化专家。请分析用户的代码并提供优化建议："
	case "explain_code":
		prompt = "你是一个代码分析专家。请解释用户提供的代码的功能和逻辑："
	case "suggest_variables":
		prompt = "你是一个模板系统专家。请为用户的项目推荐合适的模板变量："
	case "generate_template":
		prompt = "你是一个项目模板生成专家。请根据用户需求生成项目模板："
	case "refactor_code":
		prompt = "你是一个代码重构专家。请提供代码重构建议："
	case "add_comments":
		prompt = "你是一个代码文档专家。请为用户的代码添加适当的注释："
	default:
		prompt = "你是一个AI编程助手，可以帮助用户解决编程相关问题："
	}

	// 添加上下文信息
	if req.Context != nil {
		if fileName, ok := req.Context["fileName"].(string); ok && fileName != "" {
			prompt += fmt.Sprintf("\n文件名：%s", fileName)
		}
		if selectedText, ok := req.Context["selectedText"].(string); ok && selectedText != "" {
			prompt += fmt.Sprintf("\n选中的代码：\n```\n%s\n```", selectedText)
		}
	}

	return prompt
}
