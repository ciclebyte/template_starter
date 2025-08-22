package service

import (
	"context"

	apikeyApi "github.com/ciclebyte/template_starter/api/v1/apikey"
	"github.com/ciclebyte/template_starter/internal/model/entity"
)

// IApiKey API Key服务接口
type IApiKey interface {
	// ============================================================================
	// 管理员 API Key 管理
	// ============================================================================
	
	// GetApiKeys 获取API Key列表
	GetApiKeys(ctx context.Context, req *apikeyApi.GetApiKeysReq) (*apikeyApi.GetApiKeysRes, error)
	
	// CreateApiKey 创建API Key
	CreateApiKey(ctx context.Context, req *apikeyApi.CreateApiKeyReq) (*apikeyApi.CreateApiKeyRes, error)
	
	// UpdateApiKey 更新API Key
	UpdateApiKey(ctx context.Context, req *apikeyApi.UpdateApiKeyReq) (*apikeyApi.UpdateApiKeyRes, error)
	
	// DeleteApiKey 删除API Key
	DeleteApiKey(ctx context.Context, req *apikeyApi.DeleteApiKeyReq) (*apikeyApi.DeleteApiKeyRes, error)
	
	// RegenerateApiKey 重新生成API Key Secret
	RegenerateApiKey(ctx context.Context, req *apikeyApi.RegenerateApiKeyReq) (*apikeyApi.RegenerateApiKeyRes, error)
	
	// GetApiKeyLogs 获取API Key使用日志
	GetApiKeyLogs(ctx context.Context, req *apikeyApi.GetApiKeyLogsReq) (*apikeyApi.GetApiKeyLogsRes, error)

	// ============================================================================
	// 个人 API Key 管理
	// ============================================================================
	
	// GetMyApiKeys 获取我的API Key列表
	GetMyApiKeys(ctx context.Context, req *apikeyApi.GetMyApiKeysReq) (*apikeyApi.GetMyApiKeysRes, error)
	
	// CreateMyApiKey 创建我的API Key
	CreateMyApiKey(ctx context.Context, req *apikeyApi.CreateMyApiKeyReq) (*apikeyApi.CreateMyApiKeyRes, error)
	
	// UpdateMyApiKey 更新我的API Key
	UpdateMyApiKey(ctx context.Context, req *apikeyApi.UpdateMyApiKeyReq) (*apikeyApi.UpdateMyApiKeyRes, error)
	
	// DeleteMyApiKey 删除我的API Key
	DeleteMyApiKey(ctx context.Context, req *apikeyApi.DeleteMyApiKeyReq) (*apikeyApi.DeleteMyApiKeyRes, error)
	
	// RegenerateMyApiKey 重新生成我的API Key Secret
	RegenerateMyApiKey(ctx context.Context, req *apikeyApi.RegenerateMyApiKeyReq) (*apikeyApi.RegenerateMyApiKeyRes, error)

	// ============================================================================
	// API Key 验证和使用
	// ============================================================================
	
	// VerifyApiKey 验证API Key
	VerifyApiKey(ctx context.Context, keyId, keySecret string) (*entity.ApiKeys, error)
	
	// UpdateApiKeyUsage 更新API Key使用信息
	UpdateApiKeyUsage(ctx context.Context, apiKeyId int64, ip string) error
	
	// LogApiKeyUsage 记录API Key使用日志
	LogApiKeyUsage(ctx context.Context, apiKeyId, userId int64, method, path, ip, userAgent string, statusCode, responseTime int) error
}

var localApiKey IApiKey

func ApiKey() IApiKey {
	if localApiKey == nil {
		panic("implement not found for interface IApiKey, forgot register?")
	}
	return localApiKey
}

func RegisterApiKey(i IApiKey) {
	localApiKey = i
}