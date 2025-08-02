package service

import (
	"context"

	"github.com/ciclebyte/template_starter/internal/model"
)

type (
	ISprigFunctions interface {
		GetSprigFunctions(ctx context.Context) (*model.SprigFunctionsResponse, error)
	}
)

var (
	localSprigFunctions ISprigFunctions
)

func SprigFunctions() ISprigFunctions {
	if localSprigFunctions == nil {
		panic("implement not found for interface ISprigFunctions, forgot register?")
	}
	return localSprigFunctions
}

func RegisterSprigFunctions(i ISprigFunctions) {
	localSprigFunctions = i
}