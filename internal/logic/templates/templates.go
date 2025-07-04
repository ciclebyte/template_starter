package templates

import (
	"context"
	"fmt"

	api "github.com/ciclebyte/template_starter/api/v1/templates"
	dao "github.com/ciclebyte/template_starter/internal/dao"
	model "github.com/ciclebyte/template_starter/internal/model"
	do "github.com/ciclebyte/template_starter/internal/model/do"
	entity "github.com/ciclebyte/template_starter/internal/model/entity"
	service "github.com/ciclebyte/template_starter/internal/service"
	liberr "github.com/ciclebyte/template_starter/library/liberr"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

func init() {
	service.RegisterTemplates(New())
}

func New() *sTemplates {
	return &sTemplates{}
}

type sTemplates struct {
}

func (s sTemplates) List(ctx context.Context, req *api.TemplatesListReq) (total interface{}, templatesList []*model.TemplatesInfo, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		m := dao.Templates.Ctx(ctx)
		columns := dao.Templates.Columns()
		if req.Name != "" {
			m = m.Where(fmt.Sprintf("%s like ?", columns.Name), "%"+req.Name+"%")
		}
		if req.CategoryId != 0 {
			m = m.Where(columns.CategoryId+" = ?", req.CategoryId)
		}
		if req.IsFeatured != 0 {
			m = m.Where(columns.IsFeatured+" = ?", req.IsFeatured)
		}

		total, err = m.Count()
		liberr.ErrIsNil(ctx, err, "获取模板列表失败")
		orderBy := req.OrderBy
		if orderBy == "" {
			orderBy = "created_at desc"
		}
		err = m.Page(req.PageNum, req.PageSize).Order(orderBy).Scan(&templatesList)
		liberr.ErrIsNil(ctx, err, "获取模板列表失败")

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

func (s sTemplates) Add(ctx context.Context, req *api.TemplatesAddReq) (err error) {
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 判重：名称不能重复
		count, err := dao.Templates.Ctx(ctx).TX(tx).Where("name = ?", req.Name).Count()
		liberr.ErrIsNil(ctx, err, "模板名称判重失败")
		if count > 0 {
			return gerror.New("模板名称已存在")
		}

		// add
		result, err := dao.Templates.Ctx(ctx).TX(tx).Insert(do.Templates{
			Name:        req.Name,        // 模板名称
			Description: req.Description, // 模板详细描述
			CategoryId:  req.CategoryId,  // 所属分类ID
			IsFeatured:  req.IsFeatured,  // 是否推荐模板
			Logo:        req.Logo,        // 模板logo图片URL
		})
		liberr.ErrIsNil(ctx, err, "新增模板失败")
		templateId, err := result.LastInsertId()
		liberr.ErrIsNil(ctx, err, "获取模板ID失败")

		// 新增模板语言
		for _, lang := range req.Languages {
			_, err := dao.TemplateLanguages.Ctx(ctx).TX(tx).Insert(do.TemplateLanguages{
				TemplateId: templateId,
				LanguageId: lang.LanguageId,
				IsPrimary:  lang.IsPrimary,
			})
			liberr.ErrIsNil(ctx, err, "新增模板语言失败")
		}
		return nil
	})
	return
}

func (s sTemplates) Edit(ctx context.Context, req *api.TemplatesEditReq) (err error) {
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err = s.GetById(ctx, gconv.Int64(req.Id))
		liberr.ErrIsNil(ctx, err, "获取模板失败")
		// 判重：名称不能与其他模板重复
		count, err := dao.Templates.Ctx(ctx).TX(tx).Where("name = ? AND id <> ?", req.Name, req.Id).Count()
		liberr.ErrIsNil(ctx, err, "模板名称判重失败")
		if count > 0 {
			return gerror.New("模板名称已存在")
		}

		// 编辑模板主表
		_, err = dao.Templates.Ctx(ctx).TX(tx).WherePri(req.Id).Update(do.Templates{
			Id:          req.Id,          // 模板ID，自增主键
			Name:        req.Name,        // 模板名称
			Description: req.Description, // 模板详细描述
			CategoryId:  req.CategoryId,  // 所属分类ID
			IsFeatured:  req.IsFeatured,  // 是否推荐模板
			Logo:        req.Logo,        // 模板logo图片URL
		})
		liberr.ErrIsNil(ctx, err, "修改模板失败")

		// 删除原有模板语言
		_, err = dao.TemplateLanguages.Ctx(ctx).TX(tx).Where("template_id = ?", req.Id).Delete()
		liberr.ErrIsNil(ctx, err, "删除原有模板语言失败")

		// 批量插入新模板语言
		for _, lang := range req.Languages {
			_, err := dao.TemplateLanguages.Ctx(ctx).TX(tx).Insert(do.TemplateLanguages{
				TemplateId: gconv.Int64(req.Id),
				LanguageId: lang.LanguageId,
				IsPrimary:  lang.IsPrimary,
			})
			liberr.ErrIsNil(ctx, err, "新增模板语言失败")
		}
		return nil
	})
	return
}

func (s sTemplates) Delete(ctx context.Context, id int64) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = dao.Templates.Ctx(ctx).WherePri(id).Delete()
		liberr.ErrIsNil(ctx, err, "删除模板失败")
	})
	return
}

func (s sTemplates) BatchDelete(ctx context.Context, ids []int64) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = dao.Templates.Ctx(ctx).Where(dao.Templates.Columns().Id+" in(?)", ids).Delete()
		liberr.ErrIsNil(ctx, err, "批量删除模板失败")
	})
	return
}

func (s sTemplates) GetById(ctx context.Context, id int64) (res *model.TemplatesInfo, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		err = dao.Templates.Ctx(ctx).Where(fmt.Sprintf("%s=?", dao.Templates.Columns().Id), id).Scan(&res)
		liberr.ErrIsNil(ctx, err, "获取模板失败")
		// 查询并填充Languages字段
		var langs []model.TemplateLanguagesInfo
		err = dao.TemplateLanguages.Ctx(ctx).Where("template_id = ?", id).Scan(&langs)
		if err == nil && res != nil {
			res.Languages = langs
		}
	})
	return
}

func (s *sTemplates) FileTree(ctx context.Context, req *api.TemplatesFileTreeReq) (res *api.TemplatesFileTreeRes, err error) {
	res = &api.TemplatesFileTreeRes{}
	templateId := gconv.Int64(req.TemplateId)
	var files []*entity.TemplateFiles
	err = dao.TemplateFiles.Ctx(ctx).Where("template_id = ?", templateId).
		Fields("id, file_path, is_directory, parent_id, file_size, md5").
		Scan(&files)
	if err != nil {
		return
	}
	// 构建树
	idMap := make(map[int64]*api.FileTreeNode)
	for _, f := range files {
		node := &api.FileTreeNode{
			Id: f.Id, FilePath: f.FilePath, IsDirectory: f.IsDirectory,
			ParentId: f.ParentId, FileSize: f.FileSize, Md5: f.Md5,
		}
		idMap[f.Id] = node
	}
	for _, node := range idMap {
		if node.ParentId == 0 {
			res.Tree = append(res.Tree, node)
		} else if parent, ok := idMap[node.ParentId]; ok {
			parent.Children = append(parent.Children, node)
		}
	}
	return
}
