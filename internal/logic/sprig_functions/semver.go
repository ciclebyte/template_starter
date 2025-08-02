package sprig_functions

import "github.com/ciclebyte/template_starter/internal/model"

// buildSemverFunctions 构建语义版本函数分类
func (s *sSprigFunctions) buildSemverFunctions() model.SprigFunctionCategory {
	return model.SprigFunctionCategory{
		Name:        "语义版本函数",
		Description: "语义版本解析和比较函数",
		Functions: []model.SprigFunction{
			{
				Name:        "semver",
				DisplayName: "解析语义版本",
				Description: "将字符串解析为语义版本对象",
				Params: []model.SprigFunctionParam{
					{Name: "version", Type: "string", Required: true, Description: "版本字符串"},
				},
				ReturnType: "Version",
				Category:   "语义版本函数",
				Example:    `{{ $version := semver "1.2.3-alpha.1+123" }}`,
				Usage:      "semver 函数将字符串解析为语义版本。返回Version对象包含Major、Minor、Patch、Prerelease、Metadata、Original属性。如果解析失败，会导致模板执行中止并返回错误。",
				InsertText: `{{ $v := semver "1.2.3" }}{{ $v.Major }}`,
				Note:       "解析失败会中止模板执行。返回对象支持Compare方法进行版本比较",
			},
			{
				Name:        "semverCompare",
				DisplayName: "比较语义版本",
				Description: "使用约束条件比较语义版本",
				Params: []model.SprigFunctionParam{
					{Name: "constraint", Type: "string", Required: true, Description: "版本约束条件"},
					{Name: "version", Type: "string", Required: true, Description: "要比较的版本"},
				},
				ReturnType: "bool",
				Category:   "语义版本函数",
				Example:    `{{ semverCompare "^1.2.0" "1.2.3" }}`,
				Usage:      "semverCompare 函数提供更强大的版本比较功能。支持版本范围比较，如果约束匹配返回true，否则返回false。支持精确匹配、主次版本匹配等多种比较模式。",
				InsertText: `{{ semverCompare ">=1.0.0" . }}`,
				Note:       "支持复杂的版本约束语法：=、!=、>、<、>=、<=、~、^、通配符等",
			},
		},
	}
}