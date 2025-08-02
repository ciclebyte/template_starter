package builtin_functions

import (
	"github.com/gogf/gf/v2/frame/g"
)

// GetBuiltinFunctionsReq 获取内置函数列表请求
type GetBuiltinFunctionsReq struct {
	g.Meta `path:"/builtin-functions" method:"get" summary:"获取内置函数列表" tags:"内置函数"`
}

// GetBuiltinFunctionsRes 获取内置函数列表响应
type GetBuiltinFunctionsRes struct {
	Categories []BuiltinFunctionCategory `json:"categories"`
	Total      int                      `json:"total"`
}

// BuiltinFunctionCategory 内置函数分类
type BuiltinFunctionCategory struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Functions   []BuiltinFunction `json:"functions"`
}

// BuiltinFunction 内置函数
type BuiltinFunction struct {
	Name        string          `json:"name"`
	DisplayName string          `json:"display_name"`
	Description string          `json:"description"`
	Params      []FunctionParam `json:"params"`
	ReturnType  string          `json:"return_type"`
	Category    string          `json:"category"`
	Example     string          `json:"example"`
	InsertText  string          `json:"insert_text"`
}

// FunctionParam 函数参数
type FunctionParam struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
	Default     string `json:"default,omitempty"`
}