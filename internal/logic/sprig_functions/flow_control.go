package sprig_functions

import "github.com/ciclebyte/template_starter/internal/model"

// buildFlowControlFunctions 构建流程控制函数分类
func (s *sSprigFunctions) buildFlowControlFunctions() model.SprigFunctionCategory {
	return model.SprigFunctionCategory{
		Name:        "流程控制函数",
		Description: "模板流程控制和错误处理函数",
		Functions: []model.SprigFunction{
			{
				Name:        "fail",
				DisplayName: "强制失败",
				Description: "无条件返回空字符串和指定错误信息",
				Params: []model.SprigFunctionParam{
					{Name: "message", Type: "string", Required: true, Description: "错误信息"},
				},
				ReturnType: "string",
				Category:   "流程控制函数",
				Example:    `{{ fail "Please accept the end user license agreement" }}`,
				Usage:      "fail 函数无条件返回空字符串和带有指定文本的错误。这在其他条件判断确定模板渲染应该失败的场景中很有用。",
				InsertText: `{{ fail "error message" }}`,
				Note:       "会导致模板渲染停止并返回错误",
			},
		},
	}
}