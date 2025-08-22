package controller

import (
	"context"

	profileApi "github.com/ciclebyte/template_starter/api/v1/profile"
	apikeyApi "github.com/ciclebyte/template_starter/api/v1/apikey"
	service "github.com/ciclebyte/template_starter/internal/service"
)

var Profile = profileController{}

type profileController struct {
	BaseController
}

// GetProfile 获取个人资料
func (c *profileController) GetProfile(ctx context.Context, req *profileApi.GetProfileReq) (res *profileApi.GetProfileRes, err error) {
	res = new(profileApi.GetProfileRes)

	result, err := service.Profile().GetProfile(ctx, req)
	if err != nil {
		return nil, err
	}

	*res = *result
	return
}

// UpdateProfile 更新个人资料
func (c *profileController) UpdateProfile(ctx context.Context, req *profileApi.UpdateProfileReq) (res *profileApi.UpdateProfileRes, err error) {
	res = new(profileApi.UpdateProfileRes)

	_, err = service.Profile().UpdateProfile(ctx, req)
	if err != nil {
		return nil, err
	}

	return
}

// ChangePassword 修改密码
func (c *profileController) ChangePassword(ctx context.Context, req *profileApi.ChangePasswordReq) (res *profileApi.ChangePasswordRes, err error) {
	res = new(profileApi.ChangePasswordRes)

	_, err = service.Profile().ChangePassword(ctx, req)
	if err != nil {
		return nil, err
	}

	return
}

// UpdateEmail 更新邮箱
func (c *profileController) UpdateEmail(ctx context.Context, req *profileApi.UpdateEmailReq) (res *profileApi.UpdateEmailRes, err error) {
	res = new(profileApi.UpdateEmailRes)

	_, err = service.Profile().UpdateEmail(ctx, req)
	if err != nil {
		return nil, err
	}

	return
}

// UploadAvatar 上传头像
func (c *profileController) UploadAvatar(ctx context.Context, req *profileApi.UploadAvatarReq) (res *profileApi.UploadAvatarRes, err error) {
	res = new(profileApi.UploadAvatarRes)

	result, err := service.Profile().UploadAvatar(ctx, req)
	if err != nil {
		return nil, err
	}

	res.AvatarUrl = result.AvatarUrl
	return
}

// GetSecurityInfo 获取账户安全信息
func (c *profileController) GetSecurityInfo(ctx context.Context, req *profileApi.GetSecurityInfoReq) (res *profileApi.GetSecurityInfoRes, err error) {
	res = new(profileApi.GetSecurityInfoRes)

	result, err := service.Profile().GetSecurityInfo(ctx, req)
	if err != nil {
		return nil, err
	}

	res.SecurityInfo = result.SecurityInfo
	return
}

// GetLoginHistory 获取登录历史
func (c *profileController) GetLoginHistory(ctx context.Context, req *profileApi.GetLoginHistoryReq) (res *profileApi.GetLoginHistoryRes, err error) {
	res = new(profileApi.GetLoginHistoryRes)

	result, err := service.Profile().GetLoginHistory(ctx, req)
	if err != nil {
		return nil, err
	}

	*res = *result
	return
}

// ============================================================================
// 个人 API Key 管理
// ============================================================================

// GetMyApiKeys 获取我的API Key列表
func (c *profileController) GetMyApiKeys(ctx context.Context, req *apikeyApi.GetMyApiKeysReq) (res *apikeyApi.GetMyApiKeysRes, err error) {
	res = new(apikeyApi.GetMyApiKeysRes)

	result, err := service.ApiKey().GetMyApiKeys(ctx, req)
	if err != nil {
		return nil, err
	}

	*res = *result
	return
}

// CreateMyApiKey 创建我的API Key
func (c *profileController) CreateMyApiKey(ctx context.Context, req *apikeyApi.CreateMyApiKeyReq) (res *apikeyApi.CreateMyApiKeyRes, err error) {
	res = new(apikeyApi.CreateMyApiKeyRes)

	result, err := service.ApiKey().CreateMyApiKey(ctx, req)
	if err != nil {
		return nil, err
	}

	*res = *result
	return
}

// UpdateMyApiKey 更新我的API Key
func (c *profileController) UpdateMyApiKey(ctx context.Context, req *apikeyApi.UpdateMyApiKeyReq) (res *apikeyApi.UpdateMyApiKeyRes, err error) {
	res = new(apikeyApi.UpdateMyApiKeyRes)

	_, err = service.ApiKey().UpdateMyApiKey(ctx, req)
	if err != nil {
		return nil, err
	}

	return
}

// DeleteMyApiKey 删除我的API Key
func (c *profileController) DeleteMyApiKey(ctx context.Context, req *apikeyApi.DeleteMyApiKeyReq) (res *apikeyApi.DeleteMyApiKeyRes, err error) {
	res = new(apikeyApi.DeleteMyApiKeyRes)

	_, err = service.ApiKey().DeleteMyApiKey(ctx, req)
	if err != nil {
		return nil, err
	}

	return
}

// RegenerateMyApiKey 重新生成我的API Key Secret
func (c *profileController) RegenerateMyApiKey(ctx context.Context, req *apikeyApi.RegenerateMyApiKeyReq) (res *apikeyApi.RegenerateMyApiKeyRes, err error) {
	res = new(apikeyApi.RegenerateMyApiKeyRes)

	result, err := service.ApiKey().RegenerateMyApiKey(ctx, req)
	if err != nil {
		return nil, err
	}

	*res = *result
	return
}