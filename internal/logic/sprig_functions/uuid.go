package sprig_functions

import "github.com/ciclebyte/template_starter/internal/model"

// buildUUIDFunctions 构建UUID函数分类
func (s *sSprigFunctions) buildUUIDFunctions() model.SprigFunctionCategory {
	return model.SprigFunctionCategory{
		Name:        "UUID函数",
		Description: "通用唯一标识符生成函数",
		Functions: []model.SprigFunction{
			{
				Name:        "uuidv4",
				DisplayName: "生成UUIDv4",
				Description: "生成UUID v4通用唯一标识符",
				Params:      []model.SprigFunctionParam{},
				ReturnType:  "string",
				Category:    "UUID函数",
				Example:     `{{ uuidv4 }}`,
				Usage:       "uuidv4 函数生成UUID v4（随机生成）类型的新通用唯一标识符。每次调用都会返回一个全新的随机UUID。",
				InsertText:  `{{ uuidv4 }}`,
				Note:        "基于随机数生成，每次调用结果都不同",
			},
		},
	}
}