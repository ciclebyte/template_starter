package controller

import (
	"context"

	api "github.com/ciclebyte/template_starter/api/v1/profile"
	service "github.com/ciclebyte/template_starter/internal/service"
)

var Profile = profileController{}

type profileController struct {
	BaseController
}

// GetProfile 获取个人资料
func (c *profileController) GetProfile(ctx context.Context, req *api.GetProfileReq) (res *api.GetProfileRes, err error) {
	res = new(api.GetProfileRes)

	result, err := service.Profile().GetProfile(ctx, req)
	if err != nil {
		return nil, err
	}

	*res = *result
	return
}

// UpdateProfile 更新个人资料
func (c *profileController) UpdateProfile(ctx context.Context, req *api.UpdateProfileReq) (res *api.UpdateProfileRes, err error) {
	res = new(api.UpdateProfileRes)

	_, err = service.Profile().UpdateProfile(ctx, req)
	if err != nil {
		return nil, err
	}

	return
}

// ChangePassword 修改密码
func (c *profileController) ChangePassword(ctx context.Context, req *api.ChangePasswordReq) (res *api.ChangePasswordRes, err error) {
	res = new(api.ChangePasswordRes)

	_, err = service.Profile().ChangePassword(ctx, req)
	if err != nil {
		return nil, err
	}

	return
}

// UpdateEmail 更新邮箱
func (c *profileController) UpdateEmail(ctx context.Context, req *api.UpdateEmailReq) (res *api.UpdateEmailRes, err error) {
	res = new(api.UpdateEmailRes)

	_, err = service.Profile().UpdateEmail(ctx, req)
	if err != nil {
		return nil, err
	}

	return
}

// UploadAvatar 上传头像
func (c *profileController) UploadAvatar(ctx context.Context, req *api.UploadAvatarReq) (res *api.UploadAvatarRes, err error) {
	res = new(api.UploadAvatarRes)

	result, err := service.Profile().UploadAvatar(ctx, req)
	if err != nil {
		return nil, err
	}

	res.AvatarUrl = result.AvatarUrl
	return
}

// GetSecurityInfo 获取账户安全信息
func (c *profileController) GetSecurityInfo(ctx context.Context, req *api.GetSecurityInfoReq) (res *api.GetSecurityInfoRes, err error) {
	res = new(api.GetSecurityInfoRes)

	result, err := service.Profile().GetSecurityInfo(ctx, req)
	if err != nil {
		return nil, err
	}

	res.SecurityInfo = result.SecurityInfo
	return
}

// GetLoginHistory 获取登录历史
func (c *profileController) GetLoginHistory(ctx context.Context, req *api.GetLoginHistoryReq) (res *api.GetLoginHistoryRes, err error) {
	res = new(api.GetLoginHistoryRes)

	result, err := service.Profile().GetLoginHistory(ctx, req)
	if err != nil {
		return nil, err
	}

	*res = *result
	return
}