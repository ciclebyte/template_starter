package var_preset

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/ciclebyte/template_starter/api/v1/var_preset"
	"github.com/ciclebyte/template_starter/internal/dao"
	"github.com/ciclebyte/template_starter/internal/model/do"
	"github.com/ciclebyte/template_starter/internal/model/entity"
	"github.com/ciclebyte/template_starter/internal/service"
	"github.com/ciclebyte/template_starter/library/liberr"
)

type sVarPreset struct{}

func init() {
	service.RegisterVarPreset(VarPreset())
}

func VarPreset() *sVarPreset {
	return &sVarPreset{}
}

// Add 新增变量预设
func (s *sVarPreset) Add(ctx context.Context, req *var_preset.VarPresetAddReq) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// 检查名称是否已存在
		count, err := dao.VarPreset.Ctx(ctx).Where(dao.VarPreset.Columns().Name, req.Name).Count()
		liberr.ErrIsNil(ctx, err, "查询变量预设失败")
		if count > 0 {
			liberr.ErrIsNil(ctx, gerror.New("变量预设名称已存在"), "变量预设名称已存在")
		}

		// 验证JSON格式
		var schemaTest interface{}
		err = json.Unmarshal([]byte(req.SchemaJson), &schemaTest)
		liberr.ErrIsNil(ctx, err, "数据结构模板格式不正确")

		if req.DefaultDataJson != "" {
			var defaultDataTest interface{}
			err = json.Unmarshal([]byte(req.DefaultDataJson), &defaultDataTest)
			liberr.ErrIsNil(ctx, err, "默认数据格式不正确")
		}

		// 插入数据
		_, err = dao.VarPreset.Ctx(ctx).Insert(do.VarPreset{
			Name:            req.Name,
			DisplayName:     req.DisplayName,
			Description:     req.Description,
			Category:        req.Category,
			SchemaJson:      req.SchemaJson,
			DefaultDataJson: req.DefaultDataJson,
			Icon:            req.Icon,
			Sort:            req.Sort,
			Version:         req.Version,
			IsEnabled:       1,
		})
		liberr.ErrIsNil(ctx, err, "添加变量预设失败")
	})
	return
}

// BatchDel 批量删除变量预设
func (s *sVarPreset) BatchDel(ctx context.Context, req *var_preset.VarPresetBatchDelReq) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// 检查是否有系统预设
		systemCount, err := dao.VarPreset.Ctx(ctx).Where(dao.VarPreset.Columns().Id+" IN(?)", req.Ids).
			Where(dao.VarPreset.Columns().Category, "system").Count()
		liberr.ErrIsNil(ctx, err, "查询变量预设失败")
		if systemCount > 0 {
			liberr.ErrIsNil(ctx, gerror.New("系统预设不能删除"), "系统预设不能删除")
		}

		// 软删除
		_, err = dao.VarPreset.Ctx(ctx).Where(dao.VarPreset.Columns().Id+" IN(?)", req.Ids).
			Update(do.VarPreset{DeletedAt: gtime.Now()})
		liberr.ErrIsNil(ctx, err, "批量删除变量预设失败")
	})
	return
}

// Del 删除变量预设
func (s *sVarPreset) Del(ctx context.Context, req *var_preset.VarPresetDelReq) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// 查询变量预设信息
		preset, err := dao.VarPreset.Ctx(ctx).WherePri(req.Id).One()
		liberr.ErrIsNil(ctx, err, "查询变量预设失败")
		if preset.IsEmpty() {
			liberr.ErrIsNil(ctx, gerror.New("变量预设不存在"), "变量预设不存在")
		}

		// 检查是否为系统预设
		if preset["category"].String() == "system" {
			liberr.ErrIsNil(ctx, gerror.New("系统预设不能删除"), "系统预设不能删除")
		}

		// 软删除
		_, err = dao.VarPreset.Ctx(ctx).WherePri(req.Id).Update(do.VarPreset{DeletedAt: gtime.Now()})
		liberr.ErrIsNil(ctx, err, "删除变量预设失败")
	})
	return
}

// Detail 获取变量预设详情
func (s *sVarPreset) Detail(ctx context.Context, req *var_preset.VarPresetDetailReq) (info *var_preset.VarPresetDetailInfo, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		var preset *entity.VarPreset
		err = dao.VarPreset.Ctx(ctx).WherePri(req.Id).WhereNull(dao.VarPreset.Columns().DeletedAt).Scan(&preset)
		liberr.ErrIsNil(ctx, err, "查询变量预设失败")
		if preset == nil {
			liberr.ErrIsNil(ctx, gerror.New("变量预设不存在"), "变量预设不存在")
		}

		info = &var_preset.VarPresetDetailInfo{
			Id:              int64(preset.Id),
			Name:            preset.Name,
			DisplayName:     preset.DisplayName,
			Description:     gconv.String(preset.Description),
			Category:        preset.Category,
			SchemaJson:      preset.SchemaJson,
			DefaultDataJson: gconv.String(preset.DefaultDataJson),
			Icon:            gconv.String(preset.Icon),
			Sort:            preset.Sort,
			IsEnabled:       preset.IsEnabled,
			Version:         preset.Version,
			CreatedBy:       gconv.Int64(preset.CreatedBy),
			CreatedAt:       preset.CreatedAt.Format("Y-m-d H:i:s"),
			UpdatedAt:       preset.UpdatedAt.Format("Y-m-d H:i:s"),
		}
	})
	return
}

// Edit 修改变量预设
func (s *sVarPreset) Edit(ctx context.Context, req *var_preset.VarPresetEditReq) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// 查询变量预设信息
		preset, err := dao.VarPreset.Ctx(ctx).WherePri(req.Id).WhereNull(dao.VarPreset.Columns().DeletedAt).One()
		liberr.ErrIsNil(ctx, err, "查询变量预设失败")
		if preset.IsEmpty() {
			liberr.ErrIsNil(ctx, gerror.New("变量预设不存在"), "变量预设不存在")
		}

		// 检查名称是否被其他记录使用
		count, err := dao.VarPreset.Ctx(ctx).Where(dao.VarPreset.Columns().Name, req.Name).
			WhereNot(dao.VarPreset.Columns().Id, req.Id).WhereNull(dao.VarPreset.Columns().DeletedAt).Count()
		liberr.ErrIsNil(ctx, err, "查询变量预设失败")
		if count > 0 {
			liberr.ErrIsNil(ctx, gerror.New("变量预设名称已存在"), "变量预设名称已存在")
		}

		// 验证JSON格式
		var schemaTest interface{}
		err = json.Unmarshal([]byte(req.SchemaJson), &schemaTest)
		liberr.ErrIsNil(ctx, err, "数据结构模板格式不正确")

		if req.DefaultDataJson != "" {
			var defaultDataTest interface{}
			err = json.Unmarshal([]byte(req.DefaultDataJson), &defaultDataTest)
			liberr.ErrIsNil(ctx, err, "默认数据格式不正确")
		}

		// 更新数据
		_, err = dao.VarPreset.Ctx(ctx).WherePri(req.Id).Update(do.VarPreset{
			Name:            req.Name,
			DisplayName:     req.DisplayName,
			Description:     req.Description,
			Category:        req.Category,
			SchemaJson:      req.SchemaJson,
			DefaultDataJson: req.DefaultDataJson,
			Icon:            req.Icon,
			Sort:            req.Sort,
			Version:         req.Version,
			IsEnabled:       req.IsEnabled,
		})
		liberr.ErrIsNil(ctx, err, "更新变量预设失败")
	})
	return
}

// List 获取变量预设列表（分页）
func (s *sVarPreset) List(ctx context.Context, req *var_preset.VarPresetListReq) (total int64, list []*var_preset.VarPresetListInfo, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		m := dao.VarPreset.Ctx(ctx).WhereNull(dao.VarPreset.Columns().DeletedAt)

		// 搜索条件
		if req.Name != "" {
			m = m.WhereLike(dao.VarPreset.Columns().Name, "%"+req.Name+"%")
		}
		if req.Category != "" {
			m = m.Where(dao.VarPreset.Columns().Category, req.Category)
		}

		// 获取总数
		totalInt, err := m.Count()
		liberr.ErrIsNil(ctx, err, "查询变量预设总数失败")
		total = int64(totalInt)

		// 分页查询
		if req.PageNum <= 0 {
			req.PageNum = 1
		}
		if req.PageSize <= 0 {
			req.PageSize = 20
		}

		var presets []*entity.VarPreset
		err = m.Order(dao.VarPreset.Columns().Sort + " DESC," + dao.VarPreset.Columns().Id + " DESC").
			Page(req.PageNum, req.PageSize).Scan(&presets)
		liberr.ErrIsNil(ctx, err, "查询变量预设列表失败")

		// 转换数据格式
		list = make([]*var_preset.VarPresetListInfo, 0, len(presets))
		for _, preset := range presets {
			list = append(list, &var_preset.VarPresetListInfo{
				Id:          int64(preset.Id),
				Name:        preset.Name,
				DisplayName: preset.DisplayName,
				Description: gconv.String(preset.Description),
				Category:    preset.Category,
				Icon:        gconv.String(preset.Icon),
				Sort:        preset.Sort,
				IsEnabled:   preset.IsEnabled,
				Version:     preset.Version,
				CreatedAt:   preset.CreatedAt.Format("Y-m-d H:i:s"),
				UpdatedAt:   preset.UpdatedAt.Format("Y-m-d H:i:s"),
			})
		}
	})
	return
}

// All 获取所有变量预设（不分页）
func (s *sVarPreset) All(ctx context.Context, req *var_preset.VarPresetAllReq) (list []*var_preset.VarPresetListInfo, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		m := dao.VarPreset.Ctx(ctx).WhereNull(dao.VarPreset.Columns().DeletedAt)

		// 过滤条件
		if req.Category != "" {
			m = m.Where(dao.VarPreset.Columns().Category, req.Category)
		}
		if req.IsEnabled > 0 {
			m = m.Where(dao.VarPreset.Columns().IsEnabled, req.IsEnabled)
		}

		var presets []*entity.VarPreset
		err = m.Order(dao.VarPreset.Columns().Sort + " DESC," + dao.VarPreset.Columns().Id + " DESC").Scan(&presets)
		liberr.ErrIsNil(ctx, err, "查询变量预设列表失败")

		// 转换数据格式
		list = make([]*var_preset.VarPresetListInfo, 0, len(presets))
		for _, preset := range presets {
			list = append(list, &var_preset.VarPresetListInfo{
				Id:          int64(preset.Id),
				Name:        preset.Name,
				DisplayName: preset.DisplayName,
				Description: gconv.String(preset.Description),
				Category:    preset.Category,
				Icon:        gconv.String(preset.Icon),
				Sort:        preset.Sort,
				IsEnabled:   preset.IsEnabled,
				Version:     preset.Version,
				CreatedAt:   preset.CreatedAt.Format("Y-m-d H:i:s"),
				UpdatedAt:   preset.UpdatedAt.Format("Y-m-d H:i:s"),
			})
		}
	})
	return
}

// Toggle 启用/禁用变量预设
func (s *sVarPreset) Toggle(ctx context.Context, req *var_preset.VarPresetToggleReq) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = s.GetById(ctx, req.Id)
		liberr.ErrIsNil(ctx, err, "变量预设不存在")

		_, err = dao.VarPreset.Ctx(ctx).WherePri(req.Id).Update(do.VarPreset{
			IsEnabled: req.IsEnabled,
		})
		liberr.ErrIsNil(ctx, err, "更新变量预设状态失败")
	})
	return
}

// GetById 根据ID获取变量预设
func (s *sVarPreset) GetById(ctx context.Context, id int64) (preset *entity.VarPreset, err error) {
	err = dao.VarPreset.Ctx(ctx).WherePri(id).WhereNull(dao.VarPreset.Columns().DeletedAt).Scan(&preset)
	if err != nil {
		return nil, err
	}
	if preset == nil {
		return nil, gerror.New("变量预设不存在")
	}
	return
}

// GetTemplateVarPresets 获取模板关联的变量预设列表
func (s *sVarPreset) GetTemplateVarPresets(ctx context.Context, req *var_preset.TemplateVarPresetsReq) (list []*var_preset.VarPresetListInfo, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		var presets []*entity.VarPreset
		err = dao.VarPreset.Ctx(ctx).
			LeftJoin("template_var_presets tvp", "tvp.preset_id = var_preset.id").
			Where("tvp.template_id", req.TemplateId).
			WhereNull(dao.VarPreset.Columns().DeletedAt).
			Order(dao.VarPreset.Columns().Sort + " DESC," + dao.VarPreset.Columns().Id + " DESC").
			Scan(&presets)
		liberr.ErrIsNil(ctx, err, "查询模板变量预设失败")

		// 转换数据格式
		list = make([]*var_preset.VarPresetListInfo, 0, len(presets))
		for _, preset := range presets {
			list = append(list, &var_preset.VarPresetListInfo{
				Id:          int64(preset.Id),
				Name:        preset.Name,
				DisplayName: preset.DisplayName,
				Description: gconv.String(preset.Description),
				Category:    preset.Category,
				Icon:        gconv.String(preset.Icon),
				Sort:        preset.Sort,
				IsEnabled:   preset.IsEnabled,
				Version:     preset.Version,
				CreatedAt:   preset.CreatedAt.Format("Y-m-d H:i:s"),
				UpdatedAt:   preset.UpdatedAt.Format("Y-m-d H:i:s"),
			})
		}
	})
	return
}

// SetTemplateVarPresets 设置模板关联的变量预设
func (s *sVarPreset) SetTemplateVarPresets(ctx context.Context, req *var_preset.TemplateVarPresetsSetReq) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		err = dao.TemplateVarPresets.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			// 删除现有关联
			_, err = dao.TemplateVarPresets.Ctx(ctx).TX(tx).
				Where(dao.TemplateVarPresets.Columns().TemplateId, req.TemplateId).Delete()
			if err != nil {
				return err
			}

			// 添加新关联
			if len(req.PresetIds) > 0 {
				var data []do.TemplateVarPresets
				for _, presetId := range req.PresetIds {
					data = append(data, do.TemplateVarPresets{
						TemplateId: req.TemplateId,
						PresetId:   presetId,
					})
				}
				_, err = dao.TemplateVarPresets.Ctx(ctx).TX(tx).Insert(data)
				if err != nil {
					return err
				}
			}

			return nil
		})
		liberr.ErrIsNil(ctx, err, "设置模板变量预设失败")
	})
	return
}

// AddTemplateVarPresets 添加模板变量预设关联
func (s *sVarPreset) AddTemplateVarPresets(ctx context.Context, req *var_preset.TemplateVarPresetsAddReq) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		var data []do.TemplateVarPresets
		for _, presetId := range req.PresetIds {
			// 检查是否已存在
			count, err := dao.TemplateVarPresets.Ctx(ctx).
				Where(dao.TemplateVarPresets.Columns().TemplateId, req.TemplateId).
				Where(dao.TemplateVarPresets.Columns().PresetId, presetId).Count()
			liberr.ErrIsNil(ctx, err, "查询模板变量预设关联失败")

			if count == 0 {
				data = append(data, do.TemplateVarPresets{
					TemplateId: req.TemplateId,
					PresetId:   presetId,
				})
			}
		}

		if len(data) > 0 {
			_, err = dao.TemplateVarPresets.Ctx(ctx).Insert(data)
			liberr.ErrIsNil(ctx, err, "添加模板变量预设关联失败")
		}
	})
	return
}

// RemoveTemplateVarPresets 移除模板变量预设关联
func (s *sVarPreset) RemoveTemplateVarPresets(ctx context.Context, req *var_preset.TemplateVarPresetsRemoveReq) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = dao.TemplateVarPresets.Ctx(ctx).
			Where(dao.TemplateVarPresets.Columns().TemplateId, req.TemplateId).
			Where(dao.TemplateVarPresets.Columns().PresetId+" IN(?)", req.PresetIds).Delete()
		liberr.ErrIsNil(ctx, err, "移除模板变量预设关联失败")
	})
	return
}