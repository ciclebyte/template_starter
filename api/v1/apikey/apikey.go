package apikey

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// API Key 数据结构
type ApiKey struct {
	Id          int64       `json:"id" description:"API Key ID"`
	UserId      int64       `json:"userId" description:"用户ID"`
	Name        string      `json:"name" description:"API Key 名称"`
	KeyId       string      `json:"keyId" description:"API Key ID (公开标识)"`
	KeySecret   string      `json:"keySecret,omitempty" description:"API Key Secret (仅创建时返回)"`
	Permissions []string    `json:"permissions" description:"权限列表"`
	LastUsedAt  *gtime.Time `json:"lastUsedAt" description:"最后使用时间"`
	LastUsedIp  string      `json:"lastUsedIp" description:"最后使用IP"`
	ExpiresAt   *gtime.Time `json:"expiresAt" description:"过期时间"`
	Status      int         `json:"status" description:"状态: 0-禁用, 1-启用"`
	CreatedAt   *gtime.Time `json:"createdAt" description:"创建时间"`
	UpdatedAt   *gtime.Time `json:"updatedAt" description:"更新时间"`
}

// API Key 使用日志
type ApiKeyLog struct {
	Id           int64       `json:"id" description:"日志ID"`
	ApiKeyId     int64       `json:"apiKeyId" description:"API Key ID"`
	UserId       int64       `json:"userId" description:"用户ID"`
	Method       string      `json:"method" description:"HTTP方法"`
	Path         string      `json:"path" description:"请求路径"`
	Ip           string      `json:"ip" description:"请求IP"`
	UserAgent    string      `json:"userAgent" description:"User Agent"`
	StatusCode   int         `json:"statusCode" description:"响应状态码"`
	ResponseTime int         `json:"responseTime" description:"响应时间(毫秒)"`
	CreatedAt    *gtime.Time `json:"createdAt" description:"创建时间"`
}

// ============================================================================
// API Key 管理接口
// ============================================================================

// GetApiKeysReq 获取API Key列表请求
type GetApiKeysReq struct {
	g.Meta `path:"/api-keys" method:"get" summary:"获取API Key列表" tags:"ApiKey"`
	Search string `json:"search" description:"搜索关键词"`
	Status *int   `json:"status" description:"状态筛选"`
	Page   int    `json:"page" d:"1" description:"页码"`
	Size   int    `json:"size" d:"20" description:"每页数量"`
}

type GetApiKeysRes struct {
	g.Meta `mime:"application/json"`
	List   []ApiKey `json:"list" description:"API Key列表"`
	Total  int      `json:"total" description:"总数"`
	Page   int      `json:"page" description:"当前页"`
	Size   int      `json:"size" description:"每页数量"`
}

// CreateApiKeyReq 创建API Key请求
type CreateApiKeyReq struct {
	g.Meta      `path:"/api-keys" method:"post" summary:"创建API Key" tags:"ApiKey"`
	Name        string      `json:"name" v:"required|length:1,100" description:"API Key 名称"`
	Permissions []string    `json:"permissions" description:"权限列表"`
	ExpiresAt   *gtime.Time `json:"expiresAt" description:"过期时间"`
}

type CreateApiKeyRes struct {
	g.Meta    `mime:"application/json"`
	ApiKey    ApiKey `json:"apiKey" description:"创建的API Key信息"`
	KeySecret string `json:"keySecret" description:"API Key Secret (仅此次返回)"`
}

// UpdateApiKeyReq 更新API Key请求
type UpdateApiKeyReq struct {
	g.Meta      `path:"/api-keys/{id}" method:"put" summary:"更新API Key" tags:"ApiKey"`
	Id          int64       `json:"id" in:"path" v:"required" description:"API Key ID"`
	Name        string      `json:"name" v:"required|length:1,100" description:"API Key 名称"`
	Permissions []string    `json:"permissions" description:"权限列表"`
	ExpiresAt   *gtime.Time `json:"expiresAt" description:"过期时间"`
	Status      int         `json:"status" v:"in:0,1" description:"状态: 0-禁用, 1-启用"`
}

type UpdateApiKeyRes struct {
	g.Meta `mime:"application/json"`
}

// DeleteApiKeyReq 删除API Key请求
type DeleteApiKeyReq struct {
	g.Meta `path:"/api-keys/{id}" method:"delete" summary:"删除API Key" tags:"ApiKey"`
	Id     int64 `json:"id" in:"path" v:"required" description:"API Key ID"`
}

type DeleteApiKeyRes struct {
	g.Meta `mime:"application/json"`
}

// RegenerateApiKeyReq 重新生成API Key Secret请求
type RegenerateApiKeyReq struct {
	g.Meta `path:"/api-keys/{id}/regenerate" method:"post" summary:"重新生成API Key Secret" tags:"ApiKey"`
	Id     int64 `json:"id" in:"path" v:"required" description:"API Key ID"`
}

type RegenerateApiKeyRes struct {
	g.Meta    `mime:"application/json"`
	KeySecret string `json:"keySecret" description:"新的API Key Secret"`
}

// ============================================================================
// API Key 使用日志接口
// ============================================================================

// GetApiKeyLogsReq 获取API Key使用日志请求
type GetApiKeyLogsReq struct {
	g.Meta   `path:"/api-keys/{id}/logs" method:"get" summary:"获取API Key使用日志" tags:"ApiKey"`
	Id       int64  `json:"id" in:"path" v:"required" description:"API Key ID"`
	Method   string `json:"method" description:"HTTP方法筛选"`
	Path     string `json:"path" description:"路径筛选"`
	Page     int    `json:"page" d:"1" description:"页码"`
	Size     int    `json:"size" d:"20" description:"每页数量"`
}

type GetApiKeyLogsRes struct {
	g.Meta `mime:"application/json"`
	List   []ApiKeyLog `json:"list" description:"使用日志列表"`
	Total  int         `json:"total" description:"总数"`
	Page   int         `json:"page" description:"当前页"`
	Size   int         `json:"size" description:"每页数量"`
}

// ============================================================================
// 个人 API Key 管理接口 (在Profile中)
// ============================================================================

// GetMyApiKeysReq 获取我的API Key列表请求
type GetMyApiKeysReq struct {
	g.Meta `path:"/profile/api-keys" method:"get" summary:"获取我的API Key列表" tags:"Profile"`
	Search string `json:"search" description:"搜索关键词"`
	Status *int   `json:"status" description:"状态筛选"`
	Page   int    `json:"page" d:"1" description:"页码"`
	Size   int    `json:"size" d:"20" description:"每页数量"`
}

type GetMyApiKeysRes struct {
	g.Meta `mime:"application/json"`
	List   []ApiKey `json:"list" description:"API Key列表"`
	Total  int      `json:"total" description:"总数"`
	Page   int      `json:"page" description:"当前页"`
	Size   int      `json:"size" description:"每页数量"`
}

// CreateMyApiKeyReq 创建我的API Key请求
type CreateMyApiKeyReq struct {
	g.Meta      `path:"/profile/api-keys" method:"post" summary:"创建我的API Key" tags:"Profile"`
	Name        string      `json:"name" v:"required|length:1,100" description:"API Key 名称"`
	Permissions []string    `json:"permissions" description:"权限列表"`
	ExpiresAt   *gtime.Time `json:"expiresAt" description:"过期时间"`
}

type CreateMyApiKeyRes struct {
	g.Meta    `mime:"application/json"`
	ApiKey    ApiKey `json:"apiKey" description:"创建的API Key信息"`
	KeySecret string `json:"keySecret" description:"API Key Secret (仅此次返回)"`
}

// UpdateMyApiKeyReq 更新我的API Key请求
type UpdateMyApiKeyReq struct {
	g.Meta      `path:"/profile/api-keys/{id}" method:"put" summary:"更新我的API Key" tags:"Profile"`
	Id          int64       `json:"id" in:"path" v:"required" description:"API Key ID"`
	Name        string      `json:"name" v:"required|length:1,100" description:"API Key 名称"`
	Permissions []string    `json:"permissions" description:"权限列表"`
	ExpiresAt   *gtime.Time `json:"expiresAt" description:"过期时间"`
	Status      int         `json:"status" v:"in:0,1" description:"状态: 0-禁用, 1-启用"`
}

type UpdateMyApiKeyRes struct {
	g.Meta `mime:"application/json"`
}

// DeleteMyApiKeyReq 删除我的API Key请求
type DeleteMyApiKeyReq struct {
	g.Meta `path:"/profile/api-keys/{id}" method:"delete" summary:"删除我的API Key" tags:"Profile"`
	Id     int64 `json:"id" in:"path" v:"required" description:"API Key ID"`
}

type DeleteMyApiKeyRes struct {
	g.Meta `mime:"application/json"`
}

// RegenerateMyApiKeyReq 重新生成我的API Key Secret请求
type RegenerateMyApiKeyReq struct {
	g.Meta `path:"/profile/api-keys/{id}/regenerate" method:"post" summary:"重新生成我的API Key Secret" tags:"Profile"`
	Id     int64 `json:"id" in:"path" v:"required" description:"API Key ID"`
}

type RegenerateMyApiKeyRes struct {
	g.Meta    `mime:"application/json"`
	KeySecret string `json:"keySecret" description:"新的API Key Secret"`
}