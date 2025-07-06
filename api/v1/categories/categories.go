package template

import (
	commonApi "github.com/ciclebyte/template_starter/api/v1/common"
	model "github.com/ciclebyte/template_starter/internal/model"
	"github.com/gogf/gf/v2/frame/g"
)

type CategoriesAddReq struct {
	g.Meta      `path:"/categories/add" method:"post" tags:"分类" summary:"分类-新增"`
	Name        string `json:"name" v:"required#分类名称，唯一不能为空"`
	Description string `json:"description" v:"required#分类描述不能为空"`
	Icon        string `json:"icon"`
	Sort        int    `json:"sort" v:"required#数字越大越靠前不能为空"`
}

type CategoriesAddRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

type CategoriesDelReq struct {
	g.Meta `path:"/categories/del" method:"delete" tags:"分类" summary:"分类-删除"`
	Id     int `json:"id" v:"required#id不能为空"`
}

type CategoriesDelRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

type CategoriesBatchDelReq struct {
	g.Meta `path:"/categories/batchdel" method:"delete" tags:"分类" summary:"分类-批量删除"`
	Ids    []int `json:"id" v:"required#id不能为空"`
}

type CategoriesBatchDelRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

type CategoriesEditReq struct {
	g.Meta      `path:"/categories/edit" method:"put" tags:"分类" summary:"分类-修改"`
	Id          int    `json:"id" v:"required#分类ID，自增主键不能为空"`
	Name        string `json:"name" v:"required#分类名称，唯一不能为空"`
	Description string `json:"description" v:"required#分类描述不能为空"`
	Icon        string `json:"icon"`
	Sort        int    `json:"sort" v:"required#数字越大越靠前不能为空"`
}

type CategoriesEditRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

type CategoriesListReq struct {
	g.Meta `path:"/categories/list" method:"get" tags:"分类" summary:"分类-列表"`
	commonApi.PageReq
	Name string `json:"name"`
}

type CategoriesListRes struct {
	g.Meta `mime:"application/json" example:"string"`
	commonApi.ListRes
	CategoriesList []*model.CategoriesInfo `json:"categoriesList"`
}

type CategoriesDetailReq struct {
	g.Meta `path:"/categories/detail" method:"get" tags:"分类" summary:"分类-详情"`
	Id     int `json:"id" v:"required#id不能为空"`
}

type CategoriesDetailRes struct {
	g.Meta `mime:"application/json" example:"string"`
	*model.CategoriesInfo
}
