package service

import (
	"context"

	systemConfigApi "github.com/ciclebyte/template_starter/api/v1/system_config"
)

type ISystemConfigService interface {
	// 基础CRUD
	Add(ctx context.Context, req *systemConfigApi.SystemConfigAddReq) (*systemConfigApi.SystemConfigAddRes, error)
	List(ctx context.Context, req *systemConfigApi.SystemConfigListReq) (*systemConfigApi.SystemConfigListRes, error)
	Detail(ctx context.Context, req *systemConfigApi.SystemConfigDetailReq) (*systemConfigApi.SystemConfigDetailRes, error)
	Edit(ctx context.Context, req *systemConfigApi.SystemConfigEditReq) (*systemConfigApi.SystemConfigEditRes, error)
	Del(ctx context.Context, req *systemConfigApi.SystemConfigDelReq) (*systemConfigApi.SystemConfigDelRes, error)
	BatchDel(ctx context.Context, req *systemConfigApi.SystemConfigBatchDelReq) (*systemConfigApi.SystemConfigBatchDelRes, error)
	
	// 配置操作
	GetByKey(ctx context.Context, req *systemConfigApi.SystemConfigGetByKeyReq) (*systemConfigApi.SystemConfigGetByKeyRes, error)
	GetByGroup(ctx context.Context, req *systemConfigApi.SystemConfigGetByGroupReq) (*systemConfigApi.SystemConfigGetByGroupRes, error)
	BatchUpdate(ctx context.Context, req *systemConfigApi.SystemConfigBatchUpdateReq) (*systemConfigApi.SystemConfigBatchUpdateRes, error)
	GetPublic(ctx context.Context, req *systemConfigApi.SystemConfigPublicReq) (*systemConfigApi.SystemConfigPublicRes, error)
	Reset(ctx context.Context, req *systemConfigApi.SystemConfigResetReq) (*systemConfigApi.SystemConfigResetRes, error)
	Validate(ctx context.Context, req *systemConfigApi.SystemConfigValidateReq) (*systemConfigApi.SystemConfigValidateRes, error)
	
	// 便捷方法
	GetConfigValue(ctx context.Context, configKey string) (interface{}, error)
	GetConfigString(ctx context.Context, configKey string, defaultValue ...string) (string, error)
	GetConfigInt(ctx context.Context, configKey string, defaultValue ...int) (int, error)
	GetConfigBool(ctx context.Context, configKey string, defaultValue ...bool) (bool, error)
	GetConfigJson(ctx context.Context, configKey string, result interface{}) error
	SetConfigValue(ctx context.Context, configKey string, value interface{}) error
}

var localSystemConfigService ISystemConfigService

func SystemConfig() ISystemConfigService {
	if localSystemConfigService == nil {
		panic("implement not found for interface ISystemConfigService, forgot register?")
	}
	return localSystemConfigService
}

func RegisterSystemConfigService(i ISystemConfigService) {
	localSystemConfigService = i
}