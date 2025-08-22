package controller

import (
	"context"

	"github.com/ciclebyte/template_starter/api/v1/permission"
	"github.com/ciclebyte/template_starter/internal/service"
)

// Permission 权限管理控制器
var Permission = cPermission{}

type cPermission struct{}

// ============================================================================
// 权限管理
// ============================================================================

// ListPermissions 获取权限列表
func (c *cPermission) ListPermissions(ctx context.Context, req *permission.ListPermissionsReq) (res *permission.ListPermissionsRes, err error) {
	return service.Permission().ListPermissions(ctx, req)
}

// CreatePermission 创建权限
func (c *cPermission) CreatePermission(ctx context.Context, req *permission.CreatePermissionReq) (res *permission.CreatePermissionRes, err error) {
	return service.Permission().CreatePermission(ctx, req)
}

// UpdatePermission 更新权限
func (c *cPermission) UpdatePermission(ctx context.Context, req *permission.UpdatePermissionReq) (res *permission.UpdatePermissionRes, err error) {
	return service.Permission().UpdatePermission(ctx, req)
}

// DeletePermission 删除权限
func (c *cPermission) DeletePermission(ctx context.Context, req *permission.DeletePermissionReq) (res *permission.DeletePermissionRes, err error) {
	return service.Permission().DeletePermission(ctx, req)
}

// GetPermission 获取权限详情
func (c *cPermission) GetPermission(ctx context.Context, req *permission.GetPermissionReq) (res *permission.GetPermissionRes, err error) {
	return service.Permission().GetPermission(ctx, req)
}

// ============================================================================
// 角色管理
// ============================================================================

// ListRoles 获取角色列表
func (c *cPermission) ListRoles(ctx context.Context, req *permission.ListRolesReq) (res *permission.ListRolesRes, err error) {
	return service.Permission().ListRoles(ctx, req)
}

// CreateRole 创建角色
func (c *cPermission) CreateRole(ctx context.Context, req *permission.CreateRoleReq) (res *permission.CreateRoleRes, err error) {
	return service.Permission().CreateRole(ctx, req)
}

// UpdateRole 更新角色
func (c *cPermission) UpdateRole(ctx context.Context, req *permission.UpdateRoleReq) (res *permission.UpdateRoleRes, err error) {
	return service.Permission().UpdateRole(ctx, req)
}

// DeleteRole 删除角色
func (c *cPermission) DeleteRole(ctx context.Context, req *permission.DeleteRoleReq) (res *permission.DeleteRoleRes, err error) {
	return service.Permission().DeleteRole(ctx, req)
}

// GetRole 获取角色详情
func (c *cPermission) GetRole(ctx context.Context, req *permission.GetRoleReq) (res *permission.GetRoleRes, err error) {
	return service.Permission().GetRole(ctx, req)
}

// AssignRolePermissions 分配角色权限
func (c *cPermission) AssignRolePermissions(ctx context.Context, req *permission.AssignRolePermissionsReq) (res *permission.AssignRolePermissionsRes, err error) {
	return service.Permission().AssignRolePermissions(ctx, req)
}

// ============================================================================
// 用户角色管理
// ============================================================================

// ListUserRoles 获取用户角色列表
func (c *cPermission) ListUserRoles(ctx context.Context, req *permission.ListUserRolesReq) (res *permission.ListUserRolesRes, err error) {
	return service.Permission().ListUserRoles(ctx, req)
}

// AssignUserRoles 分配用户角色
func (c *cPermission) AssignUserRoles(ctx context.Context, req *permission.AssignUserRolesReq) (res *permission.AssignUserRolesRes, err error) {
	return service.Permission().AssignUserRoles(ctx, req)
}

// RemoveUserRole 移除用户角色
func (c *cPermission) RemoveUserRole(ctx context.Context, req *permission.RemoveUserRoleReq) (res *permission.RemoveUserRoleRes, err error) {
	return service.Permission().RemoveUserRole(ctx, req)
}