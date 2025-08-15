package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TemplateVariablePresets 模板订阅预设变量关联表
type TemplateVariablePresets struct {
	g.Meta     `orm:"table:template_variable_presets, do:true"`
	Id         interface{} // 模板变量预设ID，自增主键
	TemplateId interface{} // 所属模板ID
	PresetId   interface{} // 预设ID
	CreatedAt  *gtime.Time // 创建时间
	UpdatedAt  *gtime.Time // 更新时间
}