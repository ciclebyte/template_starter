package ai

import (
	"context"
	"time"

	aiApi "github.com/ciclebyte/template_starter/api/v1/ai"
	"github.com/ciclebyte/template_starter/internal/service"
	"github.com/ciclebyte/template_starter/library/libAI"
	"github.com/ciclebyte/template_starter/library/libConfig"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type sAI struct{}

func init() {
	service.RegisterAIService(New())
}

func New() *sAI {
	return &sAI{}
}

// TestConnection 测试AI连接
func (s *sAI) TestConnection(ctx context.Context, req *aiApi.TestConnectionReq) (*aiApi.TestConnectionRes, error) {
	g.Log().Debug(ctx, "AI.TestConnection called")

	res := &aiApi.TestConnectionRes{
		Success: false,
		Message: "连接失败",
	}

	// 获取AI配置
	config, err := libConfig.GetAIConfig(ctx)
	if err != nil {
		res.Message = "获取AI配置失败: " + err.Error()
		return res, nil
	}
	
	// 调试日志：显示实际读取的配置
	g.Log().Debug(ctx, "AI配置详情", "provider", config.Provider, "openai_model", config.OpenAI.Model, "claude_model", config.Claude.Model)

	if !config.Enabled {
		res.Message = "AI功能未启用"
		return res, nil
	}

	res.Provider = config.Provider

	// 创建AI客户端
	client, err := libAI.NewAIClient(ctx)
	if err != nil {
		res.Message = "创建AI客户端失败: " + err.Error()
		return res, nil
	}

	// 记录开始时间
	startTime := time.Now()

	// 测试连接
	err = client.TestConnection(ctx)
	
	// 计算延迟
	res.Latency = int(time.Since(startTime).Milliseconds())

	if err != nil {
		res.Message = "连接测试失败: " + err.Error()
		return res, nil
	}

	// 设置成功信息
	res.Success = true
	res.Message = "连接测试成功"
	
	// 显示实际配置的模型
	if config.Provider == "openai" {
		res.Model = config.OpenAI.Model
	} else if config.Provider == "claude" {
		res.Model = config.Claude.Model
	} else {
		// 对于其他provider，也显示OpenAI模型配置（可能是moonshot等兼容API）
		res.Model = config.OpenAI.Model
	}

	return res, nil
}

// GenerateTemplate 生成模板
func (s *sAI) GenerateTemplate(ctx context.Context, req *aiApi.GenerateTemplateReq) (*aiApi.GenerateTemplateRes, error) {
	g.Log().Debug(ctx, "AI.GenerateTemplate called with req:", req)

	// 检查AI功能是否启用
	config, err := libConfig.GetAIConfig(ctx)
	if err != nil {
		return nil, err
	}

	if !config.Enabled || !config.Features.TemplateGeneration {
		return nil, gerror.New("AI模板生成功能未启用")
	}

	// 创建AI客户端
	client, err := libAI.NewAIClient(ctx)
	if err != nil {
		return nil, err
	}

	// 构建请求
	aiReq := &libAI.TemplateGenerateRequest{
		Description:   req.Description,
		ProjectType:   req.ProjectType,
		TechStack:     req.TechStack,
		Variables:     req.Variables,
		Features:      req.Features,
	}

	// 调用AI生成
	aiRes, err := client.GenerateTemplate(ctx, aiReq)
	if err != nil {
		return nil, err
	}

	// 转换响应格式
	res := &aiApi.GenerateTemplateRes{
		Instructions:  aiRes.Instructions,
		EstimatedTime: aiRes.EstimatedTime,
	}

	// 转换文件结构
	for _, file := range aiRes.ProjectStructure {
		res.ProjectStructure = append(res.ProjectStructure, aiApi.FileInfo{
			Path:        file.Path,
			Content:     file.Content,
			IsDirectory: file.IsDirectory,
			Description: file.Description,
		})
	}

	// 转换变量信息
	for _, variable := range aiRes.Variables {
		res.Variables = append(res.Variables, aiApi.VariableInfo{
			Name:         variable.Name,
			Type:         variable.Type,
			Description:  variable.Description,
			DefaultValue: variable.DefaultValue,
			Required:     variable.Required,
			Options:      variable.Options,
		})
	}

	return res, nil
}

// SuggestVariables 建议变量
func (s *sAI) SuggestVariables(ctx context.Context, req *aiApi.SuggestVariablesReq) (*aiApi.SuggestVariablesRes, error) {
	g.Log().Debug(ctx, "AI.SuggestVariables called with req:", req)

	// 检查AI功能是否启用
	config, err := libConfig.GetAIConfig(ctx)
	if err != nil {
		return nil, err
	}

	if !config.Enabled || !config.Features.VariableSuggestion {
		return nil, gerror.New("AI变量建议功能未启用")
	}

	// 创建AI客户端
	client, err := libAI.NewAIClient(ctx)
	if err != nil {
		return nil, err
	}

	// 构建请求
	aiReq := &libAI.VariableSuggestRequest{
		ProjectType: req.ProjectType,
		TechStack:   req.TechStack,
		Description: req.Description,
	}

	// 调用AI建议
	aiRes, err := client.SuggestVariables(ctx, aiReq)
	if err != nil {
		return nil, err
	}

	// 转换响应格式
	res := &aiApi.SuggestVariablesRes{}
	for _, variable := range aiRes.Variables {
		res.Variables = append(res.Variables, aiApi.VariableInfo{
			Name:         variable.Name,
			Type:         variable.Type,
			Description:  variable.Description,
			DefaultValue: variable.DefaultValue,
			Required:     variable.Required,
			Options:      variable.Options,
		})
	}

	return res, nil
}

// Chat 统一AI聊天接口
func (s *sAI) Chat(ctx context.Context, req *aiApi.ChatReq) (*aiApi.ChatRes, error) {
	g.Log().Debug(ctx, "AI.Chat called with action:", req.Action)

	// 检查AI功能是否启用
	config, err := libConfig.GetAIConfig(ctx)
	if err != nil {
		return nil, err
	}

	if !config.Enabled {
		return nil, gerror.New("AI功能未启用")
	}

	// 创建AI客户端
	client, err := libAI.NewAIClient(ctx)
	if err != nil {
		return nil, err
	}

	// 构建请求
	aiReq := &libAI.ChatRequest{
		Action:      req.Action,
		Context:     req.Context,
		UserInput:   req.UserInput,
		Preferences: req.Preferences,
		ChatHistory: convertChatHistory(req.ChatHistory),
	}

	// 调用AI聊天
	aiRes, err := client.Chat(ctx, aiReq)
	if err != nil {
		return nil, err
	}

	// 转换响应格式
	res := &aiApi.ChatRes{
		Content:  aiRes.Content,
		Metadata: aiRes.Metadata,
	}

	// 转换建议
	for _, suggestion := range aiRes.Suggestions {
		res.Suggestions = append(res.Suggestions, aiApi.ChatSuggestion{
			Type:        suggestion.Type,
			Name:        suggestion.Name,
			Description: suggestion.Description,
			Code:        suggestion.Code,
			Confidence:  suggestion.Confidence,
			Priority:    suggestion.Priority,
		})
	}

	return res, nil
}

// ChatStream 流式AI聊天接口
func (s *sAI) ChatStream(ctx context.Context, req *aiApi.ChatReq, r *ghttp.Request) {
	g.Log().Debug(ctx, "AI.ChatStream called with action:", req.Action)

	// 设置SSE响应头
	r.Response.Header().Set("Content-Type", "text/event-stream")
	r.Response.Header().Set("Cache-Control", "no-cache")
	r.Response.Header().Set("Connection", "keep-alive")
	r.Response.Header().Set("Access-Control-Allow-Origin", "*")

	// 发送流式数据的辅助函数
	writeSSE := func(data *aiApi.ChatStreamData) {
		jsonData, _ := gjson.Marshal(data)
		r.Response.Write("data: " + string(jsonData) + "\n\n")
		r.Response.Flush()
	}

	// 检查AI功能是否启用
	config, err := libConfig.GetAIConfig(ctx)
	if err != nil {
		writeSSE(&aiApi.ChatStreamData{
			Type:    "error",
			Content: "获取AI配置失败: " + err.Error(),
			Done:    true,
		})
		return
	}

	if !config.Enabled {
		writeSSE(&aiApi.ChatStreamData{
			Type:    "error",
			Content: "AI功能未启用",
			Done:    true,
		})
		return
	}

	// 创建AI客户端
	client, err := libAI.NewAIClient(ctx)
	if err != nil {
		writeSSE(&aiApi.ChatStreamData{
			Type:    "error",
			Content: "创建AI客户端失败: " + err.Error(),
			Done:    true,
		})
		return
	}

	// 构建请求
	aiReq := &libAI.ChatRequest{
		Action:      req.Action,
		Context:     req.Context,
		UserInput:   req.UserInput,
		Preferences: req.Preferences,
		ChatHistory: convertChatHistory(req.ChatHistory),
	}

	// 发送开始信号
	writeSSE(&aiApi.ChatStreamData{
		Type:    "start",
		Content: "开始生成AI响应...",
		Done:    false,
	})

	// 调用AI聊天 - 这里暂时使用非流式调用，后续可以改进为真正的流式
	aiRes, err := client.Chat(ctx, aiReq)
	if err != nil {
		writeSSE(&aiApi.ChatStreamData{
			Type:    "error",
			Content: "AI调用失败: " + err.Error(),
			Done:    true,
		})
		return
	}

	// 模拟流式输出（逐字输出）
	content := aiRes.Content
	chunkSize := 20 // 每次发送的字符数
	
	for i := 0; i < len(content); i += chunkSize {
		end := i + chunkSize
		if end > len(content) {
			end = len(content)
		}
		
		chunk := content[i:end]
		writeSSE(&aiApi.ChatStreamData{
			Type:    "chunk",
			Content: chunk,
			Done:    false,
		})
		
		// 模拟打字效果
		time.Sleep(50 * time.Millisecond)
	}

	// 发送建议和元数据
	var suggestions []aiApi.ChatSuggestion
	for _, suggestion := range aiRes.Suggestions {
		suggestions = append(suggestions, aiApi.ChatSuggestion{
			Type:        suggestion.Type,
			Name:        suggestion.Name,
			Description: suggestion.Description,
			Code:        suggestion.Code,
			Confidence:  suggestion.Confidence,
			Priority:    suggestion.Priority,
		})
	}

	if len(suggestions) > 0 {
		writeSSE(&aiApi.ChatStreamData{
			Type:        "suggestions",
			Suggestions: suggestions,
			Done:        false,
		})
	}

	// 发送元数据
	writeSSE(&aiApi.ChatStreamData{
		Type:     "metadata",
		Metadata: aiRes.Metadata,
		Done:     false,
	})

	// 发送完成信号  
	writeSSE(&aiApi.ChatStreamData{
		Type:    "done",
		Content: "",
		Done:    true,
	})
}

// convertChatHistory 转换聊天历史格式
func convertChatHistory(apiHistory []aiApi.ChatMessage) []libAI.ChatMessage {
	var result []libAI.ChatMessage
	for _, msg := range apiHistory {
		result = append(result, libAI.ChatMessage{
			Role:      msg.Role,
			Content:   msg.Content,
			Timestamp: msg.Timestamp,
		})
	}
	return result
}