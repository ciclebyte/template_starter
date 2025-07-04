// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TemplateFilesDao is the data access object for table template_files.
type TemplateFilesDao struct {
	table   string               // table is the underlying table name of the DAO.
	group   string               // group is the database configuration group name of current DAO.
	columns TemplateFilesColumns // columns contains all the column names of Table for convenient usage.
}

// TemplateFilesColumns defines and stores column names for table template_files.
type TemplateFilesColumns struct {
	Id          string // 文件ID，自增主键
	TemplateId  string // 所属模板ID
	FilePath    string // 文件路径（相对路径）
	FileContent string // 文件内容
	FileSize    string // 文件大小（字节）
	IsDirectory string // 是否为目录
	Md5         string // md5
	Sort        string // 排序
	ParentId    string // 父目录ID，如果是文件则指向所属目录
	CreatedAt   string // 记录创建时间
	UpdatedAt   string //
}

// templateFilesColumns holds the columns for table template_files.
var templateFilesColumns = TemplateFilesColumns{
	Id:          "id",
	TemplateId:  "template_id",
	FilePath:    "file_path",
	FileContent: "file_content",
	FileSize:    "file_size",
	IsDirectory: "is_directory",
	Md5:         "md5",
	Sort:        "sort",
	ParentId:    "parent_id",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// NewTemplateFilesDao creates and returns a new DAO object for table data access.
func NewTemplateFilesDao() *TemplateFilesDao {
	return &TemplateFilesDao{
		group:   "default",
		table:   "template_files",
		columns: templateFilesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *TemplateFilesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *TemplateFilesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *TemplateFilesDao) Columns() TemplateFilesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *TemplateFilesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *TemplateFilesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *TemplateFilesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
