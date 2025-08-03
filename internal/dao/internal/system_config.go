// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SystemConfigDao is the data access object for table system_config.
type SystemConfigDao struct {
	table   string              // table is the underlying table name of the DAO.
	group   string              // group is the database configuration group name of current DAO.
	columns SystemConfigColumns // columns contains all the column names of Table for convenient usage.
}

// SystemConfigColumns defines and stores column names for table system_config.
type SystemConfigColumns struct {
	Id             string // 配置ID，自增主键
	ConfigKey      string // 配置键名，如 ai.openai.api_key
	ConfigValue    string // 配置值，支持JSON格式存储复杂配置
	ConfigGroup    string // 配置分组，如 ai, system, template, ui
	ConfigType     string // 配置类型：string, number, boolean, json, array
	DisplayName    string // 显示名称，用于前端展示
	Description    string // 配置描述和说明
	IsPublic       string // 是否为公开配置（前端可见）
	IsRequired     string // 是否为必填配置
	DefaultValue   string // 默认值
	ValidationRule string // 验证规则，JSON格式
	SortOrder      string // 排序顺序
	Status         string // 状态：1启用，0禁用
	CreatedAt      string // 创建时间
	UpdatedAt      string // 更新时间
}

// systemConfigColumns holds the columns for table system_config.
var systemConfigColumns = SystemConfigColumns{
	Id:             "id",
	ConfigKey:      "config_key",
	ConfigValue:    "config_value",
	ConfigGroup:    "config_group",
	ConfigType:     "config_type",
	DisplayName:    "display_name",
	Description:    "description",
	IsPublic:       "is_public",
	IsRequired:     "is_required",
	DefaultValue:   "default_value",
	ValidationRule: "validation_rule",
	SortOrder:      "sort_order",
	Status:         "status",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// NewSystemConfigDao creates and returns a new DAO object for table data access.
func NewSystemConfigDao() *SystemConfigDao {
	return &SystemConfigDao{
		group:   "default",
		table:   "system_config",
		columns: systemConfigColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *SystemConfigDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *SystemConfigDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *SystemConfigDao) Columns() SystemConfigColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *SystemConfigDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *SystemConfigDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *SystemConfigDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
