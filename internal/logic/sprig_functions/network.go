package sprig_functions

import "github.com/ciclebyte/template_starter/internal/model"

// buildNetworkFunctions 构建网络函数分类
func (s *sSprigFunctions) buildNetworkFunctions() model.SprigFunctionCategory {
	return model.SprigFunctionCategory{
		Name:        "网络函数",
		Description: "网络操作和处理函数",
		Functions: []model.SprigFunction{
			{
				Name:        "getHostByName",
				DisplayName: "域名解析",
				Description: "根据域名获取IP地址",
				Params: []model.SprigFunctionParam{
					{Name: "hostname", Type: "string", Required: true, Description: "域名"},
				},
				ReturnType: "string",
				Category:   "网络函数",
				Example:    `{{ getHostByName "www.google.com" }}`,
				Usage:      "getHostByName 函数接收一个域名，返回对应的IP地址。例如 getHostByName \"www.google.com\" 会返回 www.google.com 对应的IP地址。",
				InsertText: `{{ getHostByName "example.com" }}`,
				Note:       "⚠️ 网络操作可能失败，某些环境中可能不可用",
			},
		},
	}
}