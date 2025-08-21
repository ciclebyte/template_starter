// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UsersDao is the data access object for table users.
type UsersDao struct {
	table   string       // table is the underlying table name of the DAO.
	group   string       // group is the database configuration group name of current DAO.
	columns UsersColumns // columns contains all the column names of Table for convenient usage.
}

// UsersColumns defines and stores column names for table users.
type UsersColumns struct {
	Id             string // 用户ID
	Username       string // 用户名
	Email          string // 邮箱
	PasswordHash   string // 密码哈希
	Nickname       string // 昵称
	Avatar         string // 头像URL
	Phone          string // 手机号
	Status         string // 状态：0=禁用，1=正常
	EmailVerified  string // 邮箱验证状态
	LastLoginAt    string // 最后登录时间
	LastLoginIp    string // 最后登录IP
	LoginCount     string // 登录次数
	OrganizationId string // 所属组织ID（暂时保留，第三阶段使用）
	CreatedAt      string //
	UpdatedAt      string //
}

// usersColumns holds the columns for table users.
var usersColumns = UsersColumns{
	Id:             "id",
	Username:       "username",
	Email:          "email",
	PasswordHash:   "password_hash",
	Nickname:       "nickname",
	Avatar:         "avatar",
	Phone:          "phone",
	Status:         "status",
	EmailVerified:  "email_verified",
	LastLoginAt:    "last_login_at",
	LastLoginIp:    "last_login_ip",
	LoginCount:     "login_count",
	OrganizationId: "organization_id",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// NewUsersDao creates and returns a new DAO object for table data access.
func NewUsersDao() *UsersDao {
	return &UsersDao{
		group:   "default",
		table:   "users",
		columns: usersColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *UsersDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *UsersDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *UsersDao) Columns() UsersColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *UsersDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *UsersDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *UsersDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
