package model

import "github.com/gogf/gf/v2/os/gtime"

// TagsInfo 标签信息结构体
type TagsInfo struct {
	Id          int64       `orm:"id"  json:"id"`                    // 标签ID
	Name        string      `orm:"name"  json:"name"`                // 标签名称
	Description string      `orm:"description"  json:"description"`  // 标签描述
	Sort        int         `orm:"sort"  json:"sort"`                // 排序权重
	CreatedAt   *gtime.Time `orm:"created_at"  json:"createdAt"`     // 创建时间
	UpdatedAt   *gtime.Time `orm:"updated_at"  json:"updatedAt"`     // 更新时间
	DeletedAt   *gtime.Time `orm:"deleted_at"  json:"deletedAt"`     // 删除时间
}

// TemplateTagsInfo 模板标签关联信息结构体
type TemplateTagsInfo struct {
	Id         int64       `orm:"id"  json:"id"`                 // 关联ID
	TemplateId int64       `orm:"template_id"  json:"templateId"` // 模板ID
	TagId      int64       `orm:"tag_id"  json:"tagId"`         // 标签ID
	CreatedAt  *gtime.Time `orm:"created_at"  json:"createdAt"` // 创建时间
	UpdatedAt  *gtime.Time `orm:"updated_at"  json:"updatedAt"` // 更新时间
}

// TagWithTemplateCount 标签信息（包含模板数量）
type TagWithTemplateCount struct {
	TagsInfo
	TemplateCount int `json:"templateCount"` // 关联模板数量
}

// TemplateWithTags 模板信息（包含标签列表）
type TemplateWithTags struct {
	TemplateId   int64      `json:"templateId"`   // 模板ID
	TemplateName string     `json:"templateName"` // 模板名称
	Tags         []TagsInfo `json:"tags"`         // 关联的标签列表
}