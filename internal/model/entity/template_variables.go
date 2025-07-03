// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TemplateVariables is the golang structure for table template_variables.
type TemplateVariables struct {
	Id              int64       `json:"id"              description:"变量ID，自增主键"`
	TemplateId      int64       `json:"templateId"      description:"所属模板ID"`
	Name            string      `json:"name"            description:"变量名称"`
	Description     string      `json:"description"     description:"变量描述"`
	DefaultValue    string      `json:"defaultValue"    description:"变量默认值"`
	IsRequired      int         `json:"isRequired"      description:"是否为必填变量"`
	ValidationRegex string      `json:"validationRegex" description:"变量值验证正则表达式"`
	Sort            int         `json:"sort"            description:"排序"`
	CreatedAt       *gtime.Time `json:"createdAt"       description:"创建时间"`
}
