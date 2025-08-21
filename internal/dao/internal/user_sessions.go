// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserSessionsDao is the data access object for table user_sessions.
type UserSessionsDao struct {
	table   string              // table is the underlying table name of the DAO.
	group   string              // group is the database configuration group name of current DAO.
	columns UserSessionsColumns // columns contains all the column names of Table for convenient usage.
}

// UserSessionsColumns defines and stores column names for table user_sessions.
type UserSessionsColumns struct {
	Id        string // 会话ID
	UserId    string // 用户ID
	IpAddress string // IP地址
	UserAgent string // 用户代理
	Data      string // 会话数据
	ExpiresAt string // 过期时间
	CreatedAt string //
	UpdatedAt string //
}

// userSessionsColumns holds the columns for table user_sessions.
var userSessionsColumns = UserSessionsColumns{
	Id:        "id",
	UserId:    "user_id",
	IpAddress: "ip_address",
	UserAgent: "user_agent",
	Data:      "data",
	ExpiresAt: "expires_at",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewUserSessionsDao creates and returns a new DAO object for table data access.
func NewUserSessionsDao() *UserSessionsDao {
	return &UserSessionsDao{
		group:   "default",
		table:   "user_sessions",
		columns: userSessionsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *UserSessionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *UserSessionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *UserSessionsDao) Columns() UserSessionsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *UserSessionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *UserSessionsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *UserSessionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
