package service

import (
	"context"

	api "github.com/ciclebyte/template_starter/api/v1/var_preset"
)

type IVarPreset interface {
	// 变量预设基础操作
	Add(ctx context.Context, req *api.VarPresetAddReq) (err error)
	BatchDel(ctx context.Context, req *api.VarPresetBatchDelReq) (err error)
	Del(ctx context.Context, req *api.VarPresetDelReq) (err error)
	Detail(ctx context.Context, req *api.VarPresetDetailReq) (info *api.VarPresetDetailInfo, err error)
	Edit(ctx context.Context, req *api.VarPresetEditReq) (err error)
	List(ctx context.Context, req *api.VarPresetListReq) (total int64, list []*api.VarPresetListInfo, err error)
	All(ctx context.Context, req *api.VarPresetAllReq) (list []*api.VarPresetListInfo, err error)
	Toggle(ctx context.Context, req *api.VarPresetToggleReq) (err error)

	// 模板变量预设关联操作
	GetTemplateVarPresets(ctx context.Context, req *api.TemplateVarPresetsReq) (list []*api.VarPresetListInfo, err error)
	SetTemplateVarPresets(ctx context.Context, req *api.TemplateVarPresetsSetReq) (err error)
	AddTemplateVarPresets(ctx context.Context, req *api.TemplateVarPresetsAddReq) (err error)
	RemoveTemplateVarPresets(ctx context.Context, req *api.TemplateVarPresetsRemoveReq) (err error)
}

var localVarPreset IVarPreset

func VarPreset() IVarPreset {
	if localVarPreset == nil {
		panic("implement not found for interface IVarPreset, forgot register?")
	}
	return localVarPreset
}

func RegisterVarPreset(i IVarPreset) {
	localVarPreset = i
}