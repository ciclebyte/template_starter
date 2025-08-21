// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PermissionsDao is the data access object for table permissions.
type PermissionsDao struct {
	table   string             // table is the underlying table name of the DAO.
	group   string             // group is the database configuration group name of current DAO.
	columns PermissionsColumns // columns contains all the column names of Table for convenient usage.
}

// PermissionsColumns defines and stores column names for table permissions.
type PermissionsColumns struct {
	Id          string // 权限ID
	Name        string // 权限名称
	Code        string // 权限编码
	Resource    string // 资源类型
	Action      string // 操作类型
	Description string // 权限描述
	ParentId    string // 父权限ID
	SortOrder   string // 排序
	CreatedAt   string //
}

// permissionsColumns holds the columns for table permissions.
var permissionsColumns = PermissionsColumns{
	Id:          "id",
	Name:        "name",
	Code:        "code",
	Resource:    "resource",
	Action:      "action",
	Description: "description",
	ParentId:    "parent_id",
	SortOrder:   "sort_order",
	CreatedAt:   "created_at",
}

// NewPermissionsDao creates and returns a new DAO object for table data access.
func NewPermissionsDao() *PermissionsDao {
	return &PermissionsDao{
		group:   "default",
		table:   "permissions",
		columns: permissionsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *PermissionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *PermissionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *PermissionsDao) Columns() PermissionsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *PermissionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *PermissionsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *PermissionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
