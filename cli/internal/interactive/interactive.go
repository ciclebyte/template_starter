package interactive

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ciclebyte/template_starter/cli/internal/client"
)

// Collector 交互式变量收集器
type Collector struct {
	reader *bufio.Reader
}

// NewCollector 创建新的收集器
func NewCollector() *Collector {
	return &Collector{
		reader: bufio.NewReader(os.Stdin),
	}
}

// CollectVariables 交互式收集变量值
func (c *Collector) CollectVariables(tmpl *client.Template) (map[string]interface{}, error) {
	variables := make(map[string]interface{})
	
	fmt.Printf("🚀 正在为模板 '%s' 配置变量\n", tmpl.Name)
	fmt.Println("📝 请根据提示输入变量值（直接回车使用默认值）")
	fmt.Println()
	
	for _, variable := range tmpl.Variables {
		value, err := c.collectVariable(variable)
		if err != nil {
			return nil, err
		}
		variables[variable.Name] = value
	}
	
	return variables, nil
}

// collectVariable 收集单个变量值
func (c *Collector) collectVariable(variable client.TemplateVariable) (interface{}, error) {
	// 显示变量信息
	fmt.Printf("📌 %s", variable.Name)
	if variable.Required {
		fmt.Print(" (必需)")
	}
	fmt.Println()
	
	if variable.Description != "" {
		fmt.Printf("   📄 %s\n", variable.Description)
	}
	
	if variable.DefaultValue != nil {
		fmt.Printf("   🔧 默认值: %v\n", variable.DefaultValue)
	}
	
	if len(variable.Options) > 0 {
		fmt.Printf("   🎯 可选值: %s\n", strings.Join(variable.Options, ", "))
	}
	
	// 根据类型收集值
	switch variable.Type {
	case "string":
		return c.collectString(variable)
	case "boolean":
		return c.collectBoolean(variable)
	case "number":
		return c.collectNumber(variable)
	default:
		return c.collectString(variable)
	}
}

// collectString 收集字符串值
func (c *Collector) collectString(variable client.TemplateVariable) (string, error) {
	for {
		fmt.Printf("   ✏️  输入值: ")
		input, err := c.reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		
		input = strings.TrimSpace(input)
		
		// 如果为空且有默认值，使用默认值
		if input == "" && variable.DefaultValue != nil {
			if defaultStr, ok := variable.DefaultValue.(string); ok {
				return defaultStr, nil
			}
		}
		
		// 如果为空且是必需的，要求重新输入
		if input == "" && variable.Required {
			fmt.Println("   ❌ 此变量为必需，请输入值")
			continue
		}
		
		// 如果有选项限制，检查是否在选项中
		if len(variable.Options) > 0 && input != "" {
			valid := false
			for _, option := range variable.Options {
				if input == option {
					valid = true
					break
				}
			}
			if !valid {
				fmt.Printf("   ❌ 无效选项，请选择: %s\n", strings.Join(variable.Options, ", "))
				continue
			}
		}
		
		return input, nil
	}
}

// collectBoolean 收集布尔值
func (c *Collector) collectBoolean(variable client.TemplateVariable) (bool, error) {
	for {
		fmt.Printf("   ✏️  输入值 (y/n): ")
		input, err := c.reader.ReadString('\n')
		if err != nil {
			return false, err
		}
		
		input = strings.TrimSpace(strings.ToLower(input))
		
		// 如果为空且有默认值，使用默认值
		if input == "" && variable.DefaultValue != nil {
			if defaultBool, ok := variable.DefaultValue.(bool); ok {
				return defaultBool, nil
			}
		}
		
		switch input {
		case "y", "yes", "true", "1":
			return true, nil
		case "n", "no", "false", "0":
			return false, nil
		case "":
			if variable.Required {
				fmt.Println("   ❌ 此变量为必需，请输入 y 或 n")
				continue
			}
			return false, nil
		default:
			fmt.Println("   ❌ 请输入 y/yes/true 或 n/no/false")
			continue
		}
	}
}

// collectNumber 收集数字值
func (c *Collector) collectNumber(variable client.TemplateVariable) (float64, error) {
	for {
		fmt.Printf("   ✏️  输入数字: ")
		input, err := c.reader.ReadString('\n')
		if err != nil {
			return 0, err
		}
		
		input = strings.TrimSpace(input)
		
		// 如果为空且有默认值，使用默认值
		if input == "" && variable.DefaultValue != nil {
			switch defaultVal := variable.DefaultValue.(type) {
			case float64:
				return defaultVal, nil
			case int:
				return float64(defaultVal), nil
			}
		}
		
		// 如果为空且是必需的，要求重新输入
		if input == "" && variable.Required {
			fmt.Println("   ❌ 此变量为必需，请输入数字")
			continue
		}
		
		if input == "" {
			return 0, nil
		}
		
		// 尝试解析数字
		if value, err := strconv.ParseFloat(input, 64); err == nil {
			return value, nil
		}
		
		fmt.Println("   ❌ 请输入有效的数字")
	}
}

// ConfirmGeneration 确认生成项目
func (c *Collector) ConfirmGeneration(projectName string, templateName string, variables map[string]interface{}) (bool, error) {
	fmt.Println()
	fmt.Println("📋 配置摘要:")
	fmt.Printf("   🎯 项目名称: %s\n", projectName)
	fmt.Printf("   📦 模板: %s\n", templateName)
	fmt.Println("   🔧 变量配置:")
	
	for key, value := range variables {
		fmt.Printf("      %s = %v\n", key, value)
	}
	
	fmt.Println()
	fmt.Print("🤔 确认生成项目吗？(y/n): ")
	
	input, err := c.reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes", nil
}