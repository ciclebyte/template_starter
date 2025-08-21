// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserSessions is the golang structure of table user_sessions for DAO operations like Where/Data.
type UserSessions struct {
	g.Meta    `orm:"table:user_sessions, do:true"`
	Id        interface{} // 会话ID
	UserId    interface{} // 用户ID
	IpAddress interface{} // IP地址
	UserAgent interface{} // 用户代理
	Data      interface{} // 会话数据
	ExpiresAt *gtime.Time // 过期时间
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
