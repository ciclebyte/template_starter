// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Users is the golang structure for table users.
type Users struct {
	Id             int64       `json:"id"             description:"用户ID"`
	Username       string      `json:"username"       description:"用户名"`
	Email          string      `json:"email"          description:"邮箱"`
	PasswordHash   string      `json:"passwordHash"   description:"密码哈希"`
	Nickname       string      `json:"nickname"       description:"昵称"`
	Avatar         string      `json:"avatar"         description:"头像URL"`
	Phone          string      `json:"phone"          description:"手机号"`
	Status         int         `json:"status"         description:"状态：0=禁用，1=正常"`
	EmailVerified  int         `json:"emailVerified"  description:"邮箱验证状态"`
	LastLoginAt    *gtime.Time `json:"lastLoginAt"    description:"最后登录时间"`
	LastLoginIp    string      `json:"lastLoginIp"    description:"最后登录IP"`
	LoginCount     int         `json:"loginCount"     description:"登录次数"`
	OrganizationId int64       `json:"organizationId" description:"所属组织ID（暂时保留，第三阶段使用）"`
	CreatedAt      *gtime.Time `json:"createdAt"      description:""`
	UpdatedAt      *gtime.Time `json:"updatedAt"      description:""`
}
