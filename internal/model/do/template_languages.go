// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TemplateLanguages is the golang structure of table template_languages for DAO operations like Where/Data.
type TemplateLanguages struct {
	g.Meta     `orm:"table:template_languages, do:true"`
	Id         interface{} // 关联ID，自增主键
	TemplateId interface{} // 关联的模板ID
	LanguageId interface{} // 关联的语言ID
	IsPrimary  interface{} // 是否主要语言
	CreatedAt  *gtime.Time // 记录创建时间
}
