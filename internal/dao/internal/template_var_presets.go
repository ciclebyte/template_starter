// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TemplateVarPresetsDao is the data access object for table template_var_presets.
type TemplateVarPresetsDao struct {
	table   string                       // table is the underlying table name of the DAO.
	group   string                       // group is the database configuration group name of current DAO.
	columns TemplateVarPresetsColumns    // columns contains all the column names of Table for convenient usage.
}

// TemplateVarPresetsColumns defines and stores column names for table template_var_presets.
type TemplateVarPresetsColumns struct {
	Id         string // 模板变量预设ID
	TemplateId string // 模板ID
	PresetId   string // 预设ID
	CreatedAt  string // 创建时间
	UpdatedAt  string // 更新时间
}

// templateVarPresetsColumns holds the columns for table template_var_presets.
var templateVarPresetsColumns = TemplateVarPresetsColumns{
	Id:         "id",
	TemplateId: "template_id",
	PresetId:   "preset_id",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
}

// NewTemplateVarPresetsDao creates and returns a new DAO object for table data access.
func NewTemplateVarPresetsDao() *TemplateVarPresetsDao {
	return &TemplateVarPresetsDao{
		group:   "default",
		table:   "template_var_presets",
		columns: templateVarPresetsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *TemplateVarPresetsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *TemplateVarPresetsDao) Table() string {
	return dao.table
}

// Columns returns the columns of current dao.
func (dao *TemplateVarPresetsDao) Columns() TemplateVarPresetsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *TemplateVarPresetsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context.
func (dao *TemplateVarPresetsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
func (dao *TemplateVarPresetsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}