package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// 这些变量在构建时通过 ldflags 设置
	Version   = "dev"
	GitCommit = "unknown"
	BuildTime = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Long: `显示Template Starter CLI的版本信息。

包括版本号、Git提交哈希和构建时间。`,
	Run: func(cmd *cobra.Command, args []string) {
		short, _ := cmd.Flags().GetBool("short")
		
		if short {
			fmt.Println(Version)
		} else {
			fmt.Printf("Template Starter CLI\n")
			fmt.Printf("版本:     %s\n", Version)
			fmt.Printf("Git提交:  %s\n", GitCommit)
			fmt.Printf("构建时间: %s\n", BuildTime)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Flags
	versionCmd.Flags().BoolP("short", "s", false, "只显示版本号")
}