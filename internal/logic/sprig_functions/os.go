package sprig_functions

import "github.com/ciclebyte/template_starter/internal/model"

// buildOSFunctions 构建操作系统函数分类
func (s *sSprigFunctions) buildOSFunctions() model.SprigFunctionCategory {
	return model.SprigFunctionCategory{
		Name:        "操作系统函数",
		Description: "操作系统环境变量和系统信息函数",
		Functions: []model.SprigFunction{
			{
				Name:        "env",
				DisplayName: "读取环境变量",
				Description: "读取操作系统环境变量",
				Params: []model.SprigFunctionParam{
					{Name: "variable", Type: "string", Required: true, Description: "环境变量名"},
				},
				ReturnType: "string",
				Category:   "操作系统函数",
				Example:    `{{ env "HOME" }}`,
				Usage:      "env 函数读取操作系统环境变量。例如 env \"HOME\" 返回HOME环境变量的值。",
				InsertText: `{{ env "VARIABLE_NAME" }}`,
				Note:       "⚠️ 警告：如果使用不当，这些函数可能导致信息泄露。某些Sprig实现（如Kubernetes Helm）出于安全原因不提供这些函数",
			},
			{
				Name:        "expandenv",
				DisplayName: "展开环境变量",
				Description: "在字符串中替换环境变量",
				Params: []model.SprigFunctionParam{
					{Name: "string", Type: "string", Required: true, Description: "包含环境变量的字符串"},
				},
				ReturnType: "string",
				Category:   "操作系统函数",
				Example:    `{{ expandenv "Your path is set to $PATH" }}`,
				Usage:      "expandenv 函数在字符串中替换环境变量。例如 expandenv \"Your path is set to $PATH\" 会将$PATH替换为实际的PATH环境变量值。",
				InsertText: `{{ expandenv "String with $VARIABLE" }}`,
				Note:       "⚠️ 安全警告：可能导致信息泄露，某些环境中不可用",
			},
		},
	}
}