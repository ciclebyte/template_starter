// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ApiKeyLogs is the golang structure for table api_key_logs.
type ApiKeyLogs struct {
	Id           uint64      `json:"id"           description:"日志ID"`
	ApiKeyId     uint64      `json:"apiKeyId"     description:"API Key ID"`
	UserId       uint64      `json:"userId"       description:"用户ID"`
	Method       string      `json:"method"       description:"HTTP方法"`
	Path         string      `json:"path"         description:"请求路径"`
	Ip           string      `json:"ip"           description:"请求IP"`
	UserAgent    string      `json:"userAgent"    description:"User Agent"`
	StatusCode   int         `json:"statusCode"   description:"响应状态码"`
	ResponseTime int         `json:"responseTime" description:"响应时间(毫秒)"`
	CreatedAt    *gtime.Time `json:"createdAt"    description:"创建时间"`
}
