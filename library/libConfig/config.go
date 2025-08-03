package libConfig

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/ciclebyte/template_starter/internal/model"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

// ConfigManager 配置管理器
type ConfigManager struct {
	cache      map[string]*ConfigItem
	cacheMutex sync.RWMutex
	cacheTTL   time.Duration
}

// ConfigItem 配置项
type ConfigItem struct {
	Value     interface{}
	Type      string
	ExpiresAt time.Time
}

var (
	defaultManager *ConfigManager
	once           sync.Once
)

// GetManager 获取默认配置管理器实例
func GetManager() *ConfigManager {
	once.Do(func() {
		defaultManager = &ConfigManager{
			cache:    make(map[string]*ConfigItem),
			cacheTTL: 5 * time.Minute, // 默认5分钟缓存
		}
	})
	return defaultManager
}

// GetString 获取字符串配置
func GetString(ctx context.Context, key string, defaultValue ...string) string {
	manager := GetManager()
	value, err := manager.Get(ctx, key, "string")
	if err != nil || value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return ""
	}
	return gconv.String(value)
}

// GetInt 获取整数配置
func GetInt(ctx context.Context, key string, defaultValue ...int) int {
	manager := GetManager()
	value, err := manager.Get(ctx, key, "number")
	if err != nil || value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return gconv.Int(value)
}

// GetBool 获取布尔配置
func GetBool(ctx context.Context, key string, defaultValue ...bool) bool {
	manager := GetManager()
	value, err := manager.Get(ctx, key, "boolean")
	if err != nil || value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return false
	}
	return gconv.Bool(value)
}

// GetJson 获取JSON配置并解析到指定结构
func GetJson(ctx context.Context, key string, result interface{}) error {
	manager := GetManager()
	value, err := manager.Get(ctx, key, "json")
	if err != nil {
		return err
	}
	
	if value == nil {
		return fmt.Errorf("config %s not found", key)
	}
	
	jsonStr := gconv.String(value)
	return json.Unmarshal([]byte(jsonStr), result)
}

// GetArray 获取数组配置
func GetArray(ctx context.Context, key string) ([]interface{}, error) {
	manager := GetManager()
	value, err := manager.Get(ctx, key, "array")
	if err != nil {
		return nil, err
	}
	
	if value == nil {
		return []interface{}{}, nil
	}
	
	jsonStr := gconv.String(value)
	var result []interface{}
	err = json.Unmarshal([]byte(jsonStr), &result)
	return result, err
}

// GetAIConfig 获取AI相关配置
func GetAIConfig(ctx context.Context) (*model.AIConfig, error) {
	config := &model.AIConfig{}
	
	// 基础配置
	config.Enabled = GetBool(ctx, "ai.enabled", false)
	config.Provider = GetString(ctx, "ai.provider", "openai")
	
	// OpenAI配置
	config.OpenAI.APIKey = GetString(ctx, "ai.openai.api_key")
	config.OpenAI.BaseURL = GetString(ctx, "ai.openai.base_url", "https://api.openai.com")
	config.OpenAI.Model = GetString(ctx, "ai.openai.model", "gpt-4")
	config.OpenAI.MaxTokens = GetInt(ctx, "ai.openai.max_tokens", 4000)
	
	// 解析温度参数
	tempStr := GetString(ctx, "ai.openai.temperature", "0.7")
	if temp, err := strconv.ParseFloat(tempStr, 64); err == nil {
		config.OpenAI.Temperature = temp
	}
	
	// Claude配置
	config.Claude.APIKey = GetString(ctx, "ai.claude.api_key")
	config.Claude.BaseURL = GetString(ctx, "ai.claude.base_url", "https://api.anthropic.com")
	config.Claude.Model = GetString(ctx, "ai.claude.model", "claude-3-sonnet-20240229")
	
	// 功能开关
	config.Features.TemplateGeneration = GetBool(ctx, "ai.template_generation.enabled", true)
	config.Features.CodeReview = GetBool(ctx, "ai.code_review.enabled", false)
	config.Features.VariableSuggestion = GetBool(ctx, "ai.variable_suggestion.enabled", true)
	config.Features.AutoDocumentation = GetBool(ctx, "ai.auto_documentation.enabled", false)
	
	// 使用限制
	config.RateLimit.RequestsPerHour = GetInt(ctx, "ai.rate_limit.requests_per_hour", 100)
	config.RateLimit.RequestsPerDay = GetInt(ctx, "ai.rate_limit.requests_per_day", 1000)
	config.MaxFileSize = GetInt(ctx, "ai.max_file_size", 50)
	
	// 支持的语言
	if err := GetJson(ctx, "ai.supported_languages", &config.SupportedLanguages); err != nil {
		// 默认支持的语言
		config.SupportedLanguages = []string{"javascript", "typescript", "python", "java", "go", "vue", "react"}
	}
	
	return config, nil
}

// Set 设置配置值
func Set(ctx context.Context, key string, value interface{}) error {
	manager := GetManager()
	return manager.Set(ctx, key, value)
}

// Get 获取配置值（内部方法）
func (m *ConfigManager) Get(ctx context.Context, key string, expectedType string) (interface{}, error) {
	// 先从缓存获取
	m.cacheMutex.RLock()
	if item, exists := m.cache[key]; exists && time.Now().Before(item.ExpiresAt) {
		m.cacheMutex.RUnlock()
		return item.Value, nil
	}
	m.cacheMutex.RUnlock()
	
	// 从数据库获取
	var config model.SystemConfigInfo
	err := g.DB().Model("system_config").
		Where("config_key = ? AND status = 1", key).
		Scan(&config)
	
	if err != nil {
		g.Log().Error(ctx, "Failed to get config:", key, err)
		return nil, err
	}
	
	if config.Id == 0 {
		// 配置不存在，返回默认值
		return nil, nil
	}
	
	// 类型转换
	var value interface{}
	switch config.ConfigType {
	case "string":
		value = config.ConfigValue
	case "number":
		if val, err := strconv.Atoi(config.ConfigValue); err == nil {
			value = val
		} else {
			value = 0
		}
	case "boolean":
		value = config.ConfigValue == "true" || config.ConfigValue == "1"
	case "json", "array":
		value = config.ConfigValue
	default:
		value = config.ConfigValue
	}
	
	// 缓存结果
	m.cacheMutex.Lock()
	m.cache[key] = &ConfigItem{
		Value:     value,
		Type:      config.ConfigType,
		ExpiresAt: time.Now().Add(m.cacheTTL),
	}
	m.cacheMutex.Unlock()
	
	return value, nil
}

// Set 设置配置值（内部方法）
func (m *ConfigManager) Set(ctx context.Context, key string, value interface{}) error {
	// 转换值为字符串
	var valueStr string
	switch v := value.(type) {
	case string:
		valueStr = v
	case bool:
		if v {
			valueStr = "true"
		} else {
			valueStr = "false"
		}
	case int, int64, float64:
		valueStr = gconv.String(v)
	default:
		// 复杂类型转JSON
		if jsonBytes, err := json.Marshal(v); err == nil {
			valueStr = string(jsonBytes)
		} else {
			return err
		}
	}
	
	// 更新数据库
	_, err := g.DB().Model("system_config").
		Where("config_key = ?", key).
		Update(g.Map{"config_value": valueStr, "updated_at": time.Now()})
	
	if err != nil {
		g.Log().Error(ctx, "Failed to set config:", key, err)
		return err
	}
	
	// 清除缓存
	m.cacheMutex.Lock()
	delete(m.cache, key)
	m.cacheMutex.Unlock()
	
	return nil
}

// ClearCache 清除所有缓存
func (m *ConfigManager) ClearCache() {
	m.cacheMutex.Lock()
	m.cache = make(map[string]*ConfigItem)
	m.cacheMutex.Unlock()
}

// GetGroupConfigs 获取分组配置
func GetGroupConfigs(ctx context.Context, group string, includePrivate bool) (map[string]interface{}, error) {
	var configs []model.SystemConfigInfo
	query := g.DB().Model("system_config").
		Where("config_group = ? AND status = 1", group)
	
	if !includePrivate {
		query = query.Where("is_public = 1")
	}
	
	err := query.OrderAsc("sort_order").Scan(&configs)
	if err != nil {
		return nil, err
	}
	
	result := make(map[string]interface{})
	for _, config := range configs {
		var value interface{}
		switch config.ConfigType {
		case "string":
			value = config.ConfigValue
		case "number":
			value = gconv.Int(config.ConfigValue)
		case "boolean":
			value = gconv.Bool(config.ConfigValue)
		case "json", "array":
			var jsonValue interface{}
			if err := json.Unmarshal([]byte(config.ConfigValue), &jsonValue); err == nil {
				value = jsonValue
			} else {
				value = config.ConfigValue
			}
		default:
			value = config.ConfigValue
		}
		result[config.ConfigKey] = value
	}
	
	return result, nil
}