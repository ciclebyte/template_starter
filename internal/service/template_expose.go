package service

import (
	"context"

	api "github.com/ciclebyte/template_starter/api/v1/template_expose"
)

type ITemplateExpose interface {
	Get(ctx context.Context, req *api.TemplateExposeGetReq) (info *api.TemplateExposeInfo, err error)
	Set(ctx context.Context, req *api.TemplateExposeSetReq) (err error)
	Del(ctx context.Context, req *api.TemplateExposeDelReq) (err error)
	Versions(ctx context.Context, req *api.TemplateExposeVersionsReq) (versions []*api.TemplateExposeVersionInfo, err error)
}

var localTemplateExpose ITemplateExpose

func TemplateExpose() ITemplateExpose {
	if localTemplateExpose == nil {
		panic("implement not found for interface ITemplateExpose, forgot register?")
	}
	return localTemplateExpose
}

func RegisterTemplateExpose(i ITemplateExpose) {
	localTemplateExpose = i
}