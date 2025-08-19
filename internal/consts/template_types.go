package consts

// 模板类型常量定义
const (
	// TemplateTypeBasic 基础模板 - 代码片段、配置文件、文档等普通模板
	TemplateTypeBasic = "basic"
	
	// TemplateTypeScaffold 脚手架模板 - 完整项目结构，可与CLI结合快速搭建项目
	TemplateTypeScaffold = "scaffold"
	
	// TemplateTypeDataDriven 数据驱动模板 - 基于外部数据源生成的模板
	TemplateTypeDataDriven = "data_driven"
)

// TemplateTypeLabels 模板类型标签映射
var TemplateTypeLabels = map[string]string{
	TemplateTypeBasic:      "基础模板",
	TemplateTypeScaffold:   "脚手架模板", 
	TemplateTypeDataDriven: "数据驱动模板",
}

// TemplateTypeDescriptions 模板类型描述
var TemplateTypeDescriptions = map[string]string{
	TemplateTypeBasic:      "适用于代码片段、配置文件、文档等普通模板",
	TemplateTypeScaffold:   "完整的项目结构模板，可与CLI工具结合快速搭建项目",
	TemplateTypeDataDriven: "基于外部数据源（数据库、API等）动态生成的模板",
}

// GetTemplateTypeLabel 获取模板类型标签
func GetTemplateTypeLabel(templateType string) string {
	if label, exists := TemplateTypeLabels[templateType]; exists {
		return label
	}
	return "未知类型"
}

// GetTemplateTypeDescription 获取模板类型描述
func GetTemplateTypeDescription(templateType string) string {
	if desc, exists := TemplateTypeDescriptions[templateType]; exists {
		return desc
	}
	return ""
}

// IsValidTemplateType 验证模板类型是否有效
func IsValidTemplateType(templateType string) bool {
	_, exists := TemplateTypeLabels[templateType]
	return exists
}

// GetAllTemplateTypes 获取所有模板类型
func GetAllTemplateTypes() []string {
	return []string{TemplateTypeBasic, TemplateTypeScaffold, TemplateTypeDataDriven}
}