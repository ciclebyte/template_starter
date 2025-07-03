// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// LanguagesDao is the data access object for table languages.
type LanguagesDao struct {
	table   string           // table is the underlying table name of the DAO.
	group   string           // group is the database configuration group name of current DAO.
	columns LanguagesColumns // columns contains all the column names of Table for convenient usage.
}

// LanguagesColumns defines and stores column names for table languages.
type LanguagesColumns struct {
	Id          string // 语言ID，自增主键
	Name        string // 语言名称（如JavaScript、Python）
	DisplayName string // 语言显示名称
	Code        string // 语言代码（如js、py）
	Icon        string // 语言图标标识或URL
	Color       string // 语言代表色（十六进制）
	IsPopular   string // 是否热门语言
	CreatedAt   string // 记录创建时间
	UpdatedAt   string // 记录最后更新时间
}

// languagesColumns holds the columns for table languages.
var languagesColumns = LanguagesColumns{
	Id:          "id",
	Name:        "name",
	DisplayName: "display_name",
	Code:        "code",
	Icon:        "icon",
	Color:       "color",
	IsPopular:   "is_popular",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// NewLanguagesDao creates and returns a new DAO object for table data access.
func NewLanguagesDao() *LanguagesDao {
	return &LanguagesDao{
		group:   "default",
		table:   "languages",
		columns: languagesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *LanguagesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *LanguagesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *LanguagesDao) Columns() LanguagesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *LanguagesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *LanguagesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *LanguagesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
