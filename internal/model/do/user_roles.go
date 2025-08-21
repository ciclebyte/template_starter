// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserRoles is the golang structure of table user_roles for DAO operations like Where/Data.
type UserRoles struct {
	g.Meta     `orm:"table:user_roles, do:true"`
	Id         interface{} //
	UserId     interface{} // 用户ID
	RoleId     interface{} // 角色ID
	AssignedBy interface{} // 分配者ID
	AssignedAt *gtime.Time // 分配时间
	ExpiresAt  *gtime.Time // 过期时间
}
