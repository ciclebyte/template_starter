// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SystemConfig is the golang structure for table system_config.
type SystemConfig struct {
	Id             int64       `json:"id"             description:"配置ID，自增主键"`
	ConfigKey      string      `json:"configKey"      description:"配置键名，如 ai.openai.api_key"`
	ConfigValue    string      `json:"configValue"    description:"配置值，支持JSON格式存储复杂配置"`
	ConfigGroup    string      `json:"configGroup"    description:"配置分组，如 ai, system, template, ui"`
	ConfigType     string      `json:"configType"     description:"配置类型：string, number, boolean, json, array"`
	DisplayName    string      `json:"displayName"    description:"显示名称，用于前端展示"`
	Description    string      `json:"description"    description:"配置描述和说明"`
	IsPublic       int         `json:"isPublic"       description:"是否为公开配置（前端可见）"`
	IsRequired     int         `json:"isRequired"     description:"是否为必填配置"`
	DefaultValue   string      `json:"defaultValue"   description:"默认值"`
	ValidationRule string      `json:"validationRule" description:"验证规则，JSON格式"`
	SortOrder      int         `json:"sortOrder"      description:"排序顺序"`
	Status         int         `json:"status"         description:"状态：1启用，0禁用"`
	CreatedAt      *gtime.Time `json:"createdAt"      description:"创建时间"`
	UpdatedAt      *gtime.Time `json:"updatedAt"      description:"更新时间"`
}
