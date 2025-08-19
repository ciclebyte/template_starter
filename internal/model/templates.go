package model

type TemplatesInfo struct {
	Id           int64                   `orm:"id"  json:"id"`                           // 模板ID，自增主键
	Name         string                  `orm:"name"  json:"name"`                       // 模板名称
	Description  string                  `orm:"description"  json:"description"`         // 模板详细描述
	Introduction string                  `orm:"introduction"  json:"introduction"`       // 模板详细介绍，支持Markdown格式
	CategoryId   int                     `orm:"category_id"  json:"categoryId"`          // 所属分类ID
	IsFeatured   int                     `orm:"is_featured"  json:"isFeatured"`          // 是否推荐模板
	TemplateType string                  `orm:"template_type"  json:"templateType"`      // 模板类型：basic=基础模板，scaffold=脚手架模板，data_driven=数据驱动模板
	TypeConfig   string                  `orm:"type_config"  json:"typeConfig"`          // 类型配置，JSON格式
	Logo         string                  `orm:"logo"  json:"logo"`                       // 模板logo图片URL
	Icon         string                  `orm:"icon"  json:"icon"`                       // 模板图标名称
	Languages    []TemplateLanguagesInfo `json:"languages"`                              // 模板支持的语言
}
