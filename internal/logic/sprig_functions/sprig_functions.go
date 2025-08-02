package sprig_functions

import (
	"context"

	"github.com/ciclebyte/template_starter/internal/model"
	"github.com/ciclebyte/template_starter/internal/service"
)

type sSprigFunctions struct{}

func init() {
	service.RegisterSprigFunctions(New())
}

func New() *sSprigFunctions {
	return &sSprigFunctions{}
}

// GetSprigFunctions 获取所有Sprig函数
func (s *sSprigFunctions) GetSprigFunctions(ctx context.Context) (*model.SprigFunctionsResponse, error) {
	categories := s.buildSprigFunctionCategories()

	total := 0
	for _, category := range categories {
		total += len(category.Functions)
	}

	return &model.SprigFunctionsResponse{
		Categories: categories,
		Total:      total,
	}, nil
}

// buildSprigFunctionCategories 构建Sprig函数分类（硬编码，易于扩展）
func (s *sSprigFunctions) buildSprigFunctionCategories() []model.SprigFunctionCategory {
	return []model.SprigFunctionCategory{
		// 字符串函数分类
		s.buildStringFunctions(),
		// 字符串切片函数分类
		s.buildStringSliceFunctions(),
		// 日期函数分类
		s.buildDateFunctions(),
		// 整数数学函数分类
		s.buildIntegerMathFunctions(),
		// 整数切片函数分类
		s.buildIntegerSliceFunctions(),
		// 浮点数数学函数分类
		s.buildFloatMathFunctions(),
		// 默认值函数分类
		s.buildDefaultsFunctions(),
		// 编码函数分类
		s.buildEncodingFunctions(),
		// 列表函数分类
		s.buildListsFunctions(),
		// 字典函数分类
		s.buildDictionariesFunctions(),
		// 类型转换函数分类
		s.buildTypeConversionFunctions(),
		// 路径和文件路径函数分类
		s.buildPathFilepathFunctions(),
		// 流程控制函数分类
		s.buildFlowControlFunctions(),
		// UUID函数分类
		s.buildUUIDFunctions(),
		// 操作系统函数分类
		s.buildOSFunctions(),
		// 语义版本函数分类
		s.buildSemverFunctions(),
		// 反射函数分类
		s.buildReflectionFunctions(),
		// 加密和安全函数分类
		s.buildCryptoFunctions(),
		// 网络函数分类
		s.buildNetworkFunctions(),
		// URL函数分类
		s.buildURLFunctions(),
	}
}
