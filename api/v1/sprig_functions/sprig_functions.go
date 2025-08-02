package sprig_functions

import (
	"github.com/gogf/gf/v2/frame/g"
)

// GetSprigFunctionsReq 获取Sprig函数列表请求
type GetSprigFunctionsReq struct {
	g.Meta `path:"/sprig-functions" method:"get" summary:"获取Sprig函数列表" tags:"Sprig函数"`
}

// GetSprigFunctionsRes 获取Sprig函数列表响应
type GetSprigFunctionsRes struct {
	Categories []SprigFunctionCategory `json:"categories"`
	Total      int                     `json:"total"`
}

// SprigFunctionCategory Sprig函数分类
type SprigFunctionCategory struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Functions   []SprigFunction `json:"functions"`
}

// SprigFunction Sprig函数
type SprigFunction struct {
	Name        string                `json:"name"`
	DisplayName string                `json:"display_name"`
	Description string                `json:"description"`
	Params      []SprigFunctionParam  `json:"params"`
	ReturnType  string                `json:"return_type"`
	Category    string                `json:"category"`
	Example     string                `json:"example"`
	Usage       string                `json:"usage"`
	InsertText  string                `json:"insert_text"`
	Note        string                `json:"note,omitempty"`
	Aliases     []string              `json:"aliases,omitempty"`
}

// SprigFunctionParam Sprig函数参数
type SprigFunctionParam struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
	Default     string `json:"default,omitempty"`
}