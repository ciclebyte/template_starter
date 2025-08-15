package controller

import (
	"context"
	v1 "github.com/ciclebyte/template_starter/api/v1/template_variable_presets"
	"github.com/ciclebyte/template_starter/internal/service"
)

var TemplateVariablePresets = cTemplateVariablePresets{}

type cTemplateVariablePresets struct{}

// SubscribePreset 订阅预设变量到模板
func (c *cTemplateVariablePresets) SubscribePreset(ctx context.Context, req *v1.SubscribePresetReq) (res *v1.SubscribePresetRes, err error) {
	return service.TemplateVariablePresets().SubscribePreset(ctx, req)
}

// GetSubscribedPresets 获取模板订阅的预设变量列表
func (c *cTemplateVariablePresets) GetSubscribedPresets(ctx context.Context, req *v1.GetSubscribedPresetsReq) (res *v1.GetSubscribedPresetsRes, err error) {
	return service.TemplateVariablePresets().GetSubscribedPresets(ctx, req)
}

// UnsubscribePreset 取消订阅预设变量
func (c *cTemplateVariablePresets) UnsubscribePreset(ctx context.Context, req *v1.UnsubscribePresetReq) (res *v1.UnsubscribePresetRes, err error) {
	return service.TemplateVariablePresets().UnsubscribePreset(ctx, req)
}


// GetAvailablePresets 获取可用的预设变量列表
func (c *cTemplateVariablePresets) GetAvailablePresets(ctx context.Context, req *v1.GetAvailablePresetsReq) (res *v1.GetAvailablePresetsRes, err error) {
	return service.TemplateVariablePresets().GetAvailablePresets(ctx, req)
}