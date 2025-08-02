package controller

import (
	"context"

	v1 "github.com/ciclebyte/template_starter/api/v1/sprig_functions"
	"github.com/ciclebyte/template_starter/internal/service"
)

var SprigFunctions = cSprigFunctions{}

type cSprigFunctions struct{}

// GetSprigFunctions 获取Sprig函数列表
func (c *cSprigFunctions) GetSprigFunctions(ctx context.Context, req *v1.GetSprigFunctionsReq) (res *v1.GetSprigFunctionsRes, err error) {
	res = new(v1.GetSprigFunctionsRes)
	
	data, err := service.SprigFunctions().GetSprigFunctions(ctx)
	if err != nil {
		return
	}

	// 转换为API响应格式
	categories := make([]v1.SprigFunctionCategory, len(data.Categories))
	for i, cat := range data.Categories {
		functions := make([]v1.SprigFunction, len(cat.Functions))
		for j, fn := range cat.Functions {
			params := make([]v1.SprigFunctionParam, len(fn.Params))
			for k, param := range fn.Params {
				params[k] = v1.SprigFunctionParam{
					Name:        param.Name,
					Type:        param.Type,
					Required:    param.Required,
					Description: param.Description,
					Default:     param.Default,
				}
			}
			functions[j] = v1.SprigFunction{
				Name:        fn.Name,
				DisplayName: fn.DisplayName,
				Description: fn.Description,
				Params:      params,
				ReturnType:  fn.ReturnType,
				Category:    fn.Category,
				Example:     fn.Example,
				Usage:       fn.Usage,
				InsertText:  fn.InsertText,
				Note:        fn.Note,
				Aliases:     fn.Aliases,
			}
		}
		categories[i] = v1.SprigFunctionCategory{
			Name:         cat.Name,
			Description:  cat.Description,
			DocumentFile: cat.DocumentFile,
			Functions:    functions,
		}
	}

	res.Categories = categories
	res.Total = data.Total
	return
}