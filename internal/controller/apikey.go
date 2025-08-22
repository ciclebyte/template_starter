package controller

import (
	"context"

	api "github.com/ciclebyte/template_starter/api/v1/apikey"
	service "github.com/ciclebyte/template_starter/internal/service"
)

var ApiKey = apiKeyController{}

type apiKeyController struct {
	BaseController
}

// ============================================================================
// 管理员 API Key 管理
// ============================================================================

// GetApiKeys 获取API Key列表
func (c *apiKeyController) GetApiKeys(ctx context.Context, req *api.GetApiKeysReq) (res *api.GetApiKeysRes, err error) {
	res = new(api.GetApiKeysRes)

	result, err := service.ApiKey().GetApiKeys(ctx, req)
	if err != nil {
		return nil, err
	}

	*res = *result
	return
}

// CreateApiKey 创建API Key
func (c *apiKeyController) CreateApiKey(ctx context.Context, req *api.CreateApiKeyReq) (res *api.CreateApiKeyRes, err error) {
	res = new(api.CreateApiKeyRes)

	result, err := service.ApiKey().CreateApiKey(ctx, req)
	if err != nil {
		return nil, err
	}

	*res = *result
	return
}

// UpdateApiKey 更新API Key
func (c *apiKeyController) UpdateApiKey(ctx context.Context, req *api.UpdateApiKeyReq) (res *api.UpdateApiKeyRes, err error) {
	res = new(api.UpdateApiKeyRes)

	_, err = service.ApiKey().UpdateApiKey(ctx, req)
	if err != nil {
		return nil, err
	}

	return
}

// DeleteApiKey 删除API Key
func (c *apiKeyController) DeleteApiKey(ctx context.Context, req *api.DeleteApiKeyReq) (res *api.DeleteApiKeyRes, err error) {
	res = new(api.DeleteApiKeyRes)

	_, err = service.ApiKey().DeleteApiKey(ctx, req)
	if err != nil {
		return nil, err
	}

	return
}

// RegenerateApiKey 重新生成API Key Secret
func (c *apiKeyController) RegenerateApiKey(ctx context.Context, req *api.RegenerateApiKeyReq) (res *api.RegenerateApiKeyRes, err error) {
	res = new(api.RegenerateApiKeyRes)

	result, err := service.ApiKey().RegenerateApiKey(ctx, req)
	if err != nil {
		return nil, err
	}

	*res = *result
	return
}

// GetApiKeyLogs 获取API Key使用日志
func (c *apiKeyController) GetApiKeyLogs(ctx context.Context, req *api.GetApiKeyLogsReq) (res *api.GetApiKeyLogsRes, err error) {
	res = new(api.GetApiKeyLogsRes)

	result, err := service.ApiKey().GetApiKeyLogs(ctx, req)
	if err != nil {
		return nil, err
	}

	*res = *result
	return
}