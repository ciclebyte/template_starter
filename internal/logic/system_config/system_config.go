package system_config

import (
	"context"
	"encoding/json"
	"fmt"

	systemConfigApi "github.com/ciclebyte/template_starter/api/v1/system_config"
	"github.com/ciclebyte/template_starter/internal/dao"
	"github.com/ciclebyte/template_starter/internal/model"
	"github.com/ciclebyte/template_starter/internal/service"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type sSystemConfig struct{}

func init() {
	service.RegisterSystemConfigService(New())
}

func New() *sSystemConfig {
	return &sSystemConfig{}
}

// Add 新增配置
func (s *sSystemConfig) Add(ctx context.Context, req *systemConfigApi.SystemConfigAddReq) (*systemConfigApi.SystemConfigAddRes, error) {
	g.Log().Debug(ctx, "SystemConfig.Add called with req:", req)

	// 检查配置键是否已存在
	count, err := dao.SystemConfig.Ctx(ctx).Where("config_key", req.ConfigKey).Count()
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, fmt.Errorf("配置键 %s 已存在", req.ConfigKey)
	}

	// 插入配置
	_, err = dao.SystemConfig.Ctx(ctx).Insert(req)
	if err != nil {
		return nil, err
	}

	return &systemConfigApi.SystemConfigAddRes{}, nil
}

// List 配置列表
func (s *sSystemConfig) List(ctx context.Context, req *systemConfigApi.SystemConfigListReq) (*systemConfigApi.SystemConfigListRes, error) {
	g.Log().Debug(ctx, "SystemConfig.List called with req:", req)

	res := &systemConfigApi.SystemConfigListRes{}
	
	// 构建查询条件
	query := dao.SystemConfig.Ctx(ctx)
	
	if req.ConfigGroup != "" {
		query = query.Where("config_group", req.ConfigGroup)
	}
	
	if req.IsPublic != nil {
		query = query.Where("is_public", *req.IsPublic)
	}
	
	if req.Status != nil {
		query = query.Where("status", *req.Status)
	}

	// 获取总数
	total, err := query.Count()
	if err != nil {
		return nil, err
	}
	res.Total = total

	// 分页查询
	var configs []*model.SystemConfigInfo
	err = query.
		OrderAsc("config_group").
		OrderAsc("sort_order").
		OrderAsc("id").
		Limit(req.PageReq.PageSize).
		Offset((req.PageReq.PageNum - 1) * req.PageReq.PageSize).
		Scan(&configs)
	if err != nil {
		return nil, err
	}

	res.SystemConfigList = configs
	return res, nil
}

// Detail 配置详情
func (s *sSystemConfig) Detail(ctx context.Context, req *systemConfigApi.SystemConfigDetailReq) (*systemConfigApi.SystemConfigDetailRes, error) {
	g.Log().Debug(ctx, "SystemConfig.Detail called with req:", req)

	var config *model.SystemConfigInfo
	err := dao.SystemConfig.Ctx(ctx).Where("id", req.Id).Scan(&config)
	if err != nil {
		return nil, err
	}

	if config == nil {
		return nil, fmt.Errorf("配置不存在")
	}

	return &systemConfigApi.SystemConfigDetailRes{
		SystemConfigInfo: config,
	}, nil
}

// Edit 编辑配置
func (s *sSystemConfig) Edit(ctx context.Context, req *systemConfigApi.SystemConfigEditReq) (*systemConfigApi.SystemConfigEditRes, error) {
	g.Log().Debug(ctx, "SystemConfig.Edit called with req:", req)

	// 检查配置是否存在
	count, err := dao.SystemConfig.Ctx(ctx).Where("id", req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, fmt.Errorf("配置不存在")
	}

	// 更新配置
	_, err = dao.SystemConfig.Ctx(ctx).Where("id", req.Id).Update(req)
	if err != nil {
		return nil, err
	}

	return &systemConfigApi.SystemConfigEditRes{}, nil
}

// Del 删除配置
func (s *sSystemConfig) Del(ctx context.Context, req *systemConfigApi.SystemConfigDelReq) (*systemConfigApi.SystemConfigDelRes, error) {
	g.Log().Debug(ctx, "SystemConfig.Del called with req:", req)

	_, err := dao.SystemConfig.Ctx(ctx).Where("id", req.Id).Delete()
	if err != nil {
		return nil, err
	}

	return &systemConfigApi.SystemConfigDelRes{}, nil
}

// BatchDel 批量删除配置
func (s *sSystemConfig) BatchDel(ctx context.Context, req *systemConfigApi.SystemConfigBatchDelReq) (*systemConfigApi.SystemConfigBatchDelRes, error) {
	g.Log().Debug(ctx, "SystemConfig.BatchDel called with req:", req)

	_, err := dao.SystemConfig.Ctx(ctx).WhereIn("id", req.Ids).Delete()
	if err != nil {
		return nil, err
	}

	return &systemConfigApi.SystemConfigBatchDelRes{}, nil
}

// GetByKey 根据键名获取配置
func (s *sSystemConfig) GetByKey(ctx context.Context, req *systemConfigApi.SystemConfigGetByKeyReq) (*systemConfigApi.SystemConfigGetByKeyRes, error) {
	g.Log().Debug(ctx, "SystemConfig.GetByKey called with req:", req)

	var config model.SystemConfigInfo
	err := dao.SystemConfig.Ctx(ctx).
		Where("config_key", req.ConfigKey).
		Where("status", 1).
		Scan(&config)
	if err != nil {
		return nil, err
	}

	if config.Id == 0 {
		return nil, fmt.Errorf("配置 %s 不存在", req.ConfigKey)
	}

	return &systemConfigApi.SystemConfigGetByKeyRes{
		ConfigKey:   config.ConfigKey,
		ConfigValue: config.ConfigValue,
		ConfigType:  config.ConfigType,
	}, nil
}

// GetByGroup 根据分组获取配置
func (s *sSystemConfig) GetByGroup(ctx context.Context, req *systemConfigApi.SystemConfigGetByGroupReq) (*systemConfigApi.SystemConfigGetByGroupRes, error) {
	g.Log().Debug(ctx, "SystemConfig.GetByGroup called with req:", req)

	// 构建查询条件
	query := dao.SystemConfig.Ctx(ctx).
		Where("config_group", req.ConfigGroup).
		Where("status", 1)

	if req.IsPublic != nil {
		query = query.Where("is_public", *req.IsPublic)
	}

	// 获取配置列表
	var configs []model.SystemConfigInfo
	err := query.OrderAsc("sort_order").Scan(&configs)
	if err != nil {
		return nil, err
	}

	// 转换为键值对
	configMap := make(map[string]interface{})
	for _, config := range configs {
		value := s.parseConfigValue(config.ConfigValue, config.ConfigType)
		configMap[config.ConfigKey] = value
	}

	return &systemConfigApi.SystemConfigGetByGroupRes{
		Configs: configMap,
		GroupInfo: struct {
			Group       string `json:"group"`
			DisplayName string `json:"displayName"`
			Description string `json:"description"`
		}{
			Group:       req.ConfigGroup,
			DisplayName: s.getGroupDisplayName(req.ConfigGroup),
			Description: s.getGroupDescription(req.ConfigGroup),
		},
	}, nil
}

// BatchUpdate 批量更新配置
func (s *sSystemConfig) BatchUpdate(ctx context.Context, req *systemConfigApi.SystemConfigBatchUpdateReq) (*systemConfigApi.SystemConfigBatchUpdateRes, error) {
	g.Log().Debug(ctx, "SystemConfig.BatchUpdate called with req:", req)

	res := &systemConfigApi.SystemConfigBatchUpdateRes{
		UpdatedCount: 0,
		FailedKeys:   []string{},
	}

	for _, config := range req.Configs {
		_, err := dao.SystemConfig.Ctx(ctx).
			Where("config_key", config.ConfigKey).
			Where("status", 1).
			Update(g.Map{"config_value": config.ConfigValue})
		
		if err != nil {
			res.FailedKeys = append(res.FailedKeys, config.ConfigKey)
			g.Log().Warning(ctx, "Failed to update config:", config.ConfigKey, err)
		} else {
			res.UpdatedCount++
		}
	}

	return res, nil
}

// GetPublic 获取公开配置
func (s *sSystemConfig) GetPublic(ctx context.Context, req *systemConfigApi.SystemConfigPublicReq) (*systemConfigApi.SystemConfigPublicRes, error) {
	g.Log().Debug(ctx, "SystemConfig.GetPublic called with req:", req)

	query := dao.SystemConfig.Ctx(ctx).
		Where("is_public", 1).
		Where("status", 1)

	if req.ConfigGroup != "" {
		query = query.Where("config_group", req.ConfigGroup)
	}

	var configs []model.SystemConfigInfo
	err := query.OrderAsc("config_group").OrderAsc("sort_order").Scan(&configs)
	if err != nil {
		return nil, err
	}

	// 转换为键值对
	configMap := make(map[string]interface{})
	for _, config := range configs {
		value := s.parseConfigValue(config.ConfigValue, config.ConfigType)
		configMap[config.ConfigKey] = value
	}

	return &systemConfigApi.SystemConfigPublicRes{
		Configs: configMap,
	}, nil
}

// Reset 重置配置
func (s *sSystemConfig) Reset(ctx context.Context, req *systemConfigApi.SystemConfigResetReq) (*systemConfigApi.SystemConfigResetRes, error) {
	g.Log().Debug(ctx, "SystemConfig.Reset called with req:", req)

	var config model.SystemConfigInfo
	err := dao.SystemConfig.Ctx(ctx).
		Where("config_key", req.ConfigKey).
		Scan(&config)
	if err != nil {
		return nil, err
	}

	if config.Id == 0 {
		return nil, fmt.Errorf("配置 %s 不存在", req.ConfigKey)
	}

	// 重置为默认值
	_, err = dao.SystemConfig.Ctx(ctx).
		Where("config_key", req.ConfigKey).
		Update(g.Map{"config_value": config.DefaultValue})
	if err != nil {
		return nil, err
	}

	return &systemConfigApi.SystemConfigResetRes{}, nil
}

// Validate 验证配置值
func (s *sSystemConfig) Validate(ctx context.Context, req *systemConfigApi.SystemConfigValidateReq) (*systemConfigApi.SystemConfigValidateRes, error) {
	g.Log().Debug(ctx, "SystemConfig.Validate called with req:", req)

	var config model.SystemConfigInfo
	err := dao.SystemConfig.Ctx(ctx).
		Where("config_key", req.ConfigKey).
		Scan(&config)
	if err != nil {
		return nil, err
	}

	if config.Id == 0 {
		return &systemConfigApi.SystemConfigValidateRes{
			IsValid: false,
			Message: fmt.Sprintf("配置 %s 不存在", req.ConfigKey),
		}, nil
	}

	// 基础类型验证
	isValid, message := s.validateConfigValue(req.ConfigValue, config.ConfigType, config.ValidationRule)

	return &systemConfigApi.SystemConfigValidateRes{
		IsValid: isValid,
		Message: message,
	}, nil
}

// 便捷方法实现
func (s *sSystemConfig) GetConfigValue(ctx context.Context, configKey string) (interface{}, error) {
	var config model.SystemConfigInfo
	err := dao.SystemConfig.Ctx(ctx).
		Where("config_key", configKey).
		Where("status", 1).
		Scan(&config)
	if err != nil {
		return nil, err
	}

	if config.Id == 0 {
		return nil, nil
	}

	return s.parseConfigValue(config.ConfigValue, config.ConfigType), nil
}

func (s *sSystemConfig) GetConfigString(ctx context.Context, configKey string, defaultValue ...string) (string, error) {
	value, err := s.GetConfigValue(ctx, configKey)
	if err != nil || value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0], nil
		}
		return "", err
	}
	return gconv.String(value), nil
}

func (s *sSystemConfig) GetConfigInt(ctx context.Context, configKey string, defaultValue ...int) (int, error) {
	value, err := s.GetConfigValue(ctx, configKey)
	if err != nil || value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0], nil
		}
		return 0, err
	}
	return gconv.Int(value), nil
}

func (s *sSystemConfig) GetConfigBool(ctx context.Context, configKey string, defaultValue ...bool) (bool, error) {
	value, err := s.GetConfigValue(ctx, configKey)
	if err != nil || value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0], nil
		}
		return false, err
	}
	return gconv.Bool(value), nil
}

func (s *sSystemConfig) GetConfigJson(ctx context.Context, configKey string, result interface{}) error {
	value, err := s.GetConfigValue(ctx, configKey)
	if err != nil {
		return err
	}
	
	if value == nil {
		return fmt.Errorf("配置 %s 不存在", configKey)
	}
	
	jsonStr := gconv.String(value)
	return json.Unmarshal([]byte(jsonStr), result)
}

func (s *sSystemConfig) SetConfigValue(ctx context.Context, configKey string, value interface{}) error {
	valueStr := gconv.String(value)
	if v, ok := value.(map[string]interface{}); ok {
		if jsonBytes, err := json.Marshal(v); err == nil {
			valueStr = string(jsonBytes)
		}
	} else if v, ok := value.([]interface{}); ok {
		if jsonBytes, err := json.Marshal(v); err == nil {
			valueStr = string(jsonBytes)
		}
	}

	_, err := dao.SystemConfig.Ctx(ctx).
		Where("config_key", configKey).
		Update(g.Map{"config_value": valueStr})
	return err
}

// 辅助方法
func (s *sSystemConfig) parseConfigValue(value, configType string) interface{} {
	switch configType {
	case "number":
		return gconv.Int(value)
	case "boolean":
		return gconv.Bool(value)
	case "json", "array":
		var result interface{}
		if err := json.Unmarshal([]byte(value), &result); err == nil {
			return result
		}
		return value
	default:
		return value
	}
}

func (s *sSystemConfig) validateConfigValue(value, configType, validationRule string) (bool, string) {
	// 基础类型验证
	switch configType {
	case "number":
		if gconv.Int(value) == 0 && value != "0" {
			return false, "值必须是数字"
		}
	case "boolean":
		if value != "true" && value != "false" && value != "1" && value != "0" {
			return false, "值必须是布尔类型"
		}
	case "json", "array":
		var temp interface{}
		if err := json.Unmarshal([]byte(value), &temp); err != nil {
			return false, "值必须是有效的JSON格式"
		}
	}

	// 自定义验证规则
	if validationRule != "" {
		var rule model.ConfigValidationRule
		if err := json.Unmarshal([]byte(validationRule), &rule); err == nil {
			if rule.Required && value == "" {
				return false, "值不能为空"
			}
			// 可以添加更多验证逻辑
		}
	}

	return true, "验证通过"
}

func (s *sSystemConfig) getGroupDisplayName(group string) string {
	groupNames := map[string]string{
		"ai":       "AI配置",
		"system":   "系统配置", 
		"template": "模板配置",
		"ui":       "界面配置",
	}
	if name, exists := groupNames[group]; exists {
		return name
	}
	return group
}

func (s *sSystemConfig) getGroupDescription(group string) string {
	groupDesc := map[string]string{
		"ai":       "人工智能相关功能配置",
		"system":   "系统基础信息配置",
		"template": "模板生成相关配置",
		"ui":       "用户界面和体验配置",
	}
	if desc, exists := groupDesc[group]; exists {
		return desc
	}
	return ""
}