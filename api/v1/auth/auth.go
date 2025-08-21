package auth

import (
	"github.com/ciclebyte/template_starter/internal/service"
	"github.com/gogf/gf/v2/frame/g"
)

// RegisterReq 用户注册请求
type RegisterReq struct {
	g.Meta   `path:"/auth/register" method:"post" tags:"认证" summary:"用户注册"`
	Username string `json:"username" v:"required|length:3,20#请输入用户名|用户名长度为3-20位"`
	Email    string `json:"email" v:"required|email#请输入邮箱|邮箱格式不正确"`
	Password string `json:"password" v:"required|length:6,20#请输入密码|密码长度为6-20位"`
	Nickname string `json:"nickname" v:"length:2,20#昵称长度为2-20位"`
}

type RegisterRes struct {
	g.Meta `mime:"application/json"`
	*service.RegisterRes
}

// LoginReq 用户登录请求
type LoginReq struct {
	g.Meta   `path:"/auth/login" method:"post" tags:"认证" summary:"用户登录"`
	Username string `json:"username" v:"required#请输入用户名"`
	Password string `json:"password" v:"required#请输入密码"`
}

type LoginRes struct {
	g.Meta `mime:"application/json"`
	*service.LoginRes
}

// LogoutReq 用户登出请求
type LogoutReq struct {
	g.Meta `path:"/auth/logout" method:"post" tags:"认证" summary:"用户登出"`
}

type LogoutRes struct {
	g.Meta `mime:"application/json"`
}

// UserInfoReq 用户信息请求
type UserInfoReq struct {
	g.Meta `path:"/auth/me" method:"get" tags:"认证" summary:"获取用户信息"`
}

type UserInfoRes struct {
	g.Meta `mime:"application/json"`
	User   *service.AuthUserInfo `json:"user"`
}

// RefreshTokenReq 刷新Token请求
type RefreshTokenReq struct {
	g.Meta       `path:"/auth/refresh" method:"post" tags:"认证" summary:"刷新Token"`
	RefreshToken string `json:"refresh_token" v:"required#请提供刷新令牌"`
}

type RefreshTokenRes struct {
	g.Meta `mime:"application/json"`
	*service.RefreshTokenRes
}

// CheckPermissionReq 权限检查请求
type CheckPermissionReq struct {
	g.Meta     `path:"/auth/check-permission" method:"post" tags:"认证" summary:"检查权限"`
	Permission string `json:"permission" v:"required#请输入权限代码"`
}

type CheckPermissionRes struct {
	g.Meta        `mime:"application/json"`
	HasPermission bool `json:"has_permission"`
}

// CheckRoleReq 角色检查请求
type CheckRoleReq struct {
	g.Meta `path:"/auth/check-role" method:"post" tags:"认证" summary:"检查角色"`
	Role   string `json:"role" v:"required#请输入角色代码"`
}

type CheckRoleRes struct {
	g.Meta  `mime:"application/json"`
	HasRole bool `json:"has_role"`
}