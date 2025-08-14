// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TemplateExposeFields is the golang structure of table template_expose_fields for DAO operations like Where/Data.
type TemplateExposeFields struct {
	g.Meta          `orm:"table:template_expose_fields, do:true"`
	Id              interface{} // 暴露字段ID
	TemplateId      interface{} // 模板ID
	FieldSchemaJson interface{} // 字段结构定义（支持树形嵌套）
	Version         interface{} // 版本号
	Description     interface{} // 字段说明文档
	CreatedAt       *gtime.Time // 创建时间
	UpdatedAt       *gtime.Time // 更新时间
}