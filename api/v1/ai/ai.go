package ai

import (
	"github.com/gogf/gf/v2/frame/g"
)

// AI连接测试请求
type TestConnectionReq struct {
	g.Meta `path:"/ai/testConnection" method:"post" tags:"AI服务" summary:"AI服务-测试连接"`
}

type TestConnectionRes struct {
	g.Meta    `mime:"application/json" example:"string"`
	Success   bool   `json:"success"`   // 是否连接成功
	Message   string `json:"message"`   // 连接结果消息
	Provider  string `json:"provider"`  // 使用的服务商
	Model     string `json:"model"`     // 使用的模型
	Latency   int    `json:"latency"`   // 响应延迟(毫秒)
}

// AI模板生成请求
type GenerateTemplateReq struct {
	g.Meta        `path:"/ai/generateTemplate" method:"post" tags:"AI服务" summary:"AI服务-生成模板"`
	Description   string            `json:"description" v:"required#项目描述不能为空"`   // 项目描述
	ProjectType   string            `json:"projectType" v:"required#项目类型不能为空"`   // 项目类型
	TechStack     []string          `json:"techStack"`                             // 技术栈
	Variables     map[string]string `json:"variables"`                             // 现有变量
	Features      []string          `json:"features"`                              // 功能特性
}

type GenerateTemplateRes struct {
	g.Meta           `mime:"application/json" example:"string"`
	ProjectStructure []FileInfo     `json:"projectStructure"` // 项目结构
	Variables        []VariableInfo `json:"variables"`        // 推荐变量
	Instructions     string         `json:"instructions"`     // 使用说明
	EstimatedTime    int            `json:"estimatedTime"`    // 预估完成时间(分钟)
}

// AI变量建议请求
type SuggestVariablesReq struct {
	g.Meta        `path:"/ai/suggestVariables" method:"post" tags:"AI服务" summary:"AI服务-建议变量"`
	ProjectType   string   `json:"projectType" v:"required#项目类型不能为空"`
	TechStack     []string `json:"techStack"`
	Description   string   `json:"description"`
}

type SuggestVariablesRes struct {
	g.Meta    `mime:"application/json" example:"string"`
	Variables []VariableInfo `json:"variables"`
}

// 文件信息
type FileInfo struct {
	Path        string `json:"path"`
	Content     string `json:"content"`
	IsDirectory bool   `json:"isDirectory"`
	Description string `json:"description"`
}

// AI统一聊天请求
type ChatReq struct {
	g.Meta      `path:"/ai/chat" method:"post" tags:"AI服务" summary:"AI服务-统一聊天接口"`
	Action      string                 `json:"action" v:"required#操作类型不能为空"`      // 操作类型：optimize_code, explain_code, suggest_variables, generate_template, refactor_code, add_comments, general_chat
	Context     map[string]interface{} `json:"context"`                              // 上下文信息
	UserInput   string                 `json:"userInput" v:"required#用户输入不能为空"`    // 用户输入
	Preferences map[string]interface{} `json:"preferences"`                          // 用户偏好设置
	ChatHistory []ChatMessage          `json:"chatHistory"`                          // 聊天历史
	Stream      bool                   `json:"stream"`                               // 是否使用流式响应
}

type ChatRes struct {
	g.Meta      `mime:"application/json" example:"string"`
	Content     string                 `json:"content"`     // AI回复的主要内容
	Suggestions []ChatSuggestion       `json:"suggestions"` // 建议项
	Metadata    map[string]interface{} `json:"metadata"`    // 元数据（模型信息、耗时等）
}

// 流式响应数据
type ChatStreamData struct {
	Type        string                 `json:"type"`        // chunk, metadata, suggestions, done
	Content     string                 `json:"content"`     // 内容片段
	Suggestions []ChatSuggestion       `json:"suggestions"` // 建议项（在最后发送）
	Metadata    map[string]interface{} `json:"metadata"`    // 元数据（在最后发送）
	Done        bool                   `json:"done"`        // 是否完成
}

// 聊天消息
type ChatMessage struct {
	Role      string `json:"role"`      // user, assistant
	Content   string `json:"content"`   // 消息内容
	Timestamp string `json:"timestamp"` // 时间戳
}

// 聊天建议
type ChatSuggestion struct {
	Type        string  `json:"type"`        // code, variable, file, action
	Name        string  `json:"name"`        // 建议名称
	Description string  `json:"description"` // 建议描述
	Code        string  `json:"code"`        // 建议的代码内容
	Confidence  float64 `json:"confidence"`  // 置信度
	Priority    string  `json:"priority"`    // high, medium, low
}

// 变量信息
type VariableInfo struct {
	Name         string   `json:"name"`
	Type         string   `json:"type"`         // string, number, boolean
	Description  string   `json:"description"`
	DefaultValue string   `json:"defaultValue"`
	Required     bool     `json:"required"`
	Options      []string `json:"options,omitempty"` // 选项值（如果是选择类型）
}