package service

import (
	"context"

	"github.com/ciclebyte/template_starter/api/v1/profile"
)

// ============================================================================
// 用户个人资料服务接口
// ============================================================================

type (
	IProfile interface {
		// 个人资料管理
		GetProfile(ctx context.Context, req *profile.GetProfileReq) (*profile.GetProfileRes, error)
		UpdateProfile(ctx context.Context, req *profile.UpdateProfileReq) (*profile.UpdateProfileRes, error)
		ChangePassword(ctx context.Context, req *profile.ChangePasswordReq) (*profile.ChangePasswordRes, error)
		UpdateEmail(ctx context.Context, req *profile.UpdateEmailReq) (*profile.UpdateEmailRes, error)
		UploadAvatar(ctx context.Context, req *profile.UploadAvatarReq) (*profile.UploadAvatarRes, error)
		
		// 账户安全
		GetSecurityInfo(ctx context.Context, req *profile.GetSecurityInfoReq) (*profile.GetSecurityInfoRes, error)
		GetLoginHistory(ctx context.Context, req *profile.GetLoginHistoryReq) (*profile.GetLoginHistoryRes, error)
	}
)

var (
	localProfile IProfile
)

func Profile() IProfile {
	if localProfile == nil {
		panic("implement not found for interface IProfile, forgot register?")
	}
	return localProfile
}

func RegisterProfile(i IProfile) {
	localProfile = i
}