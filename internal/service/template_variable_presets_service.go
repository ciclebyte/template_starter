package service

import (
	"context"
	v1 "github.com/ciclebyte/template_starter/api/v1/template_variable_presets"
)

type ITemplateVariablePresets interface {
	// SubscribePreset 订阅预设变量到模板
	SubscribePreset(ctx context.Context, req *v1.SubscribePresetReq) (res *v1.SubscribePresetRes, err error)
	
	// GetSubscribedPresets 获取模板订阅的预设变量列表
	GetSubscribedPresets(ctx context.Context, req *v1.GetSubscribedPresetsReq) (res *v1.GetSubscribedPresetsRes, err error)
	
	// UnsubscribePreset 取消订阅预设变量
	UnsubscribePreset(ctx context.Context, req *v1.UnsubscribePresetReq) (res *v1.UnsubscribePresetRes, err error)
	
	// GetAvailablePresets 获取可用的预设变量列表
	GetAvailablePresets(ctx context.Context, req *v1.GetAvailablePresetsReq) (res *v1.GetAvailablePresetsRes, err error)
}

var localTemplateVariablePresets ITemplateVariablePresets

func TemplateVariablePresets() ITemplateVariablePresets {
	if localTemplateVariablePresets == nil {
		panic("implement not found for interface ITemplateVariablePresets, forgot register?")
	}
	return localTemplateVariablePresets
}

func RegisterTemplateVariablePresets(i ITemplateVariablePresets) {
	localTemplateVariablePresets = i
}