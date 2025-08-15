package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TemplateVariablePresets 模板订阅预设变量关联表
type TemplateVariablePresets struct {
	Id         uint64      `json:"id"          orm:"id,primary"    description:"模板变量预设ID"`   // 模板变量预设ID，自增主键
	TemplateId uint64      `json:"template_id" orm:"template_id"   description:"模板ID"`        // 所属模板ID
	PresetId   uint64      `json:"preset_id"   orm:"preset_id"     description:"预设ID"`        // 预设ID
	CreatedAt  *gtime.Time `json:"created_at"  orm:"created_at"    description:"创建时间"`       // 创建时间
	UpdatedAt  *gtime.Time `json:"updated_at"  orm:"updated_at"    description:"更新时间"`       // 更新时间
}