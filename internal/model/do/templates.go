// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Templates is the golang structure of table templates for DAO operations like Where/Data.
type Templates struct {
	g.Meta       `orm:"table:templates, do:true"`
	Id           interface{} // 模板ID，自增主键
	Name         interface{} // 模板名称
	Description  interface{} // 模板详细描述
	CategoryId   interface{} // 所属分类ID
	IsFeatured   interface{} // 是否推荐模板
	Logo         interface{} // 模板logo图片URL
	CreatedAt    *gtime.Time // 记录创建时间
	UpdatedAt    *gtime.Time // 记录最后更新时间
	Introduction interface{} // 模板详细介绍，支持Markdown格式
	Icon         interface{} // 模板图标名称
	TemplateType interface{} // 模板类型：basic=基础模板，scaffold=脚手架模板，data_driven=数据驱动模板
	TypeConfig   interface{} // 类型相关配置，JSON格式
}
