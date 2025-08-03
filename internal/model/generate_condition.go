package model

// GenerateCondition 文件生成条件
type GenerateCondition struct {
	Enabled       bool   `json:"enabled"`       // 是否启用条件
	VariableName  string `json:"variableName"`  // 关联变量名
	ExpectedValue bool   `json:"expectedValue"` // 期望值
	Description   string `json:"description"`   // 条件描述
}