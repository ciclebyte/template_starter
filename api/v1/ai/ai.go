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

// 变量信息
type VariableInfo struct {
	Name         string   `json:"name"`
	Type         string   `json:"type"`         // string, number, boolean
	Description  string   `json:"description"`
	DefaultValue string   `json:"defaultValue"`
	Required     bool     `json:"required"`
	Options      []string `json:"options,omitempty"` // 选项值（如果是选择类型）
}