package controller

import (
	"context"

	api "github.com/ciclebyte/template_starter/api/v1/user"
	service "github.com/ciclebyte/template_starter/internal/service"
)

var User = userController{}

type userController struct {
	BaseController
}

// ListUsers 获取用户列表
func (c *userController) ListUsers(ctx context.Context, req *api.ListUsersReq) (res *api.ListUsersRes, err error) {
	res = new(api.ListUsersRes)

	result, err := service.User().ListUsers(ctx, req)
	if err != nil {
		return nil, err
	}

	*res = *result
	return
}

// CreateUser 创建用户
func (c *userController) CreateUser(ctx context.Context, req *api.CreateUserReq) (res *api.CreateUserRes, err error) {
	res = new(api.CreateUserRes)

	result, err := service.User().CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	res.Id = result.Id
	return
}

// UpdateUser 更新用户
func (c *userController) UpdateUser(ctx context.Context, req *api.UpdateUserReq) (res *api.UpdateUserRes, err error) {
	res = new(api.UpdateUserRes)

	_, err = service.User().UpdateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return
}

// DeleteUser 删除用户
func (c *userController) DeleteUser(ctx context.Context, req *api.DeleteUserReq) (res *api.DeleteUserRes, err error) {
	res = new(api.DeleteUserRes)

	_, err = service.User().DeleteUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return
}

// GetUser 获取用户详情
func (c *userController) GetUser(ctx context.Context, req *api.GetUserReq) (res *api.GetUserRes, err error) {
	res = new(api.GetUserRes)

	result, err := service.User().GetUser(ctx, req)
	if err != nil {
		return nil, err
	}

	res.UserInfo = result.UserInfo
	return
}

// ResetPassword 重置用户密码
func (c *userController) ResetPassword(ctx context.Context, req *api.ResetPasswordReq) (res *api.ResetPasswordRes, err error) {
	res = new(api.ResetPasswordRes)

	_, err = service.User().ResetPassword(ctx, req)
	if err != nil {
		return nil, err
	}

	return
}

// UpdateUserStatus 更新用户状态
func (c *userController) UpdateUserStatus(ctx context.Context, req *api.UpdateUserStatusReq) (res *api.UpdateUserStatusRes, err error) {
	res = new(api.UpdateUserStatusRes)

	_, err = service.User().UpdateUserStatus(ctx, req)
	if err != nil {
		return nil, err
	}

	return
}

// AssignUserRoles 分配用户角色
func (c *userController) AssignUserRoles(ctx context.Context, req *api.AssignUserRolesReq) (res *api.AssignUserRolesRes, err error) {
	res = new(api.AssignUserRolesRes)

	_, err = service.User().AssignUserRoles(ctx, req)
	if err != nil {
		return nil, err
	}

	return
}