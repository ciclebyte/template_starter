package service

import (
	"context"
	"text/template"

	"github.com/ciclebyte/template_starter/internal/model"
)

type (
	IBuiltinFunctions interface {
		GetBuiltinFunctions(ctx context.Context) (*model.BuiltinFunctionsResponse, error)
		GetTemplateFuncMap() template.FuncMap
	}
)

var (
	localBuiltinFunctions IBuiltinFunctions
)

func BuiltinFunctions() IBuiltinFunctions {
	if localBuiltinFunctions == nil {
		panic("implement not found for interface IBuiltinFunctions, forgot register?")
	}
	return localBuiltinFunctions
}

func RegisterBuiltinFunctions(i IBuiltinFunctions) {
	localBuiltinFunctions = i
}