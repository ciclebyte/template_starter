package service

import (
	"context"

	api "github.com/ciclebyte/template_starter/api/v1/languages"
	model "github.com/ciclebyte/template_starter/internal/model"
)

type ILanguages interface {
	List(ctx context.Context, req *api.LanguagesListReq) (total interface{}, res []*model.LanguagesInfo, err error)
	Add(ctx context.Context, req *api.LanguagesAddReq) (err error)
	Edit(ctx context.Context, req *api.LanguagesEditReq) (err error)
	Delete(ctx context.Context, id int) (err error)
	BatchDelete(ctx context.Context, ids []int) (err error)
	GetById(ctx context.Context, id int) (res *model.LanguagesInfo, err error)
}

var localLanguages ILanguages

func Languages() ILanguages {
	if localLanguages == nil {
		panic("implement not found for interface ILanguages, forgot register?")
	}
	return localLanguages
}

func RegisterLanguages(i ILanguages) {
	localLanguages = i
}
