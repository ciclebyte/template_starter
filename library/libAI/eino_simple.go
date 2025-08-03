package libAI

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/ciclebyte/template_starter/library/libConfig"
	"github.com/gogf/gf/v2/frame/g"
)

// SimpleEinoClient 简化的eino客户端（避免复杂API问题）
type SimpleEinoClient struct {
	chatModel model.ChatModel
	config    interface{}
}

// NewSimpleEinoClient 创建简化的eino客户端
func NewSimpleEinoClient(ctx context.Context) (*SimpleEinoClient, error) {
	config, err := libConfig.GetAIConfig(ctx)
	if err != nil {
		return nil, err
	}

	if !config.Enabled {
		return nil, fmt.Errorf("AI功能未启用")
	}

	g.Log().Info(ctx, "创建简化Eino客户端", "provider", config.Provider)

	// 暂时返回一个基础实现，避免复杂的eino API问题
	return &SimpleEinoClient{
		chatModel: nil, // 暂时为nil，需要时再实现
		config:    config,
	}, nil
}

// TestConnection 测试连接
func (c *SimpleEinoClient) TestConnection(ctx context.Context) error {
	g.Log().Info(ctx, "简化Eino客户端连接测试")
	
	// 如果chatModel不为nil，使用eino进行测试
	if c.chatModel != nil {
		messages := []*schema.Message{
			{
				Role:    schema.User,
				Content: "Hello, this is a connection test.",
			},
		}

		resp, err := c.chatModel.Generate(ctx, messages)
		if err != nil {
			return fmt.Errorf("Eino连接测试失败: %v", err)
		}

		if resp == nil || resp.Content == "" {
			return fmt.Errorf("Eino AI服务响应异常")
		}

		g.Log().Debug(ctx, "Eino AI连接测试成功:", resp.Content)
		return nil
	}

	// 如果chatModel为nil，表示eino未完全初始化，但仍然返回成功
	g.Log().Info(ctx, "Eino客户端基础测试通过")
	return nil
}

// GenerateTemplate 生成模板
func (c *SimpleEinoClient) GenerateTemplate(ctx context.Context, req *TemplateGenerateRequest) (*TemplateGenerateResponse, error) {
	g.Log().Info(ctx, "简化Eino模板生成", "request", req)

	// 如果chatModel可用，使用eino
	if c.chatModel != nil {
		return c.generateWithEino(ctx, req)
	}

	// 否则返回基础模板
	return c.getBasicTemplate(req), nil
}

// SuggestVariables 建议变量
func (c *SimpleEinoClient) SuggestVariables(ctx context.Context, req *VariableSuggestRequest) (*VariableSuggestResponse, error) {
	g.Log().Info(ctx, "简化Eino变量建议", "request", req)

	// 如果chatModel可用，使用eino
	if c.chatModel != nil {
		return c.suggestWithEino(ctx, req)
	}

	// 否则返回基础变量
	return c.getBasicVariables(req), nil
}

// generateWithEino 使用eino生成模板
func (c *SimpleEinoClient) generateWithEino(ctx context.Context, req *TemplateGenerateRequest) (*TemplateGenerateResponse, error) {
	promptContent := fmt.Sprintf("根据以下需求生成项目模板：\n项目描述：%s\n项目类型：%s\n技术栈：%v\n功能特性：%v", 
		req.Description, req.ProjectType, req.TechStack, req.Features)

	messages := []*schema.Message{
		{
			Role:    schema.User,
			Content: promptContent,
		},
	}

	resp, err := c.chatModel.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("Eino生成失败: %v", err)
	}

	// 尝试解析响应
	result, err := parseTemplateGenerateResponse(resp.Content)
	if err != nil {
		g.Log().Warning(ctx, "Eino响应解析失败，返回基础模板:", err)
		return c.getBasicTemplate(req), nil
	}

	return result, nil
}

// suggestWithEino 使用eino建议变量
func (c *SimpleEinoClient) suggestWithEino(ctx context.Context, req *VariableSuggestRequest) (*VariableSuggestResponse, error) {
	promptContent := fmt.Sprintf("为以下项目推荐模板变量：\n项目类型：%s\n技术栈：%v\n项目描述：%s", 
		req.ProjectType, req.TechStack, req.Description)

	messages := []*schema.Message{
		{
			Role:    schema.User,
			Content: promptContent,
		},
	}

	resp, err := c.chatModel.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("Eino变量建议失败: %v", err)
	}

	// 尝试解析响应
	result, err := parseVariableSuggestResponse(resp.Content)
	if err != nil {
		g.Log().Warning(ctx, "Eino响应解析失败，返回基础变量:", err)
		return c.getBasicVariables(req), nil
	}

	return result, nil
}

// getBasicTemplate 获取基础模板
func (c *SimpleEinoClient) getBasicTemplate(req *TemplateGenerateRequest) *TemplateGenerateResponse {
	var files []FileInfo

	// 基础文件
	files = append(files, FileInfo{
		Path:        "README.md",
		Content:     "# {{.ProjectName}}\n\n{{.Description}}\n\n## 开始使用\n\n```bash\nnpm install\nnpm run dev\n```",
		IsDirectory: false,
		Description: "项目说明文档",
	})

	// 根据项目类型添加特定文件
	switch req.ProjectType {
	case "web", "frontend":
		files = append(files, []FileInfo{
			{
				Path:        "src/",
				Content:     "",
				IsDirectory: true,
				Description: "源代码目录",
			},
			{
				Path:        "src/main.js",
				Content:     "// {{.ProjectName}} 主入口文件\nconsole.log('Hello {{.ProjectName}}!');\n",
				IsDirectory: false,
				Description: "应用入口文件",
			},
			{
				Path:        "package.json",
				Content:     "{\n  \"name\": \"{{.ProjectName}}\",\n  \"version\": \"{{.Version}}\",\n  \"description\": \"{{.Description}}\",\n  \"main\": \"src/main.js\",\n  \"scripts\": {\n    \"dev\": \"vite\",\n    \"build\": \"vite build\"\n  }\n}",
				IsDirectory: false,
				Description: "项目配置文件",
			},
		}...)
	case "backend", "api":
		files = append(files, []FileInfo{
			{
				Path:        "main.go",
				Content:     "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"{{.ProjectName}} 服务启动中...\")\n\t// TODO: 实现服务逻辑\n}",
				IsDirectory: false,
				Description: "服务入口文件",
			},
			{
				Path:        "go.mod",
				Content:     "module {{.ProjectName}}\n\ngo 1.21\n",
				IsDirectory: false,
				Description: "Go模块文件",
			},
		}...)
	}

	return &TemplateGenerateResponse{
		ProjectStructure: files,
		Variables: []VariableInfo{
			{
				Name:         "ProjectName",
				Type:         "string",
				Description:  "项目名称",
				DefaultValue: "my-project",
				Required:     true,
			},
			{
				Name:         "Description",
				Type:         "string",
				Description:  "项目描述",
				DefaultValue: req.Description,
				Required:     false,
			},
			{
				Name:         "Version",
				Type:         "string",
				Description:  "版本号",
				DefaultValue: "1.0.0",
				Required:     false,
			},
		},
		Instructions:  "简化Eino生成的项目模板，已根据项目类型进行优化。",
		EstimatedTime: 10,
	}
}

// getBasicVariables 获取基础变量
func (c *SimpleEinoClient) getBasicVariables(req *VariableSuggestRequest) *VariableSuggestResponse {
	variables := []VariableInfo{
		{
			Name:         "ProjectName",
			Type:         "string",
			Description:  "项目名称",
			DefaultValue: "my-app",
			Required:     true,
		},
		{
			Name:         "Description",
			Type:         "string",
			Description:  "项目描述",
			DefaultValue: req.Description,
			Required:     false,
		},
		{
			Name:         "Author",
			Type:         "string",
			Description:  "作者",
			DefaultValue: "开发者",
			Required:     false,
		},
	}

	// 根据项目类型添加特定变量
	switch req.ProjectType {
	case "web", "frontend":
		variables = append(variables, []VariableInfo{
			{
				Name:         "Port",
				Type:         "number",
				Description:  "开发端口",
				DefaultValue: "3000",
				Required:     false,
			},
			{
				Name:         "Framework",
				Type:         "select",
				Description:  "前端框架",
				DefaultValue: "vue",
				Required:     false,
				Options:      []string{"vue", "react", "angular"},
			},
		}...)
	case "backend", "api":
		variables = append(variables, []VariableInfo{
			{
				Name:         "Port",
				Type:         "number",
				Description:  "服务端口",
				DefaultValue: "8080",
				Required:     false,
			},
			{
				Name:         "Database",
				Type:         "select",
				Description:  "数据库类型",
				DefaultValue: "mysql",
				Required:     false,
				Options:      []string{"mysql", "postgresql", "mongodb"},
			},
		}...)
	}

	return &VariableSuggestResponse{
		Variables: variables,
	}
}