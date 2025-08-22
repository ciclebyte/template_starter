// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ApiKeys is the golang structure for table api_keys.
type ApiKeys struct {
	Id          uint64      `json:"id"          description:"API Key ID"`
	UserId      uint64      `json:"userId"      description:"用户ID"`
	Name        string      `json:"name"        description:"API Key 名称"`
	KeyId       string      `json:"keyId"       description:"API Key ID (公开标识)"`
	KeySecret   string      `json:"keySecret"   description:"API Key Secret (加密存储)"`
	Permissions string      `json:"permissions" description:"API Key 权限列表"`
	LastUsedAt  *gtime.Time `json:"lastUsedAt"  description:"最后使用时间"`
	LastUsedIp  string      `json:"lastUsedIp"  description:"最后使用IP"`
	ExpiresAt   *gtime.Time `json:"expiresAt"   description:"过期时间"`
	Status      int         `json:"status"      description:"状态: 0-禁用, 1-启用"`
	CreatedAt   *gtime.Time `json:"createdAt"   description:"创建时间"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   description:"更新时间"`
}
