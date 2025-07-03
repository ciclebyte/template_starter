package template_languages

import (
	"context"
	"fmt"

	api "github.com/ciclebyte/template_starter/api/v1/template_languages"
	dao "github.com/ciclebyte/template_starter/internal/dao"
	model "github.com/ciclebyte/template_starter/internal/model"
	do "github.com/ciclebyte/template_starter/internal/model/do"
	service "github.com/ciclebyte/template_starter/internal/service"
	liberr "github.com/ciclebyte/template_starter/library/liberr"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

func init() {
	service.RegisterTemplateLanguages(New())
}

func New() *sTemplateLanguages {
	return &sTemplateLanguages{}
}

type sTemplateLanguages struct {
}

func (s sTemplateLanguages) List(ctx context.Context, req *api.TemplateLanguagesListReq) (total interface{}, templateLanguagesList []*model.TemplateLanguagesInfo, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		m := dao.TemplateLanguages.Ctx(ctx)
		columns := dao.TemplateLanguages.Columns()
		if req.TemplateId != "" {
			m = m.Where(columns.TemplateId+" = ?", gconv.Int64(req.TemplateId))
		}
		if req.LanguageId != 0 {
			m = m.Where(columns.LanguageId+" = ?", req.LanguageId)
		}
		if req.IsPrimary != 0 {
			m = m.Where(columns.IsPrimary+" = ?", req.IsPrimary)
		}

		total, err = m.Count()
		liberr.ErrIsNil(ctx, err, "获取模板语言列表失败")
		orderBy := req.OrderBy
		if orderBy == "" {
			orderBy = "created_at desc"
		}
		err = m.Page(req.PageNum, req.PageSize).Order(orderBy).Scan(&templateLanguagesList)
		liberr.ErrIsNil(ctx, err, "获取模板语言列表失败")
	})
	return
}

func (s sTemplateLanguages) Add(ctx context.Context, req *api.TemplateLanguagesAddReq) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// TODO 查询是否已经存在

		// add
		_, err = dao.TemplateLanguages.Ctx(ctx).Insert(do.TemplateLanguages{
			TemplateId: req.TemplateId, // 关联的模板ID
			LanguageId: req.LanguageId, // 关联的语言ID
			IsPrimary:  req.IsPrimary,  // 是否主要语言
		})
		liberr.ErrIsNil(ctx, err, "新增模板语言失败")
	})
	return
}

func (s sTemplateLanguages) Edit(ctx context.Context, req *api.TemplateLanguagesEditReq) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = s.GetById(ctx, gconv.Int64(req.Id))
		liberr.ErrIsNil(ctx, err, "获取模板语言失败")
		//TODO 根据名称等查询是否存在

		//编辑
		_, err = dao.TemplateLanguages.Ctx(ctx).WherePri(req.Id).Update(do.TemplateLanguages{
			Id:         req.Id,         // 关联ID，自增主键
			TemplateId: req.TemplateId, // 关联的模板ID
			LanguageId: req.LanguageId, // 关联的语言ID
			IsPrimary:  req.IsPrimary,  // 是否主要语言
		})
		liberr.ErrIsNil(ctx, err, "修改模板语言失败")
	})
	return
}

func (s sTemplateLanguages) Delete(ctx context.Context, id int64) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = dao.TemplateLanguages.Ctx(ctx).WherePri(id).Delete()
		liberr.ErrIsNil(ctx, err, "删除模板语言失败")
	})
	return
}

func (s sTemplateLanguages) BatchDelete(ctx context.Context, ids []int64) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = dao.TemplateLanguages.Ctx(ctx).Where(dao.TemplateLanguages.Columns().Id+" in(?)", ids).Delete()
		liberr.ErrIsNil(ctx, err, "批量删除模板语言失败")
	})
	return
}

func (s sTemplateLanguages) GetById(ctx context.Context, id int64) (res *model.TemplateLanguagesInfo, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		err = dao.TemplateLanguages.Ctx(ctx).Where(fmt.Sprintf("%s=?", dao.TemplateLanguages.Columns().Id), id).Scan(&res)
		liberr.ErrIsNil(ctx, err, "获取模板语言失败")
	})
	return
}
