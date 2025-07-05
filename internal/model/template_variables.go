package model

type TemplateVariablesInfo struct {
	Id              int64  `orm:"id"  json:"id"`                            // 变量ID，自增主键
	TemplateId      int64  `orm:"template_id"  json:"templateId"`           // 所属模板ID
	Name            string `orm:"name"  json:"name"`                        // 变量名称
	VariableType    string `orm:"variable_type"  json:"variableType"`       // 变量类型
	Description     string `orm:"description"  json:"description"`          // 变量描述
	DefaultValue    string `orm:"default_value"  json:"defaultValue"`       // 变量默认值
	IsRequired      int    `orm:"is_required"  json:"isRequired"`           // 是否为必填变量
	ValidationRegex string `orm:"validation_regex"  json:"validationRegex"` // 变量值验证正则表达式
	Sort            int    `orm:"sort"  json:"sort"`                        // 显示顺序
}
