// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TemplateExposeFields is the golang structure for table template_expose_fields.
type TemplateExposeFields struct {
	Id              uint64      `json:"id"              description:"暴露字段ID"`
	TemplateId      uint64      `json:"templateId"      description:"模板ID"`
	FieldSchemaJson string      `json:"fieldSchemaJson" description:"字段结构定义（支持树形嵌套）"`
	Version         string      `json:"version"         description:"版本号"`
	Description     *string     `json:"description"     description:"字段说明文档"`
	CreatedAt       *gtime.Time `json:"createdAt"       description:"创建时间"`
	UpdatedAt       *gtime.Time `json:"updatedAt"       description:"更新时间"`
}