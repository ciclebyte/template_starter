package template_variable_presets

import (
	"github.com/gogf/gf/v2/frame/g"
)

// 订阅预设变量请求
type SubscribePresetReq struct {
	g.Meta      `path:"/templates/{templateId}/preset-variables/subscribe" tags:"TemplateVariablePresets" method:"post" summary:"订阅预设变量"`
	TemplateId  uint64   `json:"template_id" v:"required" dc:"模板ID"`
	PresetIds   []uint64 `json:"preset_ids" v:"required" dc:"预设变量ID列表"`
}

type SubscribePresetRes struct {
	g.Meta  `mime:"application/json"`
	Success bool   `json:"success" dc:"是否成功"`
	Message string `json:"message" dc:"消息"`
}

// 获取模板订阅的预设变量列表请求
type GetSubscribedPresetsReq struct {
	g.Meta     `path:"/templates/{templateId}/preset-variables" tags:"TemplateVariablePresets" method:"get" summary:"获取模板订阅的预设变量"`
	TemplateId uint64 `json:"template_id" v:"required" dc:"模板ID"`
}

type GetSubscribedPresetsRes struct {
	g.Meta `mime:"application/json"`
	Data   []SubscribedPresetItem `json:"data" dc:"订阅的预设变量列表"`
}

// 订阅的预设变量项
type SubscribedPresetItem struct {
	Id          uint64 `json:"id" dc:"关联ID"`
	TemplateId  uint64 `json:"template_id" dc:"模板ID"`
	PresetId    uint64 `json:"preset_id" dc:"预设ID"`
	PresetName  string `json:"preset_name" dc:"预设变量名称"`
	Description string `json:"description" dc:"预设描述"`
}

// 取消订阅预设变量请求
type UnsubscribePresetReq struct {
	g.Meta     `path:"/templates/{templateId}/preset-variables/{id}" tags:"TemplateVariablePresets" method:"delete" summary:"取消订阅预设变量"`
	TemplateId uint64 `json:"template_id" v:"required" dc:"模板ID"`
	Id         uint64 `json:"id" v:"required" dc:"关联ID"`
}

type UnsubscribePresetRes struct {
	g.Meta  `mime:"application/json"`
	Success bool   `json:"success" dc:"是否成功"`
	Message string `json:"message" dc:"消息"`
}


// 获取可用预设变量列表请求 (用于订阅选择)
type GetAvailablePresetsReq struct {
	g.Meta   `path:"/templates/preset-variables/available" tags:"TemplateVariablePresets" method:"get" summary:"获取可用预设变量列表"`
	PageNum  int    `json:"pageNum" v:"min:1#页码必须大于0"`
	PageSize int    `json:"pageSize" v:"min:1#每页数量必须大于0"`
	Keyword  string `json:"keyword" dc:"搜索关键词"`
}

type GetAvailablePresetsRes struct {
	g.Meta   `mime:"application/json"`
	List     []AvailablePresetItem `json:"list" dc:"预设变量列表"`
	Total    int                   `json:"total" dc:"总数"`
	PageNum  int                   `json:"pageNum" dc:"当前页"`
	PageSize int                   `json:"pageSize" dc:"每页数量"`
}


// 可用预设变量项
type AvailablePresetItem struct {
	Id          uint64 `json:"id" dc:"预设变量ID"`
	Name        string `json:"name" dc:"预设变量名称"`
	Description string `json:"description" dc:"描述"`
}