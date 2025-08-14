// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TemplateTags is the golang structure for table template_tags.
type TemplateTags struct {
	Id         uint64      `json:"id"         description:"关联ID"`
	TemplateId uint64      `json:"templateId" description:"模板ID"`
	TagId      uint64      `json:"tagId"      description:"标签ID"`
	CreatedAt  *gtime.Time `json:"createdAt"  description:"创建时间"`
	UpdatedAt  *gtime.Time `json:"updatedAt"  description:"更新时间"`
}