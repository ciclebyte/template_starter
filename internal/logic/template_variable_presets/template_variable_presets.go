package template_variable_presets

import (
	"context"
	"encoding/json"
	"fmt"
	v1 "github.com/ciclebyte/template_starter/api/v1/template_variable_presets"
	"github.com/ciclebyte/template_starter/internal/dao"
	"github.com/ciclebyte/template_starter/internal/model/do"
	"github.com/ciclebyte/template_starter/internal/model/entity"
	"github.com/ciclebyte/template_starter/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sTemplateVariablePresets struct{}

func init() {
	service.RegisterTemplateVariablePresets(New())
}

func New() *sTemplateVariablePresets {
	return &sTemplateVariablePresets{}
}

// SubscribePreset 订阅预设变量到模板
func (s *sTemplateVariablePresets) SubscribePreset(ctx context.Context, req *v1.SubscribePresetReq) (res *v1.SubscribePresetRes, err error) {
	res = &v1.SubscribePresetRes{}
	
	// 使用事务处理
	err = dao.TemplateVariablePresets.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 批量插入新的订阅关系
		var dataList []*do.TemplateVariablePresets
		currentTime := gtime.Now()
		
		for _, preset := range req.Presets {
			data := &do.TemplateVariablePresets{
				TemplateId:  req.TemplateId,
				VarPresetId: preset.VarPresetId,
				PresetPath:  preset.PresetPath,
				MappedName:  preset.MappedName,
				IsActive:    preset.IsActive,
				Sort:        preset.Sort,
				CreatedAt:   currentTime,
				UpdatedAt:   currentTime,
			}
			dataList = append(dataList, data)
		}
		
		if len(dataList) > 0 {
			err = dao.TemplateVariablePresets.BatchInsert(ctx, dataList)
			if err != nil {
				return fmt.Errorf("批量插入订阅关系失败: %w", err)
			}
		}
		
		return nil
	})
	
	if err != nil {
		res.Success = false
		res.Message = err.Error()
		return res, err
	}
	
	res.Success = true
	res.Message = "订阅成功"
	return res, nil
}

// GetSubscribedPresets 获取模板订阅的预设变量列表
func (s *sTemplateVariablePresets) GetSubscribedPresets(ctx context.Context, req *v1.GetSubscribedPresetsReq) (res *v1.GetSubscribedPresetsRes, err error) {
	res = &v1.GetSubscribedPresetsRes{}
	
	// 查询模板订阅的预设变量关联关系
	subscriptions, err := dao.TemplateVariablePresets.GetByTemplateId(ctx, req.TemplateId)
	if err != nil {
		return nil, fmt.Errorf("查询订阅关系失败: %w", err)
	}
	
	var resultList []v1.SubscribedPresetItem
	
	for _, subscription := range subscriptions {
		// 查询预设变量详情
		var varPreset *entity.VarPreset
		err = dao.VarPreset.Ctx(ctx).Where("id", subscription.VarPresetId).Scan(&varPreset)
		if err != nil {
			g.Log().Error(ctx, "查询预设变量详情失败:", err)
			continue
		}
		
		if varPreset == nil {
			continue
		}
		
		// 解析预设变量内容，查找对应路径的变量信息
		var presetData map[string]interface{}
		contentField := varPreset.SchemaJson
		if varPreset.DefaultDataJson != nil && *varPreset.DefaultDataJson != "" {
			contentField = *varPreset.DefaultDataJson
		}
		err = json.Unmarshal([]byte(contentField), &presetData)
		if err != nil {
			g.Log().Error(ctx, "解析预设变量内容失败:", err)
			continue
		}
		
		// 类型断言
		presetPath, _ := subscription.PresetPath.(string)
		mappedName, _ := subscription.MappedName.(string)
		id, _ := subscription.Id.(int64)
		templateId, _ := subscription.TemplateId.(uint64)
		varPresetId, _ := subscription.VarPresetId.(uint64)
		isActive, _ := subscription.IsActive.(int)
		sort, _ := subscription.Sort.(int)
		
		// 根据路径查找变量信息
		varInfo := s.findVariableByPath(presetData, presetPath)
		
		item := v1.SubscribedPresetItem{
			Id:            id,
			TemplateId:    templateId,
			VarPresetId:   varPresetId,
			PresetName:    varPreset.Name,
			PresetPath:    presetPath,
			MappedName:    mappedName,
			IsActive:      isActive,
			Sort:          sort,
		}
		
		// 填充变量详细信息
		if varInfo != nil {
			if displayName, ok := varInfo["displayName"].(string); ok {
				item.DisplayName = displayName
			}
			if description, ok := varInfo["description"].(string); ok {
				item.Description = description
			}
			if example, ok := varInfo["example"].(string); ok {
				item.Example = example
			}
			if insertText, ok := varInfo["insertText"].(string); ok {
				item.InsertText = insertText
			}
			if category, ok := varInfo["category"].(string); ok {
				item.Category = category
			}
		}
		
		resultList = append(resultList, item)
	}
	
	res.Data = resultList
	return res, nil
}

// UnsubscribePreset 取消订阅预设变量
func (s *sTemplateVariablePresets) UnsubscribePreset(ctx context.Context, req *v1.UnsubscribePresetReq) (res *v1.UnsubscribePresetRes, err error) {
	res = &v1.UnsubscribePresetRes{}
	
	// 删除订阅关系
	err = dao.TemplateVariablePresets.Delete(ctx, req.Id)
	if err != nil {
		res.Success = false
		res.Message = "取消订阅失败: " + err.Error()
		return res, err
	}
	
	res.Success = true
	res.Message = "取消订阅成功"
	return res, nil
}

// UpdatePresetMapping 更新预设变量映射
func (s *sTemplateVariablePresets) UpdatePresetMapping(ctx context.Context, req *v1.UpdatePresetMappingReq) (res *v1.UpdatePresetMappingRes, err error) {
	res = &v1.UpdatePresetMappingRes{}
	
	// 更新映射关系
	updateData := &do.TemplateVariablePresets{
		MappedName: req.MappedName,
		IsActive:   req.IsActive,
		Sort:       req.Sort,
		UpdatedAt:  gtime.Now(),
	}
	
	err = dao.TemplateVariablePresets.Update(ctx, req.Id, updateData)
	if err != nil {
		res.Success = false
		res.Message = "更新映射失败: " + err.Error()
		return res, err
	}
	
	res.Success = true
	res.Message = "更新成功"
	return res, nil
}

// GetAvailablePresets 获取可用的预设变量列表
func (s *sTemplateVariablePresets) GetAvailablePresets(ctx context.Context, req *v1.GetAvailablePresetsReq) (res *v1.GetAvailablePresetsRes, err error) {
	res = &v1.GetAvailablePresetsRes{}
	
	// 构建查询条件
	query := dao.VarPreset.Ctx(ctx)
	
	if req.Keyword != "" {
		query = query.WhereLike("name", "%"+req.Keyword+"%")
	}
	
	// 分页查询
	total, err := query.Count()
	if err != nil {
		return nil, fmt.Errorf("查询总数失败: %w", err)
	}
	
	offset := (req.Page - 1) * req.Size
	var varPresets []*entity.VarPreset
	err = query.Limit(offset, req.Size).OrderDesc("updated_at").Scan(&varPresets)
	if err != nil {
		return nil, fmt.Errorf("查询预设变量列表失败: %w", err)
	}
	
	var resultList []v1.AvailablePresetItem
	
	for _, preset := range varPresets {
		description := ""
		if preset.Description != nil {
			description = *preset.Description
		}
		
		item := v1.AvailablePresetItem{
			Id:          preset.Id,
			Name:        preset.Name,
			Description: description,
			Variables:   []v1.AvailableVariableItem{},
		}
		
		// 解析预设变量内容
		var presetData map[string]interface{}
		contentField := preset.SchemaJson
		if preset.DefaultDataJson != nil && *preset.DefaultDataJson != "" {
			contentField = *preset.DefaultDataJson
		}
		err = json.Unmarshal([]byte(contentField), &presetData)
		if err != nil {
			g.Log().Error(ctx, "解析预设变量内容失败:", err)
			continue
		}
		
		// 提取所有变量
		variables := s.extractVariables(presetData, "")
		item.Variables = variables
		
		resultList = append(resultList, item)
	}
	
	res.Data = v1.AvailablePresetsData{
		List:  resultList,
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
	}
	
	return res, nil
}

// findVariableByPath 根据路径查找变量信息
func (s *sTemplateVariablePresets) findVariableByPath(data map[string]interface{}, path string) map[string]interface{} {
	if path == "" {
		return nil
	}
	
	// 直接查找根级变量
	if varInfo, ok := data[path].(map[string]interface{}); ok {
		return varInfo
	}
	
	// 递归查找嵌套变量
	return s.findVariableRecursively(data, path)
}

// findVariableRecursively 递归查找变量
func (s *sTemplateVariablePresets) findVariableRecursively(data map[string]interface{}, path string) map[string]interface{} {
	for key, value := range data {
		if valueMap, ok := value.(map[string]interface{}); ok {
			// 检查当前key是否匹配
			if key == path {
				return valueMap
			}
			
			// 如果有children，递归查找
			if children, ok := valueMap["children"].(map[string]interface{}); ok {
				if result := s.findVariableRecursively(children, path); result != nil {
					return result
				}
			}
		}
	}
	return nil
}

// extractVariables 提取预设变量中的所有变量
func (s *sTemplateVariablePresets) extractVariables(data map[string]interface{}, prefix string) []v1.AvailableVariableItem {
	var variables []v1.AvailableVariableItem
	
	for key, value := range data {
		currentPath := key
		if prefix != "" {
			currentPath = prefix + "." + key
		}
		
		if valueMap, ok := value.(map[string]interface{}); ok {
			// 检查是否是变量定义
			if s.isVariableDefinition(valueMap) {
				variable := v1.AvailableVariableItem{
					Path: currentPath,
				}
				
				if displayName, ok := valueMap["displayName"].(string); ok {
					variable.DisplayName = displayName
				}
				if description, ok := valueMap["description"].(string); ok {
					variable.Description = description
				}
				if example, ok := valueMap["example"].(string); ok {
					variable.Example = example
				}
				if insertText, ok := valueMap["insertText"].(string); ok {
					variable.InsertText = insertText
				}
				if category, ok := valueMap["category"].(string); ok {
					variable.Category = category
				}
				
				variables = append(variables, variable)
			}
			
			// 如果有children，递归提取
			if children, ok := valueMap["children"].(map[string]interface{}); ok {
				childVariables := s.extractVariables(children, currentPath)
				variables = append(variables, childVariables...)
			}
		}
	}
	
	return variables
}

// isVariableDefinition 判断是否是变量定义
func (s *sTemplateVariablePresets) isVariableDefinition(data map[string]interface{}) bool {
	// 简单判断：包含type字段且不是object类型，或者没有children字段
	if typeVal, ok := data["type"].(string); ok {
		return typeVal != "object"
	}
	
	// 如果没有children字段，认为是变量定义
	_, hasChildren := data["children"]
	return !hasChildren
}