package model

type TemplatesInfo struct {
	Id           int64                   `orm:"id"  json:"id"`                     // 模板ID，自增主键
	Name         string                  `orm:"name"  json:"name"`                 // 模板名称
	Description  string                  `orm:"description"  json:"description"`   // 模板详细描述
	Introduction string                  `orm:"introduction"  json:"introduction"` // 模板详细介绍，支持Markdown格式
	CategoryId   int                     `orm:"category_id"  json:"categoryId"`    // 所属分类ID
	IsFeatured   int                     `orm:"is_featured"  json:"isFeatured"`    // 是否推荐模板
	Logo         string                  `orm:"logo"  json:"logo"`                 // 模板logo图片URL
	Languages    []TemplateLanguagesInfo `json:"languages"`                        // 模板支持的语言
}
