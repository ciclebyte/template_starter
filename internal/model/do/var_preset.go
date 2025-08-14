// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// VarPreset is the golang structure of table var_preset for DAO operations like Where/Data.
type VarPreset struct {
	g.Meta          `orm:"table:var_preset, do:true"`
	Id              interface{} // 预设ID
	Name            interface{} // 预设名称，如MySQL、Fields、OpenAPI
	DisplayName     interface{} // 显示名称
	Description     interface{} // 预设描述
	Category        interface{} // 分类：system=系统预置，custom=用户自定义
	SchemaJson      interface{} // 数据结构模板定义（JSON Schema）
	DefaultDataJson interface{} // 默认数据示例
	Icon            interface{} // 图标名称
	Sort            interface{} // 排序权重
	IsEnabled       interface{} // 是否启用
	Version         interface{} // 版本号
	CreatedBy       interface{} // 创建者ID（系统预置为空）
	CreatedAt       *gtime.Time // 创建时间
	UpdatedAt       *gtime.Time // 更新时间
}