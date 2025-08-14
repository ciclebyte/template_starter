package controller

import (
	"context"
	"github.com/ciclebyte/template_starter/api/v1/var_preset"
	"github.com/ciclebyte/template_starter/internal/service"
)

// varPresetController 变量预设控制器
type varPresetController struct{}

var VarPreset = &varPresetController{}

// Add 新增变量预设
func (c *varPresetController) Add(ctx context.Context, req *var_preset.VarPresetAddReq) (res *var_preset.VarPresetAddRes, err error) {
	res = new(var_preset.VarPresetAddRes)
	err = service.VarPreset().Add(ctx, req)
	return
}

// BatchDel 批量删除变量预设
func (c *varPresetController) BatchDel(ctx context.Context, req *var_preset.VarPresetBatchDelReq) (res *var_preset.VarPresetBatchDelRes, err error) {
	res = new(var_preset.VarPresetBatchDelRes)
	err = service.VarPreset().BatchDel(ctx, req)
	return
}

// Del 删除变量预设
func (c *varPresetController) Del(ctx context.Context, req *var_preset.VarPresetDelReq) (res *var_preset.VarPresetDelRes, err error) {
	res = new(var_preset.VarPresetDelRes)
	err = service.VarPreset().Del(ctx, req)
	return
}

// Detail 获取变量预设详情
func (c *varPresetController) Detail(ctx context.Context, req *var_preset.VarPresetDetailReq) (res *var_preset.VarPresetDetailRes, err error) {
	res = new(var_preset.VarPresetDetailRes)
	res.VarPreset, err = service.VarPreset().Detail(ctx, req)
	return
}

// Edit 修改变量预设
func (c *varPresetController) Edit(ctx context.Context, req *var_preset.VarPresetEditReq) (res *var_preset.VarPresetEditRes, err error) {
	res = new(var_preset.VarPresetEditRes)
	err = service.VarPreset().Edit(ctx, req)
	return
}

// List 获取变量预设列表（分页）
func (c *varPresetController) List(ctx context.Context, req *var_preset.VarPresetListReq) (res *var_preset.VarPresetListRes, err error) {
	res = new(var_preset.VarPresetListRes)
	res.Total, res.VarPresetsList, err = service.VarPreset().List(ctx, req)
	return
}

// All 获取所有变量预设（不分页）
func (c *varPresetController) All(ctx context.Context, req *var_preset.VarPresetAllReq) (res *var_preset.VarPresetAllRes, err error) {
	res = new(var_preset.VarPresetAllRes)
	res.VarPresetsList, err = service.VarPreset().All(ctx, req)
	return
}

// Toggle 启用/禁用变量预设
func (c *varPresetController) Toggle(ctx context.Context, req *var_preset.VarPresetToggleReq) (res *var_preset.VarPresetToggleRes, err error) {
	res = new(var_preset.VarPresetToggleRes)
	err = service.VarPreset().Toggle(ctx, req)
	return
}

// TemplateVarPresets 获取模板关联的变量预设列表
func (c *varPresetController) TemplateVarPresets(ctx context.Context, req *var_preset.TemplateVarPresetsReq) (res *var_preset.TemplateVarPresetsRes, err error) {
	res = new(var_preset.TemplateVarPresetsRes)
	res.VarPresetsList, err = service.VarPreset().GetTemplateVarPresets(ctx, req)
	return
}

// TemplateVarPresetsSet 设置模板关联的变量预设
func (c *varPresetController) TemplateVarPresetsSet(ctx context.Context, req *var_preset.TemplateVarPresetsSetReq) (res *var_preset.TemplateVarPresetsSetRes, err error) {
	res = new(var_preset.TemplateVarPresetsSetRes)
	err = service.VarPreset().SetTemplateVarPresets(ctx, req)
	return
}

// TemplateVarPresetsAdd 添加模板变量预设关联
func (c *varPresetController) TemplateVarPresetsAdd(ctx context.Context, req *var_preset.TemplateVarPresetsAddReq) (res *var_preset.TemplateVarPresetsAddRes, err error) {
	res = new(var_preset.TemplateVarPresetsAddRes)
	err = service.VarPreset().AddTemplateVarPresets(ctx, req)
	return
}

// TemplateVarPresetsRemove 移除模板变量预设关联
func (c *varPresetController) TemplateVarPresetsRemove(ctx context.Context, req *var_preset.TemplateVarPresetsRemoveReq) (res *var_preset.TemplateVarPresetsRemoveRes, err error) {
	res = new(var_preset.TemplateVarPresetsRemoveRes)
	err = service.VarPreset().RemoveTemplateVarPresets(ctx, req)
	return
}