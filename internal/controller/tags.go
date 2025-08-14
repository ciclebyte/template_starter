package controller

import (
	"context"

	api "github.com/ciclebyte/template_starter/api/v1/tags"
	consts "github.com/ciclebyte/template_starter/internal/consts"
	service "github.com/ciclebyte/template_starter/internal/service"
	"github.com/gogf/gf/v2/util/gconv"
)

var Tags = tagsController{}

type tagsController struct {
	BaseController
}

// Add 新增标签
func (c *tagsController) Add(ctx context.Context, req *api.TagsAddReq) (res *api.TagsAddRes, err error) {
	res = new(api.TagsAddRes)
	err = service.Tags().Add(ctx, req)
	return
}

// List 标签列表
func (c *tagsController) List(ctx context.Context, req *api.TagsListReq) (res *api.TagsListRes, err error) {
	res = new(api.TagsListRes)
	if req.PageSize == 0 {
		req.PageSize = consts.PageSize
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	total, tagsList, err := service.Tags().List(ctx, req)
	res.Total = total
	res.CurrentPage = req.PageNum
	res.TagsList = tagsList
	return
}

// Get 标签详情
func (c *tagsController) Get(ctx context.Context, req *api.TagsDetailReq) (res *api.TagsDetailRes, err error) {
	res = new(api.TagsDetailRes)
	tagInfo, err := service.Tags().GetById(ctx, gconv.Int64(req.Id))
	if err != nil {
		return nil, err
	}
	res.TagsInfo = tagInfo
	return
}

// Edit 修改标签
func (c *tagsController) Edit(ctx context.Context, req *api.TagsEditReq) (res *api.TagsEditRes, err error) {
	err = service.Tags().Edit(ctx, req)
	return
}

// Delete 删除标签
func (c *tagsController) Delete(ctx context.Context, req *api.TagsDelReq) (res *api.TagsDelRes, err error) {
	err = service.Tags().Delete(ctx, gconv.Int64(req.Id))
	return
}

// BatchDelete 批量删除标签
func (c *tagsController) BatchDelete(ctx context.Context, req *api.TagsBatchDelReq) (res *api.TagsBatchDelRes, err error) {
	var int64Ids []int64
	for _, id := range req.Ids {
		int64Ids = append(int64Ids, gconv.Int64(id))
	}
	err = service.Tags().BatchDelete(ctx, int64Ids)
	return
}

// All 获取所有标签（不分页）
func (c *tagsController) All(ctx context.Context, req *api.TagsAllReq) (res *api.TagsAllRes, err error) {
	res = new(api.TagsAllRes)
	tagsList, err := service.Tags().All(ctx)
	res.TagsList = tagsList
	return
}

// WithCount 获取标签及关联模板数量
func (c *tagsController) WithCount(ctx context.Context, req *api.TagsWithCountReq) (res *api.TagsWithCountRes, err error) {
	res = new(api.TagsWithCountRes)
	tagsList, err := service.Tags().WithCount(ctx)
	res.TagsList = tagsList
	return
}

// AddTemplateTags 为模板添加标签
func (c *tagsController) AddTemplateTags(ctx context.Context, req *api.TemplateTagsAddReq) (res *api.TemplateTagsAddRes, err error) {
	res = new(api.TemplateTagsAddRes)
	err = service.Tags().AddTemplateTags(ctx, gconv.Int64(req.TemplateId), req.TagIds)
	return
}

// RemoveTemplateTags 为模板移除标签
func (c *tagsController) RemoveTemplateTags(ctx context.Context, req *api.TemplateTagsRemoveReq) (res *api.TemplateTagsRemoveRes, err error) {
	res = new(api.TemplateTagsRemoveRes)
	err = service.Tags().RemoveTemplateTags(ctx, gconv.Int64(req.TemplateId), req.TagIds)
	return
}

// SetTemplateTags 设置模板标签（批量设置，覆盖原有标签）
func (c *tagsController) SetTemplateTags(ctx context.Context, req *api.TemplateTagsSetReq) (res *api.TemplateTagsSetRes, err error) {
	res = new(api.TemplateTagsSetRes)
	err = service.Tags().SetTemplateTags(ctx, gconv.Int64(req.TemplateId), req.TagIds)
	return
}

// GetTemplateTags 获取模板标签列表
func (c *tagsController) GetTemplateTags(ctx context.Context, req *api.TemplateTagsListReq) (res *api.TemplateTagsListRes, err error) {
	res = new(api.TemplateTagsListRes)
	tagsList, err := service.Tags().GetTemplateTags(ctx, gconv.Int64(req.TemplateId))
	res.TagsList = tagsList
	return
}

// GetTemplatesByTag 根据标签获取模板列表
func (c *tagsController) GetTemplatesByTag(ctx context.Context, req *api.TemplatesByTagReq) (res *api.TemplatesByTagRes, err error) {
	res = new(api.TemplatesByTagRes)
	if req.PageSize == 0 {
		req.PageSize = consts.PageSize
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	total, templatesList, err := service.Tags().GetTemplatesByTag(ctx, gconv.Int64(req.TagId), req)
	res.Total = total
	res.CurrentPage = req.PageNum
	res.TemplatesList = templatesList
	return
}