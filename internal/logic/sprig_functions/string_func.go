package sprig_functions

import "github.com/ciclebyte/template_starter/internal/model"

// buildStringFunctions 构建字符串函数分类
func (s *sSprigFunctions) buildStringFunctions() model.SprigFunctionCategory {
	return model.SprigFunctionCategory{
		Name:        "字符串函数",
		Description: "字符串操作和处理函数",
		Functions: []model.SprigFunction{
			{
				Name:        "upper",
				DisplayName: "转大写",
				Description: "将字符串转换为大写字母",
				Params: []model.SprigFunctionParam{
					{Name: "string", Type: "string", Required: true, Description: "要转换的字符串"},
				},
				ReturnType: "string",
				Category:   "字符串函数",
				Example:    `{{ "hello" | upper }}`,
				Usage:      "upper 函数将整个字符串转换为大写字母。这对于格式化输出或标准化字符串非常有用。",
				InsertText: `{{ . | upper }}`,
				Note:       "支持Unicode字符",
			},
			{
				Name:        "lower",
				DisplayName: "转小写",
				Description: "将字符串转换为小写字母",
				Params: []model.SprigFunctionParam{
					{Name: "string", Type: "string", Required: true, Description: "要转换的字符串"},
				},
				ReturnType: "string",
				Category:   "字符串函数",
				Example:    `{{ "HELLO" | lower }}`,
				Usage:      "lower 函数将整个字符串转换为小写字母。常用于标准化用户输入或生成统一格式的标识符。",
				InsertText: `{{ . | lower }}`,
				Note:       "支持Unicode字符",
			},
			{
				Name:        "title",
				DisplayName: "首字母大写",
				Description: "将字符串转换为标题格式，每个单词的首字母大写",
				Params: []model.SprigFunctionParam{
					{Name: "string", Type: "string", Required: true, Description: "要转换的字符串"},
				},
				ReturnType: "string",
				Category:   "字符串函数",
				Example:    `{{ "hello world" | title }}`,
				Usage:      "title 函数将字符串转换为标题格式，即每个单词的首字母大写，其余字母小写。适用于生成标题、姓名等格式化文本。",
				InsertText: `{{ . | title }}`,
			},
			{
				Name:        "trim",
				DisplayName: "去除空格",
				Description: "去除字符串两端的空白字符",
				Params: []model.SprigFunctionParam{
					{Name: "string", Type: "string", Required: true, Description: "要处理的字符串"},
				},
				ReturnType: "string",
				Category:   "字符串函数",
				Example:    `{{ "  hello world  " | trim }}`,
				Usage:      "trim 函数移除字符串两端的空白字符（空格、制表符、换行符等）。这在处理用户输入或清理数据时非常有用。",
				InsertText: `{{ . | trim }}`,
				Note:       "不会移除字符串中间的空白字符",
			},
			{
				Name:        "trimPrefix",
				DisplayName: "去除前缀",
				Description: "从字符串开头移除指定的前缀",
				Params: []model.SprigFunctionParam{
					{Name: "prefix", Type: "string", Required: true, Description: "要移除的前缀"},
					{Name: "string", Type: "string", Required: true, Description: "原字符串"},
				},
				ReturnType: "string",
				Category:   "字符串函数",
				Example:    `{{ trimPrefix "prefix_" "prefix_filename.txt" }}`,
				Usage:      "trimPrefix 函数仅从字符串开头移除指定的前缀。如果字符串不以该前缀开头，则原样返回。常用于去除文件名前缀或URL协议等。",
				InsertText: `{{ trimPrefix "prefix" . }}`,
			},
			{
				Name:        "trimSuffix",
				DisplayName: "去除后缀",
				Description: "从字符串末尾移除指定的后缀",
				Params: []model.SprigFunctionParam{
					{Name: "suffix", Type: "string", Required: true, Description: "要移除的后缀"},
					{Name: "string", Type: "string", Required: true, Description: "原字符串"},
				},
				ReturnType: "string",
				Category:   "字符串函数",
				Example:    `{{ trimSuffix ".txt" "document.txt" }}`,
				Usage:      "trimSuffix 函数仅从字符串末尾移除指定的后缀。如果字符串不以该后缀结尾，则原样返回。常用于去除文件扩展名或URL参数等。",
				InsertText: `{{ trimSuffix "suffix" . }}`,
			},
			{
				Name:        "replace",
				DisplayName: "替换字符串",
				Description: "将字符串中的指定内容替换为新内容",
				Params: []model.SprigFunctionParam{
					{Name: "old", Type: "string", Required: true, Description: "要替换的内容"},
					{Name: "new", Type: "string", Required: true, Description: "替换为的内容"},
					{Name: "string", Type: "string", Required: true, Description: "原字符串"},
				},
				ReturnType: "string",
				Category:   "字符串函数",
				Example:    `{{ replace "world" "Go" "hello world" }}`,
				Usage:      "replace 函数将字符串中所有出现的指定内容替换为新内容。这是一个非常常用的字符串处理函数。",
				InsertText: `{{ replace "old" "new" . }}`,
				Note:       "会替换所有匹配的内容",
			},
			{
				Name:        "contains",
				DisplayName: "包含检查",
				Description: "检查字符串是否包含指定的子字符串",
				Params: []model.SprigFunctionParam{
					{Name: "substr", Type: "string", Required: true, Description: "要查找的子字符串"},
					{Name: "string", Type: "string", Required: true, Description: "原字符串"},
				},
				ReturnType: "bool",
				Category:   "字符串函数",
				Example:    `{{ if contains "go" "golang" }}包含{{ end }}`,
				Usage:      "contains 函数检查字符串是否包含指定的子字符串，返回布尔值。常用于条件判断。",
				InsertText: `{{ contains "substr" . }}`,
			},
			{
				Name:        "hasPrefix",
				DisplayName: "前缀检查",
				Description: "检查字符串是否以指定前缀开头",
				Params: []model.SprigFunctionParam{
					{Name: "prefix", Type: "string", Required: true, Description: "要检查的前缀"},
					{Name: "string", Type: "string", Required: true, Description: "原字符串"},
				},
				ReturnType: "bool",
				Category:   "字符串函数",
				Example:    `{{ if hasPrefix "http" .URL }}是HTTP链接{{ end }}`,
				Usage:      "hasPrefix 函数检查字符串是否以指定前缀开头，返回布尔值。常用于URL协议检查、文件类型判断等。",
				InsertText: `{{ hasPrefix "prefix" . }}`,
			},
			{
				Name:        "hasSuffix",
				DisplayName: "后缀检查",
				Description: "检查字符串是否以指定后缀结尾",
				Params: []model.SprigFunctionParam{
					{Name: "suffix", Type: "string", Required: true, Description: "要检查的后缀"},
					{Name: "string", Type: "string", Required: true, Description: "原字符串"},
				},
				ReturnType: "bool",
				Category:   "字符串函数",
				Example:    `{{ if hasSuffix ".go" .Filename }}是Go文件{{ end }}`,
				Usage:      "hasSuffix 函数检查字符串是否以指定后缀结尾，返回布尔值。常用于文件扩展名检查、URL路径判断等。",
				InsertText: `{{ hasSuffix "suffix" . }}`,
			},
		},
	}
}
