package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/ciclebyte/template_starter/cli/internal/zipper"
)

// zipCmd represents the zip command
var zipCmd = &cobra.Command{
	Use:   "zip [source-path]",
	Short: "将项目打包为适合上传的zip文件",
	Long: `将项目目录打包为zip文件，自动排除不适合上传的文件和目录。

默认排除的文件和目录：
- 二进制文件
- .git 目录
- node_modules 目录
- 隐藏目录（以.开头）
- 日志文件（*.log）
- 临时文件
- IDE配置文件

示例:
  template-cli zip                    # 打包当前目录
  template-cli zip ./my-project       # 打包指定目录
  template-cli zip -o project.zip     # 指定输出文件名
  template-cli zip --include-hidden   # 包含隐藏文件
  template-cli zip --config .ziprc    # 使用自定义配置文件`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runZipCommand(cmd, args)
	},
}

// runZipCommand 执行zip命令
func runZipCommand(cmd *cobra.Command, args []string) error {
	// 获取源路径
	sourcePath := "."
	if len(args) > 0 {
		sourcePath = args[0]
	}

	// 获取参数
	outputFile, _ := cmd.Flags().GetString("output")
	includeHidden, _ := cmd.Flags().GetBool("include-hidden")
	includeBinary, _ := cmd.Flags().GetBool("include-binary")
	configFile, _ := cmd.Flags().GetString("config")
	excludePatterns, _ := cmd.Flags().GetStringSlice("exclude")
	includePatterns, _ := cmd.Flags().GetStringSlice("include")
	verbose, _ := cmd.Flags().GetBool("verbose")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	// 解析源路径
	absSourcePath, err := filepath.Abs(sourcePath)
	if err != nil {
		return fmt.Errorf("解析源路径失败: %w", err)
	}

	// 生成默认输出文件名
	if outputFile == "" {
		projectName := filepath.Base(absSourcePath)
		outputFile = projectName + ".zip"
	}

	// 解析输出文件路径
	absOutputFile, err := filepath.Abs(outputFile)
	if err != nil {
		return fmt.Errorf("解析输出文件路径失败: %w", err)
	}

	fmt.Printf("正在打包项目...\n")
	fmt.Printf("源路径: %s\n", absSourcePath)
	fmt.Printf("输出文件: %s\n", absOutputFile)

	// 创建压缩器配置
	config := &zipper.Config{
		SourcePath:      absSourcePath,
		OutputFile:      absOutputFile,
		IncludeHidden:   includeHidden,
		IncludeBinary:   includeBinary,
		ConfigFile:      configFile,
		ExcludePatterns: excludePatterns,
		IncludePatterns: includePatterns,
		Verbose:         verbose,
		DryRun:          dryRun,
	}

	// 创建压缩器
	z, err := zipper.NewZipper(config)
	if err != nil {
		return fmt.Errorf("创建压缩器失败: %w", err)
	}

	// 执行压缩
	result, err := z.Zip()
	if err != nil {
		return fmt.Errorf("压缩失败: %w", err)
	}

	// 显示结果
	if dryRun {
		fmt.Printf("\n预览模式 - 不会创建实际文件\n")
	} else {
		fmt.Printf("\n压缩完成！\n")
	}
	
	fmt.Printf("文件数量: %d\n", result.FileCount)
	fmt.Printf("目录数量: %d\n", result.DirCount)
	fmt.Printf("跳过文件: %d\n", result.SkippedCount)
	
	if !dryRun {
		fmt.Printf("压缩包大小: %s\n", formatSize(result.CompressedSize))
		fmt.Printf("原始大小: %s\n", formatSize(result.OriginalSize))
		fmt.Printf("压缩率: %.1f%%\n", float64(result.CompressedSize)*100/float64(result.OriginalSize))
	}

	if verbose && len(result.SkippedFiles) > 0 {
		fmt.Printf("\n跳过的文件:\n")
		for _, file := range result.SkippedFiles {
			fmt.Printf("  %s\n", file)
		}
	}

	return nil
}

// formatSize 格式化文件大小
func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

func init() {
	rootCmd.AddCommand(zipCmd)

	// 输出选项
	zipCmd.Flags().StringP("output", "o", "", "输出zip文件路径 (默认: 项目名.zip)")
	
	// 包含/排除选项
	zipCmd.Flags().BoolP("include-hidden", "H", false, "包含隐藏文件和目录")
	zipCmd.Flags().BoolP("include-binary", "B", false, "包含二进制文件")
	zipCmd.Flags().StringSliceP("exclude", "e", []string{}, "额外排除的文件模式")
	zipCmd.Flags().StringSliceP("include", "i", []string{}, "强制包含的文件模式")
	
	// 配置选项
	zipCmd.Flags().StringP("config", "c", "", "自定义配置文件路径")
	
	// 行为选项
	zipCmd.Flags().BoolP("verbose", "v", false, "显示详细信息")
	zipCmd.Flags().BoolP("dry-run", "n", false, "预览模式，不创建实际文件")
}