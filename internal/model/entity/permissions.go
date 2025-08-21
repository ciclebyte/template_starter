// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Permissions is the golang structure for table permissions.
type Permissions struct {
	Id          int64       `json:"id"          description:"权限ID"`
	Name        string      `json:"name"        description:"权限名称"`
	Code        string      `json:"code"        description:"权限编码"`
	Resource    string      `json:"resource"    description:"资源类型"`
	Action      string      `json:"action"      description:"操作类型"`
	Description string      `json:"description" description:"权限描述"`
	ParentId    int64       `json:"parentId"    description:"父权限ID"`
	SortOrder   int         `json:"sortOrder"   description:"排序"`
	CreatedAt   *gtime.Time `json:"createdAt"   description:""`
}
