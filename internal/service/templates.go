package service

import (
	"context"

	api "github.com/ciclebyte/template_starter/api/v1/templates"
	model "github.com/ciclebyte/template_starter/internal/model"
)

type ITemplates interface {
	List(ctx context.Context, req *api.TemplatesListReq) (total interface{}, res []*model.TemplatesInfo, err error)
	Add(ctx context.Context, req *api.TemplatesAddReq) (err error)
	Edit(ctx context.Context, req *api.TemplatesEditReq) (err error)
	Delete(ctx context.Context, id int64) (err error)
	BatchDelete(ctx context.Context, ids []int64) (err error)
	GetById(ctx context.Context, id int64) (res *model.TemplatesInfo, err error)
}

var localTemplates ITemplates

func Templates() ITemplates {
	if localTemplates == nil {
		panic("implement not found for interface ITemplates, forgot register?")
	}
	return localTemplates
}

func RegisterTemplates(i ITemplates) {
	localTemplates = i
}
