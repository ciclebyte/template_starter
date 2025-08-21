// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Permissions is the golang structure of table permissions for DAO operations like Where/Data.
type Permissions struct {
	g.Meta      `orm:"table:permissions, do:true"`
	Id          interface{} // 权限ID
	Name        interface{} // 权限名称
	Code        interface{} // 权限编码
	Resource    interface{} // 资源类型
	Action      interface{} // 操作类型
	Description interface{} // 权限描述
	ParentId    interface{} // 父权限ID
	SortOrder   interface{} // 排序
	CreatedAt   *gtime.Time //
}
