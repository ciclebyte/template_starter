package service

import (
	"context"

	api "github.com/ciclebyte/template_starter/api/v1/template_variables"
	model "github.com/ciclebyte/template_starter/internal/model"
)

type ITemplateVariables interface {
	List(ctx context.Context, req *api.TemplateVariablesListReq) (total interface{}, res []*model.TemplateVariablesInfo, err error)
	Add(ctx context.Context, req *api.TemplateVariablesAddReq) (err error)
	Edit(ctx context.Context, req *api.TemplateVariablesEditReq) (err error)
	Delete(ctx context.Context, id int64) (err error)
	BatchDelete(ctx context.Context, ids []int64) (err error)
	GetById(ctx context.Context, id int64) (res *model.TemplateVariablesInfo, err error)
}

var localTemplateVariables ITemplateVariables

func TemplateVariables() ITemplateVariables {
	if localTemplateVariables == nil {
		panic("implement not found for interface ITemplateVariables, forgot register?")
	}
	return localTemplateVariables
}

func RegisterTemplateVariables(i ITemplateVariables) {
	localTemplateVariables = i
}
