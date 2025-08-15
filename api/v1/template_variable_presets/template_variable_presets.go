package template_variable_presets

import (
	"github.com/gogf/gf/v2/frame/g"
)

// 订阅预设变量请求
type SubscribePresetReq struct {
	g.Meta      `path:"/templates/{templateId}/preset-variables/subscribe" tags:"TemplateVariablePresets" method:"post" summary:"订阅预设变量"`
	TemplateId  uint64                    `json:"template_id" v:"required" dc:"模板ID"`
	Presets     []PresetSubscriptionItem  `json:"presets" v:"required" dc:"预设变量订阅列表"`
}

// 预设变量订阅项
type PresetSubscriptionItem struct {
	VarPresetId   uint64 `json:"var_preset_id" v:"required" dc:"预设变量ID"`
	PresetPath    string `json:"preset_path" v:"required" dc:"预设变量路径"`
	MappedName    string `json:"mapped_name" v:"required" dc:"映射变量名"`
	IsActive      int    `json:"is_active" dc:"是否启用"`
	Sort          int    `json:"sort" dc:"排序"`
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
	Id            int64  `json:"id" dc:"关联ID"`
	TemplateId    uint64 `json:"template_id" dc:"模板ID"`
	VarPresetId   uint64 `json:"var_preset_id" dc:"预设变量ID"`
	PresetName    string `json:"preset_name" dc:"预设变量名称"`
	PresetPath    string `json:"preset_path" dc:"预设变量路径"`
	MappedName    string `json:"mapped_name" dc:"映射变量名"`
	IsActive      int    `json:"is_active" dc:"是否启用"`
	Sort          int    `json:"sort" dc:"排序"`
	DisplayName   string `json:"display_name" dc:"显示名称"`
	Description   string `json:"description" dc:"描述"`
	Example       string `json:"example" dc:"示例"`
	InsertText    string `json:"insert_text" dc:"插入文本"`
	Category      string `json:"category" dc:"分类"`
}

// 取消订阅预设变量请求
type UnsubscribePresetReq struct {
	g.Meta     `path:"/templates/{templateId}/preset-variables/{id}" tags:"TemplateVariablePresets" method:"delete" summary:"取消订阅预设变量"`
	TemplateId uint64 `json:"template_id" v:"required" dc:"模板ID"`
	Id         int64  `json:"id" v:"required" dc:"关联ID"`
}

type UnsubscribePresetRes struct {
	g.Meta  `mime:"application/json"`
	Success bool   `json:"success" dc:"是否成功"`
	Message string `json:"message" dc:"消息"`
}

// 更新预设变量映射请求
type UpdatePresetMappingReq struct {
	g.Meta     `path:"/templates/{templateId}/preset-variables/{id}" tags:"TemplateVariablePresets" method:"put" summary:"更新预设变量映射"`
	TemplateId uint64 `json:"template_id" v:"required" dc:"模板ID"`
	Id         int64  `json:"id" v:"required" dc:"关联ID"`
	MappedName string `json:"mapped_name" v:"required" dc:"映射变量名"`
	IsActive   int    `json:"is_active" dc:"是否启用"`
	Sort       int    `json:"sort" dc:"排序"`
}

type UpdatePresetMappingRes struct {
	g.Meta  `mime:"application/json"`
	Success bool   `json:"success" dc:"是否成功"`
	Message string `json:"message" dc:"消息"`
}

// 获取可用预设变量列表请求 (用于订阅选择)
type GetAvailablePresetsReq struct {
	g.Meta `path:"/templates/preset-variables/available" tags:"TemplateVariablePresets" method:"get" summary:"获取可用预设变量列表"`
	Page   int    `json:"page" d:"1" v:"min:1" dc:"页码"`
	Size   int    `json:"size" d:"20" v:"min:1,max:100" dc:"每页数量"`
	Keyword string `json:"keyword" dc:"搜索关键词"`
}

type GetAvailablePresetsRes struct {
	g.Meta `mime:"application/json"`
	Data   AvailablePresetsData `json:"data" dc:"可用预设变量数据"`
}

// 可用预设变量数据
type AvailablePresetsData struct {
	List  []AvailablePresetItem `json:"list" dc:"预设变量列表"`
	Total int                   `json:"total" dc:"总数"`
	Page  int                   `json:"page" dc:"当前页"`
	Size  int                   `json:"size" dc:"每页数量"`
}

// 可用预设变量项
type AvailablePresetItem struct {
	Id          uint64                    `json:"id" dc:"预设变量ID"`
	Name        string                    `json:"name" dc:"预设变量名称"`
	Description string                    `json:"description" dc:"描述"`
	Variables   []AvailableVariableItem   `json:"variables" dc:"变量列表"`
}

// 可用变量项
type AvailableVariableItem struct {
	Path        string `json:"path" dc:"变量路径"`
	DisplayName string `json:"display_name" dc:"显示名称"`
	Description string `json:"description" dc:"描述"`
	Example     string `json:"example" dc:"示例"`
	InsertText  string `json:"insert_text" dc:"插入文本"`
	Category    string `json:"category" dc:"分类"`
}