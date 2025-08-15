package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TemplateVariablePresets 模板订阅预设变量关联表
type TemplateVariablePresets struct {
	Id            int64       `json:"id"               orm:"id,primary"       description:"ID"`                      // 关联ID，自增主键
	TemplateId    uint64      `json:"template_id"      orm:"template_id"      description:"模板ID"`                   // 所属模板ID
	VarPresetId   uint64      `json:"var_preset_id"    orm:"var_preset_id"    description:"预设变量ID"`               // 预设变量ID
	PresetPath    string      `json:"preset_path"      orm:"preset_path"      description:"预设变量路径"`               // 预设中的变量路径 (如: user.name)
	MappedName    string      `json:"mapped_name"      orm:"mapped_name"      description:"映射到的模板变量名"`          // 映射到模板中的变量名
	IsActive      int         `json:"is_active"        orm:"is_active"        description:"是否启用"`                // 是否启用此映射关系
	Sort          int         `json:"sort"             orm:"sort"             description:"排序"`                   // 排序
	CreatedAt     *gtime.Time `json:"created_at"       orm:"created_at"       description:"创建时间"`                // 创建时间
	UpdatedAt     *gtime.Time `json:"updated_at"       orm:"updated_at"       description:"更新时间"`                // 更新时间
}