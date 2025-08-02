package sprig_functions

import "github.com/ciclebyte/template_starter/internal/model"

// buildURLFunctions 构建URL函数分类
func (s *sSprigFunctions) buildURLFunctions() model.SprigFunctionCategory {
	return model.SprigFunctionCategory{
		Name:        "URL函数",
		Description: "URL解析和构建函数",
		Functions: []model.SprigFunction{
			{
				Name:        "urlParse",
				DisplayName: "解析URL",
				Description: "解析URL字符串并生成包含URL各部分的字典",
				Params: []model.SprigFunctionParam{
					{Name: "url", Type: "string", Required: true, Description: "要解析的URL字符串"},
				},
				ReturnType: "map[string]interface{}",
				Category:   "URL函数",
				Example:    `{{ urlParse "http://admin:secret@server.com:8080/api?list=false#anchor" }}`,
				Usage:      "urlParse 函数解析URL字符串并生成包含URL各部分的字典。返回的字典包含scheme、host、path、query、opaque、fragment、userinfo等字段。",
				InsertText: `{{ $url := urlParse . }}{{ $url.host }}`,
				Note:       "返回字典包含所有URL组件，详见Go的net/url包文档",
			},
			{
				Name:        "urlJoin",
				DisplayName: "构建URL",
				Description: "将字典（由urlParse生成）连接生成URL字符串",
				Params: []model.SprigFunctionParam{
					{Name: "urlDict", Type: "map[string]interface{}", Required: true, Description: "包含URL各部分的字典"},
				},
				ReturnType: "string",
				Category:   "URL函数",
				Example:    `{{ urlJoin (dict "fragment" "fragment" "host" "host:80" "path" "/path" "query" "query" "scheme" "http") }}`,
				Usage:      "urlJoin 函数将字典（由urlParse生成）连接生成URL字符串。接受包含scheme、host、path、query、fragment等字段的字典。",
				InsertText: `{{ urlJoin . }}`,
				Note:       "与urlParse配对使用，用于URL的重新构建",
			},
		},
	}
}