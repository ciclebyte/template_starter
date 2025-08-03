package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/ciclebyte/template_starter/cli/internal/client"
	"github.com/ciclebyte/template_starter/cli/internal/config"
	"github.com/ciclebyte/template_starter/cli/internal/generator"
	"github.com/ciclebyte/template_starter/cli/internal/interactive"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [project-name]",
	Short: "从模板创建新项目",
	Long: `从指定的模板创建一个新项目。

示例:
  template-cli create my-app --template go-web
  template-cli create my-service --template microservice --output ./projects
  template-cli create frontend --template vue3-admin --interactive
  template-cli create  # 进入完全交互式模式`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCreateCommand(cmd, args)
	},
}

// runCreateCommand 执行创建命令
func runCreateCommand(cmd *cobra.Command, args []string) error {
	var projectName string
	if len(args) > 0 {
		projectName = args[0]
	}
	
	templateName, _ := cmd.Flags().GetString("template")
	output, _ := cmd.Flags().GetString("output")
	interactiveMode, _ := cmd.Flags().GetBool("interactive")
	configFile, _ := cmd.Flags().GetString("config")
	force, _ := cmd.Flags().GetBool("force")

	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	// 创建API客户端
	apiClient := client.NewClient(cfg.Server.URL, cfg.Server.APIKey)

	// 进入交互式模式的条件：没有项目名称或没有模板
	isInteractiveMode := projectName == "" || templateName == ""
	
	if isInteractiveMode {
		fmt.Println("进入交互式模式...")
		
		// 第一步：收集输出目录（如果没有通过flag提供）
		if output == "." {
			fmt.Println("\n第1步：设置输出目录")
			output, err = interactive.CollectOutputDirectory()
			if err != nil {
				return fmt.Errorf("收集输出目录失败: %w", err)
			}
			fmt.Printf("输出目录设置为: %s\n", output)
		}
		
		// 第二步：收集项目名称（如果没有提供）
		if projectName == "" {
			fmt.Println("\n第2步：设置项目名称")
			projectName, err = interactive.CollectProjectName()
			if err != nil {
				return fmt.Errorf("收集项目名称失败: %w", err)
			}
			fmt.Printf("项目名称设置为: %s\n", projectName)
		}
	} else {
		// 非交互式模式，但仍然需要确认输出目录
		if output == "." {
			fmt.Printf("输出目录: %s (当前目录)\n", output)
		} else {
			fmt.Printf("输出目录: %s\n", output)
		}
	}

	var selectedTemplate *client.Template

	// 如果没有指定模板，进入交互式选择
	if templateName == "" {
		if isInteractiveMode {
			fmt.Println("\n第3步：选择项目模板")
		} else {
			fmt.Println("选择项目模板...")
		}
		selectedTemplate, err = interactive.SelectTemplate(apiClient)
		if err != nil {
			return fmt.Errorf("模板选择失败: %w", err)
		}
		if selectedTemplate == nil {
			fmt.Println("已取消项目创建")
			return nil
		}
		templateName = selectedTemplate.Name
		fmt.Printf("模板选择: %s\n", selectedTemplate.Name)
	} else {
		// 获取指定模板信息
		selectedTemplate, err = apiClient.GetTemplateInfo(templateName)
		if err != nil {
			return fmt.Errorf("获取模板信息失败: %w", err)
		}
	}

	fmt.Printf("\n使用模板: %s\n", selectedTemplate.Name)
	fmt.Printf("模板描述: %s\n", selectedTemplate.Description)

	// 收集变量
	variables := make(map[string]interface{})

	if configFile != "" {
		// TODO: 从配置文件加载变量
		fmt.Printf("从配置文件加载变量: %s\n", configFile)
	} else if interactiveMode || isInteractiveMode {
		// 交互式收集变量
		if isInteractiveMode {
			fmt.Println("\n第4步：配置模板变量")
		}
		variables, err = interactive.CollectVariables(selectedTemplate)
		if err != nil {
			return fmt.Errorf("收集变量失败: %w", err)
		}

		// 确认变量配置
		confirmed, err := interactive.ConfirmVariables(variables)
		if err != nil {
			return fmt.Errorf("确认变量失败: %w", err)
		}
		if !confirmed {
			fmt.Println("已取消项目创建")
			return nil
		}
	}

	// 添加内置变量
	variables["ProjectName"] = projectName

	fmt.Printf("\n开始创建项目...\n")

	// 调用API渲染模板
	renderedFiles, err := apiClient.RenderTemplate(fmt.Sprintf("%d", selectedTemplate.ID), variables)
	if err != nil {
		return fmt.Errorf("渲染模板失败: %w", err)
	}

	// 确定输出目录
	outputDir, err := filepath.Abs(output)
	if err != nil {
		return fmt.Errorf("解析输出目录失败: %w", err)
	}

	// 创建生成器
	gen := generator.NewGenerator(outputDir, force)

	// 生成项目
	if err := gen.GenerateProject(projectName, renderedFiles); err != nil {
		return fmt.Errorf("生成项目失败: %w", err)
	}

	projectPath := filepath.Join(outputDir, projectName)
	fmt.Printf("\n🎉 项目创建成功!\n")
	fmt.Printf("📁 项目位置: %s\n", projectPath)
	fmt.Printf("🚀 可以开始开发了!\n")
	return nil
}

func init() {
	rootCmd.AddCommand(createCmd)

	// 模板标识(可选，不指定时进入交互选择)
	createCmd.Flags().StringP("template", "t", "", "模板名称或ID (可选，不指定时进入交互式选择)")

	// 可选参数
	createCmd.Flags().StringP("output", "o", ".", "输出目录")
	createCmd.Flags().BoolP("interactive", "i", false, "启用交互式变量配置")
	createCmd.Flags().StringP("config", "c", "", "变量配置文件路径")
	createCmd.Flags().BoolP("force", "f", false, "强制覆盖已存在的目录")
}
