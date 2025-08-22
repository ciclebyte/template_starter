// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ApiKeysDao is the data access object for table api_keys.
type ApiKeysDao struct {
	table   string         // table is the underlying table name of the DAO.
	group   string         // group is the database configuration group name of current DAO.
	columns ApiKeysColumns // columns contains all the column names of Table for convenient usage.
}

// ApiKeysColumns defines and stores column names for table api_keys.
type ApiKeysColumns struct {
	Id          string // API Key ID
	UserId      string // 用户ID
	Name        string // API Key 名称
	KeyId       string // API Key ID (公开标识)
	KeySecret   string // API Key Secret (加密存储)
	Permissions string // API Key 权限列表
	LastUsedAt  string // 最后使用时间
	LastUsedIp  string // 最后使用IP
	ExpiresAt   string // 过期时间
	Status      string // 状态: 0-禁用, 1-启用
	CreatedAt   string // 创建时间
	UpdatedAt   string // 更新时间
}

// apiKeysColumns holds the columns for table api_keys.
var apiKeysColumns = ApiKeysColumns{
	Id:          "id",
	UserId:      "user_id",
	Name:        "name",
	KeyId:       "key_id",
	KeySecret:   "key_secret",
	Permissions: "permissions",
	LastUsedAt:  "last_used_at",
	LastUsedIp:  "last_used_ip",
	ExpiresAt:   "expires_at",
	Status:      "status",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// NewApiKeysDao creates and returns a new DAO object for table data access.
func NewApiKeysDao() *ApiKeysDao {
	return &ApiKeysDao{
		group:   "default",
		table:   "api_keys",
		columns: apiKeysColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *ApiKeysDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *ApiKeysDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *ApiKeysDao) Columns() ApiKeysColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *ApiKeysDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *ApiKeysDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *ApiKeysDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
