package controller

import (
	"context"

	v1 "github.com/ciclebyte/template_starter/api/v1/builtin_functions"
	"github.com/ciclebyte/template_starter/internal/service"
)

var BuiltinFunctions = cBuiltinFunctions{}

type cBuiltinFunctions struct{}

// GetBuiltinFunctions 获取内置函数列表
func (c *cBuiltinFunctions) GetBuiltinFunctions(ctx context.Context, req *v1.GetBuiltinFunctionsReq) (res *v1.GetBuiltinFunctionsRes, err error) {
	res = new(v1.GetBuiltinFunctionsRes)
	
	data, err := service.BuiltinFunctions().GetBuiltinFunctions(ctx)
	if err != nil {
		return
	}

	// 转换为API响应格式
	categories := make([]v1.BuiltinFunctionCategory, len(data.Categories))
	for i, cat := range data.Categories {
		functions := make([]v1.BuiltinFunction, len(cat.Functions))
		for j, fn := range cat.Functions {
			params := make([]v1.FunctionParam, len(fn.Params))
			for k, param := range fn.Params {
				params[k] = v1.FunctionParam{
					Name:        param.Name,
					Type:        param.Type,
					Required:    param.Required,
					Description: param.Description,
					Default:     param.Default,
				}
			}
			functions[j] = v1.BuiltinFunction{
				Name:        fn.Name,
				DisplayName: fn.DisplayName,
				Description: fn.Description,
				Params:      params,
				ReturnType:  fn.ReturnType,
				Category:    fn.Category,
				Example:     fn.Example,
				InsertText:  fn.InsertText,
			}
		}
		categories[i] = v1.BuiltinFunctionCategory{
			Name:        cat.Name,
			Description: cat.Description,
			Functions:   functions,
		}
	}

	res.Categories = categories
	res.Total = data.Total
	return
}