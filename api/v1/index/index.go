package index

import (
	model "github.com/ciclebyte/template_starter/internal/model"
	"github.com/gogf/gf/v2/frame/g"
)

// 首页请求参数
type IndexReq struct {
	g.Meta        `path:"/index" method:"get" tags:"首页" summary:"首页-数据"`
	CategoryLimit int `json:"categoryLimit"` // 每个分类返回的模板数量限制，默认6个
	FeaturedLimit int `json:"featuredLimit"` // 推荐模板数量限制，默认8个
}

// 首页统计数据
type IndexStatistics struct {
	TotalTemplates  int `json:"totalTemplates"`  // 总模板数
	TotalCategories int `json:"totalCategories"` // 总分类数
	TotalLanguages  int `json:"totalLanguages"`  // 总语言数
	FeaturedCount   int `json:"featuredCount"`   // 推荐模板数
}

// 分类及其模板
type CategoryWithTemplates struct {
	*model.CategoriesInfo                        // 继承分类信息
	Templates             []*model.TemplatesInfo `json:"templates"`     // 该分类下的模板列表
	TemplateCount         int                    `json:"templateCount"` // 该分类下的模板总数
}

// 首页响应数据
type IndexRes struct {
	g.Meta `mime:"application/json" example:"string"`

	// 统计数据
	Statistics *IndexStatistics `json:"statistics"`

	// 分类及其模板
	Categories []*CategoryWithTemplates `json:"categories"`

	// 推荐模板
	FeaturedTemplates []*model.TemplatesInfo `json:"featuredTemplates"`
}
