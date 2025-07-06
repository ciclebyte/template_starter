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
	Logo         string                `json:"logo" v:"required#模板logo图片URL不能为空"`
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

// 模板变量列表响应
type TemplatesVariablesRes struct {
	g.Meta            `mime:"application/json" example:"string"`
	CustomVariables   []*model.TemplateVariablesInfo `json:"customVariables"`
	BuiltinVariables  []*BuiltinVariableInfo         `json:"builtinVariables"`
	TemplateFunctions []*TemplateFunctionInfo        `json:"templateFunctions"`
	Statistics        *VariableStatistics            `json:"statistics"`
}
