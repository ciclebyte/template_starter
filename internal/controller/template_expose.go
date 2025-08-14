package controller

import (
	"context"
	"github.com/ciclebyte/template_starter/api/v1/template_expose"
	"github.com/ciclebyte/template_starter/internal/service"
)

// templateExposeController 模板暴露字段控制器
type templateExposeController struct{}

var TemplateExpose = &templateExposeController{}

// Get 获取模板暴露字段
func (c *templateExposeController) Get(ctx context.Context, req *template_expose.TemplateExposeGetReq) (res *template_expose.TemplateExposeGetRes, err error) {
	res = new(template_expose.TemplateExposeGetRes)
	res.TemplateExpose, err = service.TemplateExpose().Get(ctx, req)
	return
}

// Set 设置模板暴露字段
func (c *templateExposeController) Set(ctx context.Context, req *template_expose.TemplateExposeSetReq) (res *template_expose.TemplateExposeSetRes, err error) {
	res = new(template_expose.TemplateExposeSetRes)
	err = service.TemplateExpose().Set(ctx, req)
	return
}

// Del 删除模板暴露字段
func (c *templateExposeController) Del(ctx context.Context, req *template_expose.TemplateExposeDelReq) (res *template_expose.TemplateExposeDelRes, err error) {
	res = new(template_expose.TemplateExposeDelRes)
	err = service.TemplateExpose().Del(ctx, req)
	return
}

// Versions 获取模板暴露字段历史版本
func (c *templateExposeController) Versions(ctx context.Context, req *template_expose.TemplateExposeVersionsReq) (res *template_expose.TemplateExposeVersionsRes, err error) {
	res = new(template_expose.TemplateExposeVersionsRes)
	res.Versions, err = service.TemplateExpose().Versions(ctx, req)
	return
}