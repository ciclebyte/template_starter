package template

import (
	commonApi "github.com/ciclebyte/template_starter/api/v1/common"
	model "github.com/ciclebyte/template_starter/internal/model"
	"github.com/gogf/gf/v2/frame/g"
)

type LanguagesAddReq struct {
	g.Meta      `path:"/languages/add" method:"post" tags:"语言" summary:"语言-新增"`
	Name        string `json:"name" v:"required#语言名称（如JavaScript、Python）不能为空"`
	DisplayName string `json:"displayName" v:"required#语言显示名称不能为空"`
	Code        string `json:"code" v:"required#语言代码（如js、py）不能为空"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	IsPopular   int    `json:"isPopular"`
}

type LanguagesAddRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

type LanguagesDelReq struct {
	g.Meta `path:"/languages/del" method:"delete" tags:"语言" summary:"语言-删除"`
	Id     int `json:"id" v:"required#id不能为空"`
}

type LanguagesDelRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

type LanguagesBatchDelReq struct {
	g.Meta `path:"/languages/batchdel" method:"delete" tags:"语言" summary:"语言-批量删除"`
	Ids    []int `json:"id" v:"required#id不能为空"`
}

type LanguagesBatchDelRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

type LanguagesEditReq struct {
	g.Meta      `path:"/languages/edit" method:"put" tags:"语言" summary:"语言-修改"`
	Id          int    `json:"id" v:"required#语言ID，自增主键不能为空"`
	Name        string `json:"name" v:"required#语言名称（如JavaScript、Python）不能为空"`
	DisplayName string `json:"displayName" v:"required#语言显示名称不能为空"`
	Code        string `json:"code" v:"required#语言代码（如js、py）不能为空"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	IsPopular   int    `json:"isPopular"`
}

type LanguagesEditRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

type LanguagesListReq struct {
	g.Meta `path:"/languages/list" method:"get" tags:"语言" summary:"语言-列表"`
	commonApi.PageReq
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Code        string `json:"code"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	IsPopular   int    `json:"isPopular"`
}

type LanguagesListRes struct {
	g.Meta `mime:"application/json" example:"string"`
	commonApi.ListRes
	LanguagesList []*model.LanguagesInfo `json:"languagesList"`
}

type LanguagesDetailReq struct {
	g.Meta `path:"/languages/detail" method:"get" tags:"语言" summary:"语言-详情"`
	Id     int `json:"id" v:"required#id不能为空"`
}

type LanguagesDetailRes struct {
	g.Meta `mime:"application/json" example:"string"`
	*model.LanguagesInfo
}
