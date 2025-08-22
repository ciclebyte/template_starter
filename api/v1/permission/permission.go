package permission

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ============================================================================
// 权限管理相关API定义
// ============================================================================

// ListPermissionsReq 权限列表请求
type ListPermissionsReq struct {
	g.Meta `path:"/permissions" method:"get" summary:"获取权限列表" tags:"权限管理"`
	Page   int    `json:"page" v:"min:1" dc:"页码，默认1"`
	Size   int    `json:"size" v:"min:1|max:1000" dc:"每页数量，默认20"`
	Search string `json:"search" dc:"搜索关键词"`
	Resource string `json:"resource" dc:"资源类型过滤"`
}

type ListPermissionsRes struct {
	List  []PermissionInfo `json:"list" dc:"权限列表"`
	Total int              `json:"total" dc:"总数"`
	Page  int              `json:"page" dc:"当前页"`
	Size  int              `json:"size" dc:"每页数量"`
}

// CreatePermissionReq 创建权限请求
type CreatePermissionReq struct {
	g.Meta      `path:"/permissions" method:"post" summary:"创建权限" tags:"权限管理"`
	Name        string `json:"name" v:"required|length:1,100" dc:"权限名称"`
	Code        string `json:"code" v:"required|length:1,100" dc:"权限代码"`
	Resource    string `json:"resource" v:"required|length:1,50" dc:"资源类型"`
	Action      string `json:"action" v:"required|length:1,50" dc:"操作类型"`
	Description string `json:"description" v:"length:0,500" dc:"权限描述"`
}

type CreatePermissionRes struct {
	Id int64 `json:"id" dc:"权限ID"`
}

// UpdatePermissionReq 更新权限请求
type UpdatePermissionReq struct {
	g.Meta      `path:"/permissions/{id}" method:"put" summary:"更新权限" tags:"权限管理"`
	Id          int64  `json:"id" v:"required|min:1" dc:"权限ID"`
	Name        string `json:"name" v:"required|length:1,100" dc:"权限名称"`
	Code        string `json:"code" v:"required|length:1,100" dc:"权限代码"`
	Resource    string `json:"resource" v:"required|length:1,50" dc:"资源类型"`
	Action      string `json:"action" v:"required|length:1,50" dc:"操作类型"`
	Description string `json:"description" v:"length:0,500" dc:"权限描述"`
}

type UpdatePermissionRes struct{}

// DeletePermissionReq 删除权限请求
type DeletePermissionReq struct {
	g.Meta `path:"/permissions/{id}" method:"delete" summary:"删除权限" tags:"权限管理"`
	Id     int64 `json:"id" v:"required|min:1" dc:"权限ID"`
}

type DeletePermissionRes struct{}

// GetPermissionReq 获取权限详情请求
type GetPermissionReq struct {
	g.Meta `path:"/permissions/{id}" method:"get" summary:"获取权限详情" tags:"权限管理"`
	Id     int64 `json:"id" v:"required|min:1" dc:"权限ID"`
}

type GetPermissionRes struct {
	*PermissionInfo
}

// ============================================================================
// 角色管理相关API定义
// ============================================================================

// ListRolesReq 角色列表请求
type ListRolesReq struct {
	g.Meta `path:"/roles" method:"get" summary:"获取角色列表" tags:"角色管理"`
	Page   int    `json:"page" v:"min:1" dc:"页码，默认1"`
	Size   int    `json:"size" v:"min:1|max:1000" dc:"每页数量，默认20"`
	Search string `json:"search" dc:"搜索关键词"`
	Status *int   `json:"status" dc:"状态过滤 0-禁用 1-启用"`
}

type ListRolesRes struct {
	List  []RoleInfo `json:"list" dc:"角色列表"`
	Total int        `json:"total" dc:"总数"`
	Page  int        `json:"page" dc:"当前页"`
	Size  int        `json:"size" dc:"每页数量"`
}

// CreateRoleReq 创建角色请求
type CreateRoleReq struct {
	g.Meta      `path:"/roles" method:"post" summary:"创建角色" tags:"角色管理"`
	Name        string  `json:"name" v:"required|length:1,100" dc:"角色名称"`
	Code        string  `json:"code" v:"required|length:1,100" dc:"角色代码"`
	Description string  `json:"description" v:"length:0,500" dc:"角色描述"`
	Status      int     `json:"status" v:"in:0,1" dc:"状态 0-禁用 1-启用"`
	Permissions []int64 `json:"permissions" dc:"权限ID列表"`
}

type CreateRoleRes struct {
	Id int64 `json:"id" dc:"角色ID"`
}

// UpdateRoleReq 更新角色请求
type UpdateRoleReq struct {
	g.Meta      `path:"/roles/{id}" method:"put" summary:"更新角色" tags:"角色管理"`
	Id          int64   `json:"id" v:"required|min:1" dc:"角色ID"`
	Name        string  `json:"name" v:"required|length:1,100" dc:"角色名称"`
	Code        string  `json:"code" v:"required|length:1,100" dc:"角色代码"`
	Description string  `json:"description" v:"length:0,500" dc:"角色描述"`
	Status      int     `json:"status" v:"in:0,1" dc:"状态 0-禁用 1-启用"`
	Permissions []int64 `json:"permissions" dc:"权限ID列表"`
}

type UpdateRoleRes struct{}

// DeleteRoleReq 删除角色请求
type DeleteRoleReq struct {
	g.Meta `path:"/roles/{id}" method:"delete" summary:"删除角色" tags:"角色管理"`
	Id     int64 `json:"id" v:"required|min:1" dc:"角色ID"`
}

type DeleteRoleRes struct{}

// GetRoleReq 获取角色详情请求
type GetRoleReq struct {
	g.Meta `path:"/roles/{id}" method:"get" summary:"获取角色详情" tags:"角色管理"`
	Id     int64 `json:"id" v:"required|min:1" dc:"角色ID"`
}

type GetRoleRes struct {
	*RoleInfo
}

// AssignRolePermissionsReq 分配角色权限请求
type AssignRolePermissionsReq struct {
	g.Meta      `path:"/roles/{id}/permissions" method:"post" summary:"分配角色权限" tags:"角色管理"`
	Id          int64   `json:"id" v:"required|min:1" dc:"角色ID"`
	Permissions []int64 `json:"permissions" v:"required" dc:"权限ID列表"`
}

type AssignRolePermissionsRes struct{}

// ============================================================================
// 用户角色管理相关API定义
// ============================================================================

// ListUserRolesReq 用户角色列表请求
type ListUserRolesReq struct {
	g.Meta `path:"/users/{userId}/roles" method:"get" summary:"获取用户角色列表" tags:"用户角色管理"`
	UserId int64 `json:"userId" v:"required|min:1" dc:"用户ID"`
}

type ListUserRolesRes struct {
	List []UserRoleInfo `json:"list" dc:"用户角色列表"`
}

// AssignUserRolesReq 分配用户角色请求
type AssignUserRolesReq struct {
	g.Meta    `path:"/users/{userId}/roles" method:"post" summary:"分配用户角色" tags:"用户角色管理"`
	UserId    int64   `json:"userId" v:"required|min:1" dc:"用户ID"`
	Roles     []int64 `json:"roles" v:"required" dc:"角色ID列表"`
	ExpiresAt string  `json:"expiresAt" dc:"过期时间，格式：2006-01-02 15:04:05"`
}

type AssignUserRolesRes struct{}

// RemoveUserRoleReq 移除用户角色请求
type RemoveUserRoleReq struct {
	g.Meta `path:"/users/{userId}/roles/{roleId}" method:"delete" summary:"移除用户角色" tags:"用户角色管理"`
	UserId int64 `json:"userId" v:"required|min:1" dc:"用户ID"`
	RoleId int64 `json:"roleId" v:"required|min:1" dc:"角色ID"`
}

type RemoveUserRoleRes struct{}

// ============================================================================
// 数据结构定义
// ============================================================================

// PermissionInfo 权限信息
type PermissionInfo struct {
	Id          int64       `json:"id" dc:"权限ID"`
	Name        string      `json:"name" dc:"权限名称"`
	Code        string      `json:"code" dc:"权限代码"`
	Resource    string      `json:"resource" dc:"资源类型"`
	Action      string      `json:"action" dc:"操作类型"`
	Description string      `json:"description" dc:"权限描述"`
	CreatedAt   *gtime.Time `json:"createdAt" dc:"创建时间"`
	UpdatedAt   *gtime.Time `json:"updatedAt" dc:"更新时间"`
}

// RoleInfo 角色信息
type RoleInfo struct {
	Id          int64            `json:"id" dc:"角色ID"`
	Name        string           `json:"name" dc:"角色名称"`
	Code        string           `json:"code" dc:"角色代码"`
	Description string           `json:"description" dc:"角色描述"`
	Status      int              `json:"status" dc:"状态 0-禁用 1-启用"`
	UserCount   int              `json:"userCount" dc:"用户数量"`
	Permissions []PermissionInfo `json:"permissions" dc:"权限列表"`
	CreatedAt   *gtime.Time      `json:"createdAt" dc:"创建时间"`
	UpdatedAt   *gtime.Time      `json:"updatedAt" dc:"更新时间"`
}

// UserRoleInfo 用户角色信息
type UserRoleInfo struct {
	Id         int64       `json:"id" dc:"用户角色ID"`
	UserId     int64       `json:"userId" dc:"用户ID"`
	RoleId     int64       `json:"roleId" dc:"角色ID"`
	RoleName   string      `json:"roleName" dc:"角色名称"`
	RoleCode   string      `json:"roleCode" dc:"角色代码"`
	AssignedBy int64       `json:"assignedBy" dc:"分配者ID"`
	ExpiresAt  *gtime.Time `json:"expiresAt" dc:"过期时间"`
	CreatedAt  *gtime.Time `json:"createdAt" dc:"创建时间"`
}