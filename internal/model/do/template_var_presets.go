// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TemplateVarPresets is the golang structure of table template_var_presets for DAO operations like Where/Data.
type TemplateVarPresets struct {
	g.Meta     `orm:"table:template_var_presets, do:true"`
	Id         interface{} // 模板变量预设ID
	TemplateId interface{} // 模板ID
	PresetId   interface{} // 预设ID
	CreatedAt  *gtime.Time // 创建时间
	UpdatedAt  *gtime.Time // 更新时间
}