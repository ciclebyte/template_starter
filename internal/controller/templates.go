package controller

import (
	"context"

	api "github.com/ciclebyte/template_starter/api/v1/templates"
	consts "github.com/ciclebyte/template_starter/internal/consts"
	service "github.com/ciclebyte/template_starter/internal/service"
	"github.com/gogf/gf/v2/util/gconv"
)

var Templates = templatesController{}

type templatesController struct {
	BaseController
}

func (c *templatesController) Add(ctx context.Context, req *api.TemplatesAddReq) (res *api.TemplatesAddRes, err error) {
	res = new(api.TemplatesAddRes)
	err = service.Templates().Add(ctx, req)
	return
}

func (c *templatesController) List(ctx context.Context, req *api.TemplatesListReq) (res *api.TemplatesListRes, err error) {
	res = new(api.TemplatesListRes)
	if req.PageSize == 0 {
		req.PageSize = consts.PageSize
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	total, templatess, err := service.Templates().List(ctx, req)
	res.Total = total
	res.CurrentPage = req.PageNum
	res.TemplatesList = templatess
	return
}

func (c *templatesController) Get(ctx context.Context, req *api.TemplatesDetailReq) (res *api.TemplatesDetailRes, err error) {
	res = new(api.TemplatesDetailRes)
	templateInfo, err := service.Templates().GetById(ctx, gconv.Int64(req.Id))
	if err != nil {
		return nil, err
	}
	res.TemplatesInfo = templateInfo
	return
}

func (c *templatesController) Edit(ctx context.Context, req *api.TemplatesEditReq) (res *api.TemplatesEditRes, err error) {
	err = service.Templates().Edit(ctx, req)
	return
}

func (c *templatesController) Delete(ctx context.Context, req *api.TemplatesDelReq) (res *api.TemplatesDelRes, err error) {
	err = service.Templates().Delete(ctx, gconv.Int64(req.Id))
	return
}

func (c *templatesController) BatchDelete(ctx context.Context, req *api.TemplatesBatchDelReq) (res *api.TemplatesBatchDelRes, err error) {
	var int64Ids []int64
	for _, id := range req.Ids {
		int64Ids = append(int64Ids, gconv.Int64(id))
	}
	err = service.Templates().BatchDelete(ctx, int64Ids)
	return
}

func (c *templatesController) GetVariables(ctx context.Context, req *api.TemplatesVariablesReq) (res *api.TemplatesVariablesRes, err error) {
	res, err = service.Templates().GetVariables(ctx, gconv.Int64(req.TemplateId))
	return
}

func (c *templatesController) AnalyzeVariables(ctx context.Context, req *api.TemplatesAnalyzeVariablesReq) (res *api.TemplatesAnalyzeVariablesRes, err error) {
	res, err = service.Templates().AnalyzeVariables(ctx, gconv.Int64(req.TemplateId))
	return
}

func (c *templatesController) GetTypes(ctx context.Context, req *api.TemplateTypesReq) (res *api.TemplateTypesRes, err error) {
	res = &api.TemplateTypesRes{
		TemplateTypes: []*api.TemplateTypeInfo{
			{Value: "basic", Label: "基础模板", Description: "适合简单的代码生成场景"},
			{Value: "scaffold", Label: "脚手架模板", Description: "适合项目初始化和架构搭建"},
			{Value: "data_driven", Label: "数据驱动模板", Description: "基于数据源动态生成代码"},
		},
	}
	return
}
