package controller

import (
	"context"

	api "github.com/ciclebyte/template_starter/api/v1/auth"
	service "github.com/ciclebyte/template_starter/internal/service"
)

var Auth = authController{}

type authController struct {
	BaseController
}

// Register 用户注册
func (c *authController) Register(ctx context.Context, req *api.RegisterReq) (res *api.RegisterRes, err error) {
	res = new(api.RegisterRes)
	
	serviceReq := &service.RegisterReq{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Nickname: req.Nickname,
	}
	
	result, err := service.Auth().Register(ctx, serviceReq)
	if err != nil {
		return nil, err
	}
	
	res.RegisterRes = result
	return
}

// Login 用户登录
func (c *authController) Login(ctx context.Context, req *api.LoginReq) (res *api.LoginRes, err error) {
	res = new(api.LoginRes)
	
	serviceReq := &service.LoginReq{
		Username: req.Username,
		Password: req.Password,
	}
	
	result, err := service.Auth().Login(ctx, serviceReq)
	if err != nil {
		return nil, err
	}
	
	res.LoginRes = result
	return
}

// Logout 用户登出
func (c *authController) Logout(ctx context.Context, req *api.LogoutReq) (res *api.LogoutRes, err error) {
	res = new(api.LogoutRes)
	
	err = service.Auth().Logout(ctx)
	return
}

// Me 获取当前用户信息
func (c *authController) Me(ctx context.Context, req *api.UserInfoReq) (res *api.UserInfoRes, err error) {
	res = new(api.UserInfoRes)
	
	user, err := service.Auth().GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	
	res.User = user
	return
}

// RefreshToken 刷新令牌
func (c *authController) RefreshToken(ctx context.Context, req *api.RefreshTokenReq) (res *api.RefreshTokenRes, err error) {
	res = new(api.RefreshTokenRes)
	
	serviceReq := &service.RefreshTokenReq{
		RefreshToken: req.RefreshToken,
	}
	
	result, err := service.Auth().RefreshToken(ctx, serviceReq)
	if err != nil {
		return nil, err
	}
	
	res.RefreshTokenRes = result
	return
}

// CheckPermission 检查权限
func (c *authController) CheckPermission(ctx context.Context, req *api.CheckPermissionReq) (res *api.CheckPermissionRes, err error) {
	res = new(api.CheckPermissionRes)
	
	// 从上下文获取用户ID
	userInfo, err := service.Auth().GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	
	hasPermission, err := service.Auth().HasPermission(ctx, userInfo.ID, req.Permission)
	if err != nil {
		return nil, err
	}
	
	res.HasPermission = hasPermission
	return
}

// CheckRole 检查角色
func (c *authController) CheckRole(ctx context.Context, req *api.CheckRoleReq) (res *api.CheckRoleRes, err error) {
	res = new(api.CheckRoleRes)
	
	// 从上下文获取用户ID
	userInfo, err := service.Auth().GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	
	hasRole, err := service.Auth().HasRole(ctx, userInfo.ID, req.Role)
	if err != nil {
		return nil, err
	}
	
	res.HasRole = hasRole
	return
}