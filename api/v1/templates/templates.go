package template

import (
	commonApi "github.com/ciclebyte/template_starter/api/v1/common"
	model "github.com/ciclebyte/template_starter/internal/model"
	"github.com/gogf/gf/v2/frame/g"
)

type TemplateLanguageReq struct {
	LanguageId int `json:"languageId"`
	IsPrimary  int `json:"isPrimary"`
}

type TemplatesAddReq struct {
	g.Meta       `path:"/templates/add" method:"post" tags:"模板" summary:"模板-新增"`
	Name         string                `json:"name" v:"required#模板名称不能为空"`
	Description  string                `json:"description" v:"required#模板详细描述不能为空"`
	Introduction string                `json:"introduction"` // 模板详细介绍，支持Markdown格式
	CategoryId   int                   `json:"categoryId" v:"required#所属分类ID不能为空"`
	IsFeatured   int                   `json:"isFeatured" v:"required#是否推荐模板不能为空"`
	Logo         string                `json:"logo"`
	Icon         string                `json:"icon"`
	Languages    []TemplateLanguageReq `json:"languages"`
}

type TemplatesAddRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

type TemplatesDelReq struct {
	g.Meta `path:"/templates/del" method:"delete" tags:"模板" summary:"模板-删除"`
	Id     interface{} `json:"id" v:"required#id不能为空"`
}

type TemplatesDelRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

type TemplatesBatchDelReq struct {
	g.Meta `path:"/templates/batchdel" method:"delete" tags:"模板" summary:"模板-批量删除"`
	Ids    []interface{} `json:"id" v:"required#id不能为空"`
}

type TemplatesBatchDelRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

type TemplatesEditReq struct {
	g.Meta       `path:"/templates/edit" method:"put" tags:"模板" summary:"模板-修改"`
	Id           interface{}           `json:"id" v:"required#模板ID，自增主键不能为空"`
	Name         string                `json:"name" v:"required#模板名称不能为空"`
	Description  string                `json:"description" v:"required#模板详细描述不能为空"`
	Introduction string                `json:"introduction"` // 模板详细介绍，支持Markdown格式
	CategoryId   int                   `json:"categoryId" v:"required#所属分类ID不能为空"`
	IsFeatured   int                   `json:"isFeatured" v:"required#是否推荐模板不能为空"`
	Logo         string                `json:"logo"`
	Icon         string                `json:"icon"`
	Languages    []TemplateLanguageReq `json:"languages"`
}

type TemplatesEditRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

type TemplatesListReq struct {
	g.Meta `path:"/templates/list" method:"get" tags:"模板" summary:"模板-列表"`
	commonApi.PageReq
	Name        string `json:"name"`
	Description string `json:"description"`
	CategoryId  int    `json:"categoryId"`
	IsFeatured  int    `json:"isFeatured"`
	Logo        string `json:"logo"`
}

type TemplatesListRes struct {
	g.Meta `mime:"application/json" example:"string"`
	commonApi.ListRes
	TemplatesList []*model.TemplatesInfo `json:"templatesList"`
}

type TemplatesDetailReq struct {
	g.Meta `path:"/templates/detail" method:"get" tags:"模板" summary:"模板-详情"`
	Id     interface{} `json:"id" v:"required#id不能为空"`
}

type TemplatesDetailRes struct {
	g.Meta        `mime:"application/json" example:"string"`
	TemplatesInfo *model.TemplatesInfo `json:"data"`
}

// 变量统计信息
type VariableStatistics struct {
	TotalCustomVariables  int `json:"totalCustomVariables"`
	TotalBuiltinVariables int `json:"totalBuiltinVariables"`
	TotalFunctions        int `json:"totalFunctions"`
	TotalFiles            int `json:"totalFiles"`
}

// 内置变量信息
type BuiltinVariableInfo struct {
	Name        string   `json:"name"`
	Label       string   `json:"label"`
	Description string   `json:"description"`
	UsageCount  int      `json:"usageCount"`
	Files       []string `json:"files"`
}

// 模板函数信息
type TemplateFunctionInfo struct {
	Name        string   `json:"name"`
	Label       string   `json:"label"`
	Description string   `json:"description"`
	UsageCount  int      `json:"usageCount"`
	Files       []string `json:"files"`
}

// 模板变量列表请求
type TemplatesVariablesReq struct {
	g.Meta     `path:"/templates/{templateId}/variables" method:"get" tags:"模板" summary:"模板-变量列表"`
	TemplateId interface{} `json:"templateId" v:"required#模板ID不能为空"`
}

// 模板变量列表响应 - 已废弃，保留空结构体以兼容
type TemplatesVariablesRes struct {
	g.Meta            `mime:"application/json" example:"string"`
	CustomVariables   []interface{}           `json:"customVariables"`
	BuiltinVariables  []*BuiltinVariableInfo  `json:"builtinVariables"`
	TemplateFunctions []*TemplateFunctionInfo `json:"templateFunctions"`
	Statistics        *VariableStatistics     `json:"statistics"`
}

// 变量分析请求
type TemplatesAnalyzeVariablesReq struct {
	g.Meta     `path:"/templates/{templateId}/analyze-variables" method:"post" tags:"模板" summary:"模板-分析变量"`
	TemplateId interface{} `json:"templateId" v:"required#模板ID不能为空"`
}

// 检测到的变量信息
type DetectedVariable struct {
	Name        string   `json:"name"`        // 变量名
	Type        string   `json:"type"`        // 推测的类型
	Files       []string `json:"files"`       // 出现的文件
	Contexts    []string `json:"contexts"`    // 使用上下文
	Suggestions string   `json:"suggestions"` // 类型建议说明
}

// 变量分析响应
type TemplatesAnalyzeVariablesRes struct {
	g.Meta              `mime:"application/json" example:"string"`
	DetectedVariables   []*DetectedVariable `json:"detectedVariables"`   // 检测到的变量
	MissingVariables    []*DetectedVariable `json:"missingVariables"`    // 缺失的变量（模板中使用但未定义）
	UnusedVariables     []string            `json:"unusedVariables"`     // 未使用的变量（已定义但模板中未使用）
	ConflictVariables   []*DetectedVariable `json:"conflictVariables"`   // 冲突的变量（类型不匹配）
	TotalVariableCount  int                 `json:"totalVariableCount"`  // 总变量数
	AnalyzedFileCount   int                 `json:"analyzedFileCount"`   // 分析的文件数
}
