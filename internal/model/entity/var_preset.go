// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// VarPreset is the golang structure for table var_preset.
type VarPreset struct {
	Id              uint64      `json:"id"              description:"预设ID"`
	Name            string      `json:"name"            description:"预设名称，如MySQL、Fields、OpenAPI"`
	DisplayName     string      `json:"displayName"     description:"显示名称"`
	Description     *string     `json:"description"     description:"预设描述"`
	Category        string      `json:"category"        description:"分类：system=系统预置，custom=用户自定义"`
	SchemaJson      string      `json:"schemaJson"      description:"数据结构模板定义（JSON Schema）"`
	DefaultDataJson *string     `json:"defaultDataJson" description:"默认数据示例"`
	Icon            *string     `json:"icon"            description:"图标名称"`
	Sort            int         `json:"sort"            description:"排序权重"`
	IsEnabled       int         `json:"isEnabled"       description:"是否启用"`
	Version         string      `json:"version"         description:"版本号"`
	CreatedBy       *uint64     `json:"createdBy"       description:"创建者ID（系统预置为空）"`
	CreatedAt       *gtime.Time `json:"createdAt"       description:"创建时间"`
	UpdatedAt       *gtime.Time `json:"updatedAt"       description:"更新时间"`
}