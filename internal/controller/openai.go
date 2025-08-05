package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	r.Response.Header().Set("X-Accel-Buffering", "no") // 禁用nginx缓冲

	g.Log().Info(ctx, "OpenAI Stream - 开始流式响应", "model", req.Model)

	// 转换为内部AI请求格式
	aiReq, err := c.convertToAIRequest(req)
	if err != nil {
		c.writeOpenAIStreamError(r, "请求转换失败: "+err.Error())
		return
	}

	// 生成响应ID和时间戳
	responseID := "chatcmpl-" + guid.S()
	created := time.Now().Unix()

	// 直接调用真正的流式AI服务，使用拦截器转换格式
	c.streamWithRealAIService(ctx, r, aiReq, responseID, created, req.Model)
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

// streamOpenAIResponse 按OpenAI标准格式输出流式响应
func (c *cOpenAI) streamOpenAIResponse(ctx context.Context, r *ghttp.Request, aiReq *aiApi.ChatReq, responseID string, created int64, model string) {
	// 发送开始chunk
	startChunk := &openaiApi.ChatCompletionStreamResponse{
		ID:      responseID,
		Object:  "chat.completion.chunk",
		Created: created,
		Model:   model,
		Choices: []openaiApi.ChatCompletionStreamChoice{
			{
				Index: 0,
				Delta: openaiApi.ChatCompletionMessage{
					Role: "assistant",
				},
				FinishReason: nil,
			},
		},
	}
	c.writeOpenAIStreamChunk(r, startChunk)

	// 创建自定义的响应处理器，将内部AI服务的流式响应转换为OpenAI格式
	c.streamWithOpenAIFormat(ctx, r, aiReq, responseID, created, model)
}

// streamWithOpenAIFormat 使用内部AI服务的真正流式响应，并转换为OpenAI格式
func (c *cOpenAI) streamWithOpenAIFormat(ctx context.Context, r *ghttp.Request, aiReq *aiApi.ChatReq, responseID string, created int64, model string) {
	// 这个函数已经被新的streamWithRealAIService替代
	// 直接调用新的实现
	c.streamWithRealAIService(ctx, r, aiReq, responseID, created, model)
}

// writeOpenAIStreamChunk 写入OpenAI格式的流式数据块
func (c *cOpenAI) writeOpenAIStreamChunk(r *ghttp.Request, chunk *openaiApi.ChatCompletionStreamResponse) {
	jsonData, _ := gjson.Marshal(chunk)
	sseData := "data: " + string(jsonData) + "\n\n"

	// 立即写入并刷新
	r.Response.Write(sseData)
	r.Response.Flush()
	g.Log().Debug(r.Context(), "OpenAI Stream - chunk已发送", "size", len(sseData))
}

// writeOpenAIStreamError 写入OpenAI格式的流式错误
func (c *cOpenAI) writeOpenAIStreamError(r *ghttp.Request, message string) {
	errorData := openaiApi.OpenAIError{
		Error: openaiApi.OpenAIErrorDetail{
			Message: message,
			Type:    "internal_error",
		},
	}
	jsonData, _ := gjson.Marshal(errorData)
	sseData := "data: " + string(jsonData) + "\n\n"
	r.Response.Write(sseData)
	r.Response.Flush()
}

// stringPtr 返回字符串指针的辅助函数
func stringPtr(s string) *string {
	return &s
}

// openAISSEInterceptor 拦截内部AI流式响应并转换为OpenAI格式
type openAISSEInterceptor struct {
	originalWriter io.Writer
	responseID     string
	created        int64
	model          string
	ctx            context.Context
}

func (o *openAISSEInterceptor) Write(data []byte) (int, error) {
	content := string(data)

	// 处理SSE格式的数据
	if strings.HasPrefix(content, "data: ") {
		jsonStr := strings.TrimPrefix(content, "data: ")
		jsonStr = strings.TrimSuffix(jsonStr, "\n\n")

		if jsonStr != "" && jsonStr != "[DONE]" {
			// 解析内部AI服务的流式数据
			var streamData aiApi.ChatStreamData
			if err := gjson.Unmarshal([]byte(jsonStr), &streamData); err == nil {
				// 只转换chunk类型的数据
				if streamData.Type == "chunk" && streamData.Content != "" {
					g.Log().Debug(o.ctx, "OpenAI Stream - 实时转换chunk", "content", streamData.Content, "length", len(streamData.Content))

					// 转换为OpenAI格式
					openaiChunk := &openaiApi.ChatCompletionStreamResponse{
						ID:      o.responseID,
						Object:  "chat.completion.chunk",
						Created: o.created,
						Model:   o.model,
						Choices: []openaiApi.ChatCompletionStreamChoice{
							{
								Index: 0,
								Delta: openaiApi.ChatCompletionMessage{
									Content: streamData.Content,
								},
								FinishReason: nil,
							},
						},
					}

					// 立即发送到原始writer
					jsonData, _ := gjson.Marshal(openaiChunk)
					sseData := "data: " + string(jsonData) + "\n\n"
					o.originalWriter.Write([]byte(sseData))

					// 立即flush
					if flusher, ok := o.originalWriter.(interface{ Flush() }); ok {
						flusher.Flush()
					}

					g.Log().Debug(o.ctx, "OpenAI Stream - chunk已实时转发", "size", len(sseData))
				}
				// 忽略其他类型的数据（start, metadata, suggestions, done等）
			}
		}
	}

	return len(data), nil
}

// streamWithRealAIService 使用真正的AI流式服务
func (c *cOpenAI) streamWithRealAIService(ctx context.Context, r *ghttp.Request, aiReq *aiApi.ChatReq, responseID string, created int64, model string) {
	g.Log().Info(ctx, "OpenAI Stream - 开始真正的流式AI服务调用")

	// 避免多个goroutine同时访问同一个response writer导致的panic
	// 直接使用同步方式调用AI服务，然后以流式方式发送
	
	// 调用AI服务获取完整响应
	aiRes, err := service.AI().Chat(ctx, aiReq)
	if err != nil {
		g.Log().Error(ctx, "AI服务调用失败:", err)
		c.writeOpenAIStreamError(r, "AI服务调用失败: "+err.Error())
		return
	}
	
	content := aiRes.Content
	if content == "" {
		content = "抱歉，我无法处理这个请求。"
	}
	
	g.Log().Info(ctx, "OpenAI Stream - 开始流式发送响应", "content_length", len(content))
	
	// 按字符流式发送，提供更好的用户体验
	runes := []rune(content)
	chunkSize := 3 // 每次发送3个字符，平衡流畅度和性能
	
	for i := 0; i < len(runes); i += chunkSize {
		// 检查上下文是否已取消
		select {
		case <-ctx.Done():
			c.writeOpenAIStreamError(r, "请求已取消")
			return
		default:
		}
		
		end := i + chunkSize
		if end > len(runes) {
			end = len(runes)
		}
		
		chunk := string(runes[i:end])
		
		g.Log().Debug(ctx, "OpenAI Stream - 发送chunk", "content", chunk, "progress", fmt.Sprintf("%d/%d", end, len(runes)))
		
		streamChunk := &openaiApi.ChatCompletionStreamResponse{
			ID:      responseID,
			Object:  "chat.completion.chunk",
			Created: created,
			Model:   model,
			Choices: []openaiApi.ChatCompletionStreamChoice{
				{
					Index: 0,
					Delta: openaiApi.ChatCompletionMessage{
						Content: chunk,
					},
					FinishReason: nil,
				},
			},
		}
		
		c.writeOpenAIStreamChunk(r, streamChunk)
		
		// 适度延迟，模拟真实的流式响应
		time.Sleep(15 * time.Millisecond)
	}
	
	g.Log().Info(ctx, "OpenAI Stream - 内容发送完成")

	// 发送完成标记
	finishChunk := &openaiApi.ChatCompletionStreamResponse{
		ID:      responseID,
		Object:  "chat.completion.chunk",
		Created: created,
		Model:   model,
		Choices: []openaiApi.ChatCompletionStreamChoice{
			{
				Index:        0,
				Delta:        openaiApi.ChatCompletionMessage{},
				FinishReason: stringPtr("stop"),
			},
		},
	}
	c.writeOpenAIStreamChunk(r, finishChunk)

	// 发送结束标记
	r.Response.Write("data: [DONE]\n\n")
	r.Response.Flush()
	g.Log().Info(ctx, "OpenAI Stream - 流式响应完成")
}
