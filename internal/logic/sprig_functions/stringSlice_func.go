package sprig_functions

import "github.com/ciclebyte/template_starter/internal/model"

// buildStringSliceFunctions 构建字符串切片函数分类
func (s *sSprigFunctions) buildStringSliceFunctions() model.SprigFunctionCategory {
	return model.SprigFunctionCategory{
		Name:        "字符串切片函数",
		Description: "字符串切片操作和处理函数",
		Functions: []model.SprigFunction{
			{
				Name:        "join",
				DisplayName: "连接字符串",
				Description: "将字符串列表连接成单个字符串",
				Params: []model.SprigFunctionParam{
					{Name: "separator", Type: "string", Required: true, Description: "连接分隔符"},
					{Name: "list", Type: "[]string", Required: true, Description: "字符串列表"},
				},
				ReturnType: "string",
				Category:   "字符串切片函数",
				Example:    `{{ list "hello" "world" | join "_" }}`,
				Usage:      "join 函数将字符串列表用指定分隔符连接成单个字符串。会尝试将非字符串类型转换为字符串值。",
				InsertText: `{{ . | join "_" }}`,
				Note:       "自动转换非字符串类型为字符串",
			},
			{
				Name:        "splitList",
				DisplayName: "分割字符串为列表",
				Description: "将字符串分割为字符串列表",
				Params: []model.SprigFunctionParam{
					{Name: "separator", Type: "string", Required: true, Description: "分割分隔符"},
					{Name: "string", Type: "string", Required: true, Description: "要分割的字符串"},
				},
				ReturnType: "[]string",
				Category:   "字符串切片函数",
				Example:    `{{ splitList "$" "foo$bar$baz" }}`,
				Usage:      "splitList 函数将字符串按指定分隔符分割为字符串列表。返回 [foo bar baz]。",
				InsertText: `{{ splitList "$" . }}`,
			},
			{
				Name:        "split",
				DisplayName: "分割字符串为字典",
				Description: "将字符串分割为带索引键的字典",
				Params: []model.SprigFunctionParam{
					{Name: "separator", Type: "string", Required: true, Description: "分割分隔符"},
					{Name: "string", Type: "string", Required: true, Description: "要分割的字符串"},
				},
				ReturnType: "map[string]string",
				Category:   "字符串切片函数",
				Example:    `{{ $a := split "$" "foo$bar$baz" }}{{ $a._0 }}`,
				Usage:      "split 函数将字符串分割为字典，使用索引键(_0, _1, _2...)。方便使用模板点号访问成员。例如产生 {_0: foo, _1: bar, _2: baz}。",
				InsertText: `{{ $a := split "$" . }}{{ $a._0 }}`,
				Note:       "使用索引键 _0, _1, _2... 访问分割后的元素",
			},
			{
				Name:        "splitn",
				DisplayName: "限制分割字符串",
				Description: "将字符串分割为字典，限制分割次数",
				Params: []model.SprigFunctionParam{
					{Name: "separator", Type: "string", Required: true, Description: "分割分隔符"},
					{Name: "n", Type: "int", Required: true, Description: "最大分割次数"},
					{Name: "string", Type: "string", Required: true, Description: "要分割的字符串"},
				},
				ReturnType: "map[string]string",
				Category:   "字符串切片函数",
				Example:    `{{ splitn "$" 2 "foo$bar$baz" }}`,
				Usage:      "splitn 函数将字符串分割为字典，但限制分割次数。例如最多分割2次，产生 {_0: foo, _1: bar$baz}。",
				InsertText: `{{ splitn "$" 2 . }}`,
				Note:       "限制分割次数，剩余部分保留在最后一个元素中",
			},
			{
				Name:        "sortAlpha",
				DisplayName: "字母排序",
				Description: "对字符串列表进行字母顺序排序",
				Params: []model.SprigFunctionParam{
					{Name: "list", Type: "[]string", Required: true, Description: "要排序的字符串列表"},
				},
				ReturnType: "[]string",
				Category:   "字符串切片函数",
				Example:    `{{ list "zebra" "apple" "bear" | sortAlpha }}`,
				Usage:      "sortAlpha 函数对字符串列表进行字母顺序（词典序）排序。不会就地排序，而是返回排序后的副本，保持列表的不可变性。",
				InsertText: `{{ . | sortAlpha }}`,
				Note:       "返回排序后的副本，不修改原列表",
			},
		},
	}
}
