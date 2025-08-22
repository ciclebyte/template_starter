// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ApiKeyLogsDao is the data access object for table api_key_logs.
type ApiKeyLogsDao struct {
	table   string            // table is the underlying table name of the DAO.
	group   string            // group is the database configuration group name of current DAO.
	columns ApiKeyLogsColumns // columns contains all the column names of Table for convenient usage.
}

// ApiKeyLogsColumns defines and stores column names for table api_key_logs.
type ApiKeyLogsColumns struct {
	Id           string // 日志ID
	ApiKeyId     string // API Key ID
	UserId       string // 用户ID
	Method       string // HTTP方法
	Path         string // 请求路径
	Ip           string // 请求IP
	UserAgent    string // User Agent
	StatusCode   string // 响应状态码
	ResponseTime string // 响应时间(毫秒)
	CreatedAt    string // 创建时间
}

// apiKeyLogsColumns holds the columns for table api_key_logs.
var apiKeyLogsColumns = ApiKeyLogsColumns{
	Id:           "id",
	ApiKeyId:     "api_key_id",
	UserId:       "user_id",
	Method:       "method",
	Path:         "path",
	Ip:           "ip",
	UserAgent:    "user_agent",
	StatusCode:   "status_code",
	ResponseTime: "response_time",
	CreatedAt:    "created_at",
}

// NewApiKeyLogsDao creates and returns a new DAO object for table data access.
func NewApiKeyLogsDao() *ApiKeyLogsDao {
	return &ApiKeyLogsDao{
		group:   "default",
		table:   "api_key_logs",
		columns: apiKeyLogsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *ApiKeyLogsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *ApiKeyLogsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *ApiKeyLogsDao) Columns() ApiKeyLogsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *ApiKeyLogsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *ApiKeyLogsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *ApiKeyLogsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
