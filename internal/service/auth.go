package service

import (
	"context"
	"github.com/ciclebyte/template_starter/library/libJWT"
)

// AuthUserInfo 认证用户信息
type AuthUserInfo struct {
	ID          int64    `json:"id"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	Nickname    string   `json:"nickname"`
	Avatar      string   `json:"avatar"`
	Status      int      `json:"status"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	LastLoginAt string   `json:"last_login_at"`
}

// RegisterReq 注册请求
type RegisterReq struct {
	Username string `json:"username" v:"required|length:3,20#请输入用户名|用户名长度为3-20位"`
	Email    string `json:"email" v:"required|email#请输入邮箱|邮箱格式不正确"`
	Password string `json:"password" v:"required|length:6,20#请输入密码|密码长度为6-20位"`
	Nickname string `json:"nickname" v:"length:2,20#昵称长度为2-20位"`
}

// RegisterRes 注册响应
type RegisterRes struct {
	*libJWT.TokenInfo
	User *AuthUserInfo `json:"user"`
}

// LoginReq 登录请求
type LoginReq struct {
	Username string `json:"username" v:"required#请输入用户名"`
	Password string `json:"password" v:"required#请输入密码"`
}

// LoginRes 登录响应
type LoginRes struct {
	*libJWT.TokenInfo
	User *AuthUserInfo `json:"user"`
}

// RefreshTokenReq 刷新Token请求
type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" v:"required#请提供刷新令牌"`
}

// RefreshTokenRes 刷新Token响应
type RefreshTokenRes struct {
	*libJWT.TokenInfo
}

// IAuth 认证服务接口
type IAuth interface {
	// 用户注册
	Register(ctx context.Context, req *RegisterReq) (*RegisterRes, error)
	
	// 用户登录
	Login(ctx context.Context, req *LoginReq) (*LoginRes, error)
	
	// 用户登出
	Logout(ctx context.Context) error
	
	// 获取当前用户信息
	GetCurrentUser(ctx context.Context) (*AuthUserInfo, error)
	
	// 刷新Token
	RefreshToken(ctx context.Context, req *RefreshTokenReq) (*RefreshTokenRes, error)
	
	// 检查用户权限
	HasPermission(ctx context.Context, userId int64, permission string) (bool, error)
	
	// 检查用户角色
	HasRole(ctx context.Context, userId int64, role string) (bool, error)
	
	// 获取用户权限列表
	GetUserPermissions(ctx context.Context, userId int64) ([]string, error)
	
	// 获取用户角色列表
	GetUserRoles(ctx context.Context, userId int64) ([]string, error)
}

var localAuth IAuth

func Auth() IAuth {
	if localAuth == nil {
		panic("implement not found for interface IAuth, forgot register?")
	}
	return localAuth
}

func RegisterAuth(i IAuth) {
	localAuth = i
}