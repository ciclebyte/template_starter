package template

import (
	commonApi "github.com/ciclebyte/template_starter/api/v1/common"
	model "github.com/ciclebyte/template_starter/internal/model"
	"github.com/gogf/gf/v2/frame/g"
)

type TemplateFilesAddReq struct {
	g.Meta      `path:"/templateFiles/add" method:"post" tags:"模板文件" summary:"模板文件-新增"`
	TemplateId  interface{} `json:"templateId" v:"required#所属模板ID不能为空"`
	FileName    string      `json:"fileName" v:"required#文件名不能为空"`
	FileContent string      `json:"fileContent"`
	FileSize    int         `json:"fileSize"`
	IsDirectory int         `json:"isDirectory" v:"required#是否为目录不能为空"`
	Md5         string      `json:"md5"`
	Sort        int         `json:"sort"`
	ParentId    int         `json:"parentId"`
}

type TemplateFilesAddRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

type TemplateFilesDelReq struct {
	g.Meta `path:"/templateFiles/del" method:"delete" tags:"模板文件" summary:"模板文件-删除"`
	Id     interface{} `json:"id" v:"required#id不能为空"`
}

type TemplateFilesDelRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

type TemplateFilesBatchDelReq struct {
	g.Meta `path:"/templateFiles/batchdel" method:"delete" tags:"模板文件" summary:"模板文件-批量删除"`
	Ids    []interface{} `json:"id" v:"required#id不能为空"`
}

type TemplateFilesBatchDelRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

type TemplateFilesEditReq struct {
	g.Meta      `path:"/templateFiles/edit" method:"put" tags:"模板文件" summary:"模板文件-修改"`
	Id          interface{} `json:"id" v:"required#文件ID，自增主键不能为空"`
	TemplateId  interface{} `json:"templateId" v:"required#所属模板ID不能为空"`
	FileName    string      `json:"fileName" v:"required#文件名不能为空"`
	FileContent string      `json:"fileContent"`
	FileSize    int         `json:"fileSize"`
	IsDirectory int         `json:"isDirectory" v:"required#是否为目录不能为空"`
	Md5         string      `json:"md5"`
	Sort        int         `json:"sort"`
	ParentId    int         `json:"parentId"`
}

type TemplateFilesEditRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

type TemplateFilesListReq struct {
	g.Meta `path:"/templateFiles/list" method:"get" tags:"模板文件" summary:"模板文件-列表"`
	commonApi.PageReq
	TemplateId  interface{} `json:"templateId" v:"required#所属模板ID不能为空"`
	FilePath    string      `json:"filePath" v:"required#文件路径（相对路径）不能为空"`
	FileName    string      `json:"fileName" v:"required#文件名不能为空"`
	FileContent string      `json:"fileContent" v:"required#文件内容不能为空"`
	FileSize    int         `json:"fileSize" v:"required#文件大小（字节）不能为空"`
	IsDirectory int         `json:"isDirectory" v:"required#是否为目录不能为空"`
	Md5         string      `json:"md5" v:"required#md5不能为空"`
	Sort        int         `json:"sort" v:"required#排序不能为空"`
	ParentId    int         `json:"parentId" v:"required#父目录ID，如果是文件则指向所属目录不能为空"`
}

type TemplateFilesListRes struct {
	g.Meta `mime:"application/json" example:"string"`
	commonApi.ListRes
	TemplateFilesList []*model.TemplateFilesInfo `json:"templateFilesList"`
}

type TemplateFilesDetailReq struct {
	g.Meta `path:"/templateFiles/detail" method:"get" tags:"模板文件" summary:"模板文件-详情"`
	Id     interface{} `json:"id" v:"required#id不能为空"`
}

type TemplateFilesDetailRes struct {
	g.Meta `mime:"application/json" example:"string"`
	*model.TemplateFilesInfo
}

type TemplatesFileTreeReq struct {
	g.Meta     `path:"/templateFiles/fileTree" method:"get" tags:"模板" summary:"模板-文件树"`
	TemplateId interface{} `json:"templateId" v:"required#模板ID不能为空"`
}

type FileTreeNode struct {
	Id          int64           `json:"id"`
	FilePath    string          `json:"filePath"`
	FileName    string          `json:"fileName"`
	IsDirectory int             `json:"isDirectory"`
	ParentId    int64           `json:"parentId"`
	FileSize    uint            `json:"fileSize"`
	Md5         string          `json:"md5"`
	Children    []*FileTreeNode `json:"children,omitempty"`
}

type TemplatesFileTreeRes struct {
	g.Meta `mime:"application/json"`
	Tree   []*FileTreeNode `json:"tree"`
}

// 仅获取文件内容
// GET /templateFiles/content?id=xxx
// 返回 { fileContent: "..." }
type TemplateFilesContentReq struct {
	g.Meta `path:"/templateFiles/content" method:"get" tags:"模板文件" summary:"模板文件-仅获取内容"`
	Id     interface{} `json:"id" v:"required#id不能为空"`
}

type TemplateFilesContentRes struct {
	g.Meta      `mime:"application/json" example:"string"`
	FileContent string `json:"fileContent"`
}

// ZIP 包上传接口
type TemplateFilesUploadZipReq struct {
	g.Meta     `path:"/templateFiles/uploadZip" method:"post" tags:"模板文件" summary:"模板文件-上传ZIP包"`
	TemplateId interface{} `json:"templateId" v:"required#所属模板ID不能为空"`
}

type TemplateFilesUploadZipRes struct {
	g.Meta       `mime:"application/json" example:"string"`
	SuccessCount int      `json:"successCount"` // 成功解压的文件数量
	FailedFiles  []string `json:"failedFiles"`  // 解压失败的文件列表
	Message      string   `json:"message"`      // 处理结果消息
}

// 重命名接口
type TemplateFilesRenameReq struct {
	g.Meta   `path:"/templateFiles/rename" method:"put" tags:"模板文件" summary:"模板文件-重命名"`
	Id       interface{} `json:"id" v:"required#文件ID不能为空"`
	FileName string      `json:"fileName" v:"required#新文件名不能为空"`
}

type TemplateFilesRenameRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

// 上传代码文件接口
type TemplateFilesUploadCodeReq struct {
	g.Meta     `path:"/templateFiles/uploadCode" method:"post" tags:"模板文件" summary:"模板文件-上传代码文件"`
	TemplateId interface{} `json:"templateId" v:"required#所属模板ID不能为空"`
	ParentId   interface{} `json:"parentId"` // 可选的父目录ID
}

type TemplateFilesUploadCodeRes struct {
	g.Meta      `mime:"application/json" example:"string"`
	FileName    string `json:"fileName"`    // 上传的文件名
	FileSize    int    `json:"fileSize"`    // 文件大小
	IsTextFile  bool   `json:"isTextFile"`  // 是否为文本文件
	FileContent string `json:"fileContent"` // 文件内容（如果是文本文件）
	Message     string `json:"message"`     // 处理结果消息
}

// 模板文件渲染接口
type TemplateFilesRenderReq struct {
	g.Meta        `path:"/templateFiles/render" method:"post" tags:"模板文件" summary:"模板文件-渲染"`
	FileId        interface{}            `json:"fileId" v:"required#文件ID不能为空"`
	TestVariables map[string]interface{} `json:"testVariables"` // 测试变量值
}

type TemplateFilesRenderRes struct {
	g.Meta      `mime:"application/json" example:"string"`
	FileId      int64                  `json:"fileId"`
	FileName    string                 `json:"fileName"`
	FileContent string                 `json:"fileContent"` // 渲染后的内容
	Variables   map[string]interface{} `json:"variables"`   // 使用的变量
}
