package service

import (
	"context"

	"github.com/ciclebyte/template_starter/api/v1/user"
)

// ============================================================================
// 用户管理服务接口
// ============================================================================

type (
	IUser interface {
		// 用户管理
		ListUsers(ctx context.Context, req *user.ListUsersReq) (*user.ListUsersRes, error)
		CreateUser(ctx context.Context, req *user.CreateUserReq) (*user.CreateUserRes, error)
		UpdateUser(ctx context.Context, req *user.UpdateUserReq) (*user.UpdateUserRes, error)
		DeleteUser(ctx context.Context, req *user.DeleteUserReq) (*user.DeleteUserRes, error)
		GetUser(ctx context.Context, req *user.GetUserReq) (*user.GetUserRes, error)
		ResetPassword(ctx context.Context, req *user.ResetPasswordReq) (*user.ResetPasswordRes, error)
		UpdateUserStatus(ctx context.Context, req *user.UpdateUserStatusReq) (*user.UpdateUserStatusRes, error)
		AssignUserRoles(ctx context.Context, req *user.AssignUserRolesReq) (*user.AssignUserRolesRes, error)
	}
)

var (
	localUser IUser
)

func User() IUser {
	if localUser == nil {
		panic("implement not found for interface IUser, forgot register?")
	}
	return localUser
}

func RegisterUser(i IUser) {
	localUser = i
}