// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TemplateExposeFieldsDao is the data access object for table template_expose_fields.
type TemplateExposeFieldsDao struct {
	table   string                         // table is the underlying table name of the DAO.
	group   string                         // group is the database configuration group name of current DAO.
	columns TemplateExposeFieldsColumns    // columns contains all the column names of Table for convenient usage.
}

// TemplateExposeFieldsColumns defines and stores column names for table template_expose_fields.
type TemplateExposeFieldsColumns struct {
	Id              string // 暴露字段ID
	TemplateId      string // 模板ID
	FieldSchemaJson string // 字段结构定义（支持树形嵌套）
	Version         string // 版本号
	Description     string // 字段说明文档
	CreatedAt       string // 创建时间
	UpdatedAt       string // 更新时间
}

// templateExposeFieldsColumns holds the columns for table template_expose_fields.
var templateExposeFieldsColumns = TemplateExposeFieldsColumns{
	Id:              "id",
	TemplateId:      "template_id",
	FieldSchemaJson: "field_schema_json",
	Version:         "version",
	Description:     "description",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
}

// NewTemplateExposeFieldsDao creates and returns a new DAO object for table data access.
func NewTemplateExposeFieldsDao() *TemplateExposeFieldsDao {
	return &TemplateExposeFieldsDao{
		group:   "default",
		table:   "template_expose_fields",
		columns: templateExposeFieldsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *TemplateExposeFieldsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *TemplateExposeFieldsDao) Table() string {
	return dao.table
}

// Columns returns the columns of current dao.
func (dao *TemplateExposeFieldsDao) Columns() TemplateExposeFieldsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *TemplateExposeFieldsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context.
func (dao *TemplateExposeFieldsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
func (dao *TemplateExposeFieldsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}