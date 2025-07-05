package controller

import (
	"context"

	api "github.com/ciclebyte/template_starter/api/v1/template_variables"
	service "github.com/ciclebyte/template_starter/internal/service"
	"github.com/gogf/gf/v2/util/gconv"
)

var TemplateVariables = templateVariablesController{}

type templateVariablesController struct {
	BaseController
}

func (c *templateVariablesController) Add(ctx context.Context, req *api.TemplateVariablesAddReq) (res *api.TemplateVariablesAddRes, err error) {
	res = new(api.TemplateVariablesAddRes)
	err = service.TemplateVariables().Add(ctx, req)
	return
}

func (c *templateVariablesController) List(ctx context.Context, req *api.TemplateVariablesListReq) (res *api.TemplateVariablesListRes, err error) {
	res = new(api.TemplateVariablesListRes)
	templateVariabless, err := service.TemplateVariables().List(ctx, req)
	res.TemplateVariablesList = templateVariabless
	return
}

func (c *templateVariablesController) Get(ctx context.Context, req *api.TemplateVariablesDetailReq) (res *api.TemplateVariablesDetailRes, err error) {
	res = new(api.TemplateVariablesDetailRes)
	service.TemplateVariables().GetById(ctx, gconv.Int64(req.Id))
	return
}

func (c *templateVariablesController) Edit(ctx context.Context, req *api.TemplateVariablesEditReq) (res *api.TemplateVariablesEditRes, err error) {
	err = service.TemplateVariables().Edit(ctx, req)
	return
}

func (c *templateVariablesController) Delete(ctx context.Context, req *api.TemplateVariablesDelReq) (res *api.TemplateVariablesDelRes, err error) {
	err = service.TemplateVariables().Delete(ctx, gconv.Int64(req.Id))
	return
}

func (c *templateVariablesController) BatchDelete(ctx context.Context, req *api.TemplateVariablesBatchDelReq) (res *api.TemplateVariablesBatchDelRes, err error) {
	var int64Ids []int64
	for _, id := range req.Ids {
		int64Ids = append(int64Ids, gconv.Int64(id))
	}
	err = service.TemplateVariables().BatchDelete(ctx, int64Ids)
	return
}
