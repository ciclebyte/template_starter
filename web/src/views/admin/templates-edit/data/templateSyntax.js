/**
 * Go模板语法分类数据
 * 包含条件语句、循环语句、变量操作、输出控制、模板控制等分类
 */
export const templateSyntaxCategories = [
  {
    name: '条件语句',
    syntaxes: [
      {
        name: 'if',
        display_name: 'if 条件',
        description: '条件判断语句',
        syntax: '{{if .condition}}...{{end}}',
        usage: '当条件为真时执行内容。条件可以是变量、比较表达式或函数调用。',
        example: '{{if .isEnabled}}启用状态{{end}}',
        insertText: '{{if .condition}}\n  content\n{{end}}',
        params: [
          { name: 'condition', type: 'bool', required: true, description: '条件表达式' }
        ]
      },
      {
        name: 'if-else',
        display_name: 'if-else 条件',
        description: '条件判断语句（带else分支）',
        syntax: '{{if .condition}}...{{else}}...{{end}}',
        usage: '当条件为真时执行第一个分支，否则执行else分支。',
        example: '{{if .isEnabled}}启用{{else}}禁用{{end}}',
        insertText: '{{if .condition}}\n  true branch\n{{else}}\n  false branch\n{{end}}',
        params: [
          { name: 'condition', type: 'bool', required: true, description: '条件表达式' }
        ]
      },
      {
        name: 'if-else-if',
        display_name: 'if-else if 条件',
        description: '多重条件判断语句',
        syntax: '{{if .condition1}}...{{else if .condition2}}...{{else}}...{{end}}',
        usage: '按顺序检查多个条件，执行第一个为真的分支。',
        example: '{{if eq .status "active"}}活跃{{else if eq .status "inactive"}}非活跃{{else}}未知{{end}}',
        insertText: '{{if .condition1}}\n  branch1\n{{else if .condition2}}\n  branch2\n{{else}}\n  default\n{{end}}',
        params: [
          { name: 'condition1', type: 'bool', required: true, description: '第一个条件' },
          { name: 'condition2', type: 'bool', required: true, description: '第二个条件' }
        ]
      }
    ]
  },
  {
    name: '循环语句',
    syntaxes: [
      {
        name: 'range',
        display_name: 'range 循环',
        description: '遍历数组、切片或映射',
        syntax: '{{range .items}}...{{end}}',
        usage: '遍历集合中的每个元素。在循环体内，. 代表当前元素。',
        example: '{{range .users}}<p>{{.name}}</p>{{end}}',
        insertText: '{{range .items}}\n  {{.}}\n{{end}}',
        params: [
          { name: 'items', type: 'array|slice|map', required: true, description: '要遍历的集合' }
        ]
      },
      {
        name: 'range-index',
        display_name: 'range 带索引',
        description: '遍历时获取索引和值',
        syntax: '{{range $index, $element := .items}}...{{end}}',
        usage: '遍历集合并获取索引（或键）和对应的值。',
        example: '{{range $i, $user := .users}}<p>{{$i}}: {{$user.name}}</p>{{end}}',
        insertText: '{{range $index, $element := .items}}\n  {{$index}}: {{$element}}\n{{end}}',
        params: [
          { name: 'items', type: 'array|slice|map', required: true, description: '要遍历的集合' }
        ]
      },
      {
        name: 'range-empty',
        display_name: 'range 带空检查',
        description: '遍历集合，支持空值处理',
        syntax: '{{range .items}}...{{else}}...{{end}}',
        usage: '当集合为空时，执行else分支。',
        example: '{{range .users}}<p>{{.name}}</p>{{else}}<p>没有用户</p>{{end}}',
        insertText: '{{range .items}}\n  {{.}}\n{{else}}\n  empty content\n{{end}}',
        params: [
          { name: 'items', type: 'array|slice|map', required: true, description: '要遍历的集合' }
        ]
      }
    ]
  },
  {
    name: '变量操作',
    syntaxes: [
      {
        name: 'with',
        display_name: 'with 作用域',
        description: '设置新的上下文作用域',
        syntax: '{{with .value}}...{{end}}',
        usage: '在with块内，. 代表指定的值。如果值为空，则不执行块内容。',
        example: '{{with .user}}<p>Hello {{.name}}</p>{{end}}',
        insertText: '{{with .value}}\n  {{.}}\n{{end}}',
        params: [
          { name: 'value', type: 'any', required: true, description: '新的上下文值' }
        ]
      },
      {
        name: 'with-else',
        display_name: 'with 带else',
        description: '设置作用域，支持空值处理',
        syntax: '{{with .value}}...{{else}}...{{end}}',
        usage: '当value不为空时执行第一个分支，否则执行else分支。',
        example: '{{with .user}}Hello {{.name}}{{else}}No user{{end}}',
        insertText: '{{with .value}}\n  {{.}}\n{{else}}\n  default content\n{{end}}',
        params: [
          { name: 'value', type: 'any', required: true, description: '上下文值' }
        ]
      },
      {
        name: 'variable',
        display_name: '变量赋值',
        description: '定义和使用变量',
        syntax: '{{$var := .value}}',
        usage: '将值赋给变量，变量名以$开头。变量在整个模板中有效。',
        example: '{{$name := .user.name}}Hello {{$name}}',
        insertText: '{{$var := .value}}',
        params: [
          { name: 'var', type: 'string', required: true, description: '变量名（以$开头）' },
          { name: 'value', type: 'any', required: true, description: '变量值' }
        ]
      }
    ]
  },
  {
    name: '输出控制',
    syntaxes: [
      {
        name: 'printf',
        display_name: 'printf 格式化',
        description: '格式化输出',
        syntax: '{{printf "%s: %d" .name .count}}',
        usage: '使用格式化字符串输出内容，类似C语言的printf函数。',
        example: '{{printf "用户: %s, 年龄: %d" .name .age}}',
        insertText: '{{printf "%s" .value}}',
        params: [
          { name: 'format', type: 'string', required: true, description: '格式化字符串' },
          { name: 'args', type: 'any...', required: false, description: '格式化参数' }
        ]
      },
      {
        name: 'print',
        display_name: 'print 输出',
        description: '简单输出（空格分隔）',
        syntax: '{{print .value1 .value2}}',
        usage: '输出多个值，用空格分隔。',
        example: '{{print "Hello" .name "!"}}',
        insertText: '{{print .value}}',
        params: [
          { name: 'values', type: 'any...', required: true, description: '要输出的值' }
        ]
      },
      {
        name: 'println',
        display_name: 'println 输出',
        description: '输出并换行',
        syntax: '{{println .value}}',
        usage: '输出值并在末尾添加换行符。',
        example: '{{println "Line 1"}}{{println "Line 2"}}',
        insertText: '{{println .value}}',
        params: [
          { name: 'values', type: 'any...', required: true, description: '要输出的值' }
        ]
      }
    ]
  },
  {
    name: '模板控制',
    syntaxes: [
      {
        name: 'template',
        display_name: 'template 引用',
        description: '引用其他模板',
        syntax: '{{template "name" .}}',
        usage: '调用已定义的模板，传递当前上下文或指定数据。',
        example: '{{template "header" .}}{{template "content" .user}}',
        insertText: '{{template "templateName" .}}',
        params: [
          { name: 'name', type: 'string', required: true, description: '模板名称' },
          { name: 'data', type: 'any', required: false, description: '传递给模板的数据' }
        ]
      },
      {
        name: 'define',
        display_name: 'define 定义',
        description: '定义可重用的模板块',
        syntax: '{{define "name"}}...{{end}}',
        usage: '定义一个命名的模板块，可以被template调用。',
        example: '{{define "header"}}<h1>{{.title}}</h1>{{end}}',
        insertText: '{{define "templateName"}}\n  content\n{{end}}',
        params: [
          { name: 'name', type: 'string', required: true, description: '模板名称' }
        ]
      },
      {
        name: 'block',
        display_name: 'block 块定义',
        description: '定义可覆盖的默认块',
        syntax: '{{block "name" .}}...{{end}}',
        usage: '定义一个可以被子模板覆盖的默认内容块。',
        example: '{{block "content" .}}<p>默认内容</p>{{end}}',
        insertText: '{{block "blockName" .}}\n  default content\n{{end}}',
        params: [
          { name: 'name', type: 'string', required: true, description: '块名称' },
          { name: 'data', type: 'any', required: false, description: '传递给块的数据' }
        ]
      }
    ]
  }
]