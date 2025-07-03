package controller

import (
	"context"

	api "github.com/ciclebyte/template_starter/api/v1/languages"
	consts "github.com/ciclebyte/template_starter/internal/consts"
	service "github.com/ciclebyte/template_starter/internal/service"
)

var Languages = languagesController{}

type languagesController struct {
	BaseController
}

func (c *languagesController) Add(ctx context.Context, req *api.LanguagesAddReq) (res *api.LanguagesAddRes, err error) {
	res = new(api.LanguagesAddRes)
	err = service.Languages().Add(ctx, req)
	return
}

func (c *languagesController) List(ctx context.Context, req *api.LanguagesListReq) (res *api.LanguagesListRes, err error) {
	res = new(api.LanguagesListRes)
	if req.PageSize == 0 {
		req.PageSize = consts.PageSize
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	total, languagess, err := service.Languages().List(ctx, req)
	res.Total = total
	res.CurrentPage = req.PageNum
	res.LanguagesList = languagess
	return
}

func (c *languagesController) Get(ctx context.Context, req *api.LanguagesDetailReq) (res *api.LanguagesDetailRes, err error) {
	res = new(api.LanguagesDetailRes)
	service.Languages().GetById(ctx, req.Id)
	return
}

func (c *languagesController) Edit(ctx context.Context, req *api.LanguagesEditReq) (res *api.LanguagesEditRes, err error) {
	err = service.Languages().Edit(ctx, req)
	return
}

func (c *languagesController) Delete(ctx context.Context, req *api.LanguagesDelReq) (res *api.LanguagesDelRes, err error) {
	err = service.Languages().Delete(ctx, req.Id)
	return
}

func (c *languagesController) BatchDelete(ctx context.Context, req *api.LanguagesBatchDelReq) (res *api.LanguagesBatchDelRes, err error) {
	err = service.Languages().BatchDelete(ctx, req.Ids)
	return
}
