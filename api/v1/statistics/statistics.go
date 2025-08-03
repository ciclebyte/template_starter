package statistics

import (
	"github.com/gogf/gf/v2/frame/g"
)

// 统计总览请求参数
type OverviewReq struct {
	g.Meta `path:"/statistics/overview" method:"get" tags:"统计分析" summary:"统计分析-总览数据"`
}

// 统计总览响应数据
type OverviewRes struct {
	g.Meta `mime:"application/json" example:"string"`

	// 基础统计
	TotalTemplates  int `json:"totalTemplates"`  // 总模板数
	TotalCategories int `json:"totalCategories"` // 总分类数
	TotalLanguages  int `json:"totalLanguages"`  // 总语言数
	TotalFiles      int `json:"totalFiles"`      // 总文件数
	TotalVariables  int `json:"totalVariables"`  // 总变量数

	// 质量指标
	FeaturedTemplates        int `json:"featuredTemplates"`        // 推荐模板数
	TemplatesWithVariables   int `json:"templatesWithVariables"`   // 包含变量的模板数
	TemplatesWithDescription int `json:"templatesWithDescription"` // 包含描述的模板数
	AvgFilesPerTemplate      int `json:"avgFilesPerTemplate"`      // 平均每个模板的文件数
}

// 分类分布请求参数
type CategoryDistributionReq struct {
	g.Meta `path:"/statistics/category-distribution" method:"get" tags:"统计分析" summary:"统计分析-分类分布"`
}

// 分类分布数据项
type CategoryDistributionItem struct {
	CategoryName  string `json:"categoryName"`  // 分类名称
	TemplateCount int    `json:"templateCount"` // 模板数量
	Percentage    int    `json:"percentage"`    // 占比百分比
}

// 分类分布响应数据
type CategoryDistributionRes struct {
	g.Meta `mime:"application/json" example:"string"`
	Items  []*CategoryDistributionItem `json:"items"`
}

// 语言流行度请求参数
type LanguagePopularityReq struct {
	g.Meta `path:"/statistics/language-popularity" method:"get" tags:"统计分析" summary:"统计分析-语言流行度"`
}

// 语言流行度数据项
type LanguagePopularityItem struct {
	LanguageName  string `json:"languageName"`  // 语言名称
	TemplateCount int    `json:"templateCount"` // 使用该语言的模板数
	Percentage    int    `json:"percentage"`    // 占比百分比
}

// 语言流行度响应数据
type LanguagePopularityRes struct {
	g.Meta `mime:"application/json" example:"string"`
	Items  []*LanguagePopularityItem `json:"items"`
}

// 模板复杂度请求参数
type TemplateComplexityReq struct {
	g.Meta `path:"/statistics/template-complexity" method:"get" tags:"统计分析" summary:"统计分析-模板复杂度"`
}

// 模板复杂度数据
type TemplateComplexityRes struct {
	g.Meta `mime:"application/json" example:"string"`

	// 按文件数分组
	SimpleTemplates  int `json:"simpleTemplates"`  // 简单模板(1-3个文件)
	MediumTemplates  int `json:"mediumTemplates"`  // 中等模板(4-10个文件)
	ComplexTemplates int `json:"complexTemplates"` // 复杂模板(10+个文件)

	// 按变量数分组
	NoVariableTemplates   int `json:"noVariableTemplates"`   // 无变量模板
	FewVariableTemplates  int `json:"fewVariableTemplates"`  // 少量变量模板(1-5个变量)
	ManyVariableTemplates int `json:"manyVariableTemplates"` // 多变量模板(5+个变量)
}

// 使用趋势请求参数  
type UsageTrendsReq struct {
	g.Meta `path:"/statistics/usage-trends" method:"get" tags:"统计分析" summary:"统计分析-使用趋势"`
	Days   int `json:"days"` // 统计天数，默认30天
}

// 使用趋势数据项
type UsageTrendItem struct {
	Date            string `json:"date"`            // 日期 YYYY-MM-DD
	TemplateCreated int    `json:"templateCreated"` // 当日创建的模板数
}

// 使用趋势响应数据
type UsageTrendsRes struct {
	g.Meta `mime:"application/json" example:"string"`
	Items  []*UsageTrendItem `json:"items"`
}