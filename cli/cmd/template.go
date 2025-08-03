package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ciclebyte/template_starter/cli/internal/client"
	"github.com/ciclebyte/template_starter/cli/internal/config"
)

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "模板管理命令",
	Long: `管理远程模板的命令集合。

子命令:
  list    - 列出可用的模板
  info    - 显示模板详细信息
  search  - 搜索模板`,
}

// templateListCmd lists available templates
var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出可用的模板",
	Long: `列出远程服务器上所有可用的模板。

示例:
  template-cli template list
  template-cli template list --category web`,
	RunE: func(cmd *cobra.Command, args []string) error {
		category, _ := cmd.Flags().GetString("category")
		
		// 加载配置
		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("加载配置失败: %w", err)
		}
		
		// 创建API客户端
		apiClient := client.NewClient(cfg.Server.URL, cfg.Server.APIKey)
		
		// 获取模板列表
		templates, err := apiClient.ListTemplates(category)
		if err != nil {
			return fmt.Errorf("获取模板列表失败: %w", err)
		}
		
		// 显示模板列表
		if len(templates) == 0 {
			fmt.Println("没有找到模板")
			return nil
		}
		
		fmt.Printf("找到 %d 个模板:\n\n", len(templates))
		for _, tmpl := range templates {
			fmt.Printf("• %s (ID: %d)\n", tmpl.Name, tmpl.ID)
			if tmpl.Description != "" {
				fmt.Printf("  %s\n", tmpl.Description)
			}
			fmt.Println()
		}
		
		return nil
	},
}

// templateInfoCmd shows template information
var templateInfoCmd = &cobra.Command{
	Use:   "info [template-name]",
	Short: "显示模板详细信息",
	Long: `显示指定模板的详细信息，包括变量、文件结构等。

示例:
  template-cli template info go-web
  template-cli template info vue3-admin --variables`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		templateName := args[0]
		showVariables, _ := cmd.Flags().GetBool("variables")
		showFiles, _ := cmd.Flags().GetBool("files")
		
		// 加载配置
		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("加载配置失败: %w", err)
		}
		
		// 创建API客户端
		apiClient := client.NewClient(cfg.Server.URL, cfg.Server.APIKey)
		
		// 获取模板信息
		var template *client.Template
		if showFiles {
			template, err = apiClient.GetTemplate(templateName)
		} else {
			template, err = apiClient.GetTemplateInfo(templateName)
		}
		if err != nil {
			return fmt.Errorf("获取模板信息失败: %w", err)
		}
		
		// 显示基本信息
		fmt.Printf("模板名称: %s\n", template.Name)
		fmt.Printf("模板ID: %d\n", template.ID)
		fmt.Printf("分类ID: %d\n", template.CategoryId)
		fmt.Printf("描述: %s\n", template.Description)
		
		// 显示变量信息
		if showVariables && len(template.Variables) > 0 {
			fmt.Printf("\n变量列表:\n")
			for _, variable := range template.Variables {
				fmt.Printf("• %s (%s)", variable.Name, variable.VariableType)
				if variable.IsRequired == 1 {
					fmt.Printf(" *必需*")
				}
				fmt.Println()
				if variable.Description != "" {
					fmt.Printf("  %s\n", variable.Description)
				}
				if variable.DefaultValue != "" {
					fmt.Printf("  默认值: %s\n", variable.DefaultValue)
				}
				fmt.Println()
			}
		}
		
		// 显示文件结构
		if showFiles && len(template.Files) > 0 {
			fmt.Printf("\n文件结构:\n")
			for _, file := range template.Files {
				if file.IsDirectory {
					fmt.Printf("📁 %s/\n", file.Path)
				} else {
					fmt.Printf("📄 %s\n", file.Path)
				}
				if file.Condition != "" {
					fmt.Printf("   条件: %s\n", file.Condition)
				}
			}
		}
		
		return nil
	},
}


// templateSearchCmd searches templates
var templateSearchCmd = &cobra.Command{
	Use:   "search [keyword]",
	Short: "搜索模板",
	Long: `根据关键词搜索模板。

示例:
  template-cli template search web
  template-cli template search vue --category frontend`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		keyword := args[0]
		category, _ := cmd.Flags().GetString("category")
		
		// 加载配置
		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("加载配置失败: %w", err)
		}
		
		// 创建API客户端
		apiClient := client.NewClient(cfg.Server.URL, cfg.Server.APIKey)
		
		// 搜索模板
		templates, err := apiClient.SearchTemplates(keyword, category)
		if err != nil {
			return fmt.Errorf("搜索模板失败: %w", err)
		}
		
		// 显示搜索结果
		if len(templates) == 0 {
			fmt.Printf("没有找到包含 '%s' 的模板", keyword)
			if category != "" {
				fmt.Printf(" (分类: %s)", category)
			}
			fmt.Println()
			return nil
		}
		
		fmt.Printf("找到 %d 个包含 '%s' 的模板", len(templates), keyword)
		if category != "" {
			fmt.Printf(" (分类: %s)", category)
		}
		fmt.Println(":\n")
		
		for _, tmpl := range templates {
			fmt.Printf("• %s (ID: %d)\n", tmpl.Name, tmpl.ID)
			if tmpl.Description != "" {
				fmt.Printf("  %s\n", tmpl.Description)
			}
			fmt.Println()
		}
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(templateCmd)

	// Add subcommands
	templateCmd.AddCommand(templateListCmd)
	templateCmd.AddCommand(templateInfoCmd)
	templateCmd.AddCommand(templateSearchCmd)

	// Flags for template list
	templateListCmd.Flags().StringP("category", "c", "", "按分类过滤")

	// Flags for template info
	templateInfoCmd.Flags().BoolP("variables", "v", false, "显示模板变量")
	templateInfoCmd.Flags().BoolP("files", "f", false, "显示文件结构")

	// Flags for template search
	templateSearchCmd.Flags().StringP("category", "c", "", "按分类搜索")
}