// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Tags is the golang structure for table tags.
type Tags struct {
	Id          uint64      `json:"id"          description:"标签ID"`
	Name        string      `json:"name"        description:"标签名称"`
	Description string      `json:"description" description:"标签描述"`
	Sort        int         `json:"sort"        description:"排序权重"`
	CreatedAt   *gtime.Time `json:"createdAt"   description:"创建时间"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   description:"更新时间"`
	DeletedAt   *gtime.Time `json:"deletedAt"   description:"删除时间(软删除)"`
}