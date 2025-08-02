package model

// SprigFunctionsResponse Sprig函数响应
type SprigFunctionsResponse struct {
	Categories []SprigFunctionCategory `json:"categories"`
	Total      int                     `json:"total"`
}

// SprigFunctionCategory Sprig函数分类
type SprigFunctionCategory struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Functions   []SprigFunction `json:"functions"`
}

// SprigFunction Sprig函数详细信息
type SprigFunction struct {
	Name        string                `json:"name"`
	DisplayName string                `json:"display_name"`
	Description string                `json:"description"`
	Params      []SprigFunctionParam  `json:"params"`
	ReturnType  string                `json:"return_type"`
	Category    string                `json:"category"`
	Example     string                `json:"example"`
	Usage       string                `json:"usage"`      // 详细用法说明
	InsertText  string                `json:"insert_text"`
	Note        string                `json:"note,omitempty"`      // 注意事项
	Aliases     []string              `json:"aliases,omitempty"`   // 别名函数
}

// SprigFunctionParam Sprig函数参数
type SprigFunctionParam struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
	Default     string `json:"default,omitempty"`
}