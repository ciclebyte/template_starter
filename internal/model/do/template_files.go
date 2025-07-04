// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TemplateFiles is the golang structure of table template_files for DAO operations like Where/Data.
type TemplateFiles struct {
	g.Meta      `orm:"table:template_files, do:true"`
	Id          interface{} // 文件ID，自增主键
	TemplateId  interface{} // 所属模板ID
	FilePath    interface{} // 文件路径（相对路径）
	FileContent interface{} // 文件内容
	FileSize    interface{} // 文件大小（字节）
	IsDirectory interface{} // 是否为目录
	Md5         interface{} // md5
	Sort        interface{} // 排序
	ParentId    interface{} // 父目录ID，如果是文件则指向所属目录
	CreatedAt   *gtime.Time // 记录创建时间
	UpdatedAt   *gtime.Time //
}
