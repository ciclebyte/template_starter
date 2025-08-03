package controller

import (
	"context"

	api "github.com/ciclebyte/template_starter/api/v1/template_files"
	consts "github.com/ciclebyte/template_starter/internal/consts"
	service "github.com/ciclebyte/template_starter/internal/service"
	"github.com/gogf/gf/v2/util/gconv"
)

var TemplateFiles = templateFilesController{}

type templateFilesController struct {
	BaseController
}

func (c *templateFilesController) Add(ctx context.Context, req *api.TemplateFilesAddReq) (res *api.TemplateFilesAddRes, err error) {
	res = new(api.TemplateFilesAddRes)
	err = service.TemplateFiles().Add(ctx, req)
	return
}

func (c *templateFilesController) List(ctx context.Context, req *api.TemplateFilesListReq) (res *api.TemplateFilesListRes, err error) {
	res = new(api.TemplateFilesListRes)
	if req.PageSize == 0 {
		req.PageSize = consts.PageSize
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	total, templateFiless, err := service.TemplateFiles().List(ctx, req)
	res.Total = total
	res.CurrentPage = req.PageNum
	res.TemplateFilesList = templateFiless
	return
}

func (c *templateFilesController) Get(ctx context.Context, req *api.TemplateFilesDetailReq) (res *api.TemplateFilesDetailRes, err error) {
	res = new(api.TemplateFilesDetailRes)
	service.TemplateFiles().GetById(ctx, gconv.Int64(req.Id))
	return
}

func (c *templateFilesController) Edit(ctx context.Context, req *api.TemplateFilesEditReq) (res *api.TemplateFilesEditRes, err error) {
	err = service.TemplateFiles().Edit(ctx, req)
	return
}

func (c *templateFilesController) Rename(ctx context.Context, req *api.TemplateFilesRenameReq) (res *api.TemplateFilesRenameRes, err error) {
	err = service.TemplateFiles().Rename(ctx, req)
	return
}

func (c *templateFilesController) Delete(ctx context.Context, req *api.TemplateFilesDelReq) (res *api.TemplateFilesDelRes, err error) {
	err = service.TemplateFiles().Delete(ctx, gconv.Int64(req.Id))
	return
}

func (c *templateFilesController) BatchDelete(ctx context.Context, req *api.TemplateFilesBatchDelReq) (res *api.TemplateFilesBatchDelRes, err error) {
	var int64Ids []int64
	for _, id := range req.Ids {
		int64Ids = append(int64Ids, gconv.Int64(id))
	}
	err = service.TemplateFiles().BatchDelete(ctx, int64Ids)
	return
}
func (c *templateFilesController) FileTree(ctx context.Context, req *api.TemplatesFileTreeReq) (res *api.TemplatesFileTreeRes, err error) {
	return service.TemplateFiles().FileTree(ctx, req)
}

func (c *templateFilesController) GetFileContent(ctx context.Context, req *api.TemplateFilesContentReq) (res *api.TemplateFilesContentRes, err error) {
	res = new(api.TemplateFilesContentRes)
	fileContent, err := service.TemplateFiles().GetFileContent(ctx, gconv.Int64(req.Id))
	if err != nil {
		return nil, err
	}
	res.FileContent = fileContent
	return
}

func (c *templateFilesController) UploadZip(ctx context.Context, req *api.TemplateFilesUploadZipReq) (res *api.TemplateFilesUploadZipRes, err error) {
	res = new(api.TemplateFilesUploadZipRes)
	successCount, failedFiles, err := service.TemplateFiles().UploadZip(ctx, gconv.Int64(req.TemplateId))
	if err != nil {
		return nil, err
	}
	res.SuccessCount = successCount
	res.FailedFiles = failedFiles
	if len(failedFiles) > 0 {
		res.Message = "部分文件解压失败"
	} else {
		res.Message = "ZIP包解压成功"
	}
	return
}

func (c *templateFilesController) UploadCode(ctx context.Context, req *api.TemplateFilesUploadCodeReq) (res *api.TemplateFilesUploadCodeRes, err error) {
	return service.TemplateFiles().UploadCode(ctx, req)
}

// Render 渲染模板文件
func (c *templateFilesController) Render(ctx context.Context, req *api.TemplateFilesRenderReq) (res *api.TemplateFilesRenderRes, err error) {
	return service.TemplateFiles().Render(ctx, req)
}

// RenderFileTree 渲染文件树
func (c *templateFilesController) RenderFileTree(ctx context.Context, req *api.TemplateFilesRenderFileTreeReq) (res *api.TemplateFilesRenderFileTreeRes, err error) {
	return service.TemplateFiles().RenderFileTree(ctx, req)
}

// DownloadZip 下载ZIP包
func (c *templateFilesController) DownloadZip(ctx context.Context, req *api.TemplateFilesDownloadZipReq) (res *api.TemplateFilesDownloadZipRes, err error) {
	err = service.TemplateFiles().DownloadZip(ctx, req)
	return
}

func (c *templateFilesController) Move(ctx context.Context, req *api.TemplateFilesMoveReq) (res *api.TemplateFilesMoveRes, err error) {
	res = new(api.TemplateFilesMoveRes)
	err = service.TemplateFiles().Move(ctx, req)
	return
}

// SetCondition 设置文件生成条件
func (c *templateFilesController) SetCondition(ctx context.Context, req *api.TemplateFilesSetConditionReq) (res *api.TemplateFilesSetConditionRes, err error) {
	res = new(api.TemplateFilesSetConditionRes)
	err = service.TemplateFiles().SetCondition(ctx, req)
	return
}

// GetCondition 获取文件生成条件
func (c *templateFilesController) GetCondition(ctx context.Context, req *api.TemplateFilesGetConditionReq) (res *api.TemplateFilesGetConditionRes, err error) {
	return service.TemplateFiles().GetCondition(ctx, req)
}
