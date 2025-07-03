package service

import (
	"context"

	api "github.com/ciclebyte/template_starter/api/v1/template_languages"
	model "github.com/ciclebyte/template_starter/internal/model"
)

type ITemplateLanguages interface {
	List(ctx context.Context, req *api.TemplateLanguagesListReq) (total interface{}, res []*model.TemplateLanguagesInfo, err error)
	Add(ctx context.Context, req *api.TemplateLanguagesAddReq) (err error)
	Edit(ctx context.Context, req *api.TemplateLanguagesEditReq) (err error)
	Delete(ctx context.Context, id int64) (err error)
	BatchDelete(ctx context.Context, ids []int64) (err error)
	GetById(ctx context.Context, id int64) (res *model.TemplateLanguagesInfo, err error)
}

var localTemplateLanguages ITemplateLanguages

func TemplateLanguages() ITemplateLanguages {
	if localTemplateLanguages == nil {
		panic("implement not found for interface ITemplateLanguages, forgot register?")
	}
	return localTemplateLanguages
}

func RegisterTemplateLanguages(i ITemplateLanguages) {
	localTemplateLanguages = i
}
