package service

import (
	"context"

	api "github.com/ciclebyte/template_starter/api/v1/tags"
	model "github.com/ciclebyte/template_starter/internal/model"
)

type ITags interface {
	// 标签基础操作
	List(ctx context.Context, req *api.TagsListReq) (total interface{}, res []*model.TagsInfo, err error)
	Add(ctx context.Context, req *api.TagsAddReq) (err error)
	Edit(ctx context.Context, req *api.TagsEditReq) (err error)
	Delete(ctx context.Context, id int64) (err error)
	BatchDelete(ctx context.Context, ids []int64) (err error)
	GetById(ctx context.Context, id int64) (res *model.TagsInfo, err error)
	All(ctx context.Context) (res []*model.TagsInfo, err error)
	WithCount(ctx context.Context) (res []*model.TagWithTemplateCount, err error)

	// 模板标签关联操作
	AddTemplateTags(ctx context.Context, templateId int64, tagIds []int64) (err error)
	RemoveTemplateTags(ctx context.Context, templateId int64, tagIds []int64) (err error)
	SetTemplateTags(ctx context.Context, templateId int64, tagIds []int64) (err error)
	GetTemplateTags(ctx context.Context, templateId int64) (res []*model.TagsInfo, err error)
	GetTemplatesByTag(ctx context.Context, tagId int64, req *api.TemplatesByTagReq) (total interface{}, res []*model.TemplatesInfo, err error)
}

var localTags ITags

func Tags() ITags {
	if localTags == nil {
		panic("implement not found for interface ITags, forgot register?")
	}
	return localTags
}

func RegisterTags(i ITags) {
	localTags = i
}