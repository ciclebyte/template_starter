package languages

import (
	"context"
	"fmt"

	api "github.com/ciclebyte/template_starter/api/v1/languages"
	dao "github.com/ciclebyte/template_starter/internal/dao"
	model "github.com/ciclebyte/template_starter/internal/model"
	do "github.com/ciclebyte/template_starter/internal/model/do"
	service "github.com/ciclebyte/template_starter/internal/service"
	liberr "github.com/ciclebyte/template_starter/library/liberr"
	"github.com/gogf/gf/v2/frame/g"
)

func init() {
	service.RegisterLanguages(New())
}

func New() *sLanguages {
	return &sLanguages{}
}

type sLanguages struct {
}

func (s sLanguages) List(ctx context.Context, req *api.LanguagesListReq) (total interface{}, languagesList []*model.LanguagesInfo, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		m := dao.Languages.Ctx(ctx)
		columns := dao.Languages.Columns()
		if req.Name != "" {
			m = m.Where(fmt.Sprintf("%s like ?", columns.Name), "%"+req.Name+"%")
		}
		if req.DisplayName != "" {
			m = m.Where(fmt.Sprintf("%s like ?", columns.DisplayName), "%"+req.DisplayName+"%")
		}
		if req.Code != "" {
			m = m.Where(columns.Code+" = ?", req.Code)
		}
		if req.IsPopular != 0 {
			m = m.Where(columns.IsPopular+" = ?", req.IsPopular)
		}

		total, err = m.Count()
		liberr.ErrIsNil(ctx, err, "获取语言列表失败")
		orderBy := req.OrderBy
		if orderBy == "" {
			orderBy = "created_at desc"
		}
		err = m.Page(req.PageNum, req.PageSize).Order(orderBy).Scan(&languagesList)
		liberr.ErrIsNil(ctx, err, "获取语言列表失败")
	})
	return
}

func (s sLanguages) Add(ctx context.Context, req *api.LanguagesAddReq) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// TODO 查询是否已经存在

		// add
		_, err = dao.Languages.Ctx(ctx).Insert(do.Languages{
			Name:        req.Name,        // 语言名称（如JavaScript、Python）
			DisplayName: req.DisplayName, // 语言显示名称
			Code:        req.Code,        // 语言代码（如js、py）
			Icon:        req.Icon,        // 语言图标标识或URL
			Color:       req.Color,       // 语言代表色（十六进制）
			IsPopular:   req.IsPopular,   // 是否热门语言
		})
		liberr.ErrIsNil(ctx, err, "新增语言失败")
	})
	return
}

func (s sLanguages) Edit(ctx context.Context, req *api.LanguagesEditReq) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = s.GetById(ctx, req.Id)
		liberr.ErrIsNil(ctx, err, "获取语言失败")
		//TODO 根据名称等查询是否存在

		//编辑
		_, err = dao.Languages.Ctx(ctx).WherePri(req.Id).Update(do.Languages{
			Id:          req.Id,          // 语言ID，自增主键
			Name:        req.Name,        // 语言名称（如JavaScript、Python）
			DisplayName: req.DisplayName, // 语言显示名称
			Code:        req.Code,        // 语言代码（如js、py）
			Icon:        req.Icon,        // 语言图标标识或URL
			Color:       req.Color,       // 语言代表色（十六进制）
			IsPopular:   req.IsPopular,   // 是否热门语言
		})
		liberr.ErrIsNil(ctx, err, "修改语言失败")
	})
	return
}

func (s sLanguages) Delete(ctx context.Context, id int) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = dao.Languages.Ctx(ctx).WherePri(id).Delete()
		liberr.ErrIsNil(ctx, err, "删除语言失败")
	})
	return
}

func (s sLanguages) BatchDelete(ctx context.Context, ids []int) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = dao.Languages.Ctx(ctx).Where(dao.Languages.Columns().Id+" in(?)", ids).Delete()
		liberr.ErrIsNil(ctx, err, "批量删除语言失败")
	})
	return
}

func (s sLanguages) GetById(ctx context.Context, id int) (res *model.LanguagesInfo, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		err = dao.Languages.Ctx(ctx).Where(fmt.Sprintf("%s=?", dao.Languages.Columns().Id), id).Scan(&res)
		liberr.ErrIsNil(ctx, err, "获取语言失败")
	})
	return
}
