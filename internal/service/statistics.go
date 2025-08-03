package service

import (
	"context"

	statisticsApi "github.com/ciclebyte/template_starter/api/v1/statistics"
)

type IStatisticsService interface {
	GetOverview(ctx context.Context, req *statisticsApi.OverviewReq) (*statisticsApi.OverviewRes, error)
	GetCategoryDistribution(ctx context.Context, req *statisticsApi.CategoryDistributionReq) (*statisticsApi.CategoryDistributionRes, error)
	GetLanguagePopularity(ctx context.Context, req *statisticsApi.LanguagePopularityReq) (*statisticsApi.LanguagePopularityRes, error)
	GetTemplateComplexity(ctx context.Context, req *statisticsApi.TemplateComplexityReq) (*statisticsApi.TemplateComplexityRes, error)
	GetUsageTrends(ctx context.Context, req *statisticsApi.UsageTrendsReq) (*statisticsApi.UsageTrendsRes, error)
}

var localStatisticsService IStatisticsService

func Statistics() IStatisticsService {
	if localStatisticsService == nil {
		panic("implement not found for interface IStatisticsService, forgot register?")
	}
	return localStatisticsService
}

func RegisterStatisticsService(i IStatisticsService) {
	localStatisticsService = i
}