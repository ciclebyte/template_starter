package index

import (
	"context"

	indexApi "github.com/ciclebyte/template_starter/api/v1/index"
	"github.com/ciclebyte/template_starter/internal/dao"
	"github.com/ciclebyte/template_starter/internal/model"
	"github.com/ciclebyte/template_starter/internal/service"
	"github.com/gogf/gf/v2/frame/g"
)

type sIndex struct{}

func init() {
	service.RegisterIndexService(New())
}

func New() *sIndex {
	return &sIndex{}
}

// GetIndexData 获取首页数据
func (s *sIndex) GetIndexData(ctx context.Context, req *indexApi.IndexReq) (*indexApi.IndexRes, error) {
	g.Log().Debug(ctx, "GetIndexData called with req:", req)

	res := &indexApi.IndexRes{}

	// 1. 获取统计数据
	statistics, err := s.getStatistics(ctx)
	if err != nil {
		return nil, err
	}
	res.Statistics = statistics

	// 2. 获取分类及其模板
	categories, err := s.getCategoriesWithTemplates(ctx, req.CategoryLimit)
	if err != nil {
		return nil, err
	}
	res.Categories = categories

	// 3. 获取推荐模板
	featuredTemplates, err := s.getFeaturedTemplates(ctx, req.FeaturedLimit)
	if err != nil {
		return nil, err
	}
	res.FeaturedTemplates = featuredTemplates

	return res, nil
}

// getStatistics 获取统计数据
func (s *sIndex) getStatistics(ctx context.Context) (*indexApi.IndexStatistics, error) {
	statistics := &indexApi.IndexStatistics{}

	// 获取总模板数
	totalTemplates, err := dao.Templates.Ctx(ctx).Count()
	if err != nil {
		return nil, err
	}
	statistics.TotalTemplates = totalTemplates

	// 获取总分类数
	totalCategories, err := dao.Categories.Ctx(ctx).Count()
	if err != nil {
		return nil, err
	}
	statistics.TotalCategories = totalCategories

	// 获取总语言数
	totalLanguages, err := dao.Languages.Ctx(ctx).Count()
	if err != nil {
		return nil, err
	}
	statistics.TotalLanguages = totalLanguages

	// 获取推荐模板数
	featuredCount, err := dao.Templates.Ctx(ctx).Where("is_featured", 1).Count()
	if err != nil {
		return nil, err
	}
	statistics.FeaturedCount = featuredCount

	return statistics, nil
}

// getCategoriesWithTemplates 获取分类及其模板
func (s *sIndex) getCategoriesWithTemplates(ctx context.Context, limit int) ([]*indexApi.CategoryWithTemplates, error) {
	// 获取所有分类
	var categories []*model.CategoriesInfo
	err := dao.Categories.Ctx(ctx).Order("sort desc").Scan(&categories)
	if err != nil {
		return nil, err
	}

	var result []*indexApi.CategoryWithTemplates

	for _, category := range categories {
		// 获取该分类下的模板
		var templates []*model.TemplatesInfo
		err := dao.Templates.Ctx(ctx).
			Where("category_id", category.Id).
			Limit(limit).
			Scan(&templates)
		if err != nil {
			g.Log().Warning(ctx, "Failed to get templates for category", category.Id, ":", err)
			continue
		}

		// 获取该分类下的模板总数
		templateCount, err := dao.Templates.Ctx(ctx).Where("category_id", category.Id).Count()
		if err != nil {
			g.Log().Warning(ctx, "Failed to get template count for category", category.Id, ":", err)
			templateCount = 0
		}

		categoryWithTemplates := &indexApi.CategoryWithTemplates{
			CategoriesInfo: category,
			Templates:      templates,
			TemplateCount:  templateCount,
		}

		result = append(result, categoryWithTemplates)
	}

	return result, nil
}

// getFeaturedTemplates 获取推荐模板
func (s *sIndex) getFeaturedTemplates(ctx context.Context, limit int) ([]*model.TemplatesInfo, error) {
	var templates []*model.TemplatesInfo
	err := dao.Templates.Ctx(ctx).
		Where("is_featured", 1).
		Limit(limit).
		Scan(&templates)
	return templates, err
}
