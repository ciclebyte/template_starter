// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// VarPresetDao is the data access object for table var_preset.
type VarPresetDao struct {
	table   string             // table is the underlying table name of the DAO.
	group   string             // group is the database configuration group name of current DAO.
	columns VarPresetColumns   // columns contains all the column names of Table for convenient usage.
}

// VarPresetColumns defines and stores column names for table var_preset.
type VarPresetColumns struct {
	Id              string // 预设ID
	Name            string // 预设名称，如MySQL、Fields、OpenAPI
	DisplayName     string // 显示名称
	Description     string // 预设描述
	Category        string // 分类：system=系统预置，custom=用户自定义
	SchemaJson      string // 数据结构模板定义（JSON Schema）
	DefaultDataJson string // 默认数据示例
	Icon            string // 图标名称
	Sort            string // 排序权重
	IsEnabled       string // 是否启用
	Version         string // 版本号
	CreatedBy       string // 创建者ID（系统预置为空）
	CreatedAt       string // 创建时间
	UpdatedAt       string // 更新时间
	DeletedAt       string // 删除时间
}

// varPresetColumns holds the columns for table var_preset.
var varPresetColumns = VarPresetColumns{
	Id:              "id",
	Name:            "name",
	DisplayName:     "display_name",
	Description:     "description",
	Category:        "category",
	SchemaJson:      "schema_json",
	DefaultDataJson: "default_data_json",
	Icon:            "icon",
	Sort:            "sort",
	IsEnabled:       "is_enabled",
	Version:         "version",
	CreatedBy:       "created_by",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
	DeletedAt:       "deleted_at",
}

// NewVarPresetDao creates and returns a new DAO object for table data access.
func NewVarPresetDao() *VarPresetDao {
	return &VarPresetDao{
		group:   "default",
		table:   "var_preset",
		columns: varPresetColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *VarPresetDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *VarPresetDao) Table() string {
	return dao.table
}

// Columns returns the columns of current dao.
func (dao *VarPresetDao) Columns() VarPresetColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *VarPresetDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context.
func (dao *VarPresetDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
func (dao *VarPresetDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}