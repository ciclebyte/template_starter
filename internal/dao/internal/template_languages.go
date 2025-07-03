// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TemplateLanguagesDao is the data access object for table template_languages.
type TemplateLanguagesDao struct {
	table   string                   // table is the underlying table name of the DAO.
	group   string                   // group is the database configuration group name of current DAO.
	columns TemplateLanguagesColumns // columns contains all the column names of Table for convenient usage.
}

// TemplateLanguagesColumns defines and stores column names for table template_languages.
type TemplateLanguagesColumns struct {
	Id         string // 关联ID，自增主键
	TemplateId string // 关联的模板ID
	LanguageId string // 关联的语言ID
	IsPrimary  string // 是否主要语言
	CreatedAt  string // 记录创建时间
}

// templateLanguagesColumns holds the columns for table template_languages.
var templateLanguagesColumns = TemplateLanguagesColumns{
	Id:         "id",
	TemplateId: "template_id",
	LanguageId: "language_id",
	IsPrimary:  "is_primary",
	CreatedAt:  "created_at",
}

// NewTemplateLanguagesDao creates and returns a new DAO object for table data access.
func NewTemplateLanguagesDao() *TemplateLanguagesDao {
	return &TemplateLanguagesDao{
		group:   "default",
		table:   "template_languages",
		columns: templateLanguagesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *TemplateLanguagesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *TemplateLanguagesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *TemplateLanguagesDao) Columns() TemplateLanguagesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *TemplateLanguagesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *TemplateLanguagesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *TemplateLanguagesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
