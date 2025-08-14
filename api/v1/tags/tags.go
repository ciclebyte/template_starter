package tags

import (
	commonApi "github.com/ciclebyte/template_starter/api/v1/common"
	model "github.com/ciclebyte/template_starter/internal/model"
	"github.com/gogf/gf/v2/frame/g"
)

// TagsAddReq 新增标签请求
type TagsAddReq struct {
	g.Meta      `path:"/tags/add" method:"post" tags:"标签" summary:"标签-新增"`
	Name        string `json:"name" v:"required|length:1,50#标签名称不能为空|标签名称长度为1-50个字符"`
	Description string `json:"description" v:"length:0,500#标签描述长度不能超过500个字符"`
	Sort        int    `json:"sort" v:"min:0#排序权重不能小于0"`
}

type TagsAddRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

// TagsDelReq 删除标签请求
type TagsDelReq struct {
	g.Meta `path:"/tags/del" method:"delete" tags:"标签" summary:"标签-删除"`
	Id     interface{} `json:"id" v:"required#标签ID不能为空"`
}

type TagsDelRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

// TagsBatchDelReq 批量删除标签请求
type TagsBatchDelReq struct {
	g.Meta `path:"/tags/batchdel" method:"delete" tags:"标签" summary:"标签-批量删除"`
	Ids    []interface{} `json:"ids" v:"required#标签ID列表不能为空"`
}

type TagsBatchDelRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

// TagsEditReq 修改标签请求
type TagsEditReq struct {
	g.Meta      `path:"/tags/edit" method:"put" tags:"标签" summary:"标签-修改"`
	Id          interface{} `json:"id" v:"required#标签ID不能为空"`
	Name        string      `json:"name" v:"required|length:1,50#标签名称不能为空|标签名称长度为1-50个字符"`
	Description string      `json:"description" v:"length:0,500#标签描述长度不能超过500个字符"`
	Sort        int         `json:"sort" v:"min:0#排序权重不能小于0"`
}

type TagsEditRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

// TagsListReq 标签列表请求
type TagsListReq struct {
	g.Meta `path:"/tags/list" method:"get" tags:"标签" summary:"标签-列表"`
	commonApi.PageReq
	Name string `json:"name"` // 按标签名称搜索
}

type TagsListRes struct {
	g.Meta   `mime:"application/json" example:"string"`
	TagsList []*model.TagsInfo `json:"tagsList"`
	commonApi.ListRes
}

// TagsDetailReq 标签详情请求
type TagsDetailReq struct {
	g.Meta `path:"/tags/detail" method:"get" tags:"标签" summary:"标签-详情"`
	Id     interface{} `json:"id" v:"required#标签ID不能为空"`
}

type TagsDetailRes struct {
	g.Meta   `mime:"application/json" example:"string"`
	TagsInfo *model.TagsInfo `json:"data"`
}

// TagsAllReq 获取所有标签请求（不分页）
type TagsAllReq struct {
	g.Meta `path:"/tags/all" method:"get" tags:"标签" summary:"标签-全部"`
}

type TagsAllRes struct {
	g.Meta   `mime:"application/json" example:"string"`
	TagsList []*model.TagsInfo `json:"tagsList"`
}

// TemplateTagsAddReq 为模板添加标签请求
type TemplateTagsAddReq struct {
	g.Meta     `path:"/templates/{templateId}/tags/add" method:"post" tags:"模板标签" summary:"模板-添加标签"`
	TemplateId interface{} `json:"templateId" v:"required#模板ID不能为空"`
	TagIds     []int64     `json:"tagIds" v:"required|min:1#标签ID列表不能为空|至少选择一个标签"`
}

type TemplateTagsAddRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

// TemplateTagsRemoveReq 为模板移除标签请求
type TemplateTagsRemoveReq struct {
	g.Meta     `path:"/templates/{templateId}/tags/remove" method:"delete" tags:"模板标签" summary:"模板-移除标签"`
	TemplateId interface{} `json:"templateId" v:"required#模板ID不能为空"`
	TagIds     []int64     `json:"tagIds" v:"required|min:1#标签ID列表不能为空|至少选择一个标签"`
}

type TemplateTagsRemoveRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

// TemplateTagsSetReq 设置模板标签请求（批量设置，覆盖原有标签）
type TemplateTagsSetReq struct {
	g.Meta     `path:"/templates/{templateId}/tags/set" method:"put" tags:"模板标签" summary:"模板-设置标签"`
	TemplateId interface{} `json:"templateId" v:"required#模板ID不能为空"`
	TagIds     []int64     `json:"tagIds"` // 标签ID列表，空数组表示清空所有标签
}

type TemplateTagsSetRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

// TemplateTagsListReq 获取模板标签列表请求
type TemplateTagsListReq struct {
	g.Meta     `path:"/templates/{templateId}/tags" method:"get" tags:"模板标签" summary:"模板-标签列表"`
	TemplateId interface{} `json:"templateId" v:"required#模板ID不能为空"`
}

type TemplateTagsListRes struct {
	g.Meta   `mime:"application/json" example:"string"`
	TagsList []*model.TagsInfo `json:"tagsList"`
}

// TemplatesByTagReq 根据标签获取模板列表请求
type TemplatesByTagReq struct {
	g.Meta `path:"/tags/{tagId}/templates" method:"get" tags:"标签模板" summary:"标签-模板列表"`
	commonApi.PageReq
	TagId interface{} `json:"tagId" v:"required#标签ID不能为空"`
}

type TemplatesByTagRes struct {
	g.Meta        `mime:"application/json" example:"string"`
	TemplatesList []*model.TemplatesInfo `json:"templatesList"`
	commonApi.ListRes
}

// TagsWithCountReq 获取标签及关联模板数量请求
type TagsWithCountReq struct {
	g.Meta `path:"/tags/with-count" method:"get" tags:"标签" summary:"标签-带数量统计"`
}

type TagsWithCountRes struct {
	g.Meta   `mime:"application/json" example:"string"`
	TagsList []*model.TagWithTemplateCount `json:"tagsList"`
}