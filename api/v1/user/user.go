package user

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ============================================================================
// 用户管理相关API定义
// ============================================================================

// UserInfo 用户信息
type UserInfo struct {
	Id           int64       `json:"id" dc:"用户ID"`
	Username     string      `json:"username" dc:"用户名"`
	Email        string      `json:"email" dc:"邮箱"`
	Nickname     string      `json:"nickname" dc:"昵称"`
	Avatar       string      `json:"avatar" dc:"头像"`
	Phone        string      `json:"phone" dc:"手机号"`
	Status       int         `json:"status" dc:"状态：0=禁用，1=正常"`
	LastLoginAt  *gtime.Time `json:"lastLoginAt" dc:"最后登录时间"`
	CreatedAt    *gtime.Time `json:"createdAt" dc:"创建时间"`
	UpdatedAt    *gtime.Time `json:"updatedAt" dc:"更新时间"`
	Roles        []RoleInfo  `json:"roles" dc:"用户角色列表"`
}

// RoleInfo 角色信息（简化版）
type RoleInfo struct {
	Id          int64  `json:"id" dc:"角色ID"`
	Name        string `json:"name" dc:"角色名称"`
	Code        string `json:"code" dc:"角色编码"`
	Description string `json:"description" dc:"角色描述"`
}

// ListUsersReq 用户列表请求
type ListUsersReq struct {
	g.Meta `path:"/users" method:"get" summary:"获取用户列表" tags:"用户管理"`
	Page   int    `json:"page" v:"min:1" dc:"页码，默认1"`
	Size   int    `json:"size" v:"min:1|max:1000" dc:"每页数量，默认20"`
	Search string `json:"search" dc:"搜索关键词（用户名、邮箱、昵称）"`
	Status *int   `json:"status" dc:"状态过滤 0-禁用 1-启用"`
	Role   string `json:"role" dc:"角色过滤"`
}

type ListUsersRes struct {
	List  []UserInfo `json:"list" dc:"用户列表"`
	Total int        `json:"total" dc:"总数"`
	Page  int        `json:"page" dc:"当前页"`
	Size  int        `json:"size" dc:"每页数量"`
}

// CreateUserReq 创建用户请求
type CreateUserReq struct {
	g.Meta   `path:"/users" method:"post" summary:"创建用户" tags:"用户管理"`
	Username string `json:"username" v:"required|length:3,30" dc:"用户名"`
	Email    string `json:"email" v:"required|email" dc:"邮箱"`
	Password string `json:"password" v:"required|length:6,32" dc:"密码"`
	Nickname string `json:"nickname" v:"length:1,50" dc:"昵称"`
	Phone    string `json:"phone" v:"phone" dc:"手机号"`
	Status   int    `json:"status" v:"in:0,1" dc:"状态：0=禁用，1=正常"`
	Roles    []int64 `json:"roles" dc:"角色ID列表"`
}

type CreateUserRes struct {
	Id int64 `json:"id" dc:"用户ID"`
}

// UpdateUserReq 更新用户请求
type UpdateUserReq struct {
	g.Meta   `path:"/users/{id}" method:"put" summary:"更新用户" tags:"用户管理"`
	Id       int64   `json:"id" v:"required" dc:"用户ID"`
	Username string  `json:"username" v:"required|length:3,30" dc:"用户名"`
	Email    string  `json:"email" v:"required|email" dc:"邮箱"`
	Nickname string  `json:"nickname" v:"length:1,50" dc:"昵称"`
	Phone    string  `json:"phone" v:"phone" dc:"手机号"`
	Status   int     `json:"status" v:"in:0,1" dc:"状态：0=禁用，1=正常"`
	Roles    []int64 `json:"roles" dc:"角色ID列表"`
}

type UpdateUserRes struct{}

// DeleteUserReq 删除用户请求
type DeleteUserReq struct {
	g.Meta `path:"/users/{id}" method:"delete" summary:"删除用户" tags:"用户管理"`
	Id     int64 `json:"id" v:"required" dc:"用户ID"`
}

type DeleteUserRes struct{}

// GetUserReq 获取用户详情请求
type GetUserReq struct {
	g.Meta `path:"/users/{id}" method:"get" summary:"获取用户详情" tags:"用户管理"`
	Id     int64 `json:"id" v:"required" dc:"用户ID"`
}

type GetUserRes struct {
	UserInfo *UserInfo `json:"user" dc:"用户信息"`
}

// ResetPasswordReq 重置密码请求
type ResetPasswordReq struct {
	g.Meta      `path:"/users/{id}/reset-password" method:"post" summary:"重置用户密码" tags:"用户管理"`
	Id          int64  `json:"id" v:"required" dc:"用户ID"`
	NewPassword string `json:"newPassword" v:"required|length:6,32" dc:"新密码"`
}

type ResetPasswordRes struct{}

// UpdateUserStatusReq 更新用户状态请求
type UpdateUserStatusReq struct {
	g.Meta `path:"/users/{id}/status" method:"put" summary:"更新用户状态" tags:"用户管理"`
	Id     int64 `json:"id" v:"required" dc:"用户ID"`
	Status int   `json:"status" v:"in:0,1" dc:"状态：0=禁用，1=正常"`
}

type UpdateUserStatusRes struct{}

// AssignUserRolesReq 分配用户角色请求
type AssignUserRolesReq struct {
	g.Meta    `path:"/users/{id}/roles" method:"put" summary:"分配用户角色" tags:"用户管理"`
	Id        int64   `json:"id" v:"required" dc:"用户ID"`
	Roles     []int64 `json:"roles" v:"required" dc:"角色ID列表"`
	ExpiresAt string  `json:"expiresAt" dc:"过期时间（可选）"`
}

type AssignUserRolesRes struct{}