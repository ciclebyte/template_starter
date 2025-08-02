package sprig_functions

import "github.com/ciclebyte/template_starter/internal/model"

// buildTypeConversionFunctions 构建类型转换函数分类
func (s *sSprigFunctions) buildTypeConversionFunctions() model.SprigFunctionCategory {
	return model.SprigFunctionCategory{
		Name:        "类型转换函数",
		Description: "数据类型转换和格式化函数",
		Functions: []model.SprigFunction{
			{
				Name:        "atoi",
				DisplayName: "字符串转整数",
				Description: "将字符串转换为整数",
				Params: []model.SprigFunctionParam{
					{Name: "string", Type: "string", Required: true, Description: "要转换的字符串"},
				},
				ReturnType: "int",
				Category:   "类型转换函数",
				Example:    `{{ atoi "123" }}`,
				Usage:      "atoi 函数将字符串转换为整数。只有 atoi 要求输入必须是特定类型（字符串）。",
				InsertText: `{{ atoi . }}`,
				Note:       "要求输入必须是字符串类型",
			},
			{
				Name:        "float64",
				DisplayName: "转换为浮点数",
				Description: "将值转换为float64类型",
				Params: []model.SprigFunctionParam{
					{Name: "value", Type: "any", Required: true, Description: "要转换的值"},
				},
				ReturnType: "float64",
				Category:   "类型转换函数",
				Example:    `{{ float64 "123.45" }}`,
				Usage:      "float64 函数将值转换为 float64 类型。会尝试从任何类型转换到目标类型。",
				InsertText: `{{ float64 . }}`,
				Note:       "会尝试从任何类型转换",
			},
			{
				Name:        "int",
				DisplayName: "转换为整数",
				Description: "将值转换为系统位宽的int类型",
				Params: []model.SprigFunctionParam{
					{Name: "value", Type: "any", Required: true, Description: "要转换的值"},
				},
				ReturnType: "int",
				Category:   "类型转换函数",
				Example:    `{{ int "123" }}`,
				Usage:      "int 函数将值转换为系统位宽的 int 类型。会尝试从任何类型转换到目标类型，例如可以将浮点数转换为整数，也可以将字符串转换为整数。",
				InsertText: `{{ int . }}`,
				Note:       "系统位宽的int类型，支持多种类型转换",
			},
			{
				Name:        "int64",
				DisplayName: "转换为64位整数",
				Description: "将值转换为int64类型",
				Params: []model.SprigFunctionParam{
					{Name: "value", Type: "any", Required: true, Description: "要转换的值"},
				},
				ReturnType: "int64",
				Category:   "类型转换函数",
				Example:    `{{ int64 "123" }}`,
				Usage:      "int64 函数将值转换为 int64 类型。会尝试从任何类型转换到目标类型，例如可以将浮点数转换为整数，也可以将字符串转换为整数。",
				InsertText: `{{ int64 . }}`,
				Note:       "64位整数类型，支持多种类型转换",
			},
			{
				Name:        "toString",
				DisplayName: "转换为字符串",
				Description: "将值转换为字符串",
				Params: []model.SprigFunctionParam{
					{Name: "value", Type: "any", Required: true, Description: "要转换的值"},
				},
				ReturnType: "string",
				Category:   "类型转换函数",
				Example:    `{{ toString 123 }}`,
				Usage:      "toString 函数将值转换为字符串。会尝试从任何类型转换到字符串类型。",
				InsertText: `{{ toString . }}`,
				Note:       "会尝试从任何类型转换为字符串",
			},
			{
				Name:        "toStrings",
				DisplayName: "转换为字符串列表",
				Description: "将列表、切片或数组转换为字符串列表",
				Params: []model.SprigFunctionParam{
					{Name: "collection", Type: "any", Required: true, Description: "要转换的集合（列表、切片或数组）"},
				},
				ReturnType: "[]string",
				Category:   "类型转换函数",
				Example:    `{{ toStrings (list 1 2 3) }}`,
				Usage:      "toStrings 函数将类似列表的集合转换为字符串切片。例如将 1 转换为 \"1\"，将 2 转换为 \"2\"，然后作为列表返回。",
				InsertText: `{{ toStrings . }}`,
				Note:       "将集合中每个元素转换为字符串",
			},
			{
				Name:        "toDecimal",
				DisplayName: "八进制转十进制",
				Description: "将Unix八进制权限转换为十进制",
				Params: []model.SprigFunctionParam{
					{Name: "octal", Type: "any", Required: true, Description: "八进制值"},
				},
				ReturnType: "int64",
				Category:   "类型转换函数",
				Example:    `{{ toDecimal 0777 }}`,
				Usage:      "toDecimal 函数将Unix八进制权限转换为十进制。例如将 0777 转换为 511，并返回 int64 值。",
				InsertText: `{{ toDecimal . }}`,
				Note:       "常用于Unix文件权限转换",
			},
		},
	}
}