package model

import "github.com/gogf/gf/v2/os/gtime"

type TemplateFilesInfo struct {
	Id          int64       `orm:"id"  json:"id"`                    // 文件ID，自增主键
	TemplateId  int64       `orm:"template_id"  json:"templateId"`   // 所属模板ID
	FilePath    string      `orm:"file_path"  json:"filePath"`       // 文件路径（相对路径）
	FileContent string      `orm:"file_content"  json:"fileContent"` // 文件内容
	FileSize    int         `orm:"file_size"  json:"fileSize"`       // 文件大小（字节）
	IsDirectory int         `orm:"is_directory"  json:"isDirectory"` // 是否为目录
	Md5         string      `orm:"md5"  json:"md5"`                  // md5
	Sort        int         `orm:"sort"  json:"sort"`                // 排序
	ParentId    int         `orm:"parent_id"  json:"parentId"`       // 父目录ID，如果是文件则指向所属目录
	CreatedAt   *gtime.Time `orm:"created_at"  json:"createdAt"`     // 记录创建时间
	UpdatedAt   *gtime.Time `orm:"updated_at"  json:"updatedAt"`     // None
}
