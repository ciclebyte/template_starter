// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserRoles is the golang structure for table user_roles.
type UserRoles struct {
	Id         int64       `json:"id"         description:""`
	UserId     int64       `json:"userId"     description:"用户ID"`
	RoleId     int64       `json:"roleId"     description:"角色ID"`
	AssignedBy int64       `json:"assignedBy" description:"分配者ID"`
	AssignedAt *gtime.Time `json:"assignedAt" description:"分配时间"`
	ExpiresAt  *gtime.Time `json:"expiresAt"  description:"过期时间"`
}
