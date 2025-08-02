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
		// 日期函数分类
		s.buildDateFunctions(),
	}
}

// buildStringFunctions 构建字符串函数分类
func (s *sSprigFunctions) buildStringFunctions() model.SprigFunctionCategory {
	return model.SprigFunctionCategory{
		Name:         "字符串函数",
		Description:  "字符串操作和处理函数",
		DocumentFile: "string.md",
		Functions: []model.SprigFunction{
			{
				Name:        "upper",
				DisplayName: "转大写",
				Description: "将字符串转换为大写字母",
				Params: []model.SprigFunctionParam{
					{Name: "string", Type: "string", Required: true, Description: "要转换的字符串"},
				},
				ReturnType: "string",
				Category:   "字符串函数",
				Example:    `{{ "hello" | upper }}`,
				Usage:      "upper 函数将整个字符串转换为大写字母。这对于格式化输出或标准化字符串非常有用。",
				InsertText: `{{ . | upper }}`,
				Note:       "支持Unicode字符",
			},
			{
				Name:        "lower",
				DisplayName: "转小写",
				Description: "将字符串转换为小写字母",
				Params: []model.SprigFunctionParam{
					{Name: "string", Type: "string", Required: true, Description: "要转换的字符串"},
				},
				ReturnType: "string",
				Category:   "字符串函数",
				Example:    `{{ "HELLO" | lower }}`,
				Usage:      "lower 函数将整个字符串转换为小写字母。常用于标准化用户输入或生成统一格式的标识符。",
				InsertText: `{{ . | lower }}`,
				Note:       "支持Unicode字符",
			},
			{
				Name:        "title",
				DisplayName: "首字母大写",
				Description: "将字符串转换为标题格式，每个单词的首字母大写",
				Params: []model.SprigFunctionParam{
					{Name: "string", Type: "string", Required: true, Description: "要转换的字符串"},
				},
				ReturnType: "string",
				Category:   "字符串函数",
				Example:    `{{ "hello world" | title }}`,
				Usage:      "title 函数将字符串转换为标题格式，即每个单词的首字母大写，其余字母小写。适用于生成标题、姓名等格式化文本。",
				InsertText: `{{ . | title }}`,
			},
			{
				Name:        "trim",
				DisplayName: "去除空格",
				Description: "去除字符串两端的空白字符",
				Params: []model.SprigFunctionParam{
					{Name: "string", Type: "string", Required: true, Description: "要处理的字符串"},
				},
				ReturnType: "string",
				Category:   "字符串函数",
				Example:    `{{ "  hello world  " | trim }}`,
				Usage:      "trim 函数移除字符串两端的空白字符（空格、制表符、换行符等）。这在处理用户输入或清理数据时非常有用。",
				InsertText: `{{ . | trim }}`,
				Note:       "不会移除字符串中间的空白字符",
			},
			{
				Name:        "trimPrefix",
				DisplayName: "去除前缀",
				Description: "从字符串开头移除指定的前缀",
				Params: []model.SprigFunctionParam{
					{Name: "prefix", Type: "string", Required: true, Description: "要移除的前缀"},
					{Name: "string", Type: "string", Required: true, Description: "原字符串"},
				},
				ReturnType: "string",
				Category:   "字符串函数",
				Example:    `{{ trimPrefix "prefix_" "prefix_filename.txt" }}`,
				Usage:      "trimPrefix 函数仅从字符串开头移除指定的前缀。如果字符串不以该前缀开头，则原样返回。常用于去除文件名前缀或URL协议等。",
				InsertText: `{{ trimPrefix "prefix" . }}`,
			},
			{
				Name:        "trimSuffix",
				DisplayName: "去除后缀",
				Description: "从字符串末尾移除指定的后缀",
				Params: []model.SprigFunctionParam{
					{Name: "suffix", Type: "string", Required: true, Description: "要移除的后缀"},
					{Name: "string", Type: "string", Required: true, Description: "原字符串"},
				},
				ReturnType: "string",
				Category:   "字符串函数",
				Example:    `{{ trimSuffix ".txt" "document.txt" }}`,
				Usage:      "trimSuffix 函数仅从字符串末尾移除指定的后缀。如果字符串不以该后缀结尾，则原样返回。常用于去除文件扩展名或URL参数等。",
				InsertText: `{{ trimSuffix "suffix" . }}`,
			},
			{
				Name:        "replace",
				DisplayName: "替换字符串",
				Description: "将字符串中的指定内容替换为新内容",
				Params: []model.SprigFunctionParam{
					{Name: "old", Type: "string", Required: true, Description: "要替换的内容"},
					{Name: "new", Type: "string", Required: true, Description: "替换为的内容"},
					{Name: "string", Type: "string", Required: true, Description: "原字符串"},
				},
				ReturnType: "string",
				Category:   "字符串函数",
				Example:    `{{ replace "world" "Go" "hello world" }}`,
				Usage:      "replace 函数将字符串中所有出现的指定内容替换为新内容。这是一个非常常用的字符串处理函数。",
				InsertText: `{{ replace "old" "new" . }}`,
				Note:       "会替换所有匹配的内容",
			},
			{
				Name:        "contains",
				DisplayName: "包含检查",
				Description: "检查字符串是否包含指定的子字符串",
				Params: []model.SprigFunctionParam{
					{Name: "substr", Type: "string", Required: true, Description: "要查找的子字符串"},
					{Name: "string", Type: "string", Required: true, Description: "原字符串"},
				},
				ReturnType: "bool",
				Category:   "字符串函数",
				Example:    `{{ if contains "go" "golang" }}包含{{ end }}`,
				Usage:      "contains 函数检查字符串是否包含指定的子字符串，返回布尔值。常用于条件判断。",
				InsertText: `{{ contains "substr" . }}`,
			},
			{
				Name:        "hasPrefix",
				DisplayName: "前缀检查",
				Description: "检查字符串是否以指定前缀开头",
				Params: []model.SprigFunctionParam{
					{Name: "prefix", Type: "string", Required: true, Description: "要检查的前缀"},
					{Name: "string", Type: "string", Required: true, Description: "原字符串"},
				},
				ReturnType: "bool",
				Category:   "字符串函数",
				Example:    `{{ if hasPrefix "http" .URL }}是HTTP链接{{ end }}`,
				Usage:      "hasPrefix 函数检查字符串是否以指定前缀开头，返回布尔值。常用于URL协议检查、文件类型判断等。",
				InsertText: `{{ hasPrefix "prefix" . }}`,
			},
			{
				Name:        "hasSuffix",
				DisplayName: "后缀检查",
				Description: "检查字符串是否以指定后缀结尾",
				Params: []model.SprigFunctionParam{
					{Name: "suffix", Type: "string", Required: true, Description: "要检查的后缀"},
					{Name: "string", Type: "string", Required: true, Description: "原字符串"},
				},
				ReturnType: "bool",
				Category:   "字符串函数",
				Example:    `{{ if hasSuffix ".go" .Filename }}是Go文件{{ end }}`,
				Usage:      "hasSuffix 函数检查字符串是否以指定后缀结尾，返回布尔值。常用于文件扩展名检查、URL路径判断等。",
				InsertText: `{{ hasSuffix "suffix" . }}`,
			},
		},
	}
}

// buildDateFunctions 构建日期函数分类
func (s *sSprigFunctions) buildDateFunctions() model.SprigFunctionCategory {
	return model.SprigFunctionCategory{
		Name:         "日期函数",
		Description:  "日期时间处理和格式化函数",
		DocumentFile: "date.md",
		Functions: []model.SprigFunction{
			{
				Name:        "now",
				DisplayName: "当前时间",
				Description: "获取当前时间并格式化输出",
				Params:      []model.SprigFunctionParam{},
				ReturnType:  "time.Time",
				Category:    "日期函数",
				Example:     `{{ now }}`,
				Usage:       "now 函数返回当前的时间对象。通常需要结合 date 函数进行格式化输出。",
				InsertText:  `{{ now }}`,
				Note:        "返回的是 time.Time 对象，需要格式化才能显示",
			},
			{
				Name:        "date",
				DisplayName: "格式化日期",
				Description: "将时间对象按指定格式进行格式化",
				Params: []model.SprigFunctionParam{
					{Name: "format", Type: "string", Required: true, Description: "时间格式字符串"},
					{Name: "time", Type: "time.Time", Required: true, Description: "时间对象"},
				},
				ReturnType: "string",
				Category:   "日期函数",
				Example:    `{{ date "2006-01-02 15:04:05" now }}`,
				Usage:      "date 函数将时间对象按照指定的格式进行格式化。Go语言使用特定的参考时间 '2006-01-02 15:04:05' 作为格式模板。",
				InsertText: `{{ date "2006-01-02" now }}`,
				Note:       "Go时间格式使用参考时间：2006-01-02 15:04:05 -0700 MST",
				Aliases:    []string{"dateFormat"},
			},
			{
				Name:        "dateInZone",
				DisplayName: "指定时区格式化",
				Description: "在指定时区格式化时间",
				Params: []model.SprigFunctionParam{
					{Name: "format", Type: "string", Required: true, Description: "时间格式字符串"},
					{Name: "time", Type: "time.Time", Required: true, Description: "时间对象"},
					{Name: "timezone", Type: "string", Required: true, Description: "时区名称，如 'Asia/Shanghai'"},
				},
				ReturnType: "string",
				Category:   "日期函数",
				Example:    `{{ dateInZone "2006-01-02 15:04:05" now "Asia/Shanghai" }}`,
				Usage:      "dateInZone 函数在指定时区格式化时间。时区名称需要使用标准的IANA时区数据库格式，如 'Asia/Shanghai'、'UTC'、'America/New_York' 等。",
				InsertText: `{{ dateInZone "2006-01-02 15:04:05" now "Asia/Shanghai" }}`,
				Note:       "时区名称必须是有效的IANA时区标识符",
			},
			{
				Name:        "duration",
				DisplayName: "时间间隔",
				Description: "创建时间间隔对象",
				Params: []model.SprigFunctionParam{
					{Name: "duration", Type: "string", Required: true, Description: "时间间隔字符串，如 '1h30m'"},
				},
				ReturnType: "time.Duration",
				Category:   "日期函数",
				Example:    `{{ duration "1h30m" }}`,
				Usage:      "duration 函数解析时间间隔字符串并返回 Duration 对象。支持的单位：ns(纳秒)、us(微秒)、ms(毫秒)、s(秒)、m(分钟)、h(小时)。",
				InsertText: `{{ duration "1h" }}`,
				Note:       "支持组合使用，如 '1h30m45s'",
			},
			{
				Name:        "durationRound",
				DisplayName: "时间间隔舍入",
				Description: "将时间间隔舍入到指定精度",
				Params: []model.SprigFunctionParam{
					{Name: "duration", Type: "time.Duration", Required: true, Description: "要舍入的时间间隔"},
					{Name: "precision", Type: "time.Duration", Required: true, Description: "舍入精度"},
				},
				ReturnType: "time.Duration",
				Category:   "日期函数",
				Example:    `{{ durationRound .Duration (duration "1m") }}`,
				Usage:      "durationRound 函数将时间间隔舍入到指定的精度。例如将秒级精度舍入到分钟级。",
				InsertText: `{{ durationRound . (duration "1m") }}`,
			},
			{
				Name:        "unixEpoch",
				DisplayName: "Unix时间戳",
				Description: "将时间转换为Unix时间戳",
				Params: []model.SprigFunctionParam{
					{Name: "time", Type: "time.Time", Required: true, Description: "时间对象"},
				},
				ReturnType: "int64",
				Category:   "日期函数",
				Example:    `{{ unixEpoch now }}`,
				Usage:      "unixEpoch 函数将时间对象转换为Unix时间戳（自1970年1月1日以来的秒数）。",
				InsertText: `{{ unixEpoch now }}`,
				Note:       "返回秒级时间戳",
			},
			{
				Name:        "dateModify",
				DisplayName: "修改日期",
				Description: "在指定时间基础上增加或减少时间",
				Params: []model.SprigFunctionParam{
					{Name: "modification", Type: "string", Required: true, Description: "时间修改字符串，如 '+1h' 或 '-30m'"},
					{Name: "time", Type: "time.Time", Required: true, Description: "基准时间对象"},
				},
				ReturnType: "time.Time",
				Category:   "日期函数",
				Example:    `{{ dateModify "+1h" now }}`,
				Usage:      "dateModify 函数在指定时间基础上增加或减少时间。修改字符串格式为 '+/-' + 数字 + 单位，如 '+1h'、'-30m'、'+7d'。",
				InsertText: `{{ dateModify "+1h" now }}`,
				Note:       "支持的单位：ns、us、ms、s、m、h",
			},
		},
	}
}