// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TemplateVariables is the golang structure of table template_variables for DAO operations like Where/Data.
type TemplateVariables struct {
	g.Meta          `orm:"table:template_variables, do:true"`
	Id              interface{} // 变量ID，自增主键
	TemplateId      interface{} // 所属模板ID
	Name            interface{} // 变量名称
	VariableType    interface{} // 变量类型
	Description     interface{} // 变量描述
	DefaultValue    interface{} // 变量默认值
	IsRequired      interface{} // 是否为必填变量
	ValidationRegex interface{} // 变量值验证正则表达式
	Sort            interface{} // 排序
	CreatedAt       *gtime.Time // 创建时间
}
