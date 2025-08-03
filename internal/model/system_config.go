package model

import (
	"github.com/ciclebyte/template_starter/internal/model/entity"
	"github.com/gogf/gf/v2/os/gtime"
)

// SystemConfigInfo 系统配置信息
type SystemConfigInfo struct {
	entity.SystemConfig
}

// SystemConfigAdd 新增系统配置
type SystemConfigAdd struct {
	ConfigKey      string `json:"configKey" v:"required#配置键名不能为空"`
	ConfigValue    string `json:"configValue"`
	ConfigGroup    string `json:"configGroup" v:"required#配置分组不能为空"`
	ConfigType     string `json:"configType" v:"required|in:string,number,boolean,json,array#配置类型不能为空|配置类型必须为string,number,boolean,json,array之一"`
	DisplayName    string `json:"displayName" v:"required#显示名称不能为空"`
	Description    string `json:"description"`
	IsPublic       int    `json:"isPublic" v:"in:0,1#是否公开必须为0或1"`
	IsRequired     int    `json:"isRequired" v:"in:0,1#是否必填必须为0或1"`
	DefaultValue   string `json:"defaultValue"`
	ValidationRule string `json:"validationRule"`
	SortOrder      int    `json:"sortOrder"`
	Status         int    `json:"status" v:"in:0,1#状态必须为0或1"`
}

// SystemConfigEdit 编辑系统配置
type SystemConfigEdit struct {
	Id             interface{} `json:"id" v:"required#配置ID不能为空"`
	ConfigValue    string      `json:"configValue"`
	ConfigGroup    string      `json:"configGroup" v:"required#配置分组不能为空"`
	ConfigType     string      `json:"configType" v:"required|in:string,number,boolean,json,array#配置类型不能为空|配置类型必须为string,number,boolean,json,array之一"`
	DisplayName    string      `json:"displayName" v:"required#显示名称不能为空"`
	Description    string      `json:"description"`
	IsPublic       int         `json:"isPublic" v:"in:0,1#是否公开必须为0或1"`
	IsRequired     int         `json:"isRequired" v:"in:0,1#是否必填必须为0或1"`
	DefaultValue   string      `json:"defaultValue"`
	ValidationRule string      `json:"validationRule"`
	SortOrder      int         `json:"sortOrder"`
	Status         int         `json:"status" v:"in:0,1#状态必须为0或1"`
}

// ConfigGroup 配置分组信息
type ConfigGroup struct {
	Group       string `json:"group"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Count       int    `json:"count"`
}

// AIConfig AI相关配置结构
type AIConfig struct {
	// 基础配置
	Enabled  bool   `json:"enabled"`
	Provider string `json:"provider"` // openai, claude, custom
	
	// OpenAI配置
	OpenAI struct {
		APIKey      string  `json:"apiKey"`
		BaseURL     string  `json:"baseUrl"`
		Model       string  `json:"model"`
		MaxTokens   int     `json:"maxTokens"`
		Temperature float64 `json:"temperature"`
	} `json:"openai"`
	
	// Claude配置
	Claude struct {
		APIKey  string `json:"apiKey"`
		BaseURL string `json:"baseUrl"`
		Model   string `json:"model"`
	} `json:"claude"`
	
	// 功能开关
	Features struct {
		TemplateGeneration  bool `json:"templateGeneration"`
		CodeReview         bool `json:"codeReview"`
		VariableSuggestion bool `json:"variableSuggestion"`
		AutoDocumentation  bool `json:"autoDocumentation"`
	} `json:"features"`
	
	// 使用限制
	RateLimit struct {
		RequestsPerHour int `json:"requestsPerHour"`
		RequestsPerDay  int `json:"requestsPerDay"`
	} `json:"rateLimit"`
	
	// 支持的语言列表
	SupportedLanguages []string `json:"supportedLanguages"`
	MaxFileSize       int      `json:"maxFileSize"` // KB
}

// SystemGeneralConfig 系统基础配置
type SystemGeneralConfig struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	LogoURL     string `json:"logoUrl"`
}

// TemplateConfig 模板相关配置
type TemplateConfig struct {
	DefaultAuthor         string `json:"defaultAuthor"`
	DefaultCompany        string `json:"defaultCompany"`
	MaxFileSize          int    `json:"maxFileSize"`          // MB
	MaxFilesPerTemplate  int    `json:"maxFilesPerTemplate"`
}

// UIConfig 界面相关配置
type UIConfig struct {
	Theme          string `json:"theme"`          // light, dark
	Language       string `json:"language"`       // zh-CN, en-US
	PageSize       int    `json:"pageSize"`
	EditorTheme    string `json:"editorTheme"`
	EditorFontSize int    `json:"editorFontSize"`
}

// ConfigValidationRule 配置验证规则
type ConfigValidationRule struct {
	Required bool        `json:"required"`
	MinValue interface{} `json:"minValue,omitempty"`
	MaxValue interface{} `json:"maxValue,omitempty"`
	Pattern  string      `json:"pattern,omitempty"`
	Options  []string    `json:"options,omitempty"`
}

// ConfigHistory 配置变更历史
type ConfigHistory struct {
	Id          int64       `json:"id"`
	ConfigKey   string      `json:"configKey"`
	OldValue    string      `json:"oldValue"`
	NewValue    string      `json:"newValue"`
	ChangeType  string      `json:"changeType"` // create, update, delete
	ChangedBy   string      `json:"changedBy"`  // 操作人
	ChangeReason string     `json:"changeReason"`
	CreatedAt   *gtime.Time `json:"createdAt"`
}