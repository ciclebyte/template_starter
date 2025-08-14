package template_expose

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/ciclebyte/template_starter/api/v1/template_expose"
	"github.com/ciclebyte/template_starter/internal/dao"
	"github.com/ciclebyte/template_starter/internal/model/do"
	"github.com/ciclebyte/template_starter/internal/model/entity"
	"github.com/ciclebyte/template_starter/internal/service"
	"github.com/ciclebyte/template_starter/library/liberr"
)

type sTemplateExpose struct{}

func init() {
	service.RegisterTemplateExpose(TemplateExpose())
}

func TemplateExpose() *sTemplateExpose {
	return &sTemplateExpose{}
}

// Get 获取模板暴露字段
func (s *sTemplateExpose) Get(ctx context.Context, req *template_expose.TemplateExposeGetReq) (info *template_expose.TemplateExposeInfo, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		m := dao.TemplateExposeFields.Ctx(ctx).Where(dao.TemplateExposeFields.Columns().TemplateId, req.TemplateId)
		
		// 如果指定了版本，则查询指定版本
		if req.Version != "" {
			m = m.Where(dao.TemplateExposeFields.Columns().Version, req.Version)
		} else {
			// 否则查询最新版本
			m = m.Order(dao.TemplateExposeFields.Columns().Id + " DESC")
		}

		var expose *entity.TemplateExposeFields
		err = m.Scan(&expose)
		liberr.ErrIsNil(ctx, err, "查询模板暴露字段失败")

		if expose != nil {
			info = &template_expose.TemplateExposeInfo{
				Id:              int64(expose.Id),
				TemplateId:      int64(expose.TemplateId),
				FieldSchemaJson: expose.FieldSchemaJson,
				Version:         expose.Version,
				Description:     gconv.String(expose.Description),
				CreatedAt:       expose.CreatedAt.Format("Y-m-d H:i:s"),
				UpdatedAt:       expose.UpdatedAt.Format("Y-m-d H:i:s"),
			}
		}
	})
	return
}

// Set 设置模板暴露字段
func (s *sTemplateExpose) Set(ctx context.Context, req *template_expose.TemplateExposeSetReq) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// 验证JSON格式
		var schemaTest interface{}
		err = json.Unmarshal([]byte(req.FieldSchemaJson), &schemaTest)
		liberr.ErrIsNil(ctx, err, "字段结构定义格式不正确")

		// 设置默认版本
		version := req.Version
		if version == "" {
			version = "1.0"
		}

		// 检查是否已存在相同版本
		count, err := dao.TemplateExposeFields.Ctx(ctx).
			Where(dao.TemplateExposeFields.Columns().TemplateId, req.TemplateId).
			Where(dao.TemplateExposeFields.Columns().Version, version).Count()
		liberr.ErrIsNil(ctx, err, "查询模板暴露字段失败")

		if count > 0 {
			// 更新现有记录
			_, err = dao.TemplateExposeFields.Ctx(ctx).
				Where(dao.TemplateExposeFields.Columns().TemplateId, req.TemplateId).
				Where(dao.TemplateExposeFields.Columns().Version, version).
				Update(do.TemplateExposeFields{
					FieldSchemaJson: req.FieldSchemaJson,
					Description:     req.Description,
				})
			liberr.ErrIsNil(ctx, err, "更新模板暴露字段失败")
		} else {
			// 插入新记录
			_, err = dao.TemplateExposeFields.Ctx(ctx).Insert(do.TemplateExposeFields{
				TemplateId:      req.TemplateId,
				FieldSchemaJson: req.FieldSchemaJson,
				Version:         version,
				Description:     req.Description,
			})
			liberr.ErrIsNil(ctx, err, "设置模板暴露字段失败")
		}
	})
	return
}

// Del 删除模板暴露字段
func (s *sTemplateExpose) Del(ctx context.Context, req *template_expose.TemplateExposeDelReq) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		m := dao.TemplateExposeFields.Ctx(ctx).Where(dao.TemplateExposeFields.Columns().TemplateId, req.TemplateId)
		
		// 如果指定了版本，则删除指定版本
		if req.Version != "" {
			m = m.Where(dao.TemplateExposeFields.Columns().Version, req.Version)
		}

		_, err = m.Delete()
		liberr.ErrIsNil(ctx, err, "删除模板暴露字段失败")
	})
	return
}

// Versions 获取模板暴露字段历史版本
func (s *sTemplateExpose) Versions(ctx context.Context, req *template_expose.TemplateExposeVersionsReq) (versions []*template_expose.TemplateExposeVersionInfo, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		var exposes []*entity.TemplateExposeFields
		err = dao.TemplateExposeFields.Ctx(ctx).
			Where(dao.TemplateExposeFields.Columns().TemplateId, req.TemplateId).
			Order(dao.TemplateExposeFields.Columns().Id + " DESC").
			Scan(&exposes)
		liberr.ErrIsNil(ctx, err, "查询模板暴露字段版本失败")

		// 转换数据格式
		versions = make([]*template_expose.TemplateExposeVersionInfo, 0, len(exposes))
		for _, expose := range exposes {
			versions = append(versions, &template_expose.TemplateExposeVersionInfo{
				Version:     expose.Version,
				Description: gconv.String(expose.Description),
				CreatedAt:   expose.CreatedAt.Format("Y-m-d H:i:s"),
				UpdatedAt:   expose.UpdatedAt.Format("Y-m-d H:i:s"),
			})
		}
	})
	return
}