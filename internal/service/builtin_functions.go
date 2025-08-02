package service

import (
	"context"
	"reflect"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/Masterminds/sprig/v3"
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
	categories := []model.BuiltinFunctionCategory{}
	
	// 1. 添加自定义函数
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
	
	// 2. 通过反射获取Sprig函数
	sprigCategories := s.getSprigFunctions()
	categories = append(categories, sprigCategories...)
	
	total := 0
	for _, category := range categories {
		total += len(category.Functions)
	}

	return &model.BuiltinFunctionsResponse{
		Categories: categories,
		Total:      total,
	}, nil
}

// getSprigFunctions 通过反射获取Sprig函数
func (s *sBuiltinFunctions) getSprigFunctions() []model.BuiltinFunctionCategory {
	sprigFuncs := sprig.FuncMap()
	
	// 函数分类映射
	categoryMap := map[string][]string{
		"字符串函数": {
			"upper", "lower", "title", "trim", "trimSuffix", "trimPrefix", 
			"replace", "split", "join", "contains", "hasPrefix", "hasSuffix",
			"repeat", "substr", "indent", "nindent", "wrap", "trunc",
		},
		"数学函数": {
			"add", "sub", "mul", "div", "mod", "max", "min", "floor", "ceil", "round",
			"addf", "subf", "mulf", "divf", "maxf", "minf",
		},
		"条件函数": {
			"default", "empty", "coalesce", "ternary", "not", "and", "or",
		},
		"随机函数": {
			"randAlphaNum", "randAlpha", "randNumeric", "randInt", "uuidv4",
		},
		"日期函数": {
			"date", "now", "ago", "toDate", "dateInZone", "duration",
		},
		"类型转换": {
			"toString", "toInt", "toFloat64", "toBool", "toJson", "fromJson", "toPrettyJson",
		},
		"列表函数": {
			"list", "append", "prepend", "first", "last", "rest", "initial", "reverse",
			"uniq", "compact", "slice", "has", "without", "sortAlpha",
		},
		"字典函数": {
			"dict", "get", "set", "unset", "hasKey", "pluck", "keys", "values", "pick", "omit",
		},
		"编码函数": {
			"b64enc", "b64dec", "urlquery", "sha1sum", "sha256sum", "md5sum",
		},
		"正则函数": {
			"regexMatch", "regexFindAll", "regexFind", "regexReplaceAll", "regexSplit",
		},
	}
	
	categories := []model.BuiltinFunctionCategory{}
	
	for categoryName, funcNames := range categoryMap {
		functions := []model.BuiltinFunction{}
		
		for _, funcName := range funcNames {
			if fn, exists := sprigFuncs[funcName]; exists {
				function := s.buildFunctionInfo(funcName, fn, categoryName)
				functions = append(functions, function)
			}
		}
		
		if len(functions) > 0 {
			// 按函数名排序
			sort.Slice(functions, func(i, j int) bool {
				return functions[i].Name < functions[j].Name
			})
			
			category := model.BuiltinFunctionCategory{
				Name:        categoryName,
				Description: s.getCategoryDescription(categoryName),
				Functions:   functions,
			}
			categories = append(categories, category)
		}
	}
	
	return categories
}

// buildFunctionInfo 构建函数信息
func (s *sBuiltinFunctions) buildFunctionInfo(name string, fn interface{}, category string) model.BuiltinFunction {
	fnType := reflect.TypeOf(fn)
	
	// 基本信息
	function := model.BuiltinFunction{
		Name:        name,
		DisplayName: s.getFunctionDisplayName(name),
		Description: s.getFunctionDescription(name),
		ReturnType:  s.getReturnType(fnType),
		Category:    category,
		Example:     s.getFunctionExample(name),
		InsertText:  s.getFunctionInsertText(name),
	}
	
	// 参数信息
	if fnType.Kind() == reflect.Func {
		function.Params = s.getFunctionParams(name, fnType)
	}
	
	return function
}

// getFunctionParams 获取函数参数信息
func (s *sBuiltinFunctions) getFunctionParams(name string, fnType reflect.Type) []model.FunctionParam {
	params := []model.FunctionParam{}
	
	// 简化处理，对于常见函数提供参数信息
	paramInfo := s.getKnownParams(name)
	if len(paramInfo) > 0 {
		return paramInfo
	}
	
	// 通过反射获取基本参数信息
	numIn := fnType.NumIn()
	for i := 0; i < numIn; i++ {
		paramType := fnType.In(i)
		param := model.FunctionParam{
			Name:        "arg" + string(rune('0'+i)),
			Type:        s.getTypeString(paramType),
			Required:    true,
			Description: "参数",
		}
		params = append(params, param)
	}
	
	return params
}

// getKnownParams 获取已知函数的参数信息
func (s *sBuiltinFunctions) getKnownParams(name string) []model.FunctionParam {
	knownParams := map[string][]model.FunctionParam{
		// 字符串函数
		"upper":      {{Name: "string", Type: "string", Required: true, Description: "要转换的字符串"}},
		"lower":      {{Name: "string", Type: "string", Required: true, Description: "要转换的字符串"}},
		"title":      {{Name: "string", Type: "string", Required: true, Description: "要转换的字符串"}},
		"trim":       {{Name: "string", Type: "string", Required: true, Description: "要处理的字符串"}},
		"trimPrefix": {
			{Name: "prefix", Type: "string", Required: true, Description: "要移除的前缀"},
			{Name: "string", Type: "string", Required: true, Description: "原字符串"},
		},
		"trimSuffix": {
			{Name: "suffix", Type: "string", Required: true, Description: "要移除的后缀"},
			{Name: "string", Type: "string", Required: true, Description: "原字符串"},
		},
		"replace": {
			{Name: "old", Type: "string", Required: true, Description: "要替换的内容"},
			{Name: "new", Type: "string", Required: true, Description: "替换为的内容"},
			{Name: "string", Type: "string", Required: true, Description: "原字符串"},
		},
		"split": {
			{Name: "separator", Type: "string", Required: true, Description: "分隔符"},
			{Name: "string", Type: "string", Required: true, Description: "要分割的字符串"},
		},
		"join": {
			{Name: "separator", Type: "string", Required: true, Description: "连接符"},
			{Name: "list", Type: "array", Required: true, Description: "字符串数组"},
		},
		"contains": {
			{Name: "substr", Type: "string", Required: true, Description: "子字符串"},
			{Name: "string", Type: "string", Required: true, Description: "原字符串"},
		},
		"hasPrefix": {
			{Name: "prefix", Type: "string", Required: true, Description: "前缀"},
			{Name: "string", Type: "string", Required: true, Description: "原字符串"},
		},
		"hasSuffix": {
			{Name: "suffix", Type: "string", Required: true, Description: "后缀"},
			{Name: "string", Type: "string", Required: true, Description: "原字符串"},
		},
		"repeat": {
			{Name: "count", Type: "int", Required: true, Description: "重复次数"},
			{Name: "string", Type: "string", Required: true, Description: "要重复的字符串"},
		},
		"substr": {
			{Name: "start", Type: "int", Required: true, Description: "起始位置"},
			{Name: "length", Type: "int", Required: true, Description: "长度"},
			{Name: "string", Type: "string", Required: true, Description: "原字符串"},
		},
		"trunc": {
			{Name: "length", Type: "int", Required: true, Description: "截断长度"},
			{Name: "string", Type: "string", Required: true, Description: "要截断的字符串"},
		},
		
		// 数学函数
		"add": {
			{Name: "a", Type: "number", Required: true, Description: "第一个数"},
			{Name: "b", Type: "number", Required: true, Description: "第二个数"},
		},
		"sub": {
			{Name: "a", Type: "number", Required: true, Description: "被减数"},
			{Name: "b", Type: "number", Required: true, Description: "减数"},
		},
		"mul": {
			{Name: "a", Type: "number", Required: true, Description: "第一个数"},
			{Name: "b", Type: "number", Required: true, Description: "第二个数"},
		},
		"div": {
			{Name: "a", Type: "number", Required: true, Description: "被除数"},
			{Name: "b", Type: "number", Required: true, Description: "除数"},
		},
		"mod": {
			{Name: "a", Type: "number", Required: true, Description: "被除数"},
			{Name: "b", Type: "number", Required: true, Description: "除数"},
		},
		"max": {
			{Name: "a", Type: "number", Required: true, Description: "第一个数"},
			{Name: "b", Type: "number", Required: true, Description: "第二个数"},
		},
		"min": {
			{Name: "a", Type: "number", Required: true, Description: "第一个数"},
			{Name: "b", Type: "number", Required: true, Description: "第二个数"},
		},
		"floor": {{Name: "number", Type: "float", Required: true, Description: "要向下取整的数"}},
		"ceil":  {{Name: "number", Type: "float", Required: true, Description: "要向上取整的数"}},
		"round": {{Name: "number", Type: "float", Required: true, Description: "要四舍五入的数"}},
		
		// 条件函数
		"default": {
			{Name: "default", Type: "any", Required: true, Description: "默认值"},
			{Name: "value", Type: "any", Required: true, Description: "要检查的值"},
		},
		"empty": {{Name: "value", Type: "any", Required: true, Description: "要检查的值"}},
		"coalesce": {
			{Name: "values", Type: "any", Required: true, Description: "值列表，返回第一个非空值"},
		},
		"ternary": {
			{Name: "trueVal", Type: "any", Required: true, Description: "条件为真时的值"},
			{Name: "falseVal", Type: "any", Required: true, Description: "条件为假时的值"},
			{Name: "condition", Type: "bool", Required: true, Description: "条件"},
		},
		
		// 随机函数
		"randAlphaNum": {{Name: "length", Type: "int", Required: true, Description: "字符串长度"}},
		"randAlpha":    {{Name: "length", Type: "int", Required: true, Description: "字符串长度"}},
		"randNumeric":  {{Name: "length", Type: "int", Required: true, Description: "字符串长度"}},
		"randInt": {
			{Name: "min", Type: "int", Required: true, Description: "最小值"},
			{Name: "max", Type: "int", Required: true, Description: "最大值"},
		},
		"uuidv4": {},
		
		// 日期函数
		"date": {
			{Name: "format", Type: "string", Required: true, Description: "时间格式"},
			{Name: "time", Type: "time", Required: false, Description: "时间对象，默认为当前时间"},
		},
		"ago": {{Name: "time", Type: "time", Required: true, Description: "时间对象"}},
		"toDate": {
			{Name: "format", Type: "string", Required: true, Description: "时间格式"},
			{Name: "string", Type: "string", Required: true, Description: "时间字符串"},
		},
		
		// 类型转换
		"toString":   {{Name: "value", Type: "any", Required: true, Description: "要转换的值"}},
		"toInt":      {{Name: "value", Type: "any", Required: true, Description: "要转换的值"}},
		"toFloat64":  {{Name: "value", Type: "any", Required: true, Description: "要转换的值"}},
		"toBool":     {{Name: "value", Type: "any", Required: true, Description: "要转换的值"}},
		"toJson":     {{Name: "value", Type: "any", Required: true, Description: "要序列化的值"}},
		"fromJson":   {{Name: "json", Type: "string", Required: true, Description: "JSON字符串"}},
		"toPrettyJson": {{Name: "value", Type: "any", Required: true, Description: "要序列化的值"}},
		
		// 列表函数
		"list":    {}, // 可变参数
		"first":   {{Name: "list", Type: "array", Required: true, Description: "列表"}},
		"last":    {{Name: "list", Type: "array", Required: true, Description: "列表"}},
		"rest":    {{Name: "list", Type: "array", Required: true, Description: "列表"}},
		"initial": {{Name: "list", Type: "array", Required: true, Description: "列表"}},
		"reverse": {{Name: "list", Type: "array", Required: true, Description: "列表"}},
		"uniq":    {{Name: "list", Type: "array", Required: true, Description: "列表"}},
		"compact": {{Name: "list", Type: "array", Required: true, Description: "列表"}},
		"append": {
			{Name: "list", Type: "array", Required: true, Description: "列表"},
			{Name: "item", Type: "any", Required: true, Description: "要添加的元素"},
		},
		"prepend": {
			{Name: "list", Type: "array", Required: true, Description: "列表"},
			{Name: "item", Type: "any", Required: true, Description: "要添加的元素"},
		},
		"has": {
			{Name: "item", Type: "any", Required: true, Description: "要查找的元素"},
			{Name: "list", Type: "array", Required: true, Description: "列表"},
		},
		"without": {
			{Name: "list", Type: "array", Required: true, Description: "列表"},
			{Name: "items", Type: "any", Required: true, Description: "要排除的元素"},
		},
		
		// 字典函数
		"dict": {}, // 可变参数
		"get": {
			{Name: "key", Type: "string", Required: true, Description: "键名"},
			{Name: "dict", Type: "map", Required: true, Description: "字典"},
		},
		"set": {
			{Name: "key", Type: "string", Required: true, Description: "键名"},
			{Name: "value", Type: "any", Required: true, Description: "值"},
			{Name: "dict", Type: "map", Required: true, Description: "字典"},
		},
		"unset": {
			{Name: "key", Type: "string", Required: true, Description: "键名"},
			{Name: "dict", Type: "map", Required: true, Description: "字典"},
		},
		"hasKey": {
			{Name: "key", Type: "string", Required: true, Description: "键名"},
			{Name: "dict", Type: "map", Required: true, Description: "字典"},
		},
		"keys":   {{Name: "dict", Type: "map", Required: true, Description: "字典"}},
		"values": {{Name: "dict", Type: "map", Required: true, Description: "字典"}},
		
		// 编码函数
		"b64enc":    {{Name: "string", Type: "string", Required: true, Description: "要编码的字符串"}},
		"b64dec":    {{Name: "string", Type: "string", Required: true, Description: "要解码的字符串"}},
		"urlquery":  {{Name: "string", Type: "string", Required: true, Description: "要编码的字符串"}},
		"sha1sum":   {{Name: "string", Type: "string", Required: true, Description: "要计算哈希的字符串"}},
		"sha256sum": {{Name: "string", Type: "string", Required: true, Description: "要计算哈希的字符串"}},
		"md5sum":    {{Name: "string", Type: "string", Required: true, Description: "要计算哈希的字符串"}},
	}
	
	if params, exists := knownParams[name]; exists {
		return params
	}
	return []model.FunctionParam{}
}

// 辅助方法
func (s *sBuiltinFunctions) getCategoryDescription(category string) string {
	descriptions := map[string]string{
		"字符串函数": "Sprig字符串处理函数",
		"数学函数":  "Sprig数学计算函数",
		"条件函数":  "Sprig条件判断函数",
		"随机函数":  "Sprig随机值生成函数",
		"日期函数":  "Sprig日期时间函数",
		"类型转换":  "Sprig类型转换函数",
		"列表函数":  "Sprig列表操作函数",
		"字典函数":  "Sprig字典操作函数",
		"编码函数":  "Sprig编码解码函数",
		"正则函数":  "Sprig正则表达式函数",
	}
	if desc, exists := descriptions[category]; exists {
		return desc
	}
	return "Sprig函数"
}

func (s *sBuiltinFunctions) getFunctionDisplayName(name string) string {
	displayNames := map[string]string{
		"upper": "转大写", "lower": "转小写", "title": "首字母大写", "trim": "去除空格",
		"add": "加法", "sub": "减法", "mul": "乘法", "div": "除法",
		"default": "默认值", "empty": "判断为空",
		"randAlphaNum": "随机字母数字", "randAlpha": "随机字母",
		"date": "格式化日期", "now": "当前时间",
	}
	if display, exists := displayNames[name]; exists {
		return display
	}
	return strings.Title(name) // 默认首字母大写
}

func (s *sBuiltinFunctions) getFunctionDescription(name string) string {
	descriptions := map[string]string{
		"upper": "将字符串转换为大写",
		"lower": "将字符串转换为小写", 
		"add": "计算两个数的和",
		"default": "如果值为空则返回默认值",
		"randAlphaNum": "生成指定长度的随机字母数字字符串",
	}
	if desc, exists := descriptions[name]; exists {
		return desc
	}
	return "Sprig " + name + " 函数"
}

func (s *sBuiltinFunctions) getFunctionExample(name string) string {
	examples := map[string]string{
		// 字符串函数
		"upper":      `{{ "hello" | upper }}`,
		"lower":      `{{ "HELLO" | lower }}`,
		"title":      `{{ "hello world" | title }}`,
		"trim":       `{{ "  hello  " | trim }}`,
		"trimPrefix": `{{ "hello world" | trimPrefix "hello " }}`,
		"trimSuffix": `{{ "hello world" | trimSuffix " world" }}`,
		"replace":    `{{ "hello world" | replace "world" "sprig" }}`,
		"split":      `{{ "a,b,c" | split "," }}`,
		"join":       `{{ list "a" "b" "c" | join "," }}`,
		"contains":   `{{ "hello world" | contains "world" }}`,
		"hasPrefix":  `{{ "hello world" | hasPrefix "hello" }}`,
		"hasSuffix":  `{{ "hello world" | hasSuffix "world" }}`,
		"repeat":     `{{ "hello" | repeat 3 }}`,
		"substr":     `{{ "hello world" | substr 0 5 }}`,
		"trunc":      `{{ "hello world" | trunc 5 }}`,
		
		// 数学函数
		"add":   `{{ add 1 2 }}`,
		"sub":   `{{ sub 5 3 }}`,
		"mul":   `{{ mul 3 4 }}`,
		"div":   `{{ div 8 2 }}`,
		"mod":   `{{ mod 10 3 }}`,
		"max":   `{{ max 1 2 }}`,
		"min":   `{{ min 1 2 }}`,
		"floor": `{{ floor 3.7 }}`,
		"ceil":  `{{ ceil 3.2 }}`,
		"round": `{{ round 3.6 }}`,
		
		// 条件函数
		"default":  `{{ .Value | default "默认值" }}`,
		"empty":    `{{ empty .Value }}`,
		"coalesce": `{{ coalesce .A .B "default" }}`,
		"ternary":  `{{ ternary "yes" "no" .Condition }}`,
		
		// 随机函数
		"randAlphaNum": `{{ randAlphaNum 8 }}`,
		"randAlpha":    `{{ randAlpha 6 }}`,
		"randNumeric":  `{{ randNumeric 4 }}`,
		"randInt":      `{{ randInt 1 100 }}`,
		"uuidv4":       `{{ uuidv4 }}`,
		
		// 日期函数
		"date":   `{{ date "2006-01-02" now }}`,
		"ago":    `{{ ago .CreatedAt }}`,
		"toDate": `{{ toDate "2006-01-02" "2023-12-25" }}`,
		
		// 类型转换
		"toString":     `{{ .Value | toString }}`,
		"toInt":        `{{ "123" | toInt }}`,
		"toFloat64":    `{{ "3.14" | toFloat64 }}`,
		"toBool":       `{{ "true" | toBool }}`,
		"toJson":       `{{ .Data | toJson }}`,
		"fromJson":     `{{ .JsonString | fromJson }}`,
		"toPrettyJson": `{{ .Data | toPrettyJson }}`,
		
		// 列表函数
		"list":    `{{ list "a" "b" "c" }}`,
		"first":   `{{ list "a" "b" "c" | first }}`,
		"last":    `{{ list "a" "b" "c" | last }}`,
		"rest":    `{{ list "a" "b" "c" | rest }}`,
		"initial": `{{ list "a" "b" "c" | initial }}`,
		"reverse": `{{ list "a" "b" "c" | reverse }}`,
		"uniq":    `{{ list "a" "b" "a" | uniq }}`,
		"compact": `{{ list "a" "" "b" | compact }}`,
		"append":  `{{ list "a" "b" | append "c" }}`,
		"prepend": `{{ list "b" "c" | prepend "a" }}`,
		"has":     `{{ list "a" "b" "c" | has "b" }}`,
		"without": `{{ list "a" "b" "c" | without "b" }}`,
		
		// 字典函数
		"dict":   `{{ dict "key1" "value1" "key2" "value2" }}`,
		"get":    `{{ .Dict | get "key" }}`,
		"set":    `{{ .Dict | set "key" "value" }}`,
		"unset":  `{{ .Dict | unset "key" }}`,
		"hasKey": `{{ .Dict | hasKey "key" }}`,
		"keys":   `{{ .Dict | keys }}`,
		"values": `{{ .Dict | values }}`,
		
		// 编码函数
		"b64enc":    `{{ "hello" | b64enc }}`,
		"b64dec":    `{{ "aGVsbG8=" | b64dec }}`,
		"urlquery":  `{{ "hello world" | urlquery }}`,
		"sha1sum":   `{{ "hello" | sha1sum }}`,
		"sha256sum": `{{ "hello" | sha256sum }}`,
		"md5sum":    `{{ "hello" | md5sum }}`,
	}
	if example, exists := examples[name]; exists {
		return example
	}
	return "{{ " + name + " }}"
}

func (s *sBuiltinFunctions) getFunctionInsertText(name string) string {
	insertTexts := map[string]string{
		// 字符串函数
		"upper":      `{{ . | upper }}`,
		"lower":      `{{ . | lower }}`,
		"title":      `{{ . | title }}`,
		"trim":       `{{ . | trim }}`,
		"trimPrefix": `{{ . | trimPrefix "prefix" }}`,
		"trimSuffix": `{{ . | trimSuffix "suffix" }}`,
		"replace":    `{{ . | replace "old" "new" }}`,
		"split":      `{{ . | split "," }}`,
		"join":       `{{ list "a" "b" "c" | join "," }}`,
		"contains":   `{{ . | contains "substr" }}`,
		"hasPrefix":  `{{ . | hasPrefix "prefix" }}`,
		"hasSuffix":  `{{ . | hasSuffix "suffix" }}`,
		"repeat":     `{{ . | repeat 3 }}`,
		"substr":     `{{ . | substr 0 5 }}`,
		"trunc":      `{{ . | trunc 10 }}`,

		// 数学函数
		"add":   `{{ add 1 2 }}`,
		"sub":   `{{ sub 5 3 }}`,
		"mul":   `{{ mul 3 4 }}`,
		"div":   `{{ div 8 2 }}`,
		"mod":   `{{ mod 10 3 }}`,
		"max":   `{{ max 1 2 }}`,
		"min":   `{{ min 1 2 }}`,
		"floor": `{{ floor 3.7 }}`,
		"ceil":  `{{ ceil 3.2 }}`,
		"round": `{{ round 3.6 }}`,

		// 条件函数
		"default":  `{{ . | default "默认值" }}`,
		"empty":    `{{ empty . }}`,
		"coalesce": `{{ coalesce .A .B "default" }}`,
		"ternary":  `{{ ternary "yes" "no" .condition }}`,

		// 随机函数
		"randAlphaNum": `{{ randAlphaNum 8 }}`,
		"randAlpha":    `{{ randAlpha 6 }}`,
		"randNumeric":  `{{ randNumeric 4 }}`,
		"randInt":      `{{ randInt 1 100 }}`,
		"uuidv4":       `{{ uuidv4 }}`,

		// 日期函数
		"date":   `{{ date "2006-01-02" now }}`,
		"ago":    `{{ ago .createdAt }}`,
		"toDate": `{{ toDate "2006-01-02" "2023-12-25" }}`,

		// 类型转换
		"toString":     `{{ . | toString }}`,
		"toInt":        `{{ . | toInt }}`,
		"toFloat64":    `{{ . | toFloat64 }}`,
		"toBool":       `{{ . | toBool }}`,
		"toJson":       `{{ . | toJson }}`,
		"fromJson":     `{{ . | fromJson }}`,
		"toPrettyJson": `{{ . | toPrettyJson }}`,

		// 列表函数
		"list":    `{{ list "a" "b" "c" }}`,
		"first":   `{{ . | first }}`,
		"last":    `{{ . | last }}`,
		"rest":    `{{ . | rest }}`,
		"initial": `{{ . | initial }}`,
		"reverse": `{{ . | reverse }}`,
		"uniq":    `{{ . | uniq }}`,
		"compact": `{{ . | compact }}`,
		"append":  `{{ . | append "item" }}`,
		"prepend": `{{ . | prepend "item" }}`,
		"has":     `{{ . | has "item" }}`,
		"without": `{{ . | without "item" }}`,

		// 字典函数
		"dict":   `{{ dict "key1" "value1" "key2" "value2" }}`,
		"get":    `{{ . | get "key" }}`,
		"set":    `{{ . | set "key" "value" }}`,
		"unset":  `{{ . | unset "key" }}`,
		"hasKey": `{{ . | hasKey "key" }}`,
		"keys":   `{{ . | keys }}`,
		"values": `{{ . | values }}`,

		// 编码函数
		"b64enc":    `{{ . | b64enc }}`,
		"b64dec":    `{{ . | b64dec }}`,
		"urlquery":  `{{ . | urlquery }}`,
		"sha1sum":   `{{ . | sha1sum }}`,
		"sha256sum": `{{ . | sha256sum }}`,
		"md5sum":    `{{ . | md5sum }}`,
	}
	if insert, exists := insertTexts[name]; exists {
		return insert
	}
	return "{{ " + name + " }}"
}

func (s *sBuiltinFunctions) getReturnType(fnType reflect.Type) string {
	if fnType.Kind() != reflect.Func || fnType.NumOut() == 0 {
		return "any"
	}
	
	returnType := fnType.Out(0)
	return s.getTypeString(returnType)
}

func (s *sBuiltinFunctions) getTypeString(t reflect.Type) string {
	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "int"
	case reflect.Float32, reflect.Float64:
		return "float"
	case reflect.Bool:
		return "bool"
	case reflect.Slice, reflect.Array:
		return "array"
	case reflect.Map:
		return "map"
	default:
		return "any"
	}
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