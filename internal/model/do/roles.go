// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Roles is the golang structure of table roles for DAO operations like Where/Data.
type Roles struct {
	g.Meta         `orm:"table:roles, do:true"`
	Id             interface{} // 角色ID
	Name           interface{} // 角色名称
	Code           interface{} // 角色编码
	Description    interface{} // 角色描述
	IsSystem       interface{} // 是否系统角色
	OrganizationId interface{} // 所属组织ID（NULL表示系统级角色）
	Status         interface{} // 状态：0=禁用，1=正常
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
}
