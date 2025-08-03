package ai

import (
	"context"
	"time"

	aiApi "github.com/ciclebyte/template_starter/api/v1/ai"
	"github.com/ciclebyte/template_starter/internal/service"
	"github.com/ciclebyte/template_starter/library/libAI"
	"github.com/ciclebyte/template_starter/library/libConfig"
	"github.com/gogf/gf/v2/frame/g"
)

type sAI struct{}

func init() {
	service.RegisterAIService(New())
}

func New() *sAI {
	return &sAI{}
}

// TestConnection 测试AI连接
func (s *sAI) TestConnection(ctx context.Context, req *aiApi.TestConnectionReq) (*aiApi.TestConnectionRes, error) {
	g.Log().Debug(ctx, "AI.TestConnection called")

	res := &aiApi.TestConnectionRes{
		Success: false,
		Message: "连接失败",
	}

	// 获取AI配置
	config, err := libConfig.GetAIConfig(ctx)
	if err != nil {
		res.Message = "获取AI配置失败: " + err.Error()
		return res, nil
	}

	if !config.Enabled {
		res.Message = "AI功能未启用"
		return res, nil
	}

	res.Provider = config.Provider

	// 创建AI客户端
	client, err := libAI.NewAIClient(ctx)
	if err != nil {
		res.Message = "创建AI客户端失败: " + err.Error()
		return res, nil
	}

	// 记录开始时间
	startTime := time.Now()

	// 测试连接
	err = client.TestConnection(ctx)
	
	// 计算延迟
	res.Latency = int(time.Since(startTime).Milliseconds())

	if err != nil {
		res.Message = "连接测试失败: " + err.Error()
		return res, nil
	}

	// 设置成功信息
	res.Success = true
	res.Message = "连接测试成功"
	
	if config.Provider == "openai" {
		res.Model = config.OpenAI.Model
	} else if config.Provider == "claude" {
		res.Model = config.Claude.Model
	}

	return res, nil
}

// GenerateTemplate 生成模板
func (s *sAI) GenerateTemplate(ctx context.Context, req *aiApi.GenerateTemplateReq) (*aiApi.GenerateTemplateRes, error) {
	g.Log().Debug(ctx, "AI.GenerateTemplate called with req:", req)

	// 检查AI功能是否启用
	config, err := libConfig.GetAIConfig(ctx)
	if err != nil {
		return nil, err
	}

	if !config.Enabled || !config.Features.TemplateGeneration {
		return nil, g.Err("AI模板生成功能未启用")
	}

	// 创建AI客户端
	client, err := libAI.NewAIClient(ctx)
	if err != nil {
		return nil, err
	}

	// 构建请求
	aiReq := &libAI.TemplateGenerateRequest{
		Description:   req.Description,
		ProjectType:   req.ProjectType,
		TechStack:     req.TechStack,
		Variables:     req.Variables,
		Features:      req.Features,
	}

	// 调用AI生成
	aiRes, err := client.GenerateTemplate(ctx, aiReq)
	if err != nil {
		return nil, err
	}

	// 转换响应格式
	res := &aiApi.GenerateTemplateRes{
		Instructions:  aiRes.Instructions,
		EstimatedTime: aiRes.EstimatedTime,
	}

	// 转换文件结构
	for _, file := range aiRes.ProjectStructure {
		res.ProjectStructure = append(res.ProjectStructure, aiApi.FileInfo{
			Path:        file.Path,
			Content:     file.Content,
			IsDirectory: file.IsDirectory,
			Description: file.Description,
		})
	}

	// 转换变量信息
	for _, variable := range aiRes.Variables {
		res.Variables = append(res.Variables, aiApi.VariableInfo{
			Name:         variable.Name,
			Type:         variable.Type,
			Description:  variable.Description,
			DefaultValue: variable.DefaultValue,
			Required:     variable.Required,
			Options:      variable.Options,
		})
	}

	return res, nil
}

// SuggestVariables 建议变量
func (s *sAI) SuggestVariables(ctx context.Context, req *aiApi.SuggestVariablesReq) (*aiApi.SuggestVariablesRes, error) {
	g.Log().Debug(ctx, "AI.SuggestVariables called with req:", req)

	// 检查AI功能是否启用
	config, err := libConfig.GetAIConfig(ctx)
	if err != nil {
		return nil, err
	}

	if !config.Enabled || !config.Features.VariableSuggestion {
		return nil, g.Err("AI变量建议功能未启用")
	}

	// 创建AI客户端
	client, err := libAI.NewAIClient(ctx)
	if err != nil {
		return nil, err
	}

	// 构建请求
	aiReq := &libAI.VariableSuggestRequest{
		ProjectType: req.ProjectType,
		TechStack:   req.TechStack,
		Description: req.Description,
	}

	// 调用AI建议
	aiRes, err := client.SuggestVariables(ctx, aiReq)
	if err != nil {
		return nil, err
	}

	// 转换响应格式
	res := &aiApi.SuggestVariablesRes{}
	for _, variable := range aiRes.Variables {
		res.Variables = append(res.Variables, aiApi.VariableInfo{
			Name:         variable.Name,
			Type:         variable.Type,
			Description:  variable.Description,
			DefaultValue: variable.DefaultValue,
			Required:     variable.Required,
			Options:      variable.Options,
		})
	}

	return res, nil
}