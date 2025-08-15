package template_variable_presets

import (
	"context"
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

		for _, presetId := range req.PresetIds {
			data := &do.TemplateVariablePresets{
				TemplateId: req.TemplateId,
				PresetId:   presetId,
				CreatedAt:  currentTime,
				UpdatedAt:  currentTime,
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
		// 类型断言
		id, _ := subscription.Id.(uint64)
		templateId, _ := subscription.TemplateId.(uint64)
		presetId, _ := subscription.PresetId.(uint64)

		// 查询预设变量详情
		var varPreset *entity.VarPreset
		err = dao.VarPreset.Ctx(ctx).Where("id", presetId).Scan(&varPreset)
		if err != nil {
			g.Log().Error(ctx, "查询预设变量详情失败:", err)
			continue
		}

		if varPreset == nil {
			continue
		}

		description := ""
		if varPreset.Description != nil {
			description = *varPreset.Description
		}

		item := v1.SubscribedPresetItem{
			Id:          id,
			TemplateId:  templateId,
			PresetId:    presetId,
			PresetName:  varPreset.Name,
			Description: description,
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
	err = dao.TemplateVariablePresets.Delete(ctx, int64(req.Id))
	if err != nil {
		res.Success = false
		res.Message = "取消订阅失败: " + err.Error()
		return res, err
	}

	res.Success = true
	res.Message = "取消订阅成功"
	return res, nil
}

// GetAvailablePresets 获取可用的预设变量列表
func (s *sTemplateVariablePresets) GetAvailablePresets(ctx context.Context, req *v1.GetAvailablePresetsReq) (res *v1.GetAvailablePresetsRes, err error) {
	res = &v1.GetAvailablePresetsRes{}

	g.Log().Info(ctx, "GetAvailablePresets请求参数:", "pageNum=", req.PageNum, "pageSize=", req.PageSize, "keyword=", req.Keyword)
	
	// 设置默认值
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

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

	offset := (req.PageNum - 1) * req.PageSize
	var varPresets []*entity.VarPreset
	err = query.Limit(offset, req.PageSize).OrderDesc("updated_at").Scan(&varPresets)
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
		}

		resultList = append(resultList, item)
	}

	res.List = resultList
	res.Total = total
	res.PageNum = req.PageNum
	res.PageSize = req.PageSize

	return res, nil
}
