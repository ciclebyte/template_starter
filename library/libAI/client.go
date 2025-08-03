package libAI

import (
	"context"
	"fmt"

	internalModel "github.com/ciclebyte/template_starter/internal/model"
	"github.com/ciclebyte/template_starter/library/libConfig"
	"github.com/gogf/gf/v2/frame/g"
)

// AIClient AI客户端接口
type AIClient interface {
	TestConnection(ctx context.Context) error
	GenerateTemplate(ctx context.Context, req *TemplateGenerateRequest) (*TemplateGenerateResponse, error)
	SuggestVariables(ctx context.Context, req *VariableSuggestRequest) (*VariableSuggestResponse, error)
}

// TemplateGenerateRequest 模板生成请求
type TemplateGenerateRequest struct {
	Description   string            `json:"description"`   // 项目描述
	ProjectType   string            `json:"projectType"`   // 项目类型
	TechStack     []string          `json:"techStack"`     // 技术栈
	Variables     map[string]string `json:"variables"`     // 变量
	Features      []string          `json:"features"`      // 功能特性
}

// TemplateGenerateResponse 模板生成响应
type TemplateGenerateResponse struct {
	ProjectStructure []FileInfo     `json:"projectStructure"` // 项目结构
	Variables        []VariableInfo `json:"variables"`        // 推荐变量
	Instructions     string         `json:"instructions"`     // 使用说明
	EstimatedTime    int            `json:"estimatedTime"`    // 预估完成时间(分钟)
}

// VariableSuggestRequest 变量建议请求
type VariableSuggestRequest struct {
	ProjectType string   `json:"projectType"`
	TechStack   []string `json:"techStack"`
	Description string   `json:"description"`
}

// VariableSuggestResponse 变量建议响应
type VariableSuggestResponse struct {
	Variables []VariableInfo `json:"variables"`
}

// FileInfo 文件信息
type FileInfo struct {
	Path        string `json:"path"`
	Content     string `json:"content"`
	IsDirectory bool   `json:"isDirectory"`
	Description string `json:"description"`
}

// VariableInfo 变量信息
type VariableInfo struct {
	Name         string   `json:"name"`
	Type         string   `json:"type"`         // string, number, boolean
	Description  string   `json:"description"`
	DefaultValue string   `json:"defaultValue"`
	Required     bool     `json:"required"`
	Options      []string `json:"options,omitempty"` // 选项值（如果是选择类型）
}

// SimpleAIClient 简化的AI客户端
type SimpleAIClient struct {
	config *internalModel.AIConfig
}

// NewAIClient 创建AI客户端
func NewAIClient(ctx context.Context) (AIClient, error) {
	config, err := libConfig.GetAIConfig(ctx)
	if err != nil {
		return nil, err
	}

	if !config.Enabled {
		return nil, fmt.Errorf("AI功能未启用")
	}

	return &SimpleAIClient{
		config: config,
	}, nil
}

// TestConnection 测试连接
func (c *SimpleAIClient) TestConnection(ctx context.Context) error {
	g.Log().Info(ctx, "AI连接测试", "provider", c.config.Provider)
	
	// TODO: 实际的AI服务连接测试
	// 这里可以添加真实的HTTP请求来测试AI服务连接
	
	return nil
}

// GenerateTemplate 生成模板
func (c *SimpleAIClient) GenerateTemplate(ctx context.Context, req *TemplateGenerateRequest) (*TemplateGenerateResponse, error) {
	g.Log().Info(ctx, "AI模板生成", "request", req)
	
	// 根据不同项目类型生成不同的模板结构
	projectStructure := c.generateProjectStructure(req)
	variables := c.generateVariables(req)
	
	return &TemplateGenerateResponse{
		ProjectStructure: projectStructure,
		Variables:        variables,
		Instructions:     c.generateInstructions(req),
		EstimatedTime:    c.estimateTime(req),
	}, nil
}

// SuggestVariables 建议变量
func (c *SimpleAIClient) SuggestVariables(ctx context.Context, req *VariableSuggestRequest) (*VariableSuggestResponse, error) {
	g.Log().Info(ctx, "AI变量建议", "request", req)
	
	variables := c.generateVariablesByType(req)
	
	return &VariableSuggestResponse{
		Variables: variables,
	}, nil
}

// generateProjectStructure 生成项目结构
func (c *SimpleAIClient) generateProjectStructure(req *TemplateGenerateRequest) []FileInfo {
	var files []FileInfo
	
	// 基础文件
	files = append(files, FileInfo{
		Path:        "README.md",
		Content:     "# {{.ProjectName}}\n\n{{.Description}}\n\n## 安装\n\n```bash\nnpm install\n```\n\n## 使用\n\n```bash\nnpm run dev\n```",
		IsDirectory: false,
		Description: "项目说明文档",
	})
	
	files = append(files, FileInfo{
		Path:        ".gitignore",
		Content:     "node_modules/\n.env\n.DS_Store\ndist/\nbuild/",
		IsDirectory: false,
		Description: "Git忽略配置",
	})
	
	// 根据项目类型添加特定文件
	switch req.ProjectType {
	case "web", "frontend":
		files = append(files, c.generateWebFiles(req)...)
	case "backend", "api":
		files = append(files, c.generateBackendFiles(req)...)
	case "mobile":
		files = append(files, c.generateMobileFiles(req)...)
	default:
		files = append(files, c.generateGenericFiles(req)...)
	}
	
	return files
}

// generateWebFiles 生成Web项目文件
func (c *SimpleAIClient) generateWebFiles(req *TemplateGenerateRequest) []FileInfo {
	files := []FileInfo{
		{
			Path:        "src/",
			Content:     "",
			IsDirectory: true,
			Description: "源代码目录",
		},
		{
			Path:        "src/main.js",
			Content:     "import { createApp } from 'vue'\nimport App from './App.vue'\n\ncreateApp(App).mount('#app')",
			IsDirectory: false,
			Description: "应用入口文件",
		},
		{
			Path:        "src/App.vue",
			Content:     "<template>\n  <div id=\"app\">\n    <h1>{{.ProjectName}}</h1>\n    <p>{{.Description}}</p>\n  </div>\n</template>\n\n<script>\nexport default {\n  name: 'App'\n}\n</script>",
			IsDirectory: false,
			Description: "主应用组件",
		},
		{
			Path:        "package.json",
			Content:     "{\n  \"name\": \"{{.ProjectName}}\",\n  \"version\": \"{{.Version}}\",\n  \"description\": \"{{.Description}}\",\n  \"main\": \"src/main.js\",\n  \"scripts\": {\n    \"dev\": \"vite\",\n    \"build\": \"vite build\",\n    \"preview\": \"vite preview\"\n  },\n  \"dependencies\": {\n    \"vue\": \"^3.0.0\"\n  },\n  \"devDependencies\": {\n    \"vite\": \"^4.0.0\"\n  }\n}",
			IsDirectory: false,
			Description: "项目配置文件",
		},
		{
			Path:        "vite.config.js",
			Content:     "import { defineConfig } from 'vite'\nimport vue from '@vitejs/plugin-vue'\n\nexport default defineConfig({\n  plugins: [vue()],\n  server: {\n    port: {{.Port}}\n  }\n})",
			IsDirectory: false,
			Description: "Vite配置文件",
		},
	}
	
	return files
}

// generateBackendFiles 生成后端项目文件
func (c *SimpleAIClient) generateBackendFiles(req *TemplateGenerateRequest) []FileInfo {
	files := []FileInfo{
		{
			Path:        "main.go",
			Content:     "package main\n\nimport (\n\t\"fmt\"\n\t\"net/http\"\n)\n\nfunc main() {\n\thttp.HandleFunc(\"/\", func(w http.ResponseWriter, r *http.Request) {\n\t\tfmt.Fprintf(w, \"Hello from {{.ProjectName}}!\")\n\t})\n\n\tfmt.Println(\"Server starting on port {{.Port}}...\")\n\thttp.ListenAndServe(\":{{.Port}}\", nil)\n}",
			IsDirectory: false,
			Description: "应用入口文件",
		},
		{
			Path:        "go.mod",
			Content:     "module {{.ProjectName}}\n\ngo 1.21\n\nrequire (\n\t// Add your dependencies here\n)",
			IsDirectory: false,
			Description: "Go模块配置",
		},
		{
			Path:        "api/",
			Content:     "",
			IsDirectory: true,
			Description: "API接口目录",
		},
		{
			Path:        "internal/",
			Content:     "",
			IsDirectory: true,
			Description: "内部代码目录",
		},
	}
	
	return files
}

// generateMobileFiles 生成移动应用文件
func (c *SimpleAIClient) generateMobileFiles(req *TemplateGenerateRequest) []FileInfo {
	files := []FileInfo{
		{
			Path:        "App.js",
			Content:     "import React from 'react';\nimport { Text, View } from 'react-native';\n\nexport default function App() {\n  return (\n    <View style={{flex: 1, justifyContent: 'center', alignItems: 'center'}}>\n      <Text>{{.ProjectName}}</Text>\n      <Text>{{.Description}}</Text>\n    </View>\n  );\n}",
			IsDirectory: false,
			Description: "应用主组件",
		},
		{
			Path:        "package.json",
			Content:     "{\n  \"name\": \"{{.ProjectName}}\",\n  \"version\": \"{{.Version}}\",\n  \"main\": \"App.js\",\n  \"scripts\": {\n    \"start\": \"expo start\",\n    \"android\": \"expo start --android\",\n    \"ios\": \"expo start --ios\"\n  },\n  \"dependencies\": {\n    \"react\": \"^18.0.0\",\n    \"react-native\": \"^0.72.0\"\n  }\n}",
			IsDirectory: false,
			Description: "项目配置文件",
		},
	}
	
	return files
}

// generateGenericFiles 生成通用项目文件
func (c *SimpleAIClient) generateGenericFiles(req *TemplateGenerateRequest) []FileInfo {
	return []FileInfo{
		{
			Path:        "src/",
			Content:     "",
			IsDirectory: true,
			Description: "源代码目录",
		},
		{
			Path:        "docs/",
			Content:     "",
			IsDirectory: true,
			Description: "文档目录",
		},
	}
}

// generateVariables 生成推荐变量
func (c *SimpleAIClient) generateVariables(req *TemplateGenerateRequest) []VariableInfo {
	variables := []VariableInfo{
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
			Name:         "Author",
			Type:         "string",
			Description:  "作者姓名",
			DefaultValue: "开发者",
			Required:     false,
		},
		{
			Name:         "Version",
			Type:         "string",
			Description:  "版本号",
			DefaultValue: "1.0.0",
			Required:     false,
		},
	}
	
	// 根据项目类型添加特定变量
	switch req.ProjectType {
	case "web", "frontend":
		variables = append(variables, VariableInfo{
			Name:         "Port",
			Type:         "number",
			Description:  "开发服务器端口",
			DefaultValue: "3000",
			Required:     false,
		})
	case "backend", "api":
		variables = append(variables, VariableInfo{
			Name:         "Port",
			Type:         "number",
			Description:  "服务端口",
			DefaultValue: "8080",
			Required:     false,
		})
	}
	
	return variables
}

// generateVariablesByType 根据类型生成变量建议
func (c *SimpleAIClient) generateVariablesByType(req *VariableSuggestRequest) []VariableInfo {
	var variables []VariableInfo
	
	// 基础变量
	variables = append(variables, []VariableInfo{
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
			Description:  "作者姓名",
			DefaultValue: "开发者",
			Required:     false,
		},
	}...)
	
	// 根据项目类型添加特定变量
	switch req.ProjectType {
	case "web", "frontend":
		variables = append(variables, []VariableInfo{
			{
				Name:         "Theme",
				Type:         "select",
				Description:  "UI主题",
				DefaultValue: "light",
				Required:     false,
				Options:      []string{"light", "dark"},
			},
			{
				Name:         "Language",
				Type:         "select",
				Description:  "开发语言",
				DefaultValue: "javascript",
				Required:     false,
				Options:      []string{"javascript", "typescript"},
			},
		}...)
		
	case "backend", "api":
		variables = append(variables, []VariableInfo{
			{
				Name:         "DatabaseType",
				Type:         "select",
				Description:  "数据库类型",
				DefaultValue: "mysql",
				Required:     false,
				Options:      []string{"mysql", "postgresql", "mongodb"},
			},
		}...)
	}
	
	return variables
}

// generateInstructions 生成使用说明
func (c *SimpleAIClient) generateInstructions(req *TemplateGenerateRequest) string {
	instructions := "## 使用说明\n\n1. 解压模板文件到目标目录\n2. 根据需要修改模板变量\n3. 运行初始化命令\n\n"
	
	switch req.ProjectType {
	case "web", "frontend":
		instructions += "## 开发步骤\n1. 安装依赖：npm install\n2. 启动开发服务器：npm run dev\n3. 构建生产版本：npm run build"
	case "backend", "api":
		instructions += "## 开发步骤\n1. 安装Go依赖：go mod tidy\n2. 运行应用：go run main.go\n3. 构建二进制：go build"
	case "mobile":
		instructions += "## 开发步骤\n1. 安装依赖：npm install\n2. 启动开发服务器：npm start\n3. 在模拟器中运行：npm run android 或 npm run ios"
	}
	
	return instructions
}

// estimateTime 估算完成时间
func (c *SimpleAIClient) estimateTime(req *TemplateGenerateRequest) int {
	baseTime := 10 // 基础时间10分钟
	
	// 根据技术栈复杂度调整时间
	techComplexity := len(req.TechStack) * 5
	featureComplexity := len(req.Features) * 3
	
	return baseTime + techComplexity + featureComplexity
}