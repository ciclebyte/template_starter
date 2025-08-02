package service

import (
	"context"
	"text/template"
	"time"

	"github.com/ciclebyte/template_starter/internal/model"
	"github.com/google/uuid"
)

type sBuiltinFunctions struct{}

func init() {
	RegisterBuiltinFunctions(BuiltinFunctions())
}

func BuiltinFunctions() *sBuiltinFunctions {
	return &sBuiltinFunctions{}
}

// GetBuiltinFunctions 获取所有内置函数
func (s *sBuiltinFunctions) GetBuiltinFunctions(ctx context.Context) (*model.BuiltinFunctionsResponse, error) {
	categories := []model.BuiltinFunctionCategory{
		{
			Name:        "时间函数",
			Description: "时间相关的处理函数",
			Functions: []model.BuiltinFunction{
				{
					Name:        "now",
					DisplayName: "当前时间",
					Description: "获取当前时间，可指定格式。默认格式：2006-01-02 15:04:05",
					Params: []model.FunctionParam{
						{
							Name:        "format",
							Type:        "string",
							Required:    false,
							Description: "时间格式，如 2006-01-02、2006-01-02 15:04:05",
							Default:     "2006-01-02 15:04:05",
						},
					},
					ReturnType: "string",
					Category:   "时间函数",
					Example:    `{{ now }} 或 {{ now "2006-01-02" }}`,
					InsertText: `{{ now }}`,
				},
			},
		},
		{
			Name:        "工具函数",
			Description: "通用工具函数",
			Functions: []model.BuiltinFunction{
				{
					Name:        "uuid",
					DisplayName: "UUID生成",
					Description: "生成一个随机的UUID字符串，格式：xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
					Params:      []model.FunctionParam{},
					ReturnType:  "string",
					Category:    "工具函数",
					Example:     `{{ uuid }}`,
					InsertText:  `{{ uuid }}`,
				},
			},
		},
	}

	total := 0
	for _, category := range categories {
		total += len(category.Functions)
	}

	return &model.BuiltinFunctionsResponse{
		Categories: categories,
		Total:      total,
	}, nil
}

// RegisterBuiltinFunctions 注册内置函数到模板引擎
func RegisterBuiltinFunctions(s *sBuiltinFunctions) template.FuncMap {
	return template.FuncMap{
		"now": func(format ...string) string {
			if len(format) > 0 && format[0] != "" {
				return time.Now().Format(format[0])
			}
			return time.Now().Format("2006-01-02 15:04:05")
		},
		"uuid": func() string {
			return uuid.New().String()
		},
	}
}

// GetTemplateFuncMap 获取模板函数映射（供模板渲染使用）
func (s *sBuiltinFunctions) GetTemplateFuncMap() template.FuncMap {
	return RegisterBuiltinFunctions(s)
}