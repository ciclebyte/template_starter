// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Users is the golang structure of table users for DAO operations like Where/Data.
type Users struct {
	g.Meta         `orm:"table:users, do:true"`
	Id             interface{} // 用户ID
	Username       interface{} // 用户名
	Email          interface{} // 邮箱
	PasswordHash   interface{} // 密码哈希
	Nickname       interface{} // 昵称
	Avatar         interface{} // 头像URL
	Phone          interface{} // 手机号
	Status         interface{} // 状态：0=禁用，1=正常
	EmailVerified  interface{} // 邮箱验证状态
	LastLoginAt    *gtime.Time // 最后登录时间
	LastLoginIp    interface{} // 最后登录IP
	LoginCount     interface{} // 登录次数
	OrganizationId interface{} // 所属组织ID（暂时保留，第三阶段使用）
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
}
