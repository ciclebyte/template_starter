package controller

import (
	"context"

	systemConfigApi "github.com/ciclebyte/template_starter/api/v1/system_config"
	"github.com/ciclebyte/template_starter/internal/service"
	"github.com/gogf/gf/v2/frame/g"
)

var SystemConfig = cSystemConfig{}

type cSystemConfig struct {
	BaseController
}

// Add 新增配置
func (c *cSystemConfig) Add(ctx context.Context, req *systemConfigApi.SystemConfigAddReq) (res *systemConfigApi.SystemConfigAddRes, err error) {
	g.Log().Debug(ctx, "SystemConfig.Add called")
	return service.SystemConfig().Add(ctx, req)
}

// List 配置列表
func (c *cSystemConfig) List(ctx context.Context, req *systemConfigApi.SystemConfigListReq) (res *systemConfigApi.SystemConfigListRes, err error) {
	g.Log().Debug(ctx, "SystemConfig.List called")
	return service.SystemConfig().List(ctx, req)
}

// Detail 配置详情
func (c *cSystemConfig) Detail(ctx context.Context, req *systemConfigApi.SystemConfigDetailReq) (res *systemConfigApi.SystemConfigDetailRes, err error) {
	g.Log().Debug(ctx, "SystemConfig.Detail called")
	return service.SystemConfig().Detail(ctx, req)
}

// Edit 编辑配置
func (c *cSystemConfig) Edit(ctx context.Context, req *systemConfigApi.SystemConfigEditReq) (res *systemConfigApi.SystemConfigEditRes, err error) {
	g.Log().Debug(ctx, "SystemConfig.Edit called")
	return service.SystemConfig().Edit(ctx, req)
}

// Del 删除配置
func (c *cSystemConfig) Del(ctx context.Context, req *systemConfigApi.SystemConfigDelReq) (res *systemConfigApi.SystemConfigDelRes, err error) {
	g.Log().Debug(ctx, "SystemConfig.Del called")
	return service.SystemConfig().Del(ctx, req)
}

// BatchDel 批量删除配置
func (c *cSystemConfig) BatchDel(ctx context.Context, req *systemConfigApi.SystemConfigBatchDelReq) (res *systemConfigApi.SystemConfigBatchDelRes, err error) {
	g.Log().Debug(ctx, "SystemConfig.BatchDel called")
	return service.SystemConfig().BatchDel(ctx, req)
}

// GetByKey 根据键名获取配置
func (c *cSystemConfig) GetByKey(ctx context.Context, req *systemConfigApi.SystemConfigGetByKeyReq) (res *systemConfigApi.SystemConfigGetByKeyRes, err error) {
	g.Log().Debug(ctx, "SystemConfig.GetByKey called")
	return service.SystemConfig().GetByKey(ctx, req)
}

// GetByGroup 根据分组获取配置
func (c *cSystemConfig) GetByGroup(ctx context.Context, req *systemConfigApi.SystemConfigGetByGroupReq) (res *systemConfigApi.SystemConfigGetByGroupRes, err error) {
	g.Log().Debug(ctx, "SystemConfig.GetByGroup called")
	return service.SystemConfig().GetByGroup(ctx, req)
}

// BatchUpdate 批量更新配置
func (c *cSystemConfig) BatchUpdate(ctx context.Context, req *systemConfigApi.SystemConfigBatchUpdateReq) (res *systemConfigApi.SystemConfigBatchUpdateRes, err error) {
	g.Log().Debug(ctx, "SystemConfig.BatchUpdate called")
	return service.SystemConfig().BatchUpdate(ctx, req)
}

// GetPublic 获取公开配置
func (c *cSystemConfig) GetPublic(ctx context.Context, req *systemConfigApi.SystemConfigPublicReq) (res *systemConfigApi.SystemConfigPublicRes, err error) {
	g.Log().Debug(ctx, "SystemConfig.GetPublic called")
	return service.SystemConfig().GetPublic(ctx, req)
}

// Reset 重置配置
func (c *cSystemConfig) Reset(ctx context.Context, req *systemConfigApi.SystemConfigResetReq) (res *systemConfigApi.SystemConfigResetRes, err error) {
	g.Log().Debug(ctx, "SystemConfig.Reset called")
	return service.SystemConfig().Reset(ctx, req)
}

// Validate 验证配置值
func (c *cSystemConfig) Validate(ctx context.Context, req *systemConfigApi.SystemConfigValidateReq) (res *systemConfigApi.SystemConfigValidateRes, err error) {
	g.Log().Debug(ctx, "SystemConfig.Validate called")
	return service.SystemConfig().Validate(ctx, req)
}