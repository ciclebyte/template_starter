// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TemplateVariablesDao is the data access object for table template_variables.
type TemplateVariablesDao struct {
	table   string                   // table is the underlying table name of the DAO.
	group   string                   // group is the database configuration group name of current DAO.
	columns TemplateVariablesColumns // columns contains all the column names of Table for convenient usage.
}

// TemplateVariablesColumns defines and stores column names for table template_variables.
type TemplateVariablesColumns struct {
	Id              string // 变量ID，自增主键
	TemplateId      string // 所属模板ID
	Name            string // 变量名称
	Description     string // 变量描述
	DefaultValue    string // 变量默认值
	IsRequired      string // 是否为必填变量
	ValidationRegex string // 变量值验证正则表达式
	DisplayOrder    string // 显示顺序，数字越小排序越靠前
	CreatedAt       string // 创建时间
}

// templateVariablesColumns holds the columns for table template_variables.
var templateVariablesColumns = TemplateVariablesColumns{
	Id:              "id",
	TemplateId:      "template_id",
	Name:            "name",
	Description:     "description",
	DefaultValue:    "default_value",
	IsRequired:      "is_required",
	ValidationRegex: "validation_regex",
	DisplayOrder:    "display_order",
	CreatedAt:       "created_at",
}

// NewTemplateVariablesDao creates and returns a new DAO object for table data access.
func NewTemplateVariablesDao() *TemplateVariablesDao {
	return &TemplateVariablesDao{
		group:   "default",
		table:   "template_variables",
		columns: templateVariablesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *TemplateVariablesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *TemplateVariablesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *TemplateVariablesDao) Columns() TemplateVariablesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *TemplateVariablesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *TemplateVariablesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *TemplateVariablesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
