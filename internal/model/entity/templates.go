// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Templates is the golang structure for table templates.
type Templates struct {
	Id          int64       `json:"id"          description:"模板ID，自增主键"`
	Name        string      `json:"name"        description:"模板名称"`
	Description string      `json:"description" description:"模板详细描述"`
	CategoryId  uint        `json:"categoryId"  description:"所属分类ID"`
	IsFeatured  int         `json:"isFeatured"  description:"是否推荐模板"`
	Logo        string      `json:"logo"        description:"模板logo图片URL"`
	CreatedAt   *gtime.Time `json:"createdAt"   description:"记录创建时间"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   description:"记录最后更新时间"`
}
