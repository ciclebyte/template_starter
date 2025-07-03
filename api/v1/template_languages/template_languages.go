package template

import (
	commonApi "github.com/ciclebyte/template_starter/api/v1/common"
	model "github.com/ciclebyte/template_starter/internal/model"
	"github.com/gogf/gf/v2/frame/g"
)

type TemplateLanguagesAddReq struct {
	g.Meta     `path:"/templateLanguages/add" method:"post" tags:"模板语言" summary:"模板语言-新增"`
	TemplateId interface{} `json:"templateId" v:"required#关联的模板ID不能为空"`
	LanguageId interface{} `json:"languageId" v:"required#关联的语言ID不能为空"`
	IsPrimary  int         `json:"isPrimary" v:"required#是否主要语言不能为空"`
}

type TemplateLanguagesAddRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

type TemplateLanguagesDelReq struct {
	g.Meta `path:"/templateLanguages/del" method:"delete" tags:"模板语言" summary:"模板语言-删除"`
	Id     interface{} `json:"id" v:"required#id不能为空"`
}

type TemplateLanguagesDelRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

type TemplateLanguagesBatchDelReq struct {
	g.Meta `path:"/templateLanguages/batchdel" method:"delete" tags:"模板语言" summary:"模板语言-批量删除"`
	Ids    []interface{} `json:"id" v:"required#id不能为空"`
}

type TemplateLanguagesBatchDelRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

type TemplateLanguagesEditReq struct {
	g.Meta     `path:"/templateLanguages/edit" method:"put" tags:"模板语言" summary:"模板语言-修改"`
	Id         interface{} `json:"id" v:"required#关联ID，自增主键不能为空"`
	TemplateId interface{} `json:"templateId" v:"required#关联的模板ID不能为空"`
	LanguageId int         `json:"languageId" v:"required#关联的语言ID不能为空"`
	IsPrimary  int         `json:"isPrimary" v:"required#是否主要语言不能为空"`
}

type TemplateLanguagesEditRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

type TemplateLanguagesListReq struct {
	g.Meta `path:"/templateLanguages/list" method:"get" tags:"模板语言" summary:"模板语言-列表"`
	commonApi.PageReq
	TemplateId interface{} `json:"templateId"`
	LanguageId int         `json:"languageId"`
	IsPrimary  int         `json:"isPrimary"`
}

type TemplateLanguagesListRes struct {
	g.Meta `mime:"application/json" example:"string"`
	commonApi.ListRes
	TemplateLanguagesList []*model.TemplateLanguagesInfo `json:"templateLanguagesList"`
}

type TemplateLanguagesDetailReq struct {
	g.Meta `path:"/templateLanguages/detail" method:"get" tags:"模板语言" summary:"模板语言-详情"`
	Id     interface{} `json:"id" v:"required#id不能为空"`
}

type TemplateLanguagesDetailRes struct {
	g.Meta `mime:"application/json" example:"string"`
	*model.TemplateLanguagesInfo
}
