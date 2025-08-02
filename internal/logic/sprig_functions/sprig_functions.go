package sprig_functions

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

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
	categories, err := s.parseSprigDocuments()
	if err != nil {
		return nil, err
	}

	total := 0
	for _, category := range categories {
		total += len(category.Functions)
	}

	return &model.SprigFunctionsResponse{
		Categories: categories,
		Total:      total,
	}, nil
}

// parseSprigDocuments 解析Sprig文档
func (s *sSprigFunctions) parseSprigDocuments() ([]model.SprigFunctionCategory, error) {
	// 定义分类映射
	categoryConfigs := []struct {
		Name         string
		Description  string
		DocumentFile string
		FileName     string
	}{
		{
			Name:         "字符串函数",
			Description:  "字符串操作和处理函数",
			DocumentFile: "string.md",
			FileName:     "string.md",
		},
		{
			Name:         "整数数学函数",
			Description:  "整数数学运算函数",
			DocumentFile: "Integer Math Functions.md",
			FileName:     "Integer Math Functions.md",
		},
		{
			Name:         "浮点数学函数",
			Description:  "浮点数数学运算函数",
			DocumentFile: "Float Math Functions.md",
			FileName:     "Float Math Functions.md",
		},
		{
			Name:         "日期函数",
			Description:  "日期时间处理函数",
			DocumentFile: "Date Functions.md",
			FileName:     "Date Functions.md",
		},
		{
			Name:         "默认值函数",
			Description:  "默认值和条件函数",
			DocumentFile: "Default Functions.md",
			FileName:     "Default Functions.md",
		},
		{
			Name:         "流程控制函数",
			Description:  "条件判断和流程控制函数",
			DocumentFile: "Flow Control Functions.md",
			FileName:     "Flow Control Functions.md",
		},
		{
			Name:         "类型转换函数",
			Description:  "数据类型转换函数",
			DocumentFile: "Type Conversion Functions.md",
			FileName:     "Type Conversion Functions.md",
		},
		{
			Name:         "列表函数",
			Description:  "列表操作和处理函数",
			DocumentFile: "Lists and List Functions.md",
			FileName:     "Lists and List Functions.md",
		},
		{
			Name:         "字典函数",
			Description:  "字典操作和处理函数",
			DocumentFile: "Dictionaries and Dict Functions.md",
			FileName:     "Dictionaries and Dict Functions.md",
		},
		{
			Name:         "编码函数",
			Description:  "编码解码函数",
			DocumentFile: "Encoding Functions.md",
			FileName:     "Encoding Functions.md",
		},
		{
			Name:         "路径函数",
			Description:  "路径和文件名处理函数",
			DocumentFile: "Path and Filepath Functions.md",
			FileName:     "Path and Filepath Functions.md",
		},
		{
			Name:         "字符串切片函数",
			Description:  "字符串数组操作函数",
			DocumentFile: "String Slice Functions.md",
			FileName:     "String Slice Functions.md",
		},
		{
			Name:         "整数切片函数",
			Description:  "整数数组操作函数",
			DocumentFile: "Integer Slice Functions.md",
			FileName:     "Integer Slice Functions.md",
		},
	}

	var categories []model.SprigFunctionCategory

	for _, config := range categoryConfigs {
		filePath := filepath.Join("docs", "sprig", config.FileName)
		functions, err := s.parseMarkdownFile(filePath, config.Name)
		if err != nil {
			// 如果文件不存在或解析失败，记录错误但继续处理其他文件
			continue
		}

		if len(functions) > 0 {
			categories = append(categories, model.SprigFunctionCategory{
				Name:         config.Name,
				Description:  config.Description,
				DocumentFile: config.DocumentFile,
				Functions:    functions,
			})
		}
	}

	// 添加高级函数分类
	advancedCategories := s.parseAdvancedFunctions()
	categories = append(categories, advancedCategories...)

	return categories, nil
}

// parseMarkdownFile 解析单个markdown文件
func (s *sSprigFunctions) parseMarkdownFile(filePath, category string) ([]model.SprigFunction, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return s.extractFunctionsFromContent(string(content), category), nil
}

// extractFunctionsFromContent 从markdown内容中提取函数信息
func (s *sSprigFunctions) extractFunctionsFromContent(content, category string) []model.SprigFunction {
	var functions []model.SprigFunction

	// 移除SimpRead的头部信息
	content = s.cleanupContent(content)

	// 按函数分割内容
	functionSections := s.splitIntoFunctionSections(content)

	for _, section := range functionSections {
		function := s.parseFunctionSection(section, category)
		if function.Name != "" {
			functions = append(functions, function)
		}
	}

	return functions
}

// cleanupContent 清理markdown内容
func (s *sSprigFunctions) cleanupContent(content string) string {
	// 移除SimpRead头部
	re := regexp.MustCompile(`(?s)^>.*?转码.*?\n\n`)
	content = re.ReplaceAllString(content, "")

	// 移除额外的空行
	re = regexp.MustCompile(`\n{3,}`)
	content = re.ReplaceAllString(content, "\n\n")

	return strings.TrimSpace(content)
}

// splitIntoFunctionSections 将内容分割为函数段落
func (s *sSprigFunctions) splitIntoFunctionSections(content string) []string {
	// 使用标题分割，标题格式可能是 "functionName\n----" 或 "# functionName"
	lines := strings.Split(content, "\n")
	var sections []string
	var currentSection strings.Builder
	var inFunction bool

	for i, line := range lines {
		line = strings.TrimSpace(line)
		
		// 检查是否是函数标题
		var nextLine string
		if i < len(lines)-1 {
			nextLine = lines[i+1]
		}
		
		if s.isFunctionTitle(line, nextLine) {
			// 保存当前函数段落
			if inFunction && currentSection.Len() > 0 {
				sections = append(sections, currentSection.String())
				currentSection.Reset()
			}
			inFunction = true
		}

		if inFunction {
			currentSection.WriteString(line + "\n")
		}
	}

	// 添加最后一个段落
	if currentSection.Len() > 0 {
		sections = append(sections, currentSection.String())
	}

	return sections
}

// isFunctionTitle 判断是否是函数标题
func (s *sSprigFunctions) isFunctionTitle(line, nextLine string) bool {
	// 检查下一行是否是横线分隔符
	if nextLine != "" && regexp.MustCompile(`^-{3,}$`).MatchString(nextLine) {
		return true
	}
	
	// 检查是否是markdown标题格式
	if regexp.MustCompile(`^#+\s+\w+`).MatchString(line) {
		return true
	}

	// 检查是否是简单的函数名（全小写字母，可能包含数字）
	if regexp.MustCompile(`^[a-z][a-zA-Z0-9]*$`).MatchString(line) && len(line) < 30 {
		return true
	}

	return false
}

// parseFunctionSection 解析单个函数段落
func (s *sSprigFunctions) parseFunctionSection(section, category string) model.SprigFunction {
	lines := strings.Split(section, "\n")
	if len(lines) == 0 {
		return model.SprigFunction{}
	}

	// 提取函数名
	functionName := s.extractFunctionName(lines[0])
	if functionName == "" {
		return model.SprigFunction{}
	}

	// 根据函数名获取详细信息
	return s.buildSprigFunction(functionName, section, category)
}

// extractFunctionName 提取函数名
func (s *sSprigFunctions) extractFunctionName(line string) string {
	line = strings.TrimSpace(line)
	
	// 移除markdown标题标记
	line = regexp.MustCompile(`^#+\s*`).ReplaceAllString(line, "")
	
	// 提取第一个单词作为函数名
	parts := strings.Fields(line)
	if len(parts) > 0 {
		name := strings.ToLower(parts[0])
		// 验证是否是有效的函数名
		if regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]*$`).MatchString(name) {
			return name
		}
	}
	
	return ""
}

// buildSprigFunction 构建Sprig函数信息
func (s *sSprigFunctions) buildSprigFunction(name, content, category string) model.SprigFunction {
	// 根据函数名获取预定义的信息
	functionInfo := s.getPredefinedFunctionInfo(name)
	
	// 设置基本信息
	function := model.SprigFunction{
		Name:        name,
		DisplayName: functionInfo.DisplayName,
		Description: functionInfo.Description,
		Params:      functionInfo.Params,
		ReturnType:  functionInfo.ReturnType,
		Category:    category,
		Example:     functionInfo.Example,
		Usage:       functionInfo.Usage,
		InsertText:  functionInfo.InsertText,
		Note:        functionInfo.Note,
		Aliases:     functionInfo.Aliases,
	}

	return function
}

// parseAdvancedFunctions 解析高级函数目录下的函数
func (s *sSprigFunctions) parseAdvancedFunctions() []model.SprigFunctionCategory {
	advancedDir := filepath.Join("docs", "sprig", "Advanced Functions")
	
	advancedConfigs := []struct {
		Name         string
		Description  string
		FileName     string
	}{
		{
			Name:         "加密安全函数",
			Description:  "加密和安全相关函数",
			FileName:     "Cryptographic and Security Functions.md",
		},
		{
			Name:         "网络函数",
			Description:  "网络相关函数",
			FileName:     "Network Functions.md",
		},
		{
			Name:         "操作系统函数",
			Description:  "操作系统相关函数",
			FileName:     "OS Functions.md",
		},
		{
			Name:         "反射函数",
			Description:  "反射和类型检查函数",
			FileName:     "Reflection Functions.md",
		},
		{
			Name:         "语义版本函数",
			Description:  "语义版本处理函数",
			FileName:     "Semantic Version Functions.md",
		},
		{
			Name:         "URL函数",
			Description:  "URL处理函数",
			FileName:     "URL Functions.md",
		},
		{
			Name:         "UUID函数",
			Description:  "UUID生成函数",
			FileName:     "UUID Functions.md",
		},
	}

	var categories []model.SprigFunctionCategory

	for _, config := range advancedConfigs {
		filePath := filepath.Join(advancedDir, config.FileName)
		functions, err := s.parseMarkdownFile(filePath, config.Name)
		if err != nil {
			continue
		}

		if len(functions) > 0 {
			categories = append(categories, model.SprigFunctionCategory{
				Name:         config.Name,
				Description:  config.Description,
				DocumentFile: "Advanced Functions/" + config.FileName,
				Functions:    functions,
			})
		}
	}

	return categories
}

// getPredefinedFunctionInfo 获取预定义的函数信息
func (s *sSprigFunctions) getPredefinedFunctionInfo(name string) model.SprigFunction {
	// 这里会包含完整的中文翻译和详细信息
	predefinedFunctions := s.getAllPredefinedFunctions()
	
	if info, exists := predefinedFunctions[name]; exists {
		return info
	}

	// 如果没有预定义信息，返回基本结构
	return model.SprigFunction{
		Name:        name,
		DisplayName: strings.Title(name),
		Description: "Sprig " + name + " 函数",
		ReturnType:  "string",
		Example:     "{{ " + name + " }}",
		InsertText:  "{{ " + name + " }}",
	}
}

// getAllPredefinedFunctions 获取所有预定义的函数信息
func (s *sSprigFunctions) getAllPredefinedFunctions() map[string]model.SprigFunction {
	return map[string]model.SprigFunction{
		// 字符串函数
		"trim": {
			Name:        "trim",
			DisplayName: "去除空格",
			Description: "去除字符串两端的空格字符",
			Params: []model.SprigFunctionParam{
				{Name: "string", Type: "string", Required: true, Description: "要处理的字符串"},
			},
			ReturnType: "string",
			Example:    `{{ "  hello  " | trim }}`,
			Usage:      "trim 函数用于移除字符串两端的空白字符（空格、制表符、换行符等）。",
			InsertText: `{{ . | trim }}`,
		},
		"trimAll": {
			Name:        "trimAll",
			DisplayName: "去除指定字符",
			Description: "从字符串两端移除指定的字符",
			Params: []model.SprigFunctionParam{
				{Name: "chars", Type: "string", Required: true, Description: "要移除的字符"},
				{Name: "string", Type: "string", Required: true, Description: "原字符串"},
			},
			ReturnType: "string",
			Example:    `{{ trimAll "$" "$5.00$" }}`,
			Usage:      "trimAll 函数从字符串的开头和结尾移除指定的字符。",
			InsertText: `{{ trimAll "char" . }}`,
		},
		"trimSuffix": {
			Name:        "trimSuffix",
			DisplayName: "去除后缀",
			Description: "从字符串末尾移除指定的后缀",
			Params: []model.SprigFunctionParam{
				{Name: "suffix", Type: "string", Required: true, Description: "要移除的后缀"},
				{Name: "string", Type: "string", Required: true, Description: "原字符串"},
			},
			ReturnType: "string",
			Example:    `{{ trimSuffix "-" "hello-" }}`,
			Usage:      "trimSuffix 函数仅从字符串末尾移除指定的后缀。如果字符串不以该后缀结尾，则原样返回。",
			InsertText: `{{ trimSuffix "suffix" . }}`,
		},
		"trimPrefix": {
			Name:        "trimPrefix",
			DisplayName: "去除前缀",
			Description: "从字符串开头移除指定的前缀",
			Params: []model.SprigFunctionParam{
				{Name: "prefix", Type: "string", Required: true, Description: "要移除的前缀"},
				{Name: "string", Type: "string", Required: true, Description: "原字符串"},
			},
			ReturnType: "string",
			Example:    `{{ trimPrefix "-" "-hello" }}`,
			Usage:      "trimPrefix 函数仅从字符串开头移除指定的前缀。如果字符串不以该前缀开头，则原样返回。",
			InsertText: `{{ trimPrefix "prefix" . }}`,
		},
		"upper": {
			Name:        "upper",
			DisplayName: "转大写",
			Description: "将字符串转换为大写字母",
			Params: []model.SprigFunctionParam{
				{Name: "string", Type: "string", Required: true, Description: "要转换的字符串"},
			},
			ReturnType: "string",
			Example:    `{{ "hello" | upper }}`,
			Usage:      "upper 函数将整个字符串转换为大写字母。",
			InsertText: `{{ . | upper }}`,
		},
		"lower": {
			Name:        "lower",
			DisplayName: "转小写",
			Description: "将字符串转换为小写字母",
			Params: []model.SprigFunctionParam{
				{Name: "string", Type: "string", Required: true, Description: "要转换的字符串"},
			},
			ReturnType: "string",
			Example:    `{{ "HELLO" | lower }}`,
			Usage:      "lower 函数将整个字符串转换为小写字母。",
			InsertText: `{{ . | lower }}`,
		},
		"title": {
			Name:        "title",
			DisplayName: "首字母大写",
			Description: "将字符串转换为标题格式（首字母大写）",
			Params: []model.SprigFunctionParam{
				{Name: "string", Type: "string", Required: true, Description: "要转换的字符串"},
			},
			ReturnType: "string",
			Example:    `{{ "hello world" | title }}`,
			Usage:      "title 函数将字符串转换为标题格式，即每个单词的首字母大写。",
			InsertText: `{{ . | title }}`,
		},
		"untitle": {
			Name:        "untitle",
			DisplayName: "取消标题格式",
			Description: "移除标题格式，将字符串转换为全小写",
			Params: []model.SprigFunctionParam{
				{Name: "string", Type: "string", Required: true, Description: "要转换的字符串"},
			},
			ReturnType: "string",
			Example:    `{{ "Hello World" | untitle }}`,
			Usage:      "untitle 函数将标题格式的字符串转换为全小写。",
			InsertText: `{{ . | untitle }}`,
		},
		"repeat": {
			Name:        "repeat",
			DisplayName: "重复字符串",
			Description: "将字符串重复指定次数",
			Params: []model.SprigFunctionParam{
				{Name: "count", Type: "int", Required: true, Description: "重复次数"},
				{Name: "string", Type: "string", Required: true, Description: "要重复的字符串"},
			},
			ReturnType: "string",
			Example:    `{{ repeat 3 "hello" }}`,
			Usage:      "repeat 函数将字符串重复指定的次数。",
			InsertText: `{{ repeat 3 . }}`,
		},
		"substr": {
			Name:        "substr",
			DisplayName: "截取子字符串",
			Description: "从字符串中截取指定位置和长度的子字符串",
			Params: []model.SprigFunctionParam{
				{Name: "start", Type: "int", Required: true, Description: "起始位置（从0开始）"},
				{Name: "end", Type: "int", Required: true, Description: "结束位置"},
				{Name: "string", Type: "string", Required: true, Description: "原字符串"},
			},
			ReturnType: "string",
			Example:    `{{ substr 0 5 "hello world" }}`,
			Usage:      "substr 函数从字符串中提取子字符串。接受三个参数：起始位置、结束位置和原字符串。",
			InsertText: `{{ substr 0 5 . }}`,
		},
		"nospace": {
			Name:        "nospace",
			DisplayName: "移除空格",
			Description: "移除字符串中的所有空白字符",
			Params: []model.SprigFunctionParam{
				{Name: "string", Type: "string", Required: true, Description: "要处理的字符串"},
			},
			ReturnType: "string",
			Example:    `{{ "hello w o r l d" | nospace }}`,
			Usage:      "nospace 函数移除字符串中的所有空白字符（空格、制表符、换行符等）。",
			InsertText: `{{ . | nospace }}`,
		},
		"trunc": {
			Name:        "trunc",
			DisplayName: "截断字符串",
			Description: "截断字符串到指定长度",
			Params: []model.SprigFunctionParam{
				{Name: "length", Type: "int", Required: true, Description: "截断长度"},
				{Name: "string", Type: "string", Required: true, Description: "要截断的字符串"},
			},
			ReturnType: "string",
			Example:    `{{ trunc 5 "hello world" }}`,
			Usage:      "trunc 函数将字符串截断到指定的长度。如果长度为负数，则从字符串末尾开始截断。",
			InsertText: `{{ trunc 10 . }}`,
		},
		"abbrev": {
			Name:        "abbrev",
			DisplayName: "缩写字符串",
			Description: "截断字符串并添加省略号",
			Params: []model.SprigFunctionParam{
				{Name: "length", Type: "int", Required: true, Description: "最大长度"},
				{Name: "string", Type: "string", Required: true, Description: "要缩写的字符串"},
			},
			ReturnType: "string",
			Example:    `{{ abbrev 5 "hello world" }}`,
			Usage:      "abbrev 函数截断字符串到指定长度并添加省略号(...)。省略号的长度计算在最大长度内。",
			InsertText: `{{ abbrev 10 . }}`,
		},
		"abbrevboth": {
			Name:        "abbrevboth",
			DisplayName: "两端缩写",
			Description: "从字符串两端进行缩写",
			Params: []model.SprigFunctionParam{
				{Name: "left", Type: "int", Required: true, Description: "左侧偏移量"},
				{Name: "length", Type: "int", Required: true, Description: "最大长度"},
				{Name: "string", Type: "string", Required: true, Description: "要缩写的字符串"},
			},
			ReturnType: "string",
			Example:    `{{ abbrevboth 5 10 "1234567890123456" }}`,
			Usage:      "abbrevboth 函数从字符串的指定位置开始，截取指定长度的内容，两端用省略号表示。",
			InsertText: `{{ abbrevboth 2 10 . }}`,
		},
		"initials": {
			Name:        "initials",
			DisplayName: "提取首字母",
			Description: "提取多个单词的首字母并组合",
			Params: []model.SprigFunctionParam{
				{Name: "string", Type: "string", Required: true, Description: "包含多个单词的字符串"},
			},
			ReturnType: "string",
			Example:    `{{ "First Template" | initials }}`,
			Usage:      "initials 函数提取字符串中每个单词的首字母并将它们组合成一个字符串。",
			InsertText: `{{ . | initials }}`,
		},

		// 数学函数
		"add": {
			Name:        "add",
			DisplayName: "加法",
			Description: "计算两个或多个数的和",
			Params: []model.SprigFunctionParam{
				{Name: "numbers", Type: "number", Required: true, Description: "要相加的数字"},
			},
			ReturnType: "number",
			Example:    `{{ add 1 2 3 }}`,
			Usage:      "add 函数计算两个或多个数字的和。可以接受任意数量的参数。",
			InsertText: `{{ add 1 2 }}`,
		},
		"add1": {
			Name:        "add1",
			DisplayName: "加一",
			Description: "给数字加1",
			Params: []model.SprigFunctionParam{
				{Name: "number", Type: "number", Required: true, Description: "要加1的数字"},
			},
			ReturnType: "number",
			Example:    `{{ 5 | add1 }}`,
			Usage:      "add1 函数给数字增加1。这是 add 函数的简化版本。",
			InsertText: `{{ . | add1 }}`,
		},
		"sub": {
			Name:        "sub",
			DisplayName: "减法",
			Description: "计算两个数的差",
			Params: []model.SprigFunctionParam{
				{Name: "a", Type: "number", Required: true, Description: "被减数"},
				{Name: "b", Type: "number", Required: true, Description: "减数"},
			},
			ReturnType: "number",
			Example:    `{{ sub 5 3 }}`,
			Usage:      "sub 函数计算第一个数减去第二个数的结果。",
			InsertText: `{{ sub 10 . }}`,
		},
		"div": {
			Name:        "div",
			DisplayName: "除法",
			Description: "计算两个数的商",
			Params: []model.SprigFunctionParam{
				{Name: "a", Type: "number", Required: true, Description: "被除数"},
				{Name: "b", Type: "number", Required: true, Description: "除数"},
			},
			ReturnType: "number",
			Example:    `{{ div 8 2 }}`,
			Usage:      "div 函数进行整数除法运算。结果为整数，不包含小数部分。",
			InsertText: `{{ div . 2 }}`,
		},
		"mod": {
			Name:        "mod",
			DisplayName: "取模",
			Description: "计算两个数相除的余数",
			Params: []model.SprigFunctionParam{
				{Name: "a", Type: "number", Required: true, Description: "被除数"},
				{Name: "b", Type: "number", Required: true, Description: "除数"},
			},
			ReturnType: "number",
			Example:    `{{ mod 10 3 }}`,
			Usage:      "mod 函数返回第一个数除以第二个数的余数。",
			InsertText: `{{ mod . 2 }}`,
		},
		"mul": {
			Name:        "mul",
			DisplayName: "乘法",
			Description: "计算两个或多个数的积",
			Params: []model.SprigFunctionParam{
				{Name: "numbers", Type: "number", Required: true, Description: "要相乘的数字"},
			},
			ReturnType: "number",
			Example:    `{{ mul 3 4 }}`,
			Usage:      "mul 函数计算两个或多个数字的乘积。可以接受任意数量的参数。",
			InsertText: `{{ mul . 2 }}`,
		},
		"max": {
			Name:        "max",
			DisplayName: "最大值",
			Description: "返回一系列数字中的最大值",
			Params: []model.SprigFunctionParam{
				{Name: "numbers", Type: "number", Required: true, Description: "数字列表"},
			},
			ReturnType: "number",
			Example:    `{{ max 1 2 3 }}`,
			Usage:      "max 函数返回给定数字中的最大值。",
			InsertText: `{{ max 1 . }}`,
		},
		"min": {
			Name:        "min",
			DisplayName: "最小值",
			Description: "返回一系列数字中的最小值",
			Params: []model.SprigFunctionParam{
				{Name: "numbers", Type: "number", Required: true, Description: "数字列表"},
			},
			ReturnType: "number",
			Example:    `{{ min 1 2 3 }}`,
			Usage:      "min 函数返回给定数字中的最小值。",
			InsertText: `{{ min 1 . }}`,
		},
		"floor": {
			Name:        "floor",
			DisplayName: "向下取整",
			Description: "返回小于或等于输入值的最大整数",
			Params: []model.SprigFunctionParam{
				{Name: "number", Type: "float", Required: true, Description: "要向下取整的数字"},
			},
			ReturnType: "float",
			Example:    `{{ floor 123.9999 }}`,
			Usage:      "floor 函数返回小于或等于给定数字的最大整数值。",
			InsertText: `{{ . | floor }}`,
		},
		"ceil": {
			Name:        "ceil",
			DisplayName: "向上取整",
			Description: "返回大于或等于输入值的最小整数",
			Params: []model.SprigFunctionParam{
				{Name: "number", Type: "float", Required: true, Description: "要向上取整的数字"},
			},
			ReturnType: "float",
			Example:    `{{ ceil 123.001 }}`,
			Usage:      "ceil 函数返回大于或等于给定数字的最小整数值。",
			InsertText: `{{ . | ceil }}`,
		},
		"round": {
			Name:        "round",
			DisplayName: "四舍五入",
			Description: "将数字四舍五入到指定的小数位数",
			Params: []model.SprigFunctionParam{
				{Name: "number", Type: "float", Required: true, Description: "要四舍五入的数字"},
				{Name: "precision", Type: "int", Required: false, Description: "小数位数", Default: "0"},
			},
			ReturnType: "float",
			Example:    `{{ round 123.555555 3 }}`,
			Usage:      "round 函数将浮点数四舍五入到指定的小数位数。",
			InsertText: `{{ round . 2 }}`,
		},

		// 随机函数
		"randAlphaNum": {
			Name:        "randAlphaNum",
			DisplayName: "随机字母数字",
			Description: "生成指定长度的随机字母数字字符串",
			Params: []model.SprigFunctionParam{
				{Name: "length", Type: "int", Required: true, Description: "字符串长度"},
			},
			ReturnType: "string",
			Example:    `{{ randAlphaNum 8 }}`,
			Usage:      "randAlphaNum 函数生成包含字母（a-z, A-Z）和数字（0-9）的随机字符串。使用加密安全的随机数生成器。",
			InsertText: `{{ randAlphaNum 8 }}`,
		},
		"randAlpha": {
			Name:        "randAlpha",
			DisplayName: "随机字母",
			Description: "生成指定长度的随机字母字符串",
			Params: []model.SprigFunctionParam{
				{Name: "length", Type: "int", Required: true, Description: "字符串长度"},
			},
			ReturnType: "string",
			Example:    `{{ randAlpha 6 }}`,
			Usage:      "randAlpha 函数生成只包含字母（a-z, A-Z）的随机字符串。使用加密安全的随机数生成器。",
			InsertText: `{{ randAlpha 6 }}`,
		},
		"randNumeric": {
			Name:        "randNumeric",
			DisplayName: "随机数字",
			Description: "生成指定长度的随机数字字符串",
			Params: []model.SprigFunctionParam{
				{Name: "length", Type: "int", Required: true, Description: "字符串长度"},
			},
			ReturnType: "string",
			Example:    `{{ randNumeric 4 }}`,
			Usage:      "randNumeric 函数生成只包含数字（0-9）的随机字符串。使用加密安全的随机数生成器。",
			InsertText: `{{ randNumeric 4 }}`,
		},
		"randAscii": {
			Name:        "randAscii",
			DisplayName: "随机ASCII",
			Description: "生成指定长度的随机ASCII字符串",
			Params: []model.SprigFunctionParam{
				{Name: "length", Type: "int", Required: true, Description: "字符串长度"},
			},
			ReturnType: "string",
			Example:    `{{ randAscii 10 }}`,
			Usage:      "randAscii 函数生成包含所有可打印ASCII字符的随机字符串。使用加密安全的随机数生成器。",
			InsertText: `{{ randAscii 10 }}`,
		},
		"randInt": {
			Name:        "randInt",
			DisplayName: "随机整数",
			Description: "生成指定范围内的随机整数",
			Params: []model.SprigFunctionParam{
				{Name: "min", Type: "int", Required: true, Description: "最小值（包含）"},
				{Name: "max", Type: "int", Required: true, Description: "最大值（不包含）"},
			},
			ReturnType: "int",
			Example:    `{{ randInt 1 100 }}`,
			Usage:      "randInt 函数生成从最小值（包含）到最大值（不包含）范围内的随机整数。",
			InsertText: `{{ randInt 1 100 }}`,
		},

		// 条件和默认值函数
		"default": {
			Name:        "default",
			DisplayName: "默认值",
			Description: "如果值为空则返回默认值",
			Params: []model.SprigFunctionParam{
				{Name: "default", Type: "any", Required: true, Description: "默认值"},
				{Name: "value", Type: "any", Required: true, Description: "要检查的值"},
			},
			ReturnType: "any",
			Example:    `{{ .Value | default "默认值" }}`,
			Usage:      "default 函数检查值是否为空（nil、空字符串、0等），如果为空则返回指定的默认值。",
			InsertText: `{{ . | default "默认值" }}`,
		},
		"empty": {
			Name:        "empty",
			DisplayName: "判断为空",
			Description: "检查值是否为空",
			Params: []model.SprigFunctionParam{
				{Name: "value", Type: "any", Required: true, Description: "要检查的值"},
			},
			ReturnType: "bool",
			Example:    `{{ empty .Value }}`,
			Usage:      "empty 函数检查值是否为空（nil、空字符串、空切片、空映射等）。",
			InsertText: `{{ empty . }}`,
		},
		"coalesce": {
			Name:        "coalesce",
			DisplayName: "合并空值",
			Description: "返回第一个非空值",
			Params: []model.SprigFunctionParam{
				{Name: "values", Type: "any", Required: true, Description: "值列表"},
			},
			ReturnType: "any",
			Example:    `{{ coalesce .A .B "default" }}`,
			Usage:      "coalesce 函数从给定的值列表中返回第一个非空值。",
			InsertText: `{{ coalesce .value1 .value2 "default" }}`,
		},
		"ternary": {
			Name:        "ternary",
			DisplayName: "三元操作",
			Description: "根据条件返回不同的值",
			Params: []model.SprigFunctionParam{
				{Name: "trueValue", Type: "any", Required: true, Description: "条件为真时的值"},
				{Name: "falseValue", Type: "any", Required: true, Description: "条件为假时的值"},
				{Name: "condition", Type: "bool", Required: true, Description: "条件"},
			},
			ReturnType: "any",
			Example:    `{{ ternary "yes" "no" .condition }}`,
			Usage:      "ternary 函数实现三元条件操作，根据条件的真假返回相应的值。",
			InsertText: `{{ ternary "true-value" "false-value" .condition }}`,
		},

		// 更多函数可以继续添加...
	}
}