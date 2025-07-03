package model

type TemplateLanguagesInfo struct {
	Id         int64 `orm:"id"  json:"id"`                  // 关联ID，自增主键
	TemplateId int64 `orm:"template_id"  json:"templateId"` // 关联的模板ID
	LanguageId int   `orm:"language_id"  json:"languageId"` // 关联的语言ID
	IsPrimary  int   `orm:"is_primary"  json:"isPrimary"`   // 是否主要语言
}
