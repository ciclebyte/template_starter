// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Roles is the golang structure for table roles.
type Roles struct {
	Id             int64       `json:"id"             description:"角色ID"`
	Name           string      `json:"name"           description:"角色名称"`
	Code           string      `json:"code"           description:"角色编码"`
	Description    string      `json:"description"    description:"角色描述"`
	IsSystem       int         `json:"isSystem"       description:"是否系统角色"`
	OrganizationId int64       `json:"organizationId" description:"所属组织ID（NULL表示系统级角色）"`
	Status         int         `json:"status"         description:"状态：0=禁用，1=正常"`
	CreatedAt      *gtime.Time `json:"createdAt"      description:""`
	UpdatedAt      *gtime.Time `json:"updatedAt"      description:""`
}
