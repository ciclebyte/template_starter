package template

import (
	commonApi "github.com/ciclebyte/template_starter/api/v1/common"
	model "github.com/ciclebyte/template_starter/internal/model"
	"github.com/gogf/gf/v2/frame/g"
)

//int64

type TemplateVariablesAddReq struct {
	g.Meta          `path:"/templateVariables/add" method:"post" tags:"模板变量" summary:"模板变量-新增"`
	TemplateId      int64  `json:"templateId" v:"required#所属模板ID不能为空"`
	Name            string `json:"name" v:"required#变量名称不能为空"`
	Description     string `json:"description" v:"required#变量描述不能为空"`
	DefaultValue    string `json:"defaultValue" v:"required#变量默认值不能为空"`
	IsRequired      int    `json:"isRequired" v:"required#是否为必填变量不能为空"`
	ValidationRegex string `json:"validationRegex" v:"required#变量值验证正则表达式不能为空"`
	Sort            int    `json:"sort" v:"required#显示顺序，数字越小排序越靠前不能为空"`
}

type TemplateVariablesAddRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

type TemplateVariablesDelReq struct {
	g.Meta `path:"/templateVariables/del" method:"delete" tags:"模板变量" summary:"模板变量-删除"`
	Id     interface{} `json:"id" v:"required#id不能为空"`
}

type TemplateVariablesDelRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

type TemplateVariablesBatchDelReq struct {
	g.Meta `path:"/templateVariables/batchdel" method:"delete" tags:"模板变量" summary:"模板变量-批量删除"`
	Ids    []interface{} `json:"id" v:"required#id不能为空"`
}

type TemplateVariablesBatchDelRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

type TemplateVariablesEditReq struct {
	g.Meta          `path:"/templateVariables/edit" method:"put" tags:"模板变量" summary:"模板变量-修改"`
	Id              interface{} `json:"id" v:"required#变量ID，自增主键不能为空"`
	TemplateId      interface{} `json:"templateId" v:"required#所属模板ID不能为空"`
	Name            string      `json:"name" v:"required#变量名称不能为空"`
	Description     string      `json:"description" v:"required#变量描述不能为空"`
	DefaultValue    string      `json:"defaultValue" v:"required#变量默认值不能为空"`
	IsRequired      int         `json:"isRequired" v:"required#是否为必填变量不能为空"`
	ValidationRegex string      `json:"validationRegex" v:"required#变量值验证正则表达式不能为空"`
	Sort            int         `json:"sort" v:"required#显示顺序"`
}

type TemplateVariablesEditRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

type TemplateVariablesListReq struct {
	g.Meta `path:"/templateVariables/list" method:"get" tags:"模板变量" summary:"模板变量-列表"`
	commonApi.PageReq
	TemplateId      interface{} `json:"templateId"`
	Name            string      `json:"name"`
	Description     string      `json:"description"`
	DefaultValue    string      `json:"defaultValue"`
	IsRequired      int         `json:"isRequired"`
	ValidationRegex string      `json:"validationRegex"`
	Sort            int         `json:"sort"`
}

type TemplateVariablesListRes struct {
	g.Meta `mime:"application/json" example:"string"`
	commonApi.ListRes
	TemplateVariablesList []*model.TemplateVariablesInfo `json:"templateVariablesList"`
}

type TemplateVariablesDetailReq struct {
	g.Meta `path:"/templateVariables/detail" method:"get" tags:"模板变量" summary:"模板变量-详情"`
	Id     uint `json:"id" v:"required#id不能为空"`
}

type TemplateVariablesDetailRes struct {
	g.Meta `mime:"application/json" example:"string"`
	*model.TemplateVariablesInfo
}
