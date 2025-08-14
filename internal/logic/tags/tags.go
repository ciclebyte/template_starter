package tags

import (
	"context"
	"fmt"

	api "github.com/ciclebyte/template_starter/api/v1/tags"
	dao "github.com/ciclebyte/template_starter/internal/dao"
	model "github.com/ciclebyte/template_starter/internal/model"
	do "github.com/ciclebyte/template_starter/internal/model/do"
	service "github.com/ciclebyte/template_starter/internal/service"
	liberr "github.com/ciclebyte/template_starter/library/liberr"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

func init() {
	service.RegisterTags(New())
}

func New() *sTags {
	return &sTags{}
}

type sTags struct {
}

// List 标签列表
func (s sTags) List(ctx context.Context, req *api.TagsListReq) (total interface{}, tagsList []*model.TagsInfo, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		m := dao.Tags.Ctx(ctx).Where("deleted_at IS NULL")
		columns := dao.Tags.Columns()
		
		if req.Name != "" {
			m = m.Where(fmt.Sprintf("%s like ?", columns.Name), "%"+req.Name+"%")
		}

		total, err = m.Count()
		liberr.ErrIsNil(ctx, err, "获取标签列表失败")
		
		orderBy := req.OrderBy
		if orderBy == "" {
			orderBy = "sort desc, created_at desc"
		}
		
		err = m.Page(req.PageNum, req.PageSize).Order(orderBy).Scan(&tagsList)
		liberr.ErrIsNil(ctx, err, "获取标签列表失败")
	})
	return
}

// Add 新增标签
func (s sTags) Add(ctx context.Context, req *api.TagsAddReq) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// 判重：名称不能重复
		count, err := dao.Tags.Ctx(ctx).Where("name = ? AND deleted_at IS NULL", req.Name).Count()
		liberr.ErrIsNil(ctx, err, "标签名称判重失败")
		if count > 0 {
			liberr.ErrIsNil(ctx, gerror.New("标签名称已存在"), "标签名称已存在")
		}

		// 新增标签
		_, err = dao.Tags.Ctx(ctx).Insert(do.Tags{
			Name:        req.Name,
			Description: req.Description,
			Sort:        req.Sort,
		})
		liberr.ErrIsNil(ctx, err, "新增标签失败")
	})
	return
}

// Edit 修改标签
func (s sTags) Edit(ctx context.Context, req *api.TagsEditReq) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// 检查标签是否存在
		_, err = s.GetById(ctx, gconv.Int64(req.Id))
		liberr.ErrIsNil(ctx, err, "获取标签失败")
		
		// 判重：名称不能与其他标签重复
		count, err := dao.Tags.Ctx(ctx).Where("name = ? AND id <> ? AND deleted_at IS NULL", req.Name, req.Id).Count()
		liberr.ErrIsNil(ctx, err, "标签名称判重失败")
		if count > 0 {
			liberr.ErrIsNil(ctx, gerror.New("标签名称已存在"), "标签名称已存在")
		}

		// 修改标签
		_, err = dao.Tags.Ctx(ctx).WherePri(req.Id).Update(do.Tags{
			Id:          req.Id,
			Name:        req.Name,
			Description: req.Description,
			Sort:        req.Sort,
		})
		liberr.ErrIsNil(ctx, err, "修改标签失败")
	})
	return
}

// Delete 删除标签（软删除）
func (s sTags) Delete(ctx context.Context, id int64) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// 检查标签是否存在
		_, err = s.GetById(ctx, id)
		liberr.ErrIsNil(ctx, err, "标签不存在")

		// 软删除标签
		_, err = dao.Tags.Ctx(ctx).WherePri(id).Update(do.Tags{
			DeletedAt: gtime.Now(),
		})
		liberr.ErrIsNil(ctx, err, "删除标签失败")

		// 删除所有关联关系
		_, err = dao.TemplateTags.Ctx(ctx).Where("tag_id = ?", id).Delete()
		liberr.ErrIsNil(ctx, err, "删除标签关联关系失败")
	})
	return
}

// BatchDelete 批量删除标签
func (s sTags) BatchDelete(ctx context.Context, ids []int64) (err error) {
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, id := range ids {
			err = s.Delete(ctx, id)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return
}

// GetById 根据ID获取标签
func (s sTags) GetById(ctx context.Context, id int64) (res *model.TagsInfo, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		err = dao.Tags.Ctx(ctx).Where("id = ? AND deleted_at IS NULL", id).Scan(&res)
		liberr.ErrIsNil(ctx, err, "获取标签失败")

		if res == nil {
			liberr.ErrIsNil(ctx, gerror.Newf("标签ID %d 不存在", id), "标签不存在")
		}
	})
	return
}

// All 获取所有标签（不分页）
func (s sTags) All(ctx context.Context) (res []*model.TagsInfo, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		err = dao.Tags.Ctx(ctx).Where("deleted_at IS NULL").Order("sort desc, created_at desc").Scan(&res)
		liberr.ErrIsNil(ctx, err, "获取所有标签失败")
	})
	return
}

// WithCount 获取标签及关联模板数量
func (s sTags) WithCount(ctx context.Context) (res []*model.TagWithTemplateCount, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// 查询标签及其关联的模板数量
		sql := `
			SELECT 
				t.id, t.name, t.description, t.sort, t.created_at, t.updated_at, t.deleted_at,
				COALESCE(COUNT(tt.template_id), 0) as template_count
			FROM tags t
			LEFT JOIN template_tags tt ON t.id = tt.tag_id
			WHERE t.deleted_at IS NULL
			GROUP BY t.id, t.name, t.description, t.sort, t.created_at, t.updated_at, t.deleted_at
			ORDER BY t.sort DESC, t.created_at DESC
		`
		
		err = dao.Tags.DB().Raw(sql).Scan(&res)
		liberr.ErrIsNil(ctx, err, "获取标签统计失败")
	})
	return
}

// AddTemplateTags 为模板添加标签
func (s sTags) AddTemplateTags(ctx context.Context, templateId int64, tagIds []int64) (err error) {
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 验证模板是否存在
		count, err := dao.Templates.Ctx(ctx).TX(tx).Where("id = ?", templateId).Count()
		liberr.ErrIsNil(ctx, err, "验证模板失败")
		if count == 0 {
			return gerror.New("模板不存在")
		}

		// 验证标签是否存在
		for _, tagId := range tagIds {
			count, err := dao.Tags.Ctx(ctx).TX(tx).Where("id = ? AND deleted_at IS NULL", tagId).Count()
			liberr.ErrIsNil(ctx, err, "验证标签失败")
			if count == 0 {
				return gerror.Newf("标签ID %d 不存在", tagId)
			}

			// 检查是否已经关联
			count, err = dao.TemplateTags.Ctx(ctx).TX(tx).Where("template_id = ? AND tag_id = ?", templateId, tagId).Count()
			liberr.ErrIsNil(ctx, err, "检查标签关联失败")
			if count == 0 {
				// 不存在则新增关联
				_, err = dao.TemplateTags.Ctx(ctx).TX(tx).Insert(do.TemplateTags{
					TemplateId: templateId,
					TagId:      tagId,
				})
				liberr.ErrIsNil(ctx, err, "新增模板标签关联失败")
			}
		}
		return nil
	})
	return
}

// RemoveTemplateTags 为模板移除标签
func (s sTags) RemoveTemplateTags(ctx context.Context, templateId int64, tagIds []int64) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// 删除关联关系
		_, err = dao.TemplateTags.Ctx(ctx).Where("template_id = ? AND tag_id IN(?)", templateId, tagIds).Delete()
		liberr.ErrIsNil(ctx, err, "移除模板标签关联失败")
	})
	return
}

// SetTemplateTags 设置模板标签（批量设置，覆盖原有标签）
func (s sTags) SetTemplateTags(ctx context.Context, templateId int64, tagIds []int64) (err error) {
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 验证模板是否存在
		count, err := dao.Templates.Ctx(ctx).TX(tx).Where("id = ?", templateId).Count()
		liberr.ErrIsNil(ctx, err, "验证模板失败")
		if count == 0 {
			return gerror.New("模板不存在")
		}

		// 删除原有的所有关联
		_, err = dao.TemplateTags.Ctx(ctx).TX(tx).Where("template_id = ?", templateId).Delete()
		liberr.ErrIsNil(ctx, err, "删除原有模板标签关联失败")

		// 如果有新标签，则批量新增
		if len(tagIds) > 0 {
			// 验证所有标签是否存在
			for _, tagId := range tagIds {
				count, err := dao.Tags.Ctx(ctx).TX(tx).Where("id = ? AND deleted_at IS NULL", tagId).Count()
				liberr.ErrIsNil(ctx, err, "验证标签失败")
				if count == 0 {
					return gerror.Newf("标签ID %d 不存在", tagId)
				}
			}

			// 批量插入新关联
			for _, tagId := range tagIds {
				_, err = dao.TemplateTags.Ctx(ctx).TX(tx).Insert(do.TemplateTags{
					TemplateId: templateId,
					TagId:      tagId,
				})
				liberr.ErrIsNil(ctx, err, "新增模板标签关联失败")
			}
		}
		return nil
	})
	return
}

// GetTemplateTags 获取模板标签列表
func (s sTags) GetTemplateTags(ctx context.Context, templateId int64) (res []*model.TagsInfo, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// 联表查询模板的标签
		sql := `
			SELECT t.id, t.name, t.description, t.sort, t.created_at, t.updated_at, t.deleted_at
			FROM tags t
			INNER JOIN template_tags tt ON t.id = tt.tag_id
			WHERE tt.template_id = ? AND t.deleted_at IS NULL
			ORDER BY t.sort DESC, t.created_at DESC
		`
		
		err = dao.Tags.DB().Raw(sql, templateId).Scan(&res)
		liberr.ErrIsNil(ctx, err, "获取模板标签失败")
	})
	return
}

// GetTemplatesByTag 根据标签获取模板列表
func (s sTags) GetTemplatesByTag(ctx context.Context, tagId int64, req *api.TemplatesByTagReq) (total interface{}, templatesList []*model.TemplatesInfo, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// 验证标签是否存在
		_, err = s.GetById(ctx, tagId)
		liberr.ErrIsNil(ctx, err, "标签不存在")

		// 联表查询标签下的模板
		sql := `
			SELECT t.id, t.name, t.description, t.introduction, t.category_id, t.is_featured, t.logo, t.created_at, t.updated_at
			FROM templates t
			INNER JOIN template_tags tt ON t.id = tt.template_id
			WHERE tt.tag_id = ?
		`
		
		// 统计总数
		countSql := `
			SELECT COUNT(*)
			FROM templates t
			INNER JOIN template_tags tt ON t.id = tt.template_id
			WHERE tt.tag_id = ?
		`
		
		total, err = dao.Templates.DB().Raw(countSql, tagId).Value()
		liberr.ErrIsNil(ctx, err, "获取模板总数失败")

		// 分页查询
		orderBy := req.OrderBy
		if orderBy == "" {
			orderBy = "t.created_at DESC"
		}
		
		sql += fmt.Sprintf(" ORDER BY %s LIMIT %d OFFSET %d", orderBy, req.PageSize, (req.PageNum-1)*req.PageSize)
		
		err = dao.Templates.DB().Raw(sql, tagId).Scan(&templatesList)
		liberr.ErrIsNil(ctx, err, "获取标签模板列表失败")

		// 查询并填充每个模板的Languages字段
		for _, tpl := range templatesList {
			var langs []model.TemplateLanguagesInfo
			err = dao.TemplateLanguages.Ctx(ctx).Where("template_id = ?", tpl.Id).Scan(&langs)
			if err == nil {
				tpl.Languages = langs
			}
		}
	})
	return
}