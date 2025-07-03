// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TemplateLanguages is the golang structure for table template_languages.
type TemplateLanguages struct {
	Id         int64       `json:"id"         description:"关联ID，自增主键"`
	TemplateId int64       `json:"templateId" description:"关联的模板ID"`
	LanguageId int         `json:"languageId" description:"关联的语言ID"`
	IsPrimary  int         `json:"isPrimary"  description:"是否主要语言"`
	CreatedAt  *gtime.Time `json:"createdAt"  description:"记录创建时间"`
}
