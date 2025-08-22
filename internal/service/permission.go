package service

import (
	"context"

	"github.com/ciclebyte/template_starter/api/v1/permission"
)

// ============================================================================
// 权限管理服务接口
// ============================================================================

type (
	IPermission interface {
		// 权限管理
		ListPermissions(ctx context.Context, req *permission.ListPermissionsReq) (*permission.ListPermissionsRes, error)
		CreatePermission(ctx context.Context, req *permission.CreatePermissionReq) (*permission.CreatePermissionRes, error)
		UpdatePermission(ctx context.Context, req *permission.UpdatePermissionReq) (*permission.UpdatePermissionRes, error)
		DeletePermission(ctx context.Context, req *permission.DeletePermissionReq) (*permission.DeletePermissionRes, error)
		GetPermission(ctx context.Context, req *permission.GetPermissionReq) (*permission.GetPermissionRes, error)

		// 角色管理
		ListRoles(ctx context.Context, req *permission.ListRolesReq) (*permission.ListRolesRes, error)
		CreateRole(ctx context.Context, req *permission.CreateRoleReq) (*permission.CreateRoleRes, error)
		UpdateRole(ctx context.Context, req *permission.UpdateRoleReq) (*permission.UpdateRoleRes, error)
		DeleteRole(ctx context.Context, req *permission.DeleteRoleReq) (*permission.DeleteRoleRes, error)
		GetRole(ctx context.Context, req *permission.GetRoleReq) (*permission.GetRoleRes, error)
		AssignRolePermissions(ctx context.Context, req *permission.AssignRolePermissionsReq) (*permission.AssignRolePermissionsRes, error)

		// 用户角色管理
		ListUserRoles(ctx context.Context, req *permission.ListUserRolesReq) (*permission.ListUserRolesRes, error)
		AssignUserRoles(ctx context.Context, req *permission.AssignUserRolesReq) (*permission.AssignUserRolesRes, error)
		RemoveUserRole(ctx context.Context, req *permission.RemoveUserRoleReq) (*permission.RemoveUserRoleRes, error)
	}
)

var (
	localPermission IPermission
)

func Permission() IPermission {
	if localPermission == nil {
		panic("implement not found for interface IPermission, forgot register?")
	}
	return localPermission
}

func RegisterPermission(i IPermission) {
	localPermission = i
}