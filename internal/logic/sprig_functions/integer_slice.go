package sprig_functions

import "github.com/ciclebyte/template_starter/internal/model"

// buildIntegerSliceFunctions 构建整数切片函数分类
func (s *sSprigFunctions) buildIntegerSliceFunctions() model.SprigFunctionCategory {
	return model.SprigFunctionCategory{
		Name:        "整数切片函数",
		Description: "整数切片生成和操作函数",
		Functions: []model.SprigFunction{
			{
				Name:        "until",
				DisplayName: "生成整数范围",
				Description: "构建从0开始到指定数字的整数范围",
				Params: []model.SprigFunctionParam{
					{Name: "count", Type: "int", Required: true, Description: "结束数字（不包含）"},
				},
				ReturnType: "[]int",
				Category:   "整数切片函数",
				Example:    `{{ until 5 }}`,
				Usage:      "until 函数构建一个整数范围。until 5 生成列表 [0, 1, 2, 3, 4]。这对于与 range 循环结合使用很有用，如 range $i, $e := until 5。",
				InsertText: `{{ until 5 }}`,
				Note:       "生成从0开始，不包含结束数字的整数列表",
			},
			{
				Name:        "untilStep",
				DisplayName: "步长整数范围",
				Description: "生成带有自定义起始值、步长和结束值的整数列表",
				Params: []model.SprigFunctionParam{
					{Name: "start", Type: "int", Required: true, Description: "起始数字"},
					{Name: "step", Type: "int", Required: true, Description: "步长"},
					{Name: "stop", Type: "int", Required: true, Description: "结束数字（不包含）"},
				},
				ReturnType: "[]int",
				Category:   "整数切片函数",
				Example:    `{{ untilStep 3 2 6 }}`,
				Usage:      "untilStep 函数类似于 until，但允许定义起始值、停止值和步长。untilStep 3 2 6 将产生 [3 5]，从3开始，每次增加2，直到等于或大于6。这类似于Python的range函数。",
				InsertText: `{{ untilStep 0 1 10 }}`,
				Note:       "类似Python的range函数，支持自定义起始值和步长",
			},
			{
				Name:        "seq",
				DisplayName: "序列生成",
				Description: "类似bash seq命令的序列生成函数",
				Params: []model.SprigFunctionParam{
					{Name: "args", Type: "...int", Required: true, Description: "参数：1个参数(end)、2个参数(start,end)或3个参数(start,step,end)"},
				},
				ReturnType: "[]int",
				Category:   "整数切片函数",
				Example:    `{{ seq 1 5 }}`,
				Usage:      "seq 函数类似bash的seq命令。1个参数(end)：生成1到end的所有整数（包含）。2个参数(start,end)：生成start到end的所有整数（包含），步长为1。3个参数(start,step,end)：生成start到end的所有整数（包含），步长为step。",
				InsertText: `{{ seq 1 10 }}`,
				Note:       "支持1-3个参数，与bash seq命令行为一致，结果包含结束值",
			},
		},
	}
}