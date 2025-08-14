// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TemplateTagsDao is the data access object for table template_tags.
type TemplateTagsDao struct {
	table   string              // table is the underlying table name of the DAO.
	group   string              // group is the database configuration group name of current DAO.
	columns TemplateTagsColumns // columns contains all the column names of Table for convenient usage.
}

// TemplateTagsColumns defines and stores column names for table template_tags.
type TemplateTagsColumns struct {
	Id         string // 关联ID
	TemplateId string // 模板ID
	TagId      string // 标签ID
	CreatedAt  string // 创建时间
	UpdatedAt  string // 更新时间
}

// templateTagsColumns holds the columns for table template_tags.
var templateTagsColumns = TemplateTagsColumns{
	Id:         "id",
	TemplateId: "template_id",
	TagId:      "tag_id",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
}

// NewTemplateTagsDao creates and returns a new DAO object for table data access.
func NewTemplateTagsDao() *TemplateTagsDao {
	return &TemplateTagsDao{
		group:   "default",
		table:   "template_tags",
		columns: templateTagsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *TemplateTagsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *TemplateTagsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *TemplateTagsDao) Columns() TemplateTagsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *TemplateTagsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *TemplateTagsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *TemplateTagsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}