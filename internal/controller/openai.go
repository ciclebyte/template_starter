package controller

import (
	"encoding/json"
	"strings"
	"time"

	aiApi "github.com/ciclebyte/template_starter/api/v1/ai"
	openaiApi "github.com/ciclebyte/template_starter/api/v1/openai"
	"github.com/ciclebyte/template_starter/internal/service"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/guid"
)

var OpenAI = cOpenAI{}

type cOpenAI struct {
	BaseController
}

// ChatCompletions OpenAI兼容的聊天完成接口
func (c *cOpenAI) ChatCompletions(r *ghttp.Request) {
	ctx := r.Context()

	// 解析请求体
	var req *openaiApi.ChatCompletionRequest
	if err := r.Parse(&req); err != nil {
		c.writeOpenAIError(r, 400, "invalid_request_error", "请求参数解析失败: "+err.Error())
		return
	}

	g.Log().Debug(ctx, "OpenAI.ChatCompletions called", "model", req.Model, "stream", req.Stream)

	// 如果是流式请求，处理流式响应
	if req.Stream {
		c.handleStreamRequest(r, req)
		return
	}

	// 处理非流式请求
	c.handleNonStreamRequest(r, req)
}

// handleStreamRequest 处理流式请求
func (c *cOpenAI) handleStreamRequest(r *ghttp.Request, req *openaiApi.ChatCompletionRequest) {
	ctx := r.Context()
	g.Log().Debug(ctx, "OpenAI.handleStreamRequest called")

	// 设置SSE响应头
	r.Response.Header().Set("Content-Type", "text/event-stream")
	r.Response.Header().Set("Cache-Control", "no-cache")
	r.Response.Header().Set("Connection", "keep-alive")
	r.Response.Header().Set("Access-Control-Allow-Origin", "*")

	// 转换为内部AI请求格式
	aiReq, err := c.convertToAIRequest(req)
	if err != nil {
		c.writeStreamError(r, "请求转换失败: "+err.Error())
		return
	}

	// 调用内部流式AI服务
	service.AI().ChatStream(ctx, aiReq, r)
}

// handleNonStreamRequest 处理非流式请求
func (c *cOpenAI) handleNonStreamRequest(r *ghttp.Request, req *openaiApi.ChatCompletionRequest) {
	ctx := r.Context()
	g.Log().Debug(ctx, "OpenAI.handleNonStreamRequest called")

	// 转换为内部AI请求格式
	aiReq, err := c.convertToAIRequest(req)
	if err != nil {
		c.writeOpenAIError(r, 400, "invalid_request_error", "请求转换失败: "+err.Error())
		return
	}

	// 调用内部AI服务
	aiRes, err := service.AI().Chat(ctx, aiReq)
	if err != nil {
		c.writeOpenAIError(r, 500, "internal_error", "AI服务调用失败: "+err.Error())
		return
	}

	// 转换为OpenAI格式响应并返回
	response := c.convertToOpenAIResponse(req, aiRes)
	r.Response.WriteJson(response)
}

// convertToAIRequest 将OpenAI请求转换为内部AI请求格式
func (c *cOpenAI) convertToAIRequest(req *openaiApi.ChatCompletionRequest) (*aiApi.ChatReq, error) {
	// 提取系统消息和用户消息
	var systemPrompt string
	var userInput string
	var chatHistory []aiApi.ChatMessage

	for _, msg := range req.Messages {
		switch msg.Role {
		case "system":
			systemPrompt = msg.Content
		case "user":
			userInput = msg.Content // 使用最后一个用户消息作为当前输入
			chatHistory = append(chatHistory, aiApi.ChatMessage{
				Role:      "user",
				Content:   msg.Content,
				Timestamp: time.Now().Format(time.RFC3339),
			})
		case "assistant":
			chatHistory = append(chatHistory, aiApi.ChatMessage{
				Role:      "assistant",
				Content:   msg.Content,
				Timestamp: time.Now().Format(time.RFC3339),
			})
		}
	}

	// 确定操作类型
	action := c.detectActionFromMessages(req.Messages, req.Functions)

	// 构建上下文
	context := map[string]interface{}{
		"model":        req.Model,
		"systemPrompt": systemPrompt,
	}

	// 如果有函数定义，添加到上下文
	if len(req.Functions) > 0 {
		context["functions"] = req.Functions
		context["function_call"] = req.FunctionCall
	}

	// 构建偏好设置
	preferences := map[string]interface{}{
		"temperature": req.Temperature,
		"max_tokens":  req.MaxTokens,
	}

	return &aiApi.ChatReq{
		Action:      action,
		Context:     context,
		UserInput:   userInput,
		Preferences: preferences,
		ChatHistory: chatHistory,
		Stream:      req.Stream,
	}, nil
}

// detectActionFromMessages 从消息和函数中检测操作类型
func (c *cOpenAI) detectActionFromMessages(messages []openaiApi.ChatCompletionMessage, functions []openaiApi.ChatFunction) string {
	// 如果有函数定义，根据函数名确定操作类型
	if len(functions) > 0 {
		for _, fn := range functions {
			switch fn.Name {
			case "suggest_variables":
				return "suggest_variables"
			case "optimize_code":
				return "optimize_code"
			case "explain_code":
				return "explain_code"
			case "generate_template":
				return "generate_template"
			}
		}
	}

	// 从消息内容推断操作类型
	for _, msg := range messages {
		if msg.Role == "user" {
			content := strings.ToLower(msg.Content)
			if strings.Contains(content, "优化") || strings.Contains(content, "optimize") {
				return "optimize_code"
			}
			if strings.Contains(content, "解释") || strings.Contains(content, "explain") {
				return "explain_code"
			}
			if strings.Contains(content, "变量") || strings.Contains(content, "variable") {
				return "suggest_variables"
			}
			if strings.Contains(content, "模板") || strings.Contains(content, "template") {
				return "generate_template"
			}
		}
	}

	// 默认为通用聊天
	return "general_chat"
}

// convertToOpenAIResponse 将内部AI响应转换为OpenAI格式
func (c *cOpenAI) convertToOpenAIResponse(req *openaiApi.ChatCompletionRequest, aiRes *aiApi.ChatRes) *openaiApi.ChatCompletionResponse {
	responseID := "chatcmpl-" + guid.S()
	created := time.Now().Unix()

	message := openaiApi.ChatCompletionMessage{
		Role:    "assistant",
		Content: aiRes.Content,
	}

	// 如果有建议，转换为function_call
	if len(aiRes.Suggestions) > 0 {
		functionResult := map[string]interface{}{
			"suggestions": aiRes.Suggestions,
		}

		arguments, _ := json.Marshal(functionResult)
		message.FunctionCall = &openaiApi.FunctionCall{
			Name:      "suggest_actions",
			Arguments: string(arguments),
		}
	}

	return &openaiApi.ChatCompletionResponse{
		ID:      responseID,
		Object:  "chat.completion",
		Created: created,
		Model:   req.Model,
		Choices: []openaiApi.ChatCompletionChoice{
			{
				Index:        0,
				Message:      message,
				FinishReason: "stop",
			},
		},
		Usage: openaiApi.ChatCompletionUsage{
			PromptTokens:     c.estimateTokens(req.Messages),
			CompletionTokens: c.estimateTokens([]openaiApi.ChatCompletionMessage{{Content: aiRes.Content}}),
			TotalTokens:      c.estimateTokens(req.Messages) + c.estimateTokens([]openaiApi.ChatCompletionMessage{{Content: aiRes.Content}}),
		},
	}
}

// estimateTokens 简单的token估算（实际应该使用专业的tokenizer）
func (c *cOpenAI) estimateTokens(messages []openaiApi.ChatCompletionMessage) int {
	totalChars := 0
	for _, msg := range messages {
		totalChars += len(msg.Content)
	}
	// 粗略估算：4个字符约等于1个token
	return totalChars / 4
}

// writeOpenAIError 写入OpenAI格式的错误响应
func (c *cOpenAI) writeOpenAIError(r *ghttp.Request, statusCode int, errorType, message string) {
	r.Response.WriteHeader(statusCode)
	r.Response.WriteJson(openaiApi.OpenAIError{
		Error: openaiApi.OpenAIErrorDetail{
			Message: message,
			Type:    errorType,
		},
	})
}

// writeStreamError 写入流式错误
func (c *cOpenAI) writeStreamError(r *ghttp.Request, message string) {
	errorData := map[string]interface{}{
		"error": map[string]interface{}{
			"message": message,
			"type":    "internal_error",
		},
	}
	jsonData, _ := gjson.Marshal(errorData)
	sseData := "data: " + string(jsonData) + "\n\n"
	r.Response.Write(sseData)
	r.Response.Flush()
}
