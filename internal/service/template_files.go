package service

import (
	"context"

	api "github.com/ciclebyte/template_starter/api/v1/template_files"
	model "github.com/ciclebyte/template_starter/internal/model"
)

type ITemplateFiles interface {
	List(ctx context.Context, req *api.TemplateFilesListReq) (total interface{}, res []*model.TemplateFilesInfo, err error)
	Add(ctx context.Context, req *api.TemplateFilesAddReq) (err error)
	Edit(ctx context.Context, req *api.TemplateFilesEditReq) (err error)
	Delete(ctx context.Context, id int64) (err error)
	BatchDelete(ctx context.Context, ids []int64) (err error)
	GetById(ctx context.Context, id int64) (res *model.TemplateFilesInfo, err error)
	FileTree(ctx context.Context, req *api.TemplatesFileTreeReq) (res *api.TemplatesFileTreeRes, err error)
	GetFileContent(ctx context.Context, id int64) (fileContent string, err error)
	UploadZip(ctx context.Context, templateId int64) (successCount int, failedFiles []string, err error)
}

var localTemplateFiles ITemplateFiles

func TemplateFiles() ITemplateFiles {
	if localTemplateFiles == nil {
		panic("implement not found for interface ITemplateFiles, forgot register?")
	}
	return localTemplateFiles
}

func RegisterTemplateFiles(i ITemplateFiles) {
	localTemplateFiles = i
}
