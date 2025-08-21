// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserSessions is the golang structure for table user_sessions.
type UserSessions struct {
	Id        string      `json:"id"        description:"会话ID"`
	UserId    int64       `json:"userId"    description:"用户ID"`
	IpAddress string      `json:"ipAddress" description:"IP地址"`
	UserAgent string      `json:"userAgent" description:"用户代理"`
	Data      string      `json:"data"      description:"会话数据"`
	ExpiresAt *gtime.Time `json:"expiresAt" description:"过期时间"`
	CreatedAt *gtime.Time `json:"createdAt" description:""`
	UpdatedAt *gtime.Time `json:"updatedAt" description:""`
}
