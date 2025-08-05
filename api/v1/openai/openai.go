package openai

import (
	"github.com/gogf/gf/v2/frame/g"
)

// OpenAI兼容的Chat Completions请求
type ChatCompletionRequest struct {
	Model          string                    `json:"model" v:"required#模型不能为空"`
	Messages       []ChatCompletionMessage   `json:"messages" v:"required#消息不能为空"`
	Stream         bool                      `json:"stream"`
	Temperature    float64                   `json:"temperature"`
	MaxTokens      int                       `json:"max_tokens"`
	Functions      []ChatFunction            `json:"functions,omitempty"`
	FunctionCall   interface{}               `json:"function_call,omitempty"` // "auto", "none", or {"name": "function_name"}
	User           string                    `json:"user,omitempty"`
	Stop           []string                  `json:"stop,omitempty"`
	FrequencyPenalty float64                 `json:"frequency_penalty,omitempty"`
	PresencePenalty  float64                 `json:"presence_penalty,omitempty"`
}

// OpenAI兼容的消息格式
type ChatCompletionMessage struct {
	Role         string            `json:"role"`    // system, user, assistant, function
	Content      string            `json:"content"` // 消息内容
	Name         string            `json:"name,omitempty"`         // 对于function消息
	FunctionCall *FunctionCall     `json:"function_call,omitempty"` // AI调用的函数
}

// 函数调用
type FunctionCall struct {
	Name      string `json:"name"`      // 函数名
	Arguments string `json:"arguments"` // JSON格式的参数
}

// 函数定义
type ChatFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// OpenAI兼容的响应格式
type ChatCompletionResponse struct {
	g.Meta  `mime:"application/json" example:"string"`
	ID      string                   `json:"id"`
	Object  string                   `json:"object"`
	Created int64                    `json:"created"`
	Model   string                   `json:"model"`
	Choices []ChatCompletionChoice   `json:"choices"`
	Usage   ChatCompletionUsage      `json:"usage"`
}

// 选择项
type ChatCompletionChoice struct {
	Index        int                     `json:"index"`
	Message      ChatCompletionMessage   `json:"message"`
	FinishReason string                  `json:"finish_reason"` // stop, length, function_call
}

// 使用情况
type ChatCompletionUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// 流式响应的数据块
type ChatCompletionStreamResponse struct {
	ID      string                        `json:"id"`
	Object  string                        `json:"object"`
	Created int64                         `json:"created"`
	Model   string                        `json:"model"`
	Choices []ChatCompletionStreamChoice  `json:"choices"`
}

// 流式选择项
type ChatCompletionStreamChoice struct {
	Index        int                    `json:"index"`
	Delta        ChatCompletionMessage  `json:"delta"` // 增量内容
	FinishReason *string                `json:"finish_reason"` // null或"stop"等
}

// 错误响应格式
type OpenAIError struct {
	Error OpenAIErrorDetail `json:"error"`
}

type OpenAIErrorDetail struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code,omitempty"`
}