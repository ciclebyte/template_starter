package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config CLI配置结构
type Config struct {
	Server ServerConfig `mapstructure:"server"`
	User   UserConfig   `mapstructure:"user"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	URL    string `mapstructure:"url"`
	APIKey string `mapstructure:"api_key"`
}

// UserConfig 用户配置
type UserConfig struct {
	Author string `mapstructure:"author"`
	Email  string `mapstructure:"email"`
}

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
	var config Config
	
	// 确保配置目录存在
	configDir, err := GetConfigDir()
	if err != nil {
		return nil, err
	}
	
	
	// 设置默认值
	setDefaults()
	
	// 读取配置
	if err := viper.ReadInConfig(); err != nil {
		// 如果配置文件不存在，创建默认配置文件
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			configFile := filepath.Join(configDir, "config.yaml")
			if err := viper.WriteConfigAs(configFile); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	
	// 解析配置到结构体
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	
	return &config, nil
}

// SaveConfig 保存配置
func SaveConfig(config *Config) error {
	// 设置配置值
	viper.Set("server.url", config.Server.URL)
	viper.Set("server.api_key", config.Server.APIKey)
	viper.Set("user.author", config.User.Author)
	viper.Set("user.email", config.User.Email)
	
	// 保存到文件
	return viper.WriteConfig()
}

// SetConfigValue 设置配置值
func SetConfigValue(key string, value interface{}) error {
	viper.Set(key, value)
	return viper.WriteConfig()
}

// GetConfigValue 获取配置值
func GetConfigValue(key string) interface{} {
	return viper.Get(key)
}

// ResetConfig 重置配置为默认值
func ResetConfig() error {
	// 清除所有设置
	for _, key := range viper.AllKeys() {
		viper.Set(key, nil)
	}
	
	// 重新设置默认值
	setDefaults()
	
	// 保存配置
	return viper.WriteConfig()
}

// setDefaults 设置默认配置值
func setDefaults() {
	// 服务器默认配置
	viper.SetDefault("server.url", "http://localhost:8000")
	viper.SetDefault("server.api_key", "")
	
	// 用户默认配置
	viper.SetDefault("user.author", "")
	viper.SetDefault("user.email", "")
}

// GetConfigDir 获取配置目录
func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	
	configDir := filepath.Join(homeDir, ".ciclebyte", "template_starter", "config")
	
	// 确保配置目录存在
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}
	
	return configDir, nil
}


// GetServerURL 获取服务器URL
func GetServerURL() string {
	return viper.GetString("server.url")
}

// GetAPIKey 获取API密钥
func GetAPIKey() string {
	return viper.GetString("server.api_key")
}