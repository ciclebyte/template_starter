// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TemplateVariablePresetsDao is the data access object for table template_variable_presets.
type TemplateVariablePresetsDao struct {
	table   string                         // table is the underlying table name of the DAO.
	group   string                         // group is the database configuration group name of current DAO.
	columns TemplateVariablePresetsColumns // columns contains all the column names of Table for convenient usage.
}

// TemplateVariablePresetsColumns defines and stores column names for table template_variable_presets.
type TemplateVariablePresetsColumns struct {
	Id            string // 关联ID，自增主键
	TemplateId    string // 所属模板ID
	VarPresetId   string // 预设变量ID
	PresetPath    string // 预设中的变量路径 (如: user.name)
	MappedName    string // 映射到模板中的变量名
	IsActive      string // 是否启用此映射关系
	Sort          string // 排序
	CreatedAt     string // 创建时间
	UpdatedAt     string // 更新时间
}

// templateVariablePresetsColumns holds the columns for table template_variable_presets.
var templateVariablePresetsColumns = TemplateVariablePresetsColumns{
	Id:            "id",
	TemplateId:    "template_id",
	VarPresetId:   "var_preset_id",
	PresetPath:    "preset_path",
	MappedName:    "mapped_name",
	IsActive:      "is_active",
	Sort:          "sort",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

// NewTemplateVariablePresetsDao creates and returns a new DAO object for table data access.
func NewTemplateVariablePresetsDao() *TemplateVariablePresetsDao {
	return &TemplateVariablePresetsDao{
		group:   "default",
		table:   "template_variable_presets",
		columns: templateVariablePresetsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *TemplateVariablePresetsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *TemplateVariablePresetsDao) Table() string {
	return dao.table
}

// Columns returns the columns of current dao.
func (dao *TemplateVariablePresetsDao) Columns() TemplateVariablePresetsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *TemplateVariablePresetsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context.
func (dao *TemplateVariablePresetsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
func (dao *TemplateVariablePresetsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}