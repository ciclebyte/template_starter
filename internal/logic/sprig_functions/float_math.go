package sprig_functions

import "github.com/ciclebyte/template_starter/internal/model"

// buildFloatMathFunctions 构建浮点数数学函数分类
func (s *sSprigFunctions) buildFloatMathFunctions() model.SprigFunctionCategory {
	return model.SprigFunctionCategory{
		Name:        "浮点数数学函数",
		Description: "浮点数数学运算和计算函数",
		Functions: []model.SprigFunction{
			{
				Name:        "addf",
				DisplayName: "浮点数相加",
				Description: "对浮点数求和",
				Params: []model.SprigFunctionParam{
					{Name: "numbers", Type: "...float64", Required: true, Description: "要相加的浮点数"},
				},
				ReturnType: "float64",
				Category:   "浮点数数学函数",
				Example:    `{{ addf 1.5 2 2 }}`,
				Usage:      "addf 函数对浮点数求和。例如 addf 1.5 2 2 返回 5.5。",
				InsertText: `{{ addf 1.5 2.5 }}`,
			},
			{
				Name:        "add1f",
				DisplayName: "浮点数加一",
				Description: "将浮点数增加1",
				Params: []model.SprigFunctionParam{
					{Name: "number", Type: "float64", Required: true, Description: "要增加的浮点数"},
				},
				ReturnType: "float64",
				Category:   "浮点数数学函数",
				Example:    `{{ add1f 2.5 }}`,
				Usage:      "add1f 函数将浮点数增加1。",
				InsertText: `{{ add1f . }}`,
			},
			{
				Name:        "subf",
				DisplayName: "浮点数相减",
				Description: "浮点数减法运算",
				Params: []model.SprigFunctionParam{
					{Name: "numbers", Type: "...float64", Required: true, Description: "要相减的浮点数"},
				},
				ReturnType: "float64",
				Category:   "浮点数数学函数",
				Example:    `{{ subf 7.5 2 3 }}`,
				Usage:      "subf 函数执行浮点数减法运算。例如 subf 7.5 2 3 等同于 7.5 - 2 - 3，返回 2.5。",
				InsertText: `{{ subf 10.5 3.2 }}`,
			},
			{
				Name:        "divf",
				DisplayName: "浮点数除法",
				Description: "执行浮点数除法",
				Params: []model.SprigFunctionParam{
					{Name: "numbers", Type: "...float64", Required: true, Description: "要相除的浮点数"},
				},
				ReturnType: "float64",
				Category:   "浮点数数学函数",
				Example:    `{{ divf 10 2 4 }}`,
				Usage:      "divf 函数执行浮点数除法运算。例如 divf 10 2 4 等同于 10 / 2 / 4，返回 1.25。",
				InsertText: `{{ divf 10.0 2.5 }}`,
			},
			{
				Name:        "mulf",
				DisplayName: "浮点数相乘",
				Description: "浮点数乘法运算",
				Params: []model.SprigFunctionParam{
					{Name: "numbers", Type: "...float64", Required: true, Description: "要相乘的浮点数"},
				},
				ReturnType: "float64",
				Category:   "浮点数数学函数",
				Example:    `{{ mulf 1.5 2 2 }}`,
				Usage:      "mulf 函数对浮点数相乘。例如 mulf 1.5 2 2 返回 6。",
				InsertText: `{{ mulf 2.5 3.0 }}`,
			},
			{
				Name:        "maxf",
				DisplayName: "浮点数最大值",
				Description: "返回一系列浮点数中的最大值",
				Params: []model.SprigFunctionParam{
					{Name: "numbers", Type: "...float64", Required: true, Description: "要比较的浮点数"},
				},
				ReturnType: "float64",
				Category:   "浮点数数学函数",
				Example:    `{{ maxf 1 2.5 3 }}`,
				Usage:      "maxf 函数返回一系列浮点数中的最大值。例如 maxf 1 2.5 3 返回 3。",
				InsertText: `{{ maxf 1.5 2.8 3.2 }}`,
			},
			{
				Name:        "minf",
				DisplayName: "浮点数最小值",
				Description: "返回一系列浮点数中的最小值",
				Params: []model.SprigFunctionParam{
					{Name: "numbers", Type: "...float64", Required: true, Description: "要比较的浮点数"},
				},
				ReturnType: "float64",
				Category:   "浮点数数学函数",
				Example:    `{{ minf 1.5 2 3 }}`,
				Usage:      "minf 函数返回一系列浮点数中的最小值。例如 minf 1.5 2 3 返回 1.5。",
				InsertText: `{{ minf 1.5 2.8 3.2 }}`,
			},
		},
	}
}