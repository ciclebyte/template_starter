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
	FileContent string      `json:"fileContent" v:"required#文件内容不能为空"`
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
	g.Meta    `path:"/templateFiles/render" method:"post" tags:"模板文件" summary:"模板文件-渲染"`
	FileId    interface{}            `json:"fileId" v:"required#文件ID不能为空"`
	Variables map[string]interface{} `json:"variables"` // 变量值
}

type TemplateFilesRenderRes struct {
	g.Meta      `mime:"application/json" example:"string"`
	FileId      int64                  `json:"fileId"`
	FileName    string                 `json:"fileName"`
	FileContent string                 `json:"fileContent"` // 渲染后的内容
	Variables   map[string]interface{} `json:"variables"`   // 使用的变量
	Success     bool                   `json:"success"`     // 渲染是否成功
	Error       *TemplateRenderError   `json:"error,omitempty"` // 渲染错误详情
}

// 模板渲染错误详情
type TemplateRenderError struct {
	Type        string `json:"type"`        // 错误类型: "parse_error", "execute_error", "variable_error"
	Message     string `json:"message"`     // 错误消息
	Line        int    `json:"line"`        // 错误行号
	Column      int    `json:"column"`      // 错误列号
	Context     string `json:"context"`     // 错误上下文
	Suggestion  string `json:"suggestion"`  // 修复建议
}

// 渲染文件树接口
type TemplateFilesRenderFileTreeReq struct {
	g.Meta     `path:"/templateFiles/renderFileTree" method:"post" tags:"模板文件" summary:"模板文件-渲染文件树"`
	TemplateId interface{}            `json:"templateId" v:"required#模板ID不能为空"`
	Variables  map[string]interface{} `json:"variables"` // 变量值
}

type TemplateFilesRenderFileTreeRes struct {
	g.Meta     `mime:"application/json" example:"string"`
	TemplateId int64                  `json:"templateId"`
	Tree       []*RenderFileInfo      `json:"tree"`       // 渲染后的文件树
	Variables  map[string]interface{} `json:"variables"`  // 使用的变量
	TotalFiles int                    `json:"totalFiles"` // 总文件数
	TotalSize  int64                  `json:"totalSize"`  // 总文件大小
}

// 渲染后的文件信息
type RenderFileInfo struct {
	Id          int64             `json:"id"`          // 文件ID
	FilePath    string            `json:"filePath"`    // 文件路径
	FileName    string            `json:"fileName"`    // 文件名
	FileContent string            `json:"fileContent"` // 渲染后的文件内容
	FileSize    int               `json:"fileSize"`    // 文件大小
	IsDirectory int               `json:"isDirectory"` // 是否为目录
	ParentId    int               `json:"parentId"`    // 父目录ID
	Children    []*RenderFileInfo `json:"children"`    // 子文件/目录
}

// 模板文件ZIP下载接口
type TemplateFilesDownloadZipReq struct {
	g.Meta     `path:"/templateFiles/downloadZip" method:"post" tags:"模板文件" summary:"模板文件-下载ZIP包"`
	TemplateId interface{}            `json:"templateId" v:"required#模板ID不能为空"`
	Variables  map[string]interface{} `json:"variables"` // 变量值
	FileName   string                 `json:"fileName"`  // 可选的ZIP文件名，默认为模板名
}

type TemplateFilesDownloadZipRes struct {
	g.Meta `mime:"application/zip" example:"string"`
	// 直接返回ZIP文件流
}

// 在现有的结构体后面添加
type TemplateFilesMoveReq struct {
	g.Meta      `path:"/templateFiles/move" method:"put" tags:"模板文件" summary:"模板文件-移动"`
	Id          interface{} `json:"id" v:"required#文件ID不能为空"`
	NewParentId interface{} `json:"newParentId"` // 新父目录ID，null或0表示移动到根目录
}

type TemplateFilesMoveRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

// 设置文件生成条件接口
type TemplateFilesSetConditionReq struct {
	g.Meta        `path:"/templateFiles/setCondition" method:"put" tags:"模板文件" summary:"模板文件-设置生成条件"`
	Id            interface{} `json:"id" v:"required#文件ID不能为空"`
	Enabled       bool        `json:"enabled"`                                    // 是否启用条件
	VariableName  string      `json:"variableName" v:"required-if:enabled,true"`  // 关联变量名
	ExpectedValue bool        `json:"expectedValue"`                              // 期望值
	Description   string      `json:"description"`                                // 条件描述
}

type TemplateFilesSetConditionRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

// 获取文件生成条件接口
type TemplateFilesGetConditionReq struct {
	g.Meta `path:"/templateFiles/getCondition" method:"get" tags:"模板文件" summary:"模板文件-获取生成条件"`
	Id     interface{} `json:"id" v:"required#文件ID不能为空"`
}

type TemplateFilesGetConditionRes struct {
	g.Meta        `mime:"application/json" example:"string"`
	Enabled       bool   `json:"enabled"`       // 是否启用条件
	VariableName  string `json:"variableName"`  // 关联变量名  
	ExpectedValue bool   `json:"expectedValue"` // 期望值
	Description   string `json:"description"`   // 条件描述
}
