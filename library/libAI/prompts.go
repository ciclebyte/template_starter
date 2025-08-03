package libAI

import (
	"fmt"
	"strings"
)

// PromptTemplates 提示词模板集合
type PromptTemplates struct {
	// 模板生成相关
	TemplateGeneration map[string]string
	// 变量建议相关
	VariableSuggestion map[string]string
	// 代码生成相关
	CodeGeneration map[string]string
	// 文档生成相关
	Documentation map[string]string
}

// GetOptimizedPrompts 获取优化的提示词模板
func GetOptimizedPrompts() *PromptTemplates {
	return &PromptTemplates{
		TemplateGeneration: getTemplateGenerationPrompts(),
		VariableSuggestion: getVariableSuggestionPrompts(),
		CodeGeneration:     getCodeGenerationPrompts(),
		Documentation:      getDocumentationPrompts(),
	}
}

// getTemplateGenerationPrompts 获取模板生成提示词
func getTemplateGenerationPrompts() map[string]string {
	defaultPrompt := "你是一个资深的软件架构师和全栈开发工程师，专门负责生成高质量的项目模板。\n\n" +
		"## 任务描述\n根据用户需求生成一个完整、实用的项目模板，包含合理的文件结构、代码内容和配置。\n\n" +
		"## 输入信息\n- **项目描述**: {{.Description}}\n- **项目类型**: {{.ProjectType}}\n- **技术栈**: {{.TechStack}}\n- **功能特性**: {{.Features}}\n" +
		"{{if .RetrievedContext}}\n- **参考资料**:\n{{.RetrievedContext}}\n{{end}}\n\n" +
		"## 输出要求\n请严格按照以下JSON格式返回结果：\n\n" +
		"```json\n{\n  \"projectStructure\": [\n    {\n      \"path\": \"文件或目录路径\",\n      \"content\": \"文件内容(目录为空字符串)\",\n      \"isDirectory\": false,\n      \"description\": \"文件用途说明\"\n    }\n  ],\n  \"variables\": [\n    {\n      \"name\": \"变量名\",\n      \"type\": \"string|number|boolean|select|textarea\",\n      \"description\": \"变量描述\",\n      \"defaultValue\": \"默认值\",\n      \"required\": true,\n      \"options\": [\"选项1\", \"选项2\"]\n    }\n  ],\n  \"instructions\": \"使用说明和注意事项\",\n  \"estimatedTime\": 15\n}\n```\n\n" +
		"## 质量要求\n1. **文件结构**：完整、符合最佳实践，包含必要的配置文件\n2. **代码质量**：包含注释、错误处理、类型定义\n3. **变量设计**：实用、有意义，帮助用户快速定制\n4. **兼容性**：考虑跨平台兼容性和依赖管理\n5. **安全性**：遵循安全编码规范，不包含敏感信息\n\n" +
		"## 特殊说明\n- 使用双大括号语法作为模板变量：{{.VariableName}}\n- 文件内容应该是可直接使用的代码\n- 预估时间以分钟为单位\n- 确保返回的是有效的JSON格式"

	webPrompt := "你是一个Web开发专家，擅长前后端开发和现代Web技术栈。\n\n" +
		"## 项目信息\n- **描述**: {{.Description}}\n- **类型**: {{.ProjectType}}\n- **技术栈**: {{.TechStack}}\n- **特性**: {{.Features}}\n\n" +
		"## Web项目特殊要求\n1. **响应式设计**：支持移动端适配\n2. **性能优化**：代码分割、懒加载、缓存策略\n3. **SEO友好**：合理的HTML结构和meta标签\n4. **无障碍访问**：符合WCAG标准\n5. **现代化工具**：包含构建工具、代码检查、测试框架\n\n" +
		"请生成包含以下结构的Web项目：\n- 入口文件和配置\n- 组件目录结构\n- 样式文件组织\n- 资源文件管理\n- 构建和部署配置\n\n返回标准JSON格式。"

	return map[string]string{
		"default": defaultPrompt,
		"web":     webPrompt,
		"backend": "你是一个后端开发专家，精通服务器端架构和API设计。根据{{.Description}}、{{.ProjectType}}、{{.TechStack}}生成后端项目模板。返回JSON格式。",
		"mobile":  "你是一个移动应用开发专家，熟悉原生和跨平台开发。根据{{.Description}}、{{.ProjectType}}、{{.TechStack}}生成移动应用项目模板。返回JSON格式。",
		"desktop": "你是一个桌面应用开发专家，精通Electron、Tauri等跨平台框架。根据{{.Description}}、{{.ProjectType}}、{{.TechStack}}生成桌面应用项目模板。返回JSON格式。",
	}
}

// getVariableSuggestionPrompts 获取变量建议提示词
func getVariableSuggestionPrompts() map[string]string {
	defaultPrompt := "你是一个模板变量设计专家，善于分析项目需求并提供实用的模板变量建议。\n\n" +
		"## 分析目标\n为以下项目推荐合适的模板变量，让用户能够快速定制项目。\n\n" +
		"## 项目信息\n- **项目类型**: {{.ProjectType}}\n- **技术栈**: {{.TechStack}}\n- **项目描述**: {{.Description}}\n\n" +
		"## 变量设计原则\n1. **实用性**：变量应该是用户经常需要修改的内容\n2. **清晰性**：变量名称和描述要清楚明确\n3. **合理默认值**：提供有意义的默认值\n4. **类型适配**：选择合适的变量类型和选项\n\n" +
		"## 输出格式\n```json\n{\n  \"variables\": [\n    {\n      \"name\": \"变量名(驼峰命名)\",\n      \"type\": \"string|number|boolean|select|textarea\",\n      \"description\": \"变量用途的详细描述\",\n      \"defaultValue\": \"合理的默认值\",\n      \"required\": true,\n      \"options\": [\"选项1\", \"选项2\"]\n    }\n  ]\n}\n```\n\n请根据项目特点推荐5-15个实用变量。"

	return map[string]string{
		"default": defaultPrompt,
		"web":     "为Web项目推荐模板变量，重点关注应用基础信息、UI配置、功能模块、部署配置、开发工具。项目信息：类型{{.ProjectType}}，技术栈{{.TechStack}}，描述{{.Description}}。返回JSON格式的变量列表。",
		"backend": "为后端项目推荐模板变量，重点关注服务基础信息、数据库配置、API配置、中间件选择、部署配置。项目信息：类型{{.ProjectType}}，技术栈{{.TechStack}}，描述{{.Description}}。返回JSON格式的变量列表。",
	}
}

// getCodeGenerationPrompts 获取代码生成提示词
func getCodeGenerationPrompts() map[string]string {
	return map[string]string{
		"function": "生成一个高质量的{{.Language}}函数：\n\n## 需求描述\n{{.Description}}\n\n请生成完整的函数实现、注释和测试代码。",
		"class":    "设计一个{{.Language}}类：\n\n## 类描述\n{{.Description}}\n\n包含构造函数、主要方法、属性定义和使用示例。",
		"api":      "设计RESTful API接口：\n\n## API描述\n{{.Description}}\n\n生成API定义、实现代码和文档。",
	}
}

// getDocumentationPrompts 获取文档生成提示词
func getDocumentationPrompts() map[string]string {
	return map[string]string{
		"readme":     "为项目生成专业的README.md文档。项目名称：{{.ProjectName}}，描述：{{.Description}}，技术栈：{{.TechStack}}。生成结构清晰、内容完整的README文档。",
		"api_doc":    "生成API文档。API信息：{{.Description}}。生成详细的API参考文档。",
		"user_guide": "编写用户使用手册。产品信息：{{.Description}}。生成易于理解的用户手册。",
	}
}

// PromptBuilder 提示词构建器
type PromptBuilder struct {
	template  string
	variables map[string]interface{}
}

// NewPromptBuilder 创建提示词构建器
func NewPromptBuilder(template string) *PromptBuilder {
	return &PromptBuilder{
		template:  template,
		variables: make(map[string]interface{}),
	}
}

// SetVariable 设置变量
func (pb *PromptBuilder) SetVariable(key string, value interface{}) *PromptBuilder {
	pb.variables[key] = value
	return pb
}

// SetVariables 批量设置变量
func (pb *PromptBuilder) SetVariables(vars map[string]interface{}) *PromptBuilder {
	for k, v := range vars {
		pb.variables[k] = v
	}
	return pb
}

// Build 构建最终提示词
func (pb *PromptBuilder) Build() string {
	result := pb.template
	
	// 简单的模板变量替换
	for key, value := range pb.variables {
		placeholder := fmt.Sprintf("{{.%s}}", key)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
	}
	
	return result
}

// OptimizePromptForModel 根据模型优化提示词
func OptimizePromptForModel(prompt string, model string) string {
	switch {
	case strings.Contains(model, "gpt-4"):
		return addAdvancedInstructions(prompt)
	case strings.Contains(model, "gpt-3.5"):
		return addExplicitInstructions(prompt)
	case strings.Contains(model, "claude"):
		return addStructuredFormat(prompt)
	default:
		return prompt
	}
}

// addAdvancedInstructions 添加高级指令
func addAdvancedInstructions(prompt string) string {
	advanced := "\n\n## 高级要求\n- 考虑代码的可维护性和扩展性\n- 包含性能优化建议\n- 提供多种实现方案对比\n- 考虑安全性和最佳实践\n"
	return prompt + advanced
}

// addExplicitInstructions 添加明确指令
func addExplicitInstructions(prompt string) string {
	explicit := "\n\n## 明确要求\n- 严格按照指定格式输出\n- 不要包含额外的解释文字\n- 确保JSON格式正确无误\n- 每个字段都必须填写\n"
	return prompt + explicit
}

// addStructuredFormat 添加结构化格式
func addStructuredFormat(prompt string) string {
	structured := "\n\n## 输出结构\n请按照以下步骤思考和输出：\n1. 分析项目需求\n2. 设计文件结构\n3. 确定关键变量\n4. 生成完整JSON\n"
	return prompt + structured
}

// ValidatePrompt 验证提示词质量
func ValidatePrompt(prompt string) []string {
	var issues []string
	
	if len(prompt) < 100 {
		issues = append(issues, "提示词过短，可能缺少必要信息")
	}
	
	if len(prompt) > 5000 {
		issues = append(issues, "提示词过长，可能影响AI理解")
	}
	
	if !strings.Contains(prompt, "JSON") {
		issues = append(issues, "缺少明确的输出格式要求")
	}
	
	if !strings.Contains(prompt, "{{.") {
		issues = append(issues, "没有使用模板变量")
	}
	
	return issues
}