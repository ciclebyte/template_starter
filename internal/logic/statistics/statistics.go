package statistics

import (
	"context"

	statisticsApi "github.com/ciclebyte/template_starter/api/v1/statistics"
	"github.com/ciclebyte/template_starter/internal/dao"
	"github.com/ciclebyte/template_starter/internal/service"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sStatistics struct{}

func init() {
	service.RegisterStatisticsService(New())
}

func New() *sStatistics {
	return &sStatistics{}
}

// GetOverview 获取统计总览数据
func (s *sStatistics) GetOverview(ctx context.Context, req *statisticsApi.OverviewReq) (*statisticsApi.OverviewRes, error) {
	g.Log().Debug(ctx, "GetOverview called")

	res := &statisticsApi.OverviewRes{}

	// 并行获取各项统计数据
	var err error

	// 基础统计
	res.TotalTemplates, err = dao.Templates.Ctx(ctx).Count()
	if err != nil {
		return nil, err
	}

	res.TotalCategories, err = dao.Categories.Ctx(ctx).Count()
	if err != nil {
		return nil, err
	}

	res.TotalLanguages, err = dao.Languages.Ctx(ctx).Count()
	if err != nil {
		return nil, err
	}

	res.TotalFiles, err = dao.TemplateFiles.Ctx(ctx).Count()
	if err != nil {
		return nil, err
	}

	// 质量指标
	res.FeaturedTemplates, err = dao.Templates.Ctx(ctx).Where("is_featured", 1).Count()
	if err != nil {
		return nil, err
	}

	// 包含描述的模板数
	res.TemplatesWithDescription, err = dao.Templates.Ctx(ctx).
		Where("introduction IS NOT NULL AND introduction != ''").
		Count()
	if err != nil {
		return nil, err
	}

	// 计算平均每个模板的文件数
	if res.TotalTemplates > 0 {
		res.AvgFilesPerTemplate = res.TotalFiles / res.TotalTemplates
	}

	return res, nil
}

// GetCategoryDistribution 获取分类分布数据
func (s *sStatistics) GetCategoryDistribution(ctx context.Context, req *statisticsApi.CategoryDistributionReq) (*statisticsApi.CategoryDistributionRes, error) {
	g.Log().Debug(ctx, "GetCategoryDistribution called")

	res := &statisticsApi.CategoryDistributionRes{}

	// 获取分类及其模板数量
	type CategoryCount struct {
		CategoryName  string `json:"category_name"`
		TemplateCount int    `json:"template_count"`
	}

	var categoryCountData []CategoryCount
	err := dao.Categories.Ctx(ctx).
		LeftJoin("templates t", "categories.id = t.category_id").
		Fields("categories.name as category_name, COUNT(t.id) as template_count").
		Group("categories.id, categories.name").
		Order("template_count DESC").
		Scan(&categoryCountData)
	if err != nil {
		return nil, err
	}

	// 计算总模板数用于百分比计算
	totalTemplates, err := dao.Templates.Ctx(ctx).Count()
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	var items []*statisticsApi.CategoryDistributionItem
	for _, data := range categoryCountData {
		percentage := 0
		if totalTemplates > 0 {
			percentage = (data.TemplateCount * 100) / totalTemplates
		}

		items = append(items, &statisticsApi.CategoryDistributionItem{
			CategoryName:  data.CategoryName,
			TemplateCount: data.TemplateCount,
			Percentage:    percentage,
		})
	}

	res.Items = items
	return res, nil
}

// GetLanguagePopularity 获取语言流行度数据
func (s *sStatistics) GetLanguagePopularity(ctx context.Context, req *statisticsApi.LanguagePopularityReq) (*statisticsApi.LanguagePopularityRes, error) {
	g.Log().Debug(ctx, "GetLanguagePopularity called")

	res := &statisticsApi.LanguagePopularityRes{}

	// 获取语言及其关联的模板数量
	type LanguageCount struct {
		LanguageName  string `json:"language_name"`
		TemplateCount int    `json:"template_count"`
	}

	var languageCountData []LanguageCount
	err := dao.Languages.Ctx(ctx).
		LeftJoin("template_languages tl", "languages.id = tl.language_id").
		Fields("languages.name as language_name, COUNT(tl.template_id) as template_count").
		Group("languages.id, languages.name").
		Order("template_count DESC").
		Scan(&languageCountData)
	if err != nil {
		return nil, err
	}

	// 计算总模板数用于百分比计算
	totalTemplates, err := dao.Templates.Ctx(ctx).Count()
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	var items []*statisticsApi.LanguagePopularityItem
	for _, data := range languageCountData {
		percentage := 0
		if totalTemplates > 0 {
			percentage = (data.TemplateCount * 100) / totalTemplates
		}

		items = append(items, &statisticsApi.LanguagePopularityItem{
			LanguageName:  data.LanguageName,
			TemplateCount: data.TemplateCount,
			Percentage:    percentage,
		})
	}

	res.Items = items
	return res, nil
}

// GetTemplateComplexity 获取模板复杂度数据
func (s *sStatistics) GetTemplateComplexity(ctx context.Context, req *statisticsApi.TemplateComplexityReq) (*statisticsApi.TemplateComplexityRes, error) {
	g.Log().Debug(ctx, "GetTemplateComplexity called")

	res := &statisticsApi.TemplateComplexityRes{}

	// 按文件数分组统计
	type FileCountGroup struct {
		FileCount     int `json:"file_count"`
		TemplateCount int `json:"template_count"`
	}

	var fileCountData []FileCountGroup
	err := dao.TemplateFiles.Ctx(ctx).
		Fields("template_id, COUNT(*) as file_count").
		Group("template_id").
		Scan(&fileCountData)
	if err != nil {
		return nil, err
	}

	// 统计不同复杂度的模板数
	for _, data := range fileCountData {
		if data.FileCount >= 1 && data.FileCount <= 3 {
			res.SimpleTemplates++
		} else if data.FileCount >= 4 && data.FileCount <= 10 {
			res.MediumTemplates++
		} else if data.FileCount > 10 {
			res.ComplexTemplates++
		}
	}

	// 按变量数分组统计
	type VariableCountGroup struct {
		VariableCount int `json:"variable_count"`
		TemplateCount int `json:"template_count"`
	}

	var variableCountData []VariableCountGroup

	// 获取无变量的模板数
	totalTemplates, err := dao.Templates.Ctx(ctx).Count()
	if err != nil {
		return nil, err
	}

	templatesWithVariables := len(variableCountData)
	res.NoVariableTemplates = totalTemplates - templatesWithVariables

	// 统计不同变量数的模板
	for _, data := range variableCountData {
		if data.VariableCount >= 1 && data.VariableCount <= 5 {
			res.FewVariableTemplates++
		} else if data.VariableCount > 5 {
			res.ManyVariableTemplates++
		}
	}

	return res, nil
}

// GetUsageTrends 获取使用趋势数据
func (s *sStatistics) GetUsageTrends(ctx context.Context, req *statisticsApi.UsageTrendsReq) (*statisticsApi.UsageTrendsRes, error) {
	g.Log().Debug(ctx, "GetUsageTrends called with days:", req.Days)

	res := &statisticsApi.UsageTrendsRes{}

	// 默认统计30天
	days := req.Days
	if days <= 0 {
		days = 30
	}

	// 计算起始日期
	endDate := gtime.Now()
	startDate := endDate.AddDate(0, 0, -days)

	// 获取指定时间范围内的模板创建数据
	type DailyCount struct {
		Date  string `json:"date"`
		Count int    `json:"count"`
	}

	var dailyCountData []DailyCount
	err := dao.Templates.Ctx(ctx).
		Where("created_at >= ? AND created_at <= ?", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Fields("DATE(created_at) as date, COUNT(*) as count").
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&dailyCountData)
	if err != nil {
		return nil, err
	}

	// 创建完整的日期序列，填补缺失的日期
	dateMap := make(map[string]int)
	for _, data := range dailyCountData {
		dateMap[data.Date] = data.Count
	}

	// 生成完整的日期序列
	var items []*statisticsApi.UsageTrendItem
	current := startDate
	for current.Before(endDate) || current.Equal(endDate) {
		dateStr := current.Format("2006-01-02")
		count := dateMap[dateStr] // 如果没有数据，默认为0

		items = append(items, &statisticsApi.UsageTrendItem{
			Date:            dateStr,
			TemplateCreated: count,
		})

		current = current.AddDate(0, 0, 1)
	}

	res.Items = items
	return res, nil
}
