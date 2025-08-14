package var_preset

import (
	"github.com/gogf/gf/v2/frame/g"
)

// 变量预设-新增
type VarPresetAddReq struct {
	g.Meta         `path:"/var-preset/add" method:"post" tags:"变量预设" summary:"变量预设-新增"`
	Name           string `json:"name" v:"required|length:1,50#预设名称不能为空|预设名称长度为1-50个字符"`
	DisplayName    string `json:"displayName" v:"required|length:1,100#显示名称不能为空|显示名称长度为1-100个字符"`
	Description    string `json:"description" v:"length:0,500#描述长度不能超过500个字符"`
	Category       string `json:"category" v:"required|in:system,custom#分类不能为空|分类必须为system或custom"`
	SchemaJson     string `json:"schemaJson" v:"required#数据结构模板不能为空"`
	DefaultDataJson string `json:"defaultDataJson"`
	Icon           string `json:"icon" v:"length:0,50#图标名称长度不能超过50个字符"`
	Sort           int    `json:"sort" v:"min:0#排序权重不能小于0"`
	Version        string `json:"version" v:"length:1,20#版本号长度为1-20个字符"`
}

type VarPresetAddRes struct{}

// 变量预设-批量删除
type VarPresetBatchDelReq struct {
	g.Meta `path:"/var-preset/batchdel" method:"delete" tags:"变量预设" summary:"变量预设-批量删除"`
	Ids    []int64 `json:"ids" v:"required|min-length:1#预设ID列表不能为空"`
}

type VarPresetBatchDelRes struct{}

// 变量预设-删除
type VarPresetDelReq struct {
	g.Meta `path:"/var-preset/del" method:"delete" tags:"变量预设" summary:"变量预设-删除"`
	Id     int64 `json:"id" v:"required|min:1#预设ID不能为空"`
}

type VarPresetDelRes struct{}

// 变量预设-详情
type VarPresetDetailReq struct {
	g.Meta `path:"/var-preset/detail" method:"get" tags:"变量预设" summary:"变量预设-详情"`
	Id     int64 `json:"id" v:"required|min:1#预设ID不能为空"`
}

type VarPresetDetailRes struct {
	VarPreset *VarPresetDetailInfo `json:"varPreset"`
}

type VarPresetDetailInfo struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	DisplayName     string `json:"displayName"`
	Description     string `json:"description"`
	Category        string `json:"category"`
	SchemaJson      string `json:"schemaJson"`
	DefaultDataJson string `json:"defaultDataJson"`
	Icon            string `json:"icon"`
	Sort            int    `json:"sort"`
	IsEnabled       int    `json:"isEnabled"`
	Version         string `json:"version"`
	CreatedBy       int64  `json:"createdBy"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

// 变量预设-修改
type VarPresetEditReq struct {
	g.Meta         `path:"/var-preset/edit" method:"put" tags:"变量预设" summary:"变量预设-修改"`
	Id             int64  `json:"id" v:"required|min:1#预设ID不能为空"`
	Name           string `json:"name" v:"required|length:1,50#预设名称不能为空|预设名称长度为1-50个字符"`
	DisplayName    string `json:"displayName" v:"required|length:1,100#显示名称不能为空|显示名称长度为1-100个字符"`
	Description    string `json:"description" v:"length:0,500#描述长度不能超过500个字符"`
	Category       string `json:"category" v:"required|in:system,custom#分类不能为空|分类必须为system或custom"`
	SchemaJson     string `json:"schemaJson" v:"required#数据结构模板不能为空"`
	DefaultDataJson string `json:"defaultDataJson"`
	Icon           string `json:"icon" v:"length:0,50#图标名称长度不能超过50个字符"`
	Sort           int    `json:"sort" v:"min:0#排序权重不能小于0"`
	Version        string `json:"version" v:"length:1,20#版本号长度为1-20个字符"`
	IsEnabled      int    `json:"isEnabled" v:"in:0,1#启用状态必须为0或1"`
}

type VarPresetEditRes struct{}

// 变量预设-列表（分页）
type VarPresetListReq struct {
	g.Meta   `path:"/var-preset/list" method:"get" tags:"变量预设" summary:"变量预设-列表"`
	PageNum  int    `json:"pageNum" v:"min:1#页码必须大于0"`
	PageSize int    `json:"pageSize" v:"min:1#每页数量必须大于0"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

type VarPresetListRes struct {
	Total          int64              `json:"total"`
	VarPresetsList []*VarPresetListInfo `json:"varPresetsList"`
}

type VarPresetListInfo struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	DisplayName     string `json:"displayName"`
	Description     string `json:"description"`
	Category        string `json:"category"`
	Icon            string `json:"icon"`
	Sort            int    `json:"sort"`
	IsEnabled       int    `json:"isEnabled"`
	Version         string `json:"version"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

// 变量预设-全部（不分页）
type VarPresetAllReq struct {
	g.Meta   `path:"/var-preset/all" method:"get" tags:"变量预设" summary:"变量预设-全部"`
	Category string `json:"category"`
	IsEnabled int   `json:"isEnabled"`
}

type VarPresetAllRes struct {
	VarPresetsList []*VarPresetListInfo `json:"varPresetsList"`
}

// 变量预设-启用/禁用
type VarPresetToggleReq struct {
	g.Meta    `path:"/var-preset/toggle" method:"put" tags:"变量预设" summary:"变量预设-启用/禁用"`
	Id        int64 `json:"id" v:"required|min:1#预设ID不能为空"`
	IsEnabled int   `json:"isEnabled" v:"in:0,1#启用状态必须为0或1"`
}

type VarPresetToggleRes struct{}

// 模板变量预设-获取模板关联的预设列表
type TemplateVarPresetsReq struct {
	g.Meta     `path:"/templates/{templateId}/var-presets" method:"get" tags:"变量预设" summary:"模板变量预设-获取列表"`
	TemplateId int64 `json:"templateId" v:"required|min:1#模板ID不能为空"`
}

type TemplateVarPresetsRes struct {
	VarPresetsList []*VarPresetListInfo `json:"varPresetsList"`
}

// 模板变量预设-设置模板关联的预设
type TemplateVarPresetsSetReq struct {
	g.Meta     `path:"/templates/{templateId}/var-presets/set" method:"put" tags:"变量预设" summary:"模板变量预设-设置关联"`
	TemplateId int64   `json:"templateId" v:"required|min:1#模板ID不能为空"`
	PresetIds  []int64 `json:"presetIds"`
}

type TemplateVarPresetsSetRes struct{}

// 模板变量预设-添加关联
type TemplateVarPresetsAddReq struct {
	g.Meta     `path:"/templates/{templateId}/var-presets/add" method:"post" tags:"变量预设" summary:"模板变量预设-添加关联"`
	TemplateId int64   `json:"templateId" v:"required|min:1#模板ID不能为空"`
	PresetIds  []int64 `json:"presetIds" v:"required|min-length:1#预设ID列表不能为空"`
}

type TemplateVarPresetsAddRes struct{}

// 模板变量预设-移除关联
type TemplateVarPresetsRemoveReq struct {
	g.Meta     `path:"/templates/{templateId}/var-presets/remove" method:"delete" tags:"变量预设" summary:"模板变量预设-移除关联"`
	TemplateId int64   `json:"templateId" v:"required|min:1#模板ID不能为空"`
	PresetIds  []int64 `json:"presetIds" v:"required|min-length:1#预设ID列表不能为空"`
}

type TemplateVarPresetsRemoveRes struct{}