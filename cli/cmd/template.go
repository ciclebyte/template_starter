package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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
	Run: func(cmd *cobra.Command, args []string) {
		category, _ := cmd.Flags().GetString("category")
		
		fmt.Printf("列出远程模板 (分类: %s)\n", category)
		
		// TODO: 实现模板列表逻辑
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
	Run: func(cmd *cobra.Command, args []string) {
		templateName := args[0]
		showVariables, _ := cmd.Flags().GetBool("variables")
		showFiles, _ := cmd.Flags().GetBool("files")
		
		fmt.Printf("模板信息: %s\n", templateName)
		fmt.Printf("显示变量: %t\n", showVariables)
		fmt.Printf("显示文件: %t\n", showFiles)
		
		// TODO: 实现模板信息显示逻辑
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
	Run: func(cmd *cobra.Command, args []string) {
		keyword := args[0]
		category, _ := cmd.Flags().GetString("category")
		
		fmt.Printf("搜索模板: %s (分类: %s)\n", keyword, category)
		
		// TODO: 实现模板搜索逻辑
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