package sprig_functions

import "github.com/ciclebyte/template_starter/internal/model"

// buildEncodingFunctions 构建编码函数分类
func (s *sSprigFunctions) buildEncodingFunctions() model.SprigFunctionCategory {
	return model.SprigFunctionCategory{
		Name:        "编码函数",
		Description: "编码和解码转换函数",
		Functions: []model.SprigFunction{
			{
				Name:        "b64enc",
				DisplayName: "Base64编码",
				Description: "使用Base64编码字符串",
				Params: []model.SprigFunctionParam{
					{Name: "string", Type: "string", Required: true, Description: "要编码的字符串"},
				},
				ReturnType: "string",
				Category:   "编码函数",
				Example:    `{{ "hello" | b64enc }}`,
				Usage:      "b64enc 函数使用Base64编码字符串。常用于数据传输和存储中的文本编码。",
				InsertText: `{{ . | b64enc }}`,
				Note:       "输出符合RFC 4648标准的Base64编码",
			},
			{
				Name:        "b64dec",
				DisplayName: "Base64解码",
				Description: "使用Base64解码字符串",
				Params: []model.SprigFunctionParam{
					{Name: "string", Type: "string", Required: true, Description: "要解码的Base64字符串"},
				},
				ReturnType: "string",
				Category:   "编码函数",
				Example:    `{{ "aGVsbG8=" | b64dec }}`,
				Usage:      "b64dec 函数使用Base64解码字符串。用于解码Base64编码的数据。",
				InsertText: `{{ . | b64dec }}`,
				Note:       "输入必须是有效的Base64编码字符串",
			},
			{
				Name:        "b32enc",
				DisplayName: "Base32编码",
				Description: "使用Base32编码字符串",
				Params: []model.SprigFunctionParam{
					{Name: "string", Type: "string", Required: true, Description: "要编码的字符串"},
				},
				ReturnType: "string",
				Category:   "编码函数",
				Example:    `{{ "hello" | b32enc }}`,
				Usage:      "b32enc 函数使用Base32编码字符串。Base32编码使用32个可打印字符，比Base64更易读但效率略低。",
				InsertText: `{{ . | b32enc }}`,
				Note:       "输出符合RFC 4648标准的Base32编码",
			},
			{
				Name:        "b32dec",
				DisplayName: "Base32解码",
				Description: "使用Base32解码字符串",
				Params: []model.SprigFunctionParam{
					{Name: "string", Type: "string", Required: true, Description: "要解码的Base32字符串"},
				},
				ReturnType: "string",
				Category:   "编码函数",
				Example:    `{{ "NBSWY3DP" | b32dec }}`,
				Usage:      "b32dec 函数使用Base32解码字符串。用于解码Base32编码的数据。",
				InsertText: `{{ . | b32dec }}`,
				Note:       "输入必须是有效的Base32编码字符串",
			},
		},
	}
}