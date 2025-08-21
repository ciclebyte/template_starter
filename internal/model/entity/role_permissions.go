// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// RolePermissions is the golang structure for table role_permissions.
type RolePermissions struct {
	Id           int64       `json:"id"           description:""`
	RoleId       int64       `json:"roleId"       description:"角色ID"`
	PermissionId int64       `json:"permissionId" description:"权限ID"`
	CreatedAt    *gtime.Time `json:"createdAt"    description:""`
}
