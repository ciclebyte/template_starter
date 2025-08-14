package template_expose

import (
	"github.com/gogf/gf/v2/frame/g"
)

// 模板暴露字段-获取
type TemplateExposeGetReq struct {
	g.Meta     `path:"/templates/{templateId}/expose" method:"get" tags:"模板暴露字段" summary:"模板暴露字段-获取"`
	TemplateId int64  `json:"templateId" v:"required|min:1#模板ID不能为空"`
	Version    string `json:"version"`
}

type TemplateExposeGetRes struct {
	TemplateExpose *TemplateExposeInfo `json:"templateExpose"`
}

type TemplateExposeInfo struct {
	Id              int64  `json:"id"`
	TemplateId      int64  `json:"templateId"`
	FieldSchemaJson string `json:"fieldSchemaJson"`
	Version         string `json:"version"`
	Description     string `json:"description"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

// 模板暴露字段-设置
type TemplateExposeSetReq struct {
	g.Meta          `path:"/templates/{templateId}/expose" method:"put" tags:"模板暴露字段" summary:"模板暴露字段-设置"`
	TemplateId      int64  `json:"templateId" v:"required|min:1#模板ID不能为空"`
	FieldSchemaJson string `json:"fieldSchemaJson" v:"required#字段结构定义不能为空"`
	Version         string `json:"version" v:"length:1,20#版本号长度为1-20个字符"`
	Description     string `json:"description" v:"length:0,1000#描述长度不能超过1000个字符"`
}

type TemplateExposeSetRes struct{}

// 模板暴露字段-删除
type TemplateExposeDelReq struct {
	g.Meta     `path:"/templates/{templateId}/expose" method:"delete" tags:"模板暴露字段" summary:"模板暴露字段-删除"`
	TemplateId int64  `json:"templateId" v:"required|min:1#模板ID不能为空"`
	Version    string `json:"version"`
}

type TemplateExposeDelRes struct{}

// 模板暴露字段-历史版本列表
type TemplateExposeVersionsReq struct {
	g.Meta     `path:"/templates/{templateId}/expose/versions" method:"get" tags:"模板暴露字段" summary:"模板暴露字段-历史版本"`
	TemplateId int64 `json:"templateId" v:"required|min:1#模板ID不能为空"`
}

type TemplateExposeVersionsRes struct {
	Versions []*TemplateExposeVersionInfo `json:"versions"`
}

type TemplateExposeVersionInfo struct {
	Version     string `json:"version"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}