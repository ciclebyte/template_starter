package system_config

import (
	commonApi "github.com/ciclebyte/template_starter/api/v1/common"
	model "github.com/ciclebyte/template_starter/internal/model"
	"github.com/gogf/gf/v2/frame/g"
)

// 获取配置列表请求参数
type SystemConfigListReq struct {
	g.Meta `path:"/systemConfig/list" method:"get" tags:"系统配置" summary:"系统配置-列表"`
	commonApi.PageReq
	ConfigGroup string `json:"configGroup"` // 配置分组筛选
	IsPublic    *int   `json:"isPublic"`    // 是否公开配置筛选
	Status      *int   `json:"status"`      // 状态筛选
}

type SystemConfigListRes struct {
	g.Meta `mime:"application/json" example:"string"`
	commonApi.ListRes
	SystemConfigList []*model.SystemConfigInfo `json:"systemConfigList"`
}

// 获取配置详情请求参数
type SystemConfigDetailReq struct {
	g.Meta `path:"/systemConfig/detail" method:"get" tags:"系统配置" summary:"系统配置-详情"`
	Id     interface{} `json:"id" v:"required#配置ID不能为空"`
}

type SystemConfigDetailRes struct {
	g.Meta `mime:"application/json" example:"string"`
	*model.SystemConfigInfo
}

// 新增配置请求参数
type SystemConfigAddReq struct {
	g.Meta         `path:"/systemConfig/add" method:"post" tags:"系统配置" summary:"系统配置-新增"`
	ConfigKey      string `json:"configKey" v:"required#配置键名不能为空"`
	ConfigValue    string `json:"configValue"`
	ConfigGroup    string `json:"configGroup" v:"required#配置分组不能为空"`
	ConfigType     string `json:"configType" v:"required|in:string,number,boolean,json,array#配置类型不能为空|配置类型必须为string,number,boolean,json,array之一"`
	DisplayName    string `json:"displayName" v:"required#显示名称不能为空"`
	Description    string `json:"description"`
	IsPublic       int    `json:"isPublic" v:"in:0,1#是否公开必须为0或1"`
	IsRequired     int    `json:"isRequired" v:"in:0,1#是否必填必须为0或1"`
	DefaultValue   string `json:"defaultValue"`
	ValidationRule string `json:"validationRule"`
	SortOrder      int    `json:"sortOrder"`
	Status         int    `json:"status" v:"in:0,1#状态必须为0或1"`
}

type SystemConfigAddRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

// 编辑配置请求参数
type SystemConfigEditReq struct {
	g.Meta         `path:"/systemConfig/edit" method:"put" tags:"系统配置" summary:"系统配置-编辑"`
	Id             interface{} `json:"id" v:"required#配置ID不能为空"`
	ConfigValue    string      `json:"configValue"`
	ConfigGroup    string      `json:"configGroup" v:"required#配置分组不能为空"`
	ConfigType     string      `json:"configType" v:"required|in:string,number,boolean,json,array#配置类型不能为空|配置类型必须为string,number,boolean,json,array之一"`
	DisplayName    string      `json:"displayName" v:"required#显示名称不能为空"`
	Description    string      `json:"description"`
	IsPublic       int         `json:"isPublic" v:"in:0,1#是否公开必须为0或1"`
	IsRequired     int         `json:"isRequired" v:"in:0,1#是否必填必须为0或1"`
	DefaultValue   string      `json:"defaultValue"`
	ValidationRule string      `json:"validationRule"`
	SortOrder      int         `json:"sortOrder"`
	Status         int         `json:"status" v:"in:0,1#状态必须为0或1"`
}

type SystemConfigEditRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

// 删除配置请求参数
type SystemConfigDelReq struct {
	g.Meta `path:"/systemConfig/del" method:"delete" tags:"系统配置" summary:"系统配置-删除"`
	Id     interface{} `json:"id" v:"required#配置ID不能为空"`
}

type SystemConfigDelRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

// 批量删除配置请求参数
type SystemConfigBatchDelReq struct {
	g.Meta `path:"/systemConfig/batchdel" method:"delete" tags:"系统配置" summary:"系统配置-批量删除"`
	Ids    []interface{} `json:"ids" v:"required#配置ID列表不能为空"`
}

type SystemConfigBatchDelRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

// 根据配置键获取配置值
type SystemConfigGetByKeyReq struct {
	g.Meta    `path:"/systemConfig/getByKey" method:"get" tags:"系统配置" summary:"系统配置-根据键名获取"`
	ConfigKey string `json:"configKey" v:"required#配置键名不能为空"`
}

type SystemConfigGetByKeyRes struct {
	g.Meta      `mime:"application/json" example:"string"`
	ConfigKey   string `json:"configKey"`
	ConfigValue string `json:"configValue"`
	ConfigType  string `json:"configType"`
}

// 根据分组获取配置
type SystemConfigGetByGroupReq struct {
	g.Meta      `path:"/systemConfig/getByGroup" method:"get" tags:"系统配置" summary:"系统配置-根据分组获取"`
	ConfigGroup string `json:"configGroup" v:"required#配置分组不能为空"`
	IsPublic    *int   `json:"isPublic"` // 可选，筛选是否公开
}

type SystemConfigGetByGroupRes struct {
	g.Meta   `mime:"application/json" example:"string"`
	Configs  map[string]interface{} `json:"configs"` // 配置键值对
	GroupInfo struct {
		Group       string `json:"group"`
		DisplayName string `json:"displayName"`
		Description string `json:"description"`
	} `json:"groupInfo"`
}

// 批量更新配置
type SystemConfigBatchUpdateReq struct {
	g.Meta  `path:"/systemConfig/batchUpdate" method:"put" tags:"系统配置" summary:"系统配置-批量更新"`
	Configs []struct {
		ConfigKey   string `json:"configKey" v:"required#配置键名不能为空"`
		ConfigValue string `json:"configValue"`
	} `json:"configs" v:"required#配置列表不能为空"`
}

type SystemConfigBatchUpdateRes struct {
	g.Meta `mime:"application/json" example:"string"`
	UpdatedCount int      `json:"updatedCount"` // 更新成功的配置数量
	FailedKeys   []string `json:"failedKeys"`   // 更新失败的配置键名
}

// 获取公开配置（前端使用）
type SystemConfigPublicReq struct {
	g.Meta      `path:"/systemConfig/public" method:"get" tags:"系统配置" summary:"系统配置-获取公开配置"`
	ConfigGroup string `json:"configGroup"` // 可选，筛选配置分组
}

type SystemConfigPublicRes struct {
	g.Meta  `mime:"application/json" example:"string"`
	Configs map[string]interface{} `json:"configs"` // 公开配置的键值对
}

// 重置配置到默认值
type SystemConfigResetReq struct {
	g.Meta    `path:"/systemConfig/reset" method:"put" tags:"系统配置" summary:"系统配置-重置配置"`
	ConfigKey string `json:"configKey" v:"required#配置键名不能为空"`
}

type SystemConfigResetRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

// 验证配置值
type SystemConfigValidateReq struct {
	g.Meta      `path:"/systemConfig/validate" method:"post" tags:"系统配置" summary:"系统配置-验证配置值"`
	ConfigKey   string `json:"configKey" v:"required#配置键名不能为空"`
	ConfigValue string `json:"configValue" v:"required#配置值不能为空"`
}

type SystemConfigValidateRes struct {
	g.Meta  `mime:"application/json" example:"string"`
	IsValid bool   `json:"isValid"` // 是否有效
	Message string `json:"message"` // 验证消息
}