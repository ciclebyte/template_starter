// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TemplateVarPresets is the golang structure for table template_var_presets.
type TemplateVarPresets struct {
	Id         uint64      `json:"id"         description:"模板变量预设ID"`
	TemplateId uint64      `json:"templateId" description:"模板ID"`
	PresetId   uint64      `json:"presetId"   description:"预设ID"`
	CreatedAt  *gtime.Time `json:"createdAt"  description:"创建时间"`
	UpdatedAt  *gtime.Time `json:"updatedAt"  description:"更新时间"`
}