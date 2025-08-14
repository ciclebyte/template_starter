// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TemplateTags is the golang structure of table template_tags for DAO operations like Where/Data.
type TemplateTags struct {
	g.Meta     `orm:"table:template_tags, do:true"`
	Id         interface{} // 关联ID
	TemplateId interface{} // 模板ID
	TagId      interface{} // 标签ID
	CreatedAt  *gtime.Time // 创建时间
	UpdatedAt  *gtime.Time // 更新时间
}