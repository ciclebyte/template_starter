package controller

import (
	"context"

	api "github.com/ciclebyte/template_starter/api/v1/template_languages"
	consts "github.com/ciclebyte/template_starter/internal/consts"
	service "github.com/ciclebyte/template_starter/internal/service"
	"github.com/gogf/gf/v2/util/gconv"
)

var TemplateLanguages = templateLanguagesController{}

type templateLanguagesController struct {
	BaseController
}

func (c *templateLanguagesController) Add(ctx context.Context, req *api.TemplateLanguagesAddReq) (res *api.TemplateLanguagesAddRes, err error) {
	res = new(api.TemplateLanguagesAddRes)
	err = service.TemplateLanguages().Add(ctx, req)
	return
}

func (c *templateLanguagesController) List(ctx context.Context, req *api.TemplateLanguagesListReq) (res *api.TemplateLanguagesListRes, err error) {
	res = new(api.TemplateLanguagesListRes)
	if req.PageSize == 0 {
		req.PageSize = consts.PageSize
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	total, templateLanguagess, err := service.TemplateLanguages().List(ctx, req)
	res.Total = total
	res.CurrentPage = req.PageNum
	res.TemplateLanguagesList = templateLanguagess
	return
}

func (c *templateLanguagesController) Get(ctx context.Context, req *api.TemplateLanguagesDetailReq) (res *api.TemplateLanguagesDetailRes, err error) {
	res = new(api.TemplateLanguagesDetailRes)
	service.TemplateLanguages().GetById(ctx, gconv.Int64(req.Id))
	return
}

func (c *templateLanguagesController) Edit(ctx context.Context, req *api.TemplateLanguagesEditReq) (res *api.TemplateLanguagesEditRes, err error) {
	err = service.TemplateLanguages().Edit(ctx, req)
	return
}

func (c *templateLanguagesController) Delete(ctx context.Context, req *api.TemplateLanguagesDelReq) (res *api.TemplateLanguagesDelRes, err error) {
	err = service.TemplateLanguages().Delete(ctx, gconv.Int64(req.Id))
	return
}

func (c *templateLanguagesController) BatchDelete(ctx context.Context, req *api.TemplateLanguagesBatchDelReq) (res *api.TemplateLanguagesBatchDelRes, err error) {
	var int64Ids []int64
	for _, id := range req.Ids {
		int64Ids = append(int64Ids, gconv.Int64(id))
	}
	err = service.TemplateLanguages().BatchDelete(ctx, int64Ids)
	return
}
