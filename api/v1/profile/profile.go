package profile

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ============================================================================
// 用户个人资料管理API定义
// ============================================================================

// UserProfile 用户个人资料
type UserProfile struct {
	Id          int64       `json:"id" dc:"用户ID"`
	Username    string      `json:"username" dc:"用户名"`
	Email       string      `json:"email" dc:"邮箱"`
	Nickname    string      `json:"nickname" dc:"昵称"`
	Avatar      string      `json:"avatar" dc:"头像"`
	Phone       string      `json:"phone" dc:"手机号"`
	Status      int         `json:"status" dc:"状态"`
	LastLoginAt *gtime.Time `json:"lastLoginAt" dc:"最后登录时间"`
	CreatedAt   *gtime.Time `json:"createdAt" dc:"创建时间"`
	UpdatedAt   *gtime.Time `json:"updatedAt" dc:"更新时间"`
}

// UserRole 用户角色信息
type UserRole struct {
	Id          int64       `json:"id" dc:"角色ID"`
	Name        string      `json:"name" dc:"角色名称"`
	Code        string      `json:"code" dc:"角色编码"`
	Description string      `json:"description" dc:"角色描述"`
	AssignedAt  *gtime.Time `json:"assignedAt" dc:"分配时间"`
	ExpiresAt   *gtime.Time `json:"expiresAt" dc:"过期时间"`
}

// UserPermission 用户权限信息
type UserPermission struct {
	Id          int64  `json:"id" dc:"权限ID"`
	Name        string `json:"name" dc:"权限名称"`
	Code        string `json:"code" dc:"权限编码"`
	Resource    string `json:"resource" dc:"资源类型"`
	Action      string `json:"action" dc:"操作类型"`
	Description string `json:"description" dc:"权限描述"`
}

// GetProfileReq 获取个人资料请求
type GetProfileReq struct {
	g.Meta `path:"/profile" method:"get" summary:"获取个人资料" tags:"个人中心"`
}

type GetProfileRes struct {
	Profile     *UserProfile      `json:"profile" dc:"用户资料"`
	Roles       []UserRole        `json:"roles" dc:"用户角色"`
	Permissions []UserPermission  `json:"permissions" dc:"用户权限"`
}

// UpdateProfileReq 更新个人资料请求
type UpdateProfileReq struct {
	g.Meta   `path:"/profile" method:"put" summary:"更新个人资料" tags:"个人中心"`
	Nickname string `json:"nickname" v:"length:1,50" dc:"昵称"`
	Avatar   string `json:"avatar" dc:"头像URL"`
	Phone    string `json:"phone" v:"phone" dc:"手机号"`
}

type UpdateProfileRes struct{}

// ChangePasswordReq 修改密码请求
type ChangePasswordReq struct {
	g.Meta      `path:"/profile/password" method:"put" summary:"修改密码" tags:"个人中心"`
	OldPassword string `json:"oldPassword" v:"required|length:6,32" dc:"原密码"`
	NewPassword string `json:"newPassword" v:"required|length:6,32" dc:"新密码"`
}

type ChangePasswordRes struct{}

// UpdateEmailReq 更新邮箱请求
type UpdateEmailReq struct {
	g.Meta   `path:"/profile/email" method:"put" summary:"更新邮箱" tags:"个人中心"`
	Email    string `json:"email" v:"required|email" dc:"新邮箱"`
	Password string `json:"password" v:"required" dc:"当前密码"`
}

type UpdateEmailRes struct{}

// UploadAvatarReq 上传头像请求
type UploadAvatarReq struct {
	g.Meta `path:"/profile/avatar" method:"post" summary:"上传头像" tags:"个人中心"`
	// 这里应该是文件上传，暂时用URL方式
	AvatarUrl string `json:"avatarUrl" v:"required|url" dc:"头像URL"`
}

type UploadAvatarRes struct {
	AvatarUrl string `json:"avatarUrl" dc:"头像URL"`
}

// GetSecurityInfoReq 获取账户安全信息请求
type GetSecurityInfoReq struct {
	g.Meta `path:"/profile/security" method:"get" summary:"获取账户安全信息" tags:"个人中心"`
}

type SecurityInfo struct {
	EmailVerified      bool        `json:"emailVerified" dc:"邮箱是否验证"`
	PhoneVerified      bool        `json:"phoneVerified" dc:"手机是否验证"`
	TwoFactorEnabled   bool        `json:"twoFactorEnabled" dc:"是否启用双因子认证"`
	LastLoginAt        *gtime.Time `json:"lastLoginAt" dc:"最后登录时间"`
	LastLoginIp        string      `json:"lastLoginIp" dc:"最后登录IP"`
	LoginCount         int         `json:"loginCount" dc:"登录次数"`
	PasswordUpdatedAt  *gtime.Time `json:"passwordUpdatedAt" dc:"密码最后更新时间"`
}

type GetSecurityInfoRes struct {
	SecurityInfo *SecurityInfo `json:"securityInfo" dc:"安全信息"`
}

// GetLoginHistoryReq 获取登录历史请求
type GetLoginHistoryReq struct {
	g.Meta `path:"/profile/login-history" method:"get" summary:"获取登录历史" tags:"个人中心"`
	Page   int `json:"page" v:"min:1" dc:"页码，默认1"`
	Size   int `json:"size" v:"min:1|max:100" dc:"每页数量，默认20"`
}

type LoginHistory struct {
	Id        string      `json:"id" dc:"记录ID"`
	LoginAt   *gtime.Time `json:"loginAt" dc:"登录时间"`
	LoginIp   string      `json:"loginIp" dc:"登录IP"`
	UserAgent string      `json:"userAgent" dc:"用户代理"`
	Location  string      `json:"location" dc:"登录地点"`
	Status    string      `json:"status" dc:"登录状态"`
}

type GetLoginHistoryRes struct {
	List  []LoginHistory `json:"list" dc:"登录历史列表"`
	Total int            `json:"total" dc:"总数"`
	Page  int            `json:"page" dc:"当前页"`
	Size  int            `json:"size" dc:"每页数量"`
}