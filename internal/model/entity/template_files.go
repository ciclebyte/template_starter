// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TemplateFiles is the golang structure for table template_files.
type TemplateFiles struct {
	Id          int64       `json:"id"          description:"文件ID，自增主键"`
	TemplateId  int64       `json:"templateId"  description:"所属模板ID"`
	FilePath    string      `json:"filePath"    description:"文件路径（相对路径）"`
	FileContent string      `json:"fileContent" description:"文件内容"`
	FileSize    uint        `json:"fileSize"    description:"文件大小（字节）"`
	IsDirectory int         `json:"isDirectory" description:"是否为目录"`
	Md5         string      `json:"md5"         description:"md5"`
	Sort        int         `json:"sort"        description:"排序"`
	ParentId    int64       `json:"parentId"    description:"父目录ID，如果是文件则指向所属目录"`
	CreatedAt   *gtime.Time `json:"createdAt"   description:"记录创建时间"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   description:""`
}
