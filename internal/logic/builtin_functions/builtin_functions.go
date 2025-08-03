package builtin_functions

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/ciclebyte/template_starter/internal/model"
	"github.com/ciclebyte/template_starter/internal/service"
)

type sBuiltinFunctions struct{}

func init() {
	service.RegisterBuiltinFunctions(New())
}

func New() *sBuiltinFunctions {
	return &sBuiltinFunctions{}
}

// GetBuiltinFunctions 获取所有内置函数（分组显示实用的自定义函数）
func (s *sBuiltinFunctions) GetBuiltinFunctions(ctx context.Context) (*model.BuiltinFunctionsResponse, error) {
	categories := []model.BuiltinFunctionCategory{
		s.buildFileOperationFunctions(),
		s.buildStringProcessingFunctions(),
		s.buildCodeGenerationFunctions(),
		s.buildProjectStructureFunctions(),
	}

	total := 0
	for _, category := range categories {
		total += len(category.Functions)
	}

	return &model.BuiltinFunctionsResponse{
		Categories: categories,
		Total:      total,
	}, nil
}

// buildFileOperationFunctions 构建文件操作函数分组
func (s *sBuiltinFunctions) buildFileOperationFunctions() model.BuiltinFunctionCategory {
	return model.BuiltinFunctionCategory{
		Name:        "文件操作函数",
		Description: "文件和目录操作相关函数",
		Functions: []model.BuiltinFunction{
			{
				Name:        "readFile",
				DisplayName: "读取文件",
				Description: "读取指定文件的内容",
				Params: []model.FunctionParam{
					{
						Name:        "filepath",
						Type:        "string",
						Required:    true,
						Description: "文件路径",
					},
				},
				ReturnType: "string",
				Category:   "文件操作函数",
				Example:    `{{ readFile "config/app.yaml" }}`,
				InsertText: `{{ readFile "path/to/file" }}`,
			},
			{
				Name:        "fileExists",
				DisplayName: "文件是否存在",
				Description: "检查指定文件是否存在",
				Params: []model.FunctionParam{
					{
						Name:        "filepath",
						Type:        "string",
						Required:    true,
						Description: "文件路径",
					},
				},
				ReturnType: "bool",
				Category:   "文件操作函数",
				Example:    `{{ if fileExists "README.md" }}文件存在{{ end }}`,
				InsertText: `{{ fileExists "path/to/file" }}`,
			},
			{
				Name:        "dirExists",
				DisplayName: "目录是否存在",
				Description: "检查指定目录是否存在",
				Params: []model.FunctionParam{
					{
						Name:        "dirpath",
						Type:        "string",
						Required:    true,
						Description: "目录路径",
					},
				},
				ReturnType: "bool",
				Category:   "文件操作函数",
				Example:    `{{ if dirExists "src" }}目录存在{{ end }}`,
				InsertText: `{{ dirExists "path/to/dir" }}`,
			},
		},
	}
}

// buildStringProcessingFunctions 构建字符串处理函数分组
func (s *sBuiltinFunctions) buildStringProcessingFunctions() model.BuiltinFunctionCategory {
	return model.BuiltinFunctionCategory{
		Name:        "字符串处理函数",
		Description: "扩展的字符串处理和格式化函数",
		Functions: []model.BuiltinFunction{
			{
				Name:        "indent",
				DisplayName: "缩进字符串",
				Description: "为字符串的每一行添加指定数量的空格缩进",
				Params: []model.FunctionParam{
					{
						Name:        "spaces",
						Type:        "int",
						Required:    true,
						Description: "缩进空格数",
					},
					{
						Name:        "text",
						Type:        "string",
						Required:    true,
						Description: "要缩进的文本",
					},
				},
				ReturnType: "string",
				Category:   "字符串处理函数",
				Example:    `{{ indent 4 "line1\nline2" }}`,
				InsertText: `{{ indent 4 . }}`,
			},
			{
				Name:        "comment",
				DisplayName: "生成注释",
				Description: "生成指定样式的注释块",
				Params: []model.FunctionParam{
					{
						Name:        "style",
						Type:        "string",
						Required:    true,
						Description: "注释样式：line(//)、block(/**/)、hash(#)",
					},
					{
						Name:        "text",
						Type:        "string",
						Required:    true,
						Description: "注释内容",
					},
				},
				ReturnType: "string",
				Category:   "字符串处理函数",
				Example:    `{{ comment "line" "这是一个注释" }}`,
				InsertText: `{{ comment "line" . }}`,
			},
			{
				Name:        "quote",
				DisplayName: "智能引号",
				Description: "为字符串添加合适的引号（双引号、单引号或反引号）",
				Params: []model.FunctionParam{
					{
						Name:        "text",
						Type:        "string",
						Required:    true,
						Description: "要加引号的文本",
					},
					{
						Name:        "style",
						Type:        "string",
						Required:    false,
						Description: "引号样式：double、single、back，默认为double",
						Default:     "double",
					},
				},
				ReturnType: "string",
				Category:   "字符串处理函数",
				Example:    `{{ quote "hello world" }}`,
				InsertText: `{{ quote . }}`,
			},
		},
	}
}

// buildCodeGenerationFunctions 构建代码生成函数分组
func (s *sBuiltinFunctions) buildCodeGenerationFunctions() model.BuiltinFunctionCategory {
	return model.BuiltinFunctionCategory{
		Name:        "代码生成函数",
		Description: "代码生成和命名转换相关函数",
		Functions: []model.BuiltinFunction{
			{
				Name:        "camelToSnake",
				DisplayName: "驼峰转下划线",
				Description: "将驼峰命名转换为下划线命名（如：UserName -> user_name）",
				Params: []model.FunctionParam{
					{
						Name:        "text",
						Type:        "string",
						Required:    true,
						Description: "驼峰命名的字符串",
					},
				},
				ReturnType: "string",
				Category:   "代码生成函数",
				Example:    `{{ camelToSnake "UserName" }}`,
				InsertText: `{{ camelToSnake . }}`,
			},
			{
				Name:        "snakeToCamel",
				DisplayName: "下划线转驼峰",
				Description: "将下划线命名转换为驼峰命名（如：user_name -> UserName）",
				Params: []model.FunctionParam{
					{
						Name:        "text",
						Type:        "string",
						Required:    true,
						Description: "下划线命名的字符串",
					},
					{
						Name:        "firstUpper",
						Type:        "bool",
						Required:    false,
						Description: "首字母是否大写，默认为true",
						Default:     "true",
					},
				},
				ReturnType: "string",
				Category:   "代码生成函数",
				Example:    `{{ snakeToCamel "user_name" }}`,
				InsertText: `{{ snakeToCamel . }}`,
			},
			{
				Name:        "pluralize",
				DisplayName: "单词复数化",
				Description: "将英文单词转换为复数形式（如：user -> users, child -> children）",
				Params: []model.FunctionParam{
					{
						Name:        "word",
						Type:        "string",
						Required:    true,
						Description: "要复数化的单词",
					},
				},
				ReturnType: "string",
				Category:   "代码生成函数",
				Example:    `{{ pluralize "user" }}`,
				InsertText: `{{ pluralize . }}`,
			},
			{
				Name:        "singularize",
				DisplayName: "单词单数化",
				Description: "将英文单词转换为单数形式（如：users -> user, children -> child）",
				Params: []model.FunctionParam{
					{
						Name:        "word",
						Type:        "string",
						Required:    true,
						Description: "要单数化的单词",
					},
				},
				ReturnType: "string",
				Category:   "代码生成函数",
				Example:    `{{ singularize "users" }}`,
				InsertText: `{{ singularize . }}`,
			},
		},
	}
}

// buildProjectStructureFunctions 构建项目结构函数分组
func (s *sBuiltinFunctions) buildProjectStructureFunctions() model.BuiltinFunctionCategory {
	return model.BuiltinFunctionCategory{
		Name:        "项目结构函数",
		Description: "项目结构和路径处理相关函数",
		Functions: []model.BuiltinFunction{
			{
				Name:        "packagePath",
				DisplayName: "生成包路径",
				Description: "根据目录结构生成Go包的导入路径",
				Params: []model.FunctionParam{
					{
						Name:        "moduleName",
						Type:        "string",
						Required:    true,
						Description: "模块名称",
					},
					{
						Name:        "packageDir",
						Type:        "string",
						Required:    true,
						Description: "包目录路径",
					},
				},
				ReturnType: "string",
				Category:   "项目结构函数",
				Example:    `{{ packagePath "github.com/user/project" "internal/service" }}`,
				InsertText: `{{ packagePath "module" "path" }}`,
			},
			{
				Name:        "relativeImport",
				DisplayName: "相对导入路径",
				Description: "生成相对于当前包的导入路径",
				Params: []model.FunctionParam{
					{
						Name:        "fromPath",
						Type:        "string",
						Required:    true,
						Description: "源路径",
					},
					{
						Name:        "toPath",
						Type:        "string",
						Required:    true,
						Description: "目标路径",
					},
				},
				ReturnType: "string",
				Category:   "项目结构函数",
				Example:    `{{ relativeImport "internal/service" "internal/model" }}`,
				InsertText: `{{ relativeImport "from" "to" }}`,
			},
			{
				Name:        "modulePrefix",
				DisplayName: "模块前缀",
				Description: "从go.mod文件获取模块前缀",
				Params: []model.FunctionParam{
					{
						Name:        "goModPath",
						Type:        "string",
						Required:    false,
						Description: "go.mod文件路径，默认为当前目录",
						Default:     "go.mod",
					},
				},
				ReturnType: "string",
				Category:   "项目结构函数",
				Example:    `{{ modulePrefix }}`,
				InsertText: `{{ modulePrefix }}`,
			},
		},
	}
}

// GetTemplateFuncMap 获取模板函数映射（供模板渲染使用）
func (s *sBuiltinFunctions) GetTemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		// 文件操作函数
		"readFile": func(filepath string) string {
			content, err := ioutil.ReadFile(filepath)
			if err != nil {
				return ""
			}
			return string(content)
		},
		"fileExists": func(filepath string) bool {
			_, err := os.Stat(filepath)
			return !os.IsNotExist(err)
		},
		"dirExists": func(dirpath string) bool {
			info, err := os.Stat(dirpath)
			return err == nil && info.IsDir()
		},

		// 字符串处理函数
		"indent": func(spaces int, text string) string {
			if spaces <= 0 {
				return text
			}
			indentStr := strings.Repeat(" ", spaces)
			lines := strings.Split(text, "\n")
			for i, line := range lines {
				if line != "" {
					lines[i] = indentStr + line
				}
			}
			return strings.Join(lines, "\n")
		},
		"comment": func(style, text string) string {
			switch style {
			case "line":
				lines := strings.Split(text, "\n")
				for i, line := range lines {
					lines[i] = "// " + line
				}
				return strings.Join(lines, "\n")
			case "block":
				return "/*\n" + text + "\n*/"
			case "hash":
				lines := strings.Split(text, "\n")
				for i, line := range lines {
					lines[i] = "# " + line
				}
				return strings.Join(lines, "\n")
			default:
				return text
			}
		},
		"quote": func(text string, style ...string) string {
			quoteStyle := "double"
			if len(style) > 0 && style[0] != "" {
				quoteStyle = style[0]
			}
			switch quoteStyle {
			case "single":
				return "'" + strings.ReplaceAll(text, "'", "\\'") + "'"
			case "back":
				return "`" + text + "`"
			default: // double
				return `"` + strings.ReplaceAll(text, `"`, `\"`) + `"`
			}
		},

		// 代码生成函数
		"camelToSnake": func(text string) string {
			re := regexp.MustCompile("([a-z0-9])([A-Z])")
			snake := re.ReplaceAllString(text, "${1}_${2}")
			return strings.ToLower(snake)
		},
		"snakeToCamel": func(text string, firstUpper ...bool) string {
			isFirstUpper := true
			if len(firstUpper) > 0 {
				isFirstUpper = firstUpper[0]
			}

			parts := strings.Split(text, "_")
			var result strings.Builder

			for i, part := range parts {
				if part == "" {
					continue
				}
				if i == 0 && !isFirstUpper {
					result.WriteString(strings.ToLower(part))
				} else {
					result.WriteString(strings.Title(strings.ToLower(part)))
				}
			}
			return result.String()
		},
		"pluralize": func(word string) string {
			// 简单的英文复数规则
			word = strings.ToLower(word)
			if strings.HasSuffix(word, "y") && len(word) > 1 && !strings.ContainsAny(string(word[len(word)-2]), "aeiou") {
				return word[:len(word)-1] + "ies"
			}
			if strings.HasSuffix(word, "s") || strings.HasSuffix(word, "x") || strings.HasSuffix(word, "z") ||
				strings.HasSuffix(word, "ch") || strings.HasSuffix(word, "sh") {
				return word + "es"
			}
			if strings.HasSuffix(word, "f") {
				return word[:len(word)-1] + "ves"
			}
			if strings.HasSuffix(word, "fe") {
				return word[:len(word)-2] + "ves"
			}
			// 不规则复数
			irregulars := map[string]string{
				"child": "children", "person": "people", "man": "men", "woman": "women",
				"foot": "feet", "tooth": "teeth", "mouse": "mice", "goose": "geese",
			}
			if plural, exists := irregulars[word]; exists {
				return plural
			}
			return word + "s"
		},
		"singularize": func(word string) string {
			// 简单的英文单数规则
			word = strings.ToLower(word)
			// 不规则单数
			irregulars := map[string]string{
				"children": "child", "people": "person", "men": "man", "women": "woman",
				"feet": "foot", "teeth": "tooth", "mice": "mouse", "geese": "goose",
			}
			if singular, exists := irregulars[word]; exists {
				return singular
			}
			if strings.HasSuffix(word, "ies") {
				return word[:len(word)-3] + "y"
			}
			if strings.HasSuffix(word, "ves") {
				if strings.HasSuffix(word[:len(word)-3], "l") || strings.HasSuffix(word[:len(word)-3], "f") {
					return word[:len(word)-3] + "f"
				}
				return word[:len(word)-3] + "fe"
			}
			if strings.HasSuffix(word, "ses") || strings.HasSuffix(word, "xes") || strings.HasSuffix(word, "zes") ||
				strings.HasSuffix(word, "ches") || strings.HasSuffix(word, "shes") {
				return word[:len(word)-2]
			}
			if strings.HasSuffix(word, "s") && len(word) > 1 {
				return word[:len(word)-1]
			}
			return word
		},

		// 项目结构函数
		"packagePath": func(moduleName, packageDir string) string {
			return filepath.Join(moduleName, packageDir)
		},
		"relativeImport": func(fromPath, toPath string) string {
			rel, err := filepath.Rel(fromPath, toPath)
			if err != nil {
				return toPath
			}
			return filepath.ToSlash(rel)
		},
		"modulePrefix": func(goModPath ...string) string {
			modPath := "go.mod"
			if len(goModPath) > 0 && goModPath[0] != "" {
				modPath = goModPath[0]
			}

			content, err := ioutil.ReadFile(modPath)
			if err != nil {
				return ""
			}

			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "module ") {
					return strings.TrimSpace(strings.TrimPrefix(line, "module"))
				}
			}
			return ""
		},
	}
}
