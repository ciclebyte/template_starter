package controller

import (
	"context"

	statisticsApi "github.com/ciclebyte/template_starter/api/v1/statistics"
	"github.com/ciclebyte/template_starter/internal/service"
	"github.com/gogf/gf/v2/frame/g"
)

var Statistics = cStatistics{}

type cStatistics struct {
	BaseController
}

// GetOverview 获取统计总览数据
func (c *cStatistics) GetOverview(ctx context.Context, req *statisticsApi.OverviewReq) (res *statisticsApi.OverviewRes, err error) {
	g.Log().Debug(ctx, "Statistics.GetOverview called")
	return service.Statistics().GetOverview(ctx, req)
}

// GetCategoryDistribution 获取分类分布数据
func (c *cStatistics) GetCategoryDistribution(ctx context.Context, req *statisticsApi.CategoryDistributionReq) (res *statisticsApi.CategoryDistributionRes, err error) {
	g.Log().Debug(ctx, "Statistics.GetCategoryDistribution called")
	return service.Statistics().GetCategoryDistribution(ctx, req)
}

// GetLanguagePopularity 获取语言流行度数据
func (c *cStatistics) GetLanguagePopularity(ctx context.Context, req *statisticsApi.LanguagePopularityReq) (res *statisticsApi.LanguagePopularityRes, err error) {
	g.Log().Debug(ctx, "Statistics.GetLanguagePopularity called")
	return service.Statistics().GetLanguagePopularity(ctx, req)
}

// GetTemplateComplexity 获取模板复杂度数据
func (c *cStatistics) GetTemplateComplexity(ctx context.Context, req *statisticsApi.TemplateComplexityReq) (res *statisticsApi.TemplateComplexityRes, err error) {
	g.Log().Debug(ctx, "Statistics.GetTemplateComplexity called")
	return service.Statistics().GetTemplateComplexity(ctx, req)
}

// GetUsageTrends 获取使用趋势数据
func (c *cStatistics) GetUsageTrends(ctx context.Context, req *statisticsApi.UsageTrendsReq) (res *statisticsApi.UsageTrendsRes, err error) {
	g.Log().Debug(ctx, "Statistics.GetUsageTrends called")
	return service.Statistics().GetUsageTrends(ctx, req)
}