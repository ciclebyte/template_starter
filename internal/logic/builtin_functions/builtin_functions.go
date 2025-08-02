package builtin_functions

import (
	"context"
	"text/template"
	"time"

	"github.com/ciclebyte/template_starter/internal/model"
	"github.com/ciclebyte/template_starter/internal/service"
	"github.com/google/uuid"
)

type sBuiltinFunctions struct{}

func init() {
	service.RegisterBuiltinFunctions(New())
}

func New() *sBuiltinFunctions {
	return &sBuiltinFunctions{}
}

// GetBuiltinFunctions 获取所有内置函数（只返回自定义函数）
func (s *sBuiltinFunctions) GetBuiltinFunctions(ctx context.Context) (*model.BuiltinFunctionsResponse, error) {
	categories := []model.BuiltinFunctionCategory{}
	
	// 添加自定义函数
	customCategory := model.BuiltinFunctionCategory{
		Name:        "自定义函数",
		Description: "项目专用的内置函数",
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
				Category:   "自定义函数",
				Example:    `{{ now }} 或 {{ now "2006-01-02" }}`,
				InsertText: `{{ now }}`,
			},
			{
				Name:        "uuid",
				DisplayName: "UUID生成",
				Description: "生成一个随机的UUID字符串，格式：xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
				Params:      []model.FunctionParam{},
				ReturnType:  "string",
				Category:    "自定义函数",
				Example:     `{{ uuid }}`,
				InsertText:  `{{ uuid }}`,
			},
		},
	}
	categories = append(categories, customCategory)
	
	total := 0
	for _, category := range categories {
		total += len(category.Functions)
	}

	return &model.BuiltinFunctionsResponse{
		Categories: categories,
		Total:      total,
	}, nil
}

// GetTemplateFuncMap 获取模板函数映射（供模板渲染使用）
func (s *sBuiltinFunctions) GetTemplateFuncMap() template.FuncMap {
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