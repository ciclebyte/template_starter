package sprig_functions

import "github.com/ciclebyte/template_starter/internal/model"

// buildReflectionFunctions 构建反射函数分类
func (s *sSprigFunctions) buildReflectionFunctions() model.SprigFunctionCategory {
	return model.SprigFunctionCategory{
		Name:        "反射函数",
		Description: "类型反射和检查函数",
		Functions: []model.SprigFunction{
			{
				Name:        "kindOf",
				DisplayName: "获取类型种类",
				Description: "返回对象的Go类型种类",
				Params: []model.SprigFunctionParam{
					{Name: "value", Type: "any", Required: true, Description: "要检查的值"},
				},
				ReturnType: "string",
				Category:   "反射函数",
				Example:    `{{ kindOf "hello" }}`,
				Usage:      "kindOf 函数返回对象的种类（kind）。例如 kindOf \"hello\" 返回 \"string\"。Go有几种原始种类，如string、slice、int64、bool等。",
				InsertText: `{{ kindOf . }}`,
				Note:       "返回Go的基本类型种类名称",
			},
			{
				Name:        "kindIs",
				DisplayName: "检查类型种类",
				Description: "验证值是否为特定种类",
				Params: []model.SprigFunctionParam{
					{Name: "kind", Type: "string", Required: true, Description: "要检查的种类名"},
					{Name: "value", Type: "any", Required: true, Description: "要检查的值"},
				},
				ReturnType: "bool",
				Category:   "反射函数",
				Example:    `{{ kindIs "string" "hello" }}`,
				Usage:      "kindIs 函数验证值是否为特定种类。例如 kindIs \"string\" \"hello\" 返回 true。适用于在if块中进行简单类型测试。",
				InsertText: `{{ kindIs "string" . }}`,
				Note:       "常用于模板条件判断中的类型检查",
			},
			{
				Name:        "typeOf",
				DisplayName: "获取类型",
				Description: "返回值的底层类型",
				Params: []model.SprigFunctionParam{
					{Name: "value", Type: "any", Required: true, Description: "要检查的值"},
				},
				ReturnType: "string",
				Category:   "反射函数",
				Example:    `{{ typeOf $foo }}`,
				Usage:      "typeOf 函数返回值的底层类型。Go有开放的类型系统，允许开发者创建自己的类型。typeOf可以获取这些自定义类型的信息。",
				InsertText: `{{ typeOf . }}`,
				Note:       "返回完整的类型信息，包括包名",
			},
			{
				Name:        "typeIs",
				DisplayName: "检查类型",
				Description: "验证值是否为特定类型",
				Params: []model.SprigFunctionParam{
					{Name: "type", Type: "string", Required: true, Description: "要检查的类型名"},
					{Name: "value", Type: "any", Required: true, Description: "要检查的值"},
				},
				ReturnType: "bool",
				Category:   "反射函数",
				Example:    `{{ typeIs "*io.Buffer" $myVal }}`,
				Usage:      "typeIs 函数类似于kindIs，但用于类型检查。例如 typeIs \"*io.Buffer\" $myVal 检查值是否为指定类型。",
				InsertText: `{{ typeIs "string" . }}`,
				Note:       "类似kindIs但针对具体类型而非种类",
			},
			{
				Name:        "typeIsLike",
				DisplayName: "类似类型检查",
				Description: "验证值是否为类似类型（会解引用指针）",
				Params: []model.SprigFunctionParam{
					{Name: "type", Type: "string", Required: true, Description: "要检查的类型名"},
					{Name: "value", Type: "any", Required: true, Description: "要检查的值"},
				},
				ReturnType: "bool",
				Category:   "反射函数",
				Example:    `{{ typeIsLike "io.Buffer" $myVal }}`,
				Usage:      "typeIsLike 函数的工作方式与typeIs相同，但它还会解引用指针。这在处理指针类型时很有用。",
				InsertText: `{{ typeIsLike "string" . }}`,
				Note:       "会自动解引用指针进行类型比较",
			},
			{
				Name:        "deepEqual",
				DisplayName: "深度相等比较",
				Description: "检查两个值是否深度相等",
				Params: []model.SprigFunctionParam{
					{Name: "value1", Type: "any", Required: true, Description: "第一个值"},
					{Name: "value2", Type: "any", Required: true, Description: "第二个值"},
				},
				ReturnType: "bool",
				Category:   "反射函数",
				Example:    `{{ deepEqual $obj1 $obj2 }}`,
				Usage:      "deepEqual 函数如果两个值\"深度相等\"则返回true。适用于非原始类型的比较（相比内置的eq函数）。使用Go的reflect.DeepEqual进行比较。",
				InsertText: `{{ deepEqual . $other }}`,
				Note:       "支持复杂类型的深度比较，比内置eq更强大",
			},
		},
	}
}