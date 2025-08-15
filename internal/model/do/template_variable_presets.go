package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TemplateVariablePresets 模板订阅预设变量关联表
type TemplateVariablePresets struct {
	g.Meta        `orm:"table:template_variable_presets, do:true"`
	Id            interface{} // 关联ID，自增主键
	TemplateId    interface{} // 所属模板ID
	VarPresetId   interface{} // 预设变量ID
	PresetPath    interface{} // 预设中的变量路径 (如: user.name)
	MappedName    interface{} // 映射到模板中的变量名
	IsActive      interface{} // 是否启用此映射关系
	Sort          interface{} // 排序
	CreatedAt     *gtime.Time // 创建时间
	UpdatedAt     *gtime.Time // 更新时间
}