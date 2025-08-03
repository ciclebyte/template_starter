package libAI

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// parseTemplateGenerateResponse 解析AI生成的模板响应
func parseTemplateGenerateResponse(content string) (*TemplateGenerateResponse, error) {
	// 提取JSON内容
	jsonContent, err := extractJSON(content)
	if err != nil {
		return nil, fmt.Errorf("提取JSON失败: %v", err)
	}

	var response TemplateGenerateResponse
	err = json.Unmarshal([]byte(jsonContent), &response)
	if err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	// 验证响应数据
	if err := validateTemplateResponse(&response); err != nil {
		return nil, fmt.Errorf("响应数据验证失败: %v", err)
	}

	return &response, nil
}

// parseVariableSuggestResponse 解析AI生成的变量建议响应
func parseVariableSuggestResponse(content string) (*VariableSuggestResponse, error) {
	jsonContent, err := extractJSON(content)
	if err != nil {
		return nil, fmt.Errorf("提取JSON失败: %v", err)
	}

	var response VariableSuggestResponse
	err = json.Unmarshal([]byte(jsonContent), &response)
	if err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	// 验证响应数据
	if err := validateVariableResponse(&response); err != nil {
		return nil, fmt.Errorf("响应数据验证失败: %v", err)
	}

	return &response, nil
}

// extractJSON 从文本中提取JSON内容
func extractJSON(content string) (string, error) {
	// 清理内容
	content = strings.TrimSpace(content)
	
	// 尝试直接解析（如果内容就是JSON）
	if strings.HasPrefix(content, "{") && strings.HasSuffix(content, "}") {
		return content, nil
	}

	// 使用正则表达式提取JSON代码块
	patterns := []string{
		// ```json ... ```
		"(?s)```json\\s*(\\{.*?\\})\\s*```",
		// ``` ... ```
		"(?s)```\\s*(\\{.*?\\})\\s*```",
		// 直接查找JSON对象
		"(?s)(\\{[^{}]*(?:\\{[^{}]*\\}[^{}]*)*\\})",
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(content)
		if len(matches) > 1 {
			jsonStr := strings.TrimSpace(matches[1])
			if isValidJSONStart(jsonStr) {
				return jsonStr, nil
			}
		}
	}

	return "", fmt.Errorf("未找到有效的JSON内容")
}

// isValidJSONStart 检查字符串是否像JSON开头
func isValidJSONStart(s string) bool {
	s = strings.TrimSpace(s)
	return strings.HasPrefix(s, "{") && strings.Contains(s, ":")
}

// validateTemplateResponse 验证模板响应数据
func validateTemplateResponse(response *TemplateGenerateResponse) error {
	if response == nil {
		return fmt.Errorf("响应数据为空")
	}

	// 验证项目结构
	if len(response.ProjectStructure) == 0 {
		return fmt.Errorf("项目结构不能为空")
	}

	for i, file := range response.ProjectStructure {
		if file.Path == "" {
			return fmt.Errorf("第%d个文件路径不能为空", i+1)
		}
		// 目录不需要内容，文件需要内容
		if !file.IsDirectory && file.Content == "" {
			// 允许某些文件为空（如配置文件）
		}
	}

	// 验证变量
	for i, variable := range response.Variables {
		if variable.Name == "" {
			return fmt.Errorf("第%d个变量名称不能为空", i+1)
		}
		if variable.Type == "" {
			variable.Type = "string" // 默认类型
		}
		// 验证类型有效性
		if !isValidVariableType(variable.Type) {
			return fmt.Errorf("第%d个变量类型无效: %s", i+1, variable.Type)
		}
	}

	// 预估时间合理性检查
	if response.EstimatedTime < 0 {
		response.EstimatedTime = 15 // 默认15分钟
	}

	return nil
}

// validateVariableResponse 验证变量响应数据
func validateVariableResponse(response *VariableSuggestResponse) error {
	if response == nil {
		return fmt.Errorf("响应数据为空")
	}

	if len(response.Variables) == 0 {
		return fmt.Errorf("变量列表不能为空")
	}

	for i, variable := range response.Variables {
		if variable.Name == "" {
			return fmt.Errorf("第%d个变量名称不能为空", i+1)
		}
		if variable.Type == "" {
			variable.Type = "string" // 默认类型
		}
		if !isValidVariableType(variable.Type) {
			return fmt.Errorf("第%d个变量类型无效: %s", i+1, variable.Type)
		}
	}

	return nil
}

// isValidVariableType 检查变量类型是否有效
func isValidVariableType(varType string) bool {
	validTypes := []string{"string", "number", "boolean", "select", "textarea"}
	for _, valid := range validTypes {
		if varType == valid {
			return true
		}
	}
	return false
}

// cleanAndParseJSON 清理并解析JSON
func cleanAndParseJSON(content string) (map[string]interface{}, error) {
	// 移除注释和多余字符
	content = removeJSONComments(content)
	
	var result map[string]interface{}
	err := json.Unmarshal([]byte(content), &result)
	return result, err
}

// removeJSONComments 移除JSON中的注释（非标准，但AI可能生成）
func removeJSONComments(content string) string {
	// 移除 // 注释
	lines := strings.Split(content, "\n")
	var cleanLines []string
	
	for _, line := range lines {
		// 查找 // 但不在字符串内
		inString := false
		escape := false
		commentStart := -1
		
		for i, ch := range line {
			if escape {
				escape = false
				continue
			}
			
			if ch == '\\' {
				escape = true
				continue
			}
			
			if ch == '"' {
				inString = !inString
				continue
			}
			
			if !inString && ch == '/' && i+1 < len(line) && line[i+1] == '/' {
				commentStart = i
				break
			}
		}
		
		if commentStart >= 0 {
			line = line[:commentStart]
		}
		
		cleanLines = append(cleanLines, line)
	}
	
	return strings.Join(cleanLines, "\n")
}

// formatJSONResponse 格式化JSON响应以提高解析成功率
func formatJSONResponse(content string) string {
	// 基本清理
	content = strings.TrimSpace(content)
	
	// 移除可能的前后文
	if idx := strings.Index(content, "{"); idx > 0 {
		content = content[idx:]
	}
	
	if idx := strings.LastIndex(content, "}"); idx >= 0 && idx < len(content)-1 {
		content = content[:idx+1]
	}
	
	return content
}