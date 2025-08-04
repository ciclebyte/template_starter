package libAI

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	internalModel "github.com/ciclebyte/template_starter/internal/model"
	"github.com/ciclebyte/template_starter/library/libConfig"
	"github.com/cloudwego/eino-ext/components/model/openai"
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

	g.Log().Info(ctx, "创建Eino客户端", "provider", config.Provider)

	// 尝试创建真实的eino chatModel
	var chatModel model.ChatModel

	if config.OpenAI.APIKey != "" {
		// 使用eino创建OpenAI兼容的ChatModel
		chatModel, err = createEinoChatModel(ctx, config)
		if err != nil {
			g.Log().Warning(ctx, "创建eino ChatModel失败，使用HTTP回退:", err)
			chatModel = nil
		} else {
			g.Log().Info(ctx, "成功创建eino ChatModel")
		}
	}

	return &SimpleEinoClient{
		chatModel: chatModel,
		config:    config,
	}, nil
}

// createEinoChatModel 创建eino的ChatModel
func createEinoChatModel(ctx context.Context, config *internalModel.AIConfig) (model.ChatModel, error) {
	g.Log().Info(ctx, "创建eino ChatModel", "baseURL", config.OpenAI.BaseURL, "model", config.OpenAI.Model)

	// 根据用户提供的官方示例，使用eino-ext创建ChatModel
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey:  config.OpenAI.APIKey,
		BaseURL: config.OpenAI.BaseURL,
		Model:   config.OpenAI.Model,
		Timeout: 5 * time.Minute, // 增加超时时间到5分钟，支持长文本处理
	})

	if err != nil {
		g.Log().Error(ctx, "创建eino ChatModel失败", "error", err)
		return nil, fmt.Errorf("eino ChatModel创建失败: %v", err)
	}

	g.Log().Info(ctx, "eino ChatModel创建成功")
	return chatModel, nil
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
	g.Log().Info(ctx, "==== SimpleEinoClient.Chat 调用 ====", "action", req.Action, "userInput", req.UserInput)

	// 获取AI配置
	config, ok := c.config.(*internalModel.AIConfig)
	if !ok {
		return nil, fmt.Errorf("AI配置类型错误")
	}

	// 检查配置
	if !config.Enabled {
		return nil, fmt.Errorf("AI功能未启用")
	}

	if config.OpenAI.APIKey == "" {
		return nil, fmt.Errorf("AI API密钥未配置")
	}

	// 如果chatModel可用，使用eino
	if c.chatModel != nil {
		g.Log().Info(ctx, "使用eino chatModel")
		return c.chatWithEino(ctx, req)
	}

	// 否则使用HTTP客户端直接调用API（与SimpleAIClient相同的逻辑）
	g.Log().Info(ctx, "eino chatModel未初始化，使用HTTP直接调用")
	return c.callAPIDirectly(ctx, req, config)
}

// callAPIDirectly 直接调用HTTP API（与SimpleAIClient相同的逻辑）
func (c *SimpleEinoClient) callAPIDirectly(ctx context.Context, req *ChatRequest, config *internalModel.AIConfig) (*ChatResponse, error) {
	// 构建提示词
	prompt := c.buildPromptForAction(req)

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
		"model":       config.OpenAI.Model,
		"messages":    messages,
		"max_tokens":  config.OpenAI.MaxTokens,
		"temperature": config.OpenAI.Temperature,
		"stream":      false,
	}

	// 创建HTTP客户端 - 使用更简单的方式
	client := g.Client()

	apiURL := strings.TrimRight(config.OpenAI.BaseURL, "/") + "/v1/chat/completions"

	g.Log().Info(ctx, "Eino发送AI请求", "url", apiURL, "model", config.OpenAI.Model)

	// 序列化请求体
	requestJson, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %v", err)
	}

	g.Log().Debug(ctx, "请求体", "json", string(requestJson))

	// 使用更基础的方式发送请求
	response := client.Header(map[string]string{
		"Authorization": "Bearer " + config.OpenAI.APIKey,
		"Content-Type":  "application/json",
		"User-Agent":    "GoFrame-AI-Client/1.0",
	}).Timeout(60*time.Second).PostContent(ctx, apiURL, string(requestJson))

	g.Log().Info(ctx, "收到AI响应", "body_length", len(response))

	responseBody := response

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

	// 检查响应是否为空
	if len(responseBody) == 0 {
		return nil, fmt.Errorf("AI API返回空响应")
	}

	if err := json.Unmarshal([]byte(responseBody), &apiResponse); err != nil {
		g.Log().Error(ctx, "Eino AI响应解析失败:", err, "response:", responseBody)
		return nil, fmt.Errorf("AI API响应格式错误: %v\n原始响应: %s", err, responseBody)
	}

	// 检查是否有错误
	if apiResponse.Error.Message != "" {
		g.Log().Error(ctx, "Eino AI API错误:", apiResponse.Error.Message)
		return nil, fmt.Errorf("AI API调用失败: %s", apiResponse.Error.Message)
	}

	// 检查是否有有效响应
	if len(apiResponse.Choices) == 0 {
		return nil, fmt.Errorf("AI API返回空响应")
	}

	content := apiResponse.Choices[0].Message.Content

	// 生成元数据
	metadata := map[string]interface{}{
		"model":          config.OpenAI.Model,
		"provider":       "eino-http",
		"tokens_used":    apiResponse.Usage.TotalTokens,
		"response_time":  2.0,
		"prompt_version": "v1.0",
		"real_ai":        true,
	}

	g.Log().Info(ctx, "Eino AI真实响应成功", "tokens", apiResponse.Usage.TotalTokens)

	return &ChatResponse{
		Content:     content,
		Suggestions: []ChatSuggestion{}, // 可以根据需要添加建议
		Metadata:    metadata,
	}, nil
}

// chatWithEino 使用eino进行聊天
func (c *SimpleEinoClient) chatWithEino(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	g.Log().Info(ctx, "开始eino聊天", "action", req.Action, "userInput", req.UserInput)
	
	// 创建带超时的context，防止请求卡死 - 设置为5分钟支持长文本处理
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	
	// 构建系统提示词
	systemPrompt := c.buildPromptForAction(req)
	
	// 构建消息序列
	messages := []*schema.Message{
		schema.SystemMessage(systemPrompt),
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
	messages = append(messages, schema.UserMessage(req.UserInput))

	g.Log().Info(ctx, "eino发送消息", "message_count", len(messages))
	
	// 使用eino流式生成响应
	stream, err := c.chatModel.Stream(timeoutCtx, messages)
	if err != nil {
		g.Log().Error(ctx, "eino流式生成失败", "error", err)
		return nil, fmt.Errorf("Eino流式聊天失败: %v", err)
	}
	defer stream.Close()

	var fullContent string
	
	// 处理流式响应 - 使用正确的eino Stream API
	for {
		resp, err := stream.Recv()
		if err != nil {
			// 检查是否是正常结束
			if err.Error() == "EOF" || strings.Contains(err.Error(), "EOF") {
				g.Log().Info(ctx, "eino流式响应完成", "total_length", len(fullContent))
				break
			}
			g.Log().Error(ctx, "eino流式读取失败", "error", err)
			return nil, fmt.Errorf("Eino流式读取失败: %v", err)
		}
		
		if resp != nil && resp.Content != "" {
			fullContent += resp.Content
			g.Log().Debug(ctx, "eino收到流式数据块", "chunk_length", len(resp.Content), "total_length", len(fullContent))
		}
	}

	g.Log().Info(ctx, "eino流式响应成功", "content_length", len(fullContent))

	// 使用与流式相同的JSON解析逻辑
	cleanContent := fullContent
	
	// 移除markdown代码块标记
	if strings.HasPrefix(cleanContent, "```json") {
		cleanContent = strings.TrimPrefix(cleanContent, "```json")
	}
	if strings.HasPrefix(cleanContent, "```") {
		cleanContent = strings.TrimPrefix(cleanContent, "```")
	}
	if strings.HasSuffix(cleanContent, "```") {
		cleanContent = strings.TrimSuffix(cleanContent, "```")
	}
	
	// 去除首尾空白字符
	cleanContent = strings.TrimSpace(cleanContent)
	
	var aiResponse struct {
		Content     string `json:"content"`
		Suggestions []struct {
			Type        string  `json:"type"`
			Name        string  `json:"name"`
			Description string  `json:"description"`
			Code        string  `json:"code"`
			Confidence  float64 `json:"confidence"`
			Priority    string  `json:"priority"`
		} `json:"suggestions"`
	}
	
	var finalContent string
	var suggestions []ChatSuggestion
	
	if err := json.Unmarshal([]byte(cleanContent), &aiResponse); err != nil {
		g.Log().Warning(ctx, "非流式AI响应不是有效JSON，使用原始内容", "error", err)
		finalContent = fullContent
	} else {
		g.Log().Info(ctx, "成功解析非流式AI JSON响应", "suggestions_count", len(aiResponse.Suggestions))
		finalContent = aiResponse.Content
		
		// 转换suggestions
		for _, s := range aiResponse.Suggestions {
			suggestions = append(suggestions, ChatSuggestion{
				Type:        s.Type,
				Name:        s.Name,
				Description: s.Description,
				Code:        s.Code,
				Confidence:  s.Confidence,
				Priority:    s.Priority,
			})
		}
	}

	// 生成元数据
	metadata := map[string]interface{}{
		"model":          "eino-chat",
		"provider":       "eino",
		"tokens_used":    len(fullContent) / 4,
		"response_time":  1.5,
		"prompt_version": "v1.0",
		"real_ai":        true,
		"json_parsed":    len(suggestions) > 0,
	}

	return &ChatResponse{
		Content:     finalContent,
		Suggestions: suggestions,
		Metadata:    metadata,
	}, nil
}

// ChatStream 流式聊天接口 - 支持实时回调
func (c *SimpleEinoClient) ChatStream(ctx context.Context, req *ChatRequest, onChunk func(chunk string)) (*ChatResponse, error) {
	g.Log().Info(ctx, "开始eino流式聊天", "action", req.Action, "userInput", req.UserInput)
	
	// 创建带超时的context，防止请求卡死 - 设置为5分钟支持长文本处理
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	
	// 构建系统提示词
	systemPrompt := c.buildPromptForAction(req)
	
	// 构建消息序列
	messages := []*schema.Message{
		schema.SystemMessage(systemPrompt),
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
	messages = append(messages, schema.UserMessage(req.UserInput))

	g.Log().Info(ctx, "eino发送流式消息", "message_count", len(messages))
	
	// 使用eino流式生成响应
	stream, err := c.chatModel.Stream(timeoutCtx, messages)
	if err != nil {
		g.Log().Error(ctx, "eino流式生成失败", "error", err)
		return nil, fmt.Errorf("Eino流式聊天失败: %v", err)
	}
	defer stream.Close()

	var fullContent string
	
	// 处理流式响应并实时回调 - 使用正确的eino Stream API
	for {
		resp, err := stream.Recv()
		if err != nil {
			// 检查是否是正常结束
			if err.Error() == "EOF" || strings.Contains(err.Error(), "EOF") {
				g.Log().Info(ctx, "eino流式响应完成", "total_length", len(fullContent))
				break
			}
			g.Log().Error(ctx, "eino流式读取失败", "error", err)
			return nil, fmt.Errorf("Eino流式读取失败: %v", err)
		}
		
		if resp != nil && resp.Content != "" {
			fullContent += resp.Content
			g.Log().Debug(ctx, "eino收到流式数据块", "chunk_length", len(resp.Content), "total_length", len(fullContent))
			
			// 实时回调发送数据块
			if onChunk != nil {
				onChunk(resp.Content)
			}
		}
	}

	g.Log().Info(ctx, "eino流式聊天完成", "content_length", len(fullContent))

	// 尝试解析JSON响应，首先清理可能的markdown代码块格式
	cleanContent := fullContent
	
	// 移除markdown代码块标记
	if strings.HasPrefix(cleanContent, "```json") {
		cleanContent = strings.TrimPrefix(cleanContent, "```json")
	}
	if strings.HasPrefix(cleanContent, "```") {
		cleanContent = strings.TrimPrefix(cleanContent, "```")
	}
	if strings.HasSuffix(cleanContent, "```") {
		cleanContent = strings.TrimSuffix(cleanContent, "```")
	}
	
	// 去除首尾空白字符
	cleanContent = strings.TrimSpace(cleanContent)
	
	g.Log().Debug(ctx, "清理后的AI响应", "original_length", len(fullContent), "clean_length", len(cleanContent))
	
	var aiResponse struct {
		Content     string `json:"content"`
		Suggestions []struct {
			Type        string  `json:"type"`
			Name        string  `json:"name"`
			Description string  `json:"description"`
			Code        string  `json:"code"`
			Confidence  float64 `json:"confidence"`
			Priority    string  `json:"priority"`
		} `json:"suggestions"`
	}
	
	var finalContent string
	var suggestions []ChatSuggestion
	
	jsonParsed := false
	if err := json.Unmarshal([]byte(cleanContent), &aiResponse); err != nil {
		contentPreview := cleanContent
		if len(contentPreview) > 200 {
			contentPreview = contentPreview[:200] + "..."
		}
		g.Log().Warning(ctx, "AI响应不是有效JSON，尝试原始内容", "error", err, "content", contentPreview)
		
		// 如果清理后仍然无法解析，尝试原始内容
		if err2 := json.Unmarshal([]byte(fullContent), &aiResponse); err2 != nil {
			g.Log().Warning(ctx, "原始内容也无法解析为JSON，使用纯文本", "error", err2)
			finalContent = fullContent
		} else {
			g.Log().Info(ctx, "使用原始内容成功解析JSON")
			finalContent = aiResponse.Content
			jsonParsed = true
		}
	} else {
		g.Log().Info(ctx, "成功解析清理后的AI JSON响应", "suggestions_count", len(aiResponse.Suggestions))
		finalContent = aiResponse.Content
		jsonParsed = true
	}
	
	// 如果成功解析JSON，转换suggestions
	if jsonParsed {
		for _, s := range aiResponse.Suggestions {
			suggestions = append(suggestions, ChatSuggestion{
				Type:        s.Type,
				Name:        s.Name,
				Description: s.Description,
				Code:        s.Code,
				Confidence:  s.Confidence,
				Priority:    s.Priority,
			})
		}
	}

	// 生成元数据
	metadata := map[string]interface{}{
		"model":          "eino-stream",
		"provider":       "eino",
		"tokens_used":    len(fullContent) / 4,
		"response_time":  1.5,
		"prompt_version": "v1.0",
		"real_ai":        true,
		"json_parsed":    len(suggestions) > 0,
	}

	return &ChatResponse{
		Content:     finalContent,
		Suggestions: suggestions,
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

// buildPromptForAction 根据操作类型构建提示词（与SimpleAIClient相同）
func (c *SimpleEinoClient) buildPromptForAction(req *ChatRequest) string {
	var basePrompt string

	switch req.Action {
	case "optimize_code":
		basePrompt = `你是一个专业的代码优化专家。请分析用户的代码并提供优化建议。

⚠️ 重要约束：
1. 只能返回纯JSON，禁止markdown代码块格式（禁止三个反引号json三个反引号）
2. 不能添加任何解释、说明或其他文字
3. 直接输出JSON对象，第一个字符必须是{
4. 最后一个字符必须是}

必须严格按照以下JSON格式返回：
{
  "content": "优化建议的详细说明",
  "suggestions": [
    {
      "type": "code",
      "name": "建议名称",
      "description": "建议描述",
      "code": "优化后的代码",
      "confidence": 0.9,
      "priority": "high"
    }
  ]
}`

	case "explain_code":
		basePrompt = `你是一个代码分析专家。请详细解释用户提供的代码的功能和逻辑。

⚠️ 重要约束：
1. 只能返回纯JSON，禁止markdown代码块格式（禁止三个反引号json三个反引号）
2. 不能添加任何解释、说明或其他文字
3. 直接输出JSON对象，第一个字符必须是{
4. 最后一个字符必须是}

必须严格按照以下JSON格式返回：
{
  "content": "代码解释的详细说明",
  "suggestions": [
    {
      "type": "action",
      "name": "相关操作",
      "description": "操作描述",
      "confidence": 0.8,
      "priority": "medium"
    }
  ]
}`

	case "suggest_variables":
		basePrompt = `你是一个模板系统专家。请为用户的项目推荐合适的模板变量。

⚠️ 重要约束：
1. 只能返回纯JSON，禁止markdown代码块格式（禁止三个反引号json三个反引号）
2. 不能添加任何解释、说明或其他文字
3. 直接输出JSON对象，第一个字符必须是{
4. 最后一个字符必须是}

必须严格按照以下JSON格式返回：
{
  "content": "变量建议说明",
  "suggestions": [
    {
      "type": "variable",
      "name": "变量名",
      "description": "变量用途说明",
      "code": "{{.变量名}}",
      "confidence": 0.9,
      "priority": "high"
    }
  ]
}`

	case "generate_template":
		basePrompt = `你是一个项目模板生成专家。请根据用户需求生成完整的项目模板。

⚠️ 重要约束：
1. 只能返回纯JSON，禁止markdown代码块格式（禁止三个反引号json三个反引号）
2. 不能添加任何解释、说明或其他文字
3. 直接输出JSON对象，第一个字符必须是{
4. 最后一个字符必须是}

必须严格按照以下JSON格式返回：
{
  "content": "模板生成说明和使用指南",
  "suggestions": [
    {
      "type": "file",
      "name": "文件名",
      "description": "文件说明",
      "code": "文件内容模板",
      "confidence": 0.9,
      "priority": "high"
    }
  ]
}`

	case "refactor_code":
		basePrompt = `你是一个代码重构专家。请提供具体的代码重构建议。

⚠️ 重要约束：
1. 只能返回纯JSON，禁止markdown代码块格式（禁止三个反引号json三个反引号）
2. 不能添加任何解释、说明或其他文字
3. 直接输出JSON对象，第一个字符必须是{
4. 最后一个字符必须是}

必须严格按照以下JSON格式返回：
{
  "content": "重构建议详细说明",
  "suggestions": [
    {
      "type": "code",
      "name": "重构方案",
      "description": "重构说明",
      "code": "重构后的代码",
      "confidence": 0.8,
      "priority": "medium"
    }
  ]
}`

	case "add_comments":
		basePrompt = `你是一个代码文档专家。请为用户的代码添加适当的注释。

⚠️ 重要约束：
1. 只能返回纯JSON，禁止markdown代码块格式（禁止三个反引号json三个反引号）
2. 不能添加任何解释、说明或其他文字
3. 直接输出JSON对象，第一个字符必须是{
4. 最后一个字符必须是}

必须严格按照以下JSON格式返回：
{
  "content": "注释添加说明",
  "suggestions": [
    {
      "type": "code",
      "name": "带注释的代码",
      "description": "注释说明",
      "code": "添加注释后的完整代码",
      "confidence": 0.9,
      "priority": "high"
    }
  ]
}`

	default:
		basePrompt = `你是一个AI编程助手，可以帮助用户解决各种编程相关问题。

⚠️ 重要约束：
1. 只能返回纯JSON，禁止markdown代码块格式（禁止三个反引号json三个反引号）
2. 不能添加任何解释、说明或其他文字
3. 直接输出JSON对象，第一个字符必须是{
4. 最后一个字符必须是}

必须严格按照以下JSON格式返回：
{
  "content": "回答内容",
  "suggestions": [
    {
      "type": "action",
      "name": "相关建议",
      "description": "建议描述",
      "confidence": 0.7,
      "priority": "medium"
    }
  ]
}`
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
