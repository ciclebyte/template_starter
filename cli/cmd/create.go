package cmd

import (
	"fmt"

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
  template-cli create frontend --template vue3-admin --interactive`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		template, _ := cmd.Flags().GetString("template")
		output, _ := cmd.Flags().GetString("output")
		interactive, _ := cmd.Flags().GetBool("interactive")
		
		fmt.Printf("创建项目: %s\n", projectName)
		fmt.Printf("使用模板: %s\n", template)
		fmt.Printf("输出目录: %s\n", output)
		fmt.Printf("交互模式: %t\n", interactive)
		
		// TODO: 实现项目创建逻辑
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// 必需的模板标识
	createCmd.Flags().StringP("template", "t", "", "模板名称或ID (必需)")
	createCmd.MarkFlagRequired("template")
	
	// 可选参数
	createCmd.Flags().StringP("output", "o", ".", "输出目录")
	createCmd.Flags().BoolP("interactive", "i", false, "启用交互式变量配置")
	createCmd.Flags().StringP("config", "c", "", "变量配置文件路径")
	createCmd.Flags().BoolP("force", "f", false, "强制覆盖已存在的目录")
}