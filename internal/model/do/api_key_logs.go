// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ApiKeyLogs is the golang structure of table api_key_logs for DAO operations like Where/Data.
type ApiKeyLogs struct {
	g.Meta       `orm:"table:api_key_logs, do:true"`
	Id           interface{} // 日志ID
	ApiKeyId     interface{} // API Key ID
	UserId       interface{} // 用户ID
	Method       interface{} // HTTP方法
	Path         interface{} // 请求路径
	Ip           interface{} // 请求IP
	UserAgent    interface{} // User Agent
	StatusCode   interface{} // 响应状态码
	ResponseTime interface{} // 响应时间(毫秒)
	CreatedAt    *gtime.Time // 创建时间
}
