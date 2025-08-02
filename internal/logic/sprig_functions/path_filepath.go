package sprig_functions

import "github.com/ciclebyte/template_starter/internal/model"

// buildPathFilepathFunctions 构建路径和文件路径函数分类
func (s *sSprigFunctions) buildPathFilepathFunctions() model.SprigFunctionCategory {
	return model.SprigFunctionCategory{
		Name:        "路径和文件路径函数",
		Description: "路径字符串处理和文件路径操作函数",
		Functions: []model.SprigFunction{
			{
				Name:        "base",
				DisplayName: "获取路径基名",
				Description: "返回路径的最后一个元素",
				Params: []model.SprigFunctionParam{
					{Name: "path", Type: "string", Required: true, Description: "路径字符串"},
				},
				ReturnType: "string",
				Category:   "路径和文件路径函数",
				Example:    `{{ base "foo/bar/baz" }}`,
				Usage:      "base 函数返回路径的最后一个元素。例如 base \"foo/bar/baz\" 返回 \"baz\"。使用斜杠(/)分隔的路径，由path包处理。",
				InsertText: `{{ base . }}`,
				Note:       "适用于Linux/MacOS文件系统路径和URI路径组件",
			},
			{
				Name:        "dir",
				DisplayName: "获取目录路径",
				Description: "返回去除最后部分后的目录路径",
				Params: []model.SprigFunctionParam{
					{Name: "path", Type: "string", Required: true, Description: "路径字符串"},
				},
				ReturnType: "string",
				Category:   "路径和文件路径函数",
				Example:    `{{ dir "foo/bar/baz" }}`,
				Usage:      "dir 函数返回目录路径，去除路径的最后部分。例如 dir \"foo/bar/baz\" 返回 \"foo/bar\"。",
				InsertText: `{{ dir . }}`,
			},
			{
				Name:        "clean",
				DisplayName: "清理路径",
				Description: "清理路径，解析相对路径元素",
				Params: []model.SprigFunctionParam{
					{Name: "path", Type: "string", Required: true, Description: "要清理的路径字符串"},
				},
				ReturnType: "string",
				Category:   "路径和文件路径函数",
				Example:    `{{ clean "foo/bar/../baz" }}`,
				Usage:      "clean 函数清理路径。例如 clean \"foo/bar/../baz\" 解析 \"..\" 并返回 \"foo/baz\"。",
				InsertText: `{{ clean . }}`,
				Note:       "解析 .. 和 . 等相对路径元素",
			},
			{
				Name:        "ext",
				DisplayName: "获取文件扩展名",
				Description: "返回文件的扩展名",
				Params: []model.SprigFunctionParam{
					{Name: "path", Type: "string", Required: true, Description: "文件路径"},
				},
				ReturnType: "string",
				Category:   "路径和文件路径函数",
				Example:    `{{ ext "foo.bar" }}`,
				Usage:      "ext 函数返回文件扩展名。例如 ext \"foo.bar\" 返回 \".bar\"。",
				InsertText: `{{ ext . }}`,
				Note:       "包含点号的完整扩展名",
			},
			{
				Name:        "isAbs",
				DisplayName: "检查绝对路径",
				Description: "检查路径是否为绝对路径",
				Params: []model.SprigFunctionParam{
					{Name: "path", Type: "string", Required: true, Description: "要检查的路径"},
				},
				ReturnType: "bool",
				Category:   "路径和文件路径函数",
				Example:    `{{ isAbs "/foo/bar" }}`,
				Usage:      "isAbs 函数检查路径是否为绝对路径。如果是绝对路径返回 true，否则返回 false。",
				InsertText: `{{ isAbs . }}`,
			},
			{
				Name:        "osBase",
				DisplayName: "获取系统路径基名",
				Description: "返回文件路径的最后一个元素（系统相关）",
				Params: []model.SprigFunctionParam{
					{Name: "filepath", Type: "string", Required: true, Description: "文件路径字符串"},
				},
				ReturnType: "string",
				Category:   "路径和文件路径函数",
				Example:    `{{ osBase "/foo/bar/baz" }}`,
				Usage:      "osBase 函数返回文件路径的最后一个元素。由path/filepath包处理，使用os.PathSeparator分隔符。在Linux上打印 \"baz\"，在Windows上同样打印 \"baz\"。",
				InsertText: `{{ osBase . }}`,
				Note:       "系统相关路径分隔符，推荐用于本地文件系统路径",
			},
			{
				Name:        "osDir",
				DisplayName: "获取系统目录路径",
				Description: "返回去除最后部分后的目录路径（系统相关）",
				Params: []model.SprigFunctionParam{
					{Name: "filepath", Type: "string", Required: true, Description: "文件路径字符串"},
				},
				ReturnType: "string",
				Category:   "路径和文件路径函数",
				Example:    `{{ osDir "/foo/bar/baz" }}`,
				Usage:      "osDir 函数返回目录路径，去除路径的最后部分。osDir \"/foo/bar/baz\" 在Linux上返回 \"/foo/bar\"，osDir \"C:\\\\foo\\\\bar\\\\baz\" 在Windows上返回 \"C:\\\\foo\\\\bar\"。",
				InsertText: `{{ osDir . }}`,
				Note:       "使用系统特定的路径分隔符",
			},
			{
				Name:        "osClean",
				DisplayName: "清理系统路径",
				Description: "清理文件路径，解析相对路径元素（系统相关）",
				Params: []model.SprigFunctionParam{
					{Name: "filepath", Type: "string", Required: true, Description: "要清理的文件路径"},
				},
				ReturnType: "string",
				Category:   "路径和文件路径函数",
				Example:    `{{ osClean "foo/bar/../baz" }}`,
				Usage:      "osClean 函数清理文件路径。解析 \"..\" 并在Linux上返回 \"foo/baz\"，在Windows上返回 \"C:\\\\foo\\\\baz\"。",
				InsertText: `{{ osClean . }}`,
				Note:       "系统相关的路径清理，解析相对路径元素",
			},
			{
				Name:        "osExt",
				DisplayName: "获取系统文件扩展名",
				Description: "返回文件的扩展名（系统相关）",
				Params: []model.SprigFunctionParam{
					{Name: "filepath", Type: "string", Required: true, Description: "文件路径"},
				},
				ReturnType: "string",
				Category:   "路径和文件路径函数",
				Example:    `{{ osExt "foo.bar" }}`,
				Usage:      "osExt 函数返回文件扩展名。在Linux和Windows上都返回 \".bar\"。",
				InsertText: `{{ osExt . }}`,
				Note:       "在各系统上行为一致，返回包含点号的扩展名",
			},
			{
				Name:        "osIsAbs",
				DisplayName: "检查系统绝对路径",
				Description: "检查文件路径是否为绝对路径（系统相关）",
				Params: []model.SprigFunctionParam{
					{Name: "filepath", Type: "string", Required: true, Description: "要检查的文件路径"},
				},
				ReturnType: "bool",
				Category:   "路径和文件路径函数",
				Example:    `{{ osIsAbs "/foo/bar" }}`,
				Usage:      "osIsAbs 函数检查文件路径是否为绝对路径。使用系统特定的规则判断绝对路径。",
				InsertText: `{{ osIsAbs . }}`,
				Note:       "使用系统特定规则，Linux用/开头，Windows用盘符开头",
			},
		},
	}
}