package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "template-cli",
	Short: "模板生成器CLI工具",
	Long: `Template Starter CLI是一个轻量级的命令行工具，用于快速创建项目模板。

特性:
  • 从远程模板库创建项目
  • 支持交互式变量配置
  • 条件文件生成
  • 轻量级远程API客户端`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ciclebyte/template_starter/config/config.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in ciclebyte directory
		configDir := filepath.Join(home, ".ciclebyte", "template_starter", "config")
		
		// 确保配置目录存在
		if err := os.MkdirAll(configDir, 0755); err != nil {
			cobra.CheckErr(err)
		}
		
		viper.AddConfigPath(configDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// 设置默认值
	setDefaultConfig()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		// 如果配置文件不存在，创建默认配置文件
		if _, ok := err.(viper.ConfigFileNotFoundError); ok && cfgFile == "" {
			// 只在使用默认配置路径时自动创建配置文件
			home, _ := os.UserHomeDir()
			configFile := filepath.Join(home, ".ciclebyte", "template_starter", "config", "config.yaml")
			if err := viper.WriteConfigAs(configFile); err == nil {
				fmt.Fprintln(os.Stderr, "Created default config file:", configFile)
			}
		}
	} else {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

// setDefaultConfig 设置默认配置值
func setDefaultConfig() {
	// 服务器默认配置
	viper.SetDefault("server.url", "http://localhost:8000")
	viper.SetDefault("server.api_key", "")
	
	// 用户默认配置
	viper.SetDefault("user.author", "")
	viper.SetDefault("user.email", "")
}