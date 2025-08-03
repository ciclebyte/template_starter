package interactive

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/ciclebyte/template_starter/cli/internal/client"
)

// CollectProjectName 收集项目名称
func CollectProjectName() (string, error) {
	validate := func(input string) error {
		input = strings.TrimSpace(input)
		if input == "" {
			return fmt.Errorf("项目名称不能为空")
		}
		// 简单验证：只允许字母、数字、连字符和下划线
		for _, r := range input {
			if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || 
			     (r >= '0' && r <= '9') || r == '-' || r == '_') {
				return fmt.Errorf("项目名称只能包含字母、数字、连字符和下划线")
			}
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "项目名称",
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

// CollectOutputDirectory 收集输出目录
func CollectOutputDirectory() (string, error) {
	validate := func(input string) error {
		input = strings.TrimSpace(input)
		if input == "" {
			return nil // 允许空值，默认为当前目录
		}
		// 这里可以添加更多路径验证逻辑
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "输出目录 (留空则使用当前目录)",
		Validate: validate,
		Default:  ".",
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	result = strings.TrimSpace(result)
	if result == "" {
		result = "."
	}

	return result, nil
}

// CollectVariables 收集模板变量
func CollectVariables(template *client.Template) (map[string]interface{}, error) {
	variables := make(map[string]interface{})

	fmt.Printf("\n配置模板变量 (%s)\n", template.Name)
	fmt.Println(strings.Repeat("-", 50))

	for _, variable := range template.Variables {
		value, err := collectSingleVariable(variable)
		if err != nil {
			return nil, fmt.Errorf("收集变量 %s 失败: %w", variable.Name, err)
		}
		variables[variable.Name] = value
	}

	return variables, nil
}

// collectSingleVariable 收集单个变量
func collectSingleVariable(variable client.TemplateVariable) (interface{}, error) {
	// 构建提示信息
	label := variable.Name
	if variable.Description != "" {
		label += fmt.Sprintf(" (%s)", variable.Description)
	}
	if variable.Required {
		label += " *"
	}

	switch variable.Type {
	case "boolean":
		return collectBooleanVariable(variable, label)
	case "number":
		return collectNumberVariable(variable, label)
	case "select":
		return collectSelectVariable(variable, label)
	default: // string
		return collectStringVariable(variable, label)
	}
}

// collectStringVariable 收集字符串变量
func collectStringVariable(variable client.TemplateVariable, label string) (string, error) {
	validate := func(input string) error {
		if variable.Required && strings.TrimSpace(input) == "" {
			return fmt.Errorf("该字段为必填项")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}

	// 设置默认值
	if variable.DefaultValue != nil {
		if defaultStr, ok := variable.DefaultValue.(string); ok {
			prompt.Default = defaultStr
		}
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

// collectBooleanVariable 收集布尔变量
func collectBooleanVariable(variable client.TemplateVariable, label string) (bool, error) {
	// 确定默认值
	defaultValue := false
	if variable.DefaultValue != nil {
		if defaultBool, ok := variable.DefaultValue.(bool); ok {
			defaultValue = defaultBool
		}
	}

	prompt := promptui.Prompt{
		Label:     label + " (y/n)",
		Default:   map[bool]string{true: "y", false: "n"}[defaultValue],
		AllowEdit: true,
		Validate: func(input string) error {
			input = strings.ToLower(strings.TrimSpace(input))
			if input != "y" && input != "n" && input != "yes" && input != "no" && 
			   input != "true" && input != "false" {
				return fmt.Errorf("请输入 y/n, yes/no, 或 true/false")
			}
			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		return false, err
	}

	result = strings.ToLower(strings.TrimSpace(result))
	return result == "y" || result == "yes" || result == "true", nil
}

// collectNumberVariable 收集数字变量
func collectNumberVariable(variable client.TemplateVariable, label string) (float64, error) {
	validate := func(input string) error {
		if variable.Required && strings.TrimSpace(input) == "" {
			return fmt.Errorf("该字段为必填项")
		}
		if input != "" {
			if _, err := strconv.ParseFloat(input, 64); err != nil {
				return fmt.Errorf("请输入有效的数字")
			}
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}

	// 设置默认值
	if variable.DefaultValue != nil {
		if defaultNum, ok := variable.DefaultValue.(float64); ok {
			prompt.Default = fmt.Sprintf("%.0f", defaultNum)
		} else if defaultInt, ok := variable.DefaultValue.(int); ok {
			prompt.Default = strconv.Itoa(defaultInt)
		}
	}

	result, err := prompt.Run()
	if err != nil {
		return 0, err
	}

	if strings.TrimSpace(result) == "" {
		return 0, nil
	}

	return strconv.ParseFloat(result, 64)
}

// collectSelectVariable 收集选择变量
func collectSelectVariable(variable client.TemplateVariable, label string) (string, error) {
	if len(variable.Options) == 0 {
		// 如果没有选项，回退到字符串输入
		return collectStringVariable(variable, label)
	}

	prompt := promptui.Select{
		Label: label,
		Items: variable.Options,
		Size:  10,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

// ConfirmVariables 确认变量配置
func ConfirmVariables(variables map[string]interface{}) (bool, error) {
	fmt.Println("\n变量配置预览:")
	fmt.Println(strings.Repeat("-", 30))
	
	for key, value := range variables {
		fmt.Printf("  %s: %v\n", key, value)
	}
	
	prompt := promptui.Prompt{
		Label:     "\n确认创建项目",
		Default:   "y",
		AllowEdit: true,
		Validate: func(input string) error {
			input = strings.ToLower(strings.TrimSpace(input))
			if input != "y" && input != "n" && input != "yes" && input != "no" {
				return fmt.Errorf("请输入 y/n 或 yes/no")
			}
			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		return false, err
	}

	result = strings.ToLower(strings.TrimSpace(result))
	return result == "y" || result == "yes", nil
}