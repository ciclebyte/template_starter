// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// RolePermissions is the golang structure of table role_permissions for DAO operations like Where/Data.
type RolePermissions struct {
	g.Meta       `orm:"table:role_permissions, do:true"`
	Id           interface{} //
	RoleId       interface{} // 角色ID
	PermissionId interface{} // 权限ID
	CreatedAt    *gtime.Time //
}
