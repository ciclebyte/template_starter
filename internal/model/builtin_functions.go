package model

// BuiltinFunction 内置函数定义
type BuiltinFunction struct {
	Name        string           `json:"name"`
	DisplayName string           `json:"display_name"`
	Description string           `json:"description"`
	Params      []FunctionParam  `json:"params"`
	ReturnType  string           `json:"return_type"`
	Category    string           `json:"category"`
	Example     string           `json:"example"`
	InsertText  string           `json:"insert_text"`
}

// FunctionParam 函数参数定义
type FunctionParam struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
	Default     string `json:"default,omitempty"`
}

// BuiltinFunctionCategory 内置函数分类
type BuiltinFunctionCategory struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Functions   []BuiltinFunction `json:"functions"`
}

// BuiltinFunctionsResponse API响应
type BuiltinFunctionsResponse struct {
	Categories []BuiltinFunctionCategory `json:"categories"`
	Total      int                      `json:"total"`
}