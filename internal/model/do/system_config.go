// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SystemConfig is the golang structure of table system_config for DAO operations like Where/Data.
type SystemConfig struct {
	g.Meta         `orm:"table:system_config, do:true"`
	Id             interface{} // 配置ID，自增主键
	ConfigKey      interface{} // 配置键名，如 ai.openai.api_key
	ConfigValue    interface{} // 配置值，支持JSON格式存储复杂配置
	ConfigGroup    interface{} // 配置分组，如 ai, system, template, ui
	ConfigType     interface{} // 配置类型：string, number, boolean, json, array
	DisplayName    interface{} // 显示名称，用于前端展示
	Description    interface{} // 配置描述和说明
	IsPublic       interface{} // 是否为公开配置（前端可见）
	IsRequired     interface{} // 是否为必填配置
	DefaultValue   interface{} // 默认值
	ValidationRule interface{} // 验证规则，JSON格式
	SortOrder      interface{} // 排序顺序
	Status         interface{} // 状态：1启用，0禁用
	CreatedAt      *gtime.Time // 创建时间
	UpdatedAt      *gtime.Time // 更新时间
}
