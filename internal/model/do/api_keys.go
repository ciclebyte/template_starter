// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ApiKeys is the golang structure of table api_keys for DAO operations like Where/Data.
type ApiKeys struct {
	g.Meta      `orm:"table:api_keys, do:true"`
	Id          interface{} // API Key ID
	UserId      interface{} // 用户ID
	Name        interface{} // API Key 名称
	KeyId       interface{} // API Key ID (公开标识)
	KeySecret   interface{} // API Key Secret (加密存储)
	Permissions interface{} // API Key 权限列表
	LastUsedAt  *gtime.Time // 最后使用时间
	LastUsedIp  interface{} // 最后使用IP
	ExpiresAt   *gtime.Time // 过期时间
	Status      interface{} // 状态: 0-禁用, 1-启用
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
}
