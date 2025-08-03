package template_variables

import (
	"context"
	"fmt"

	api "github.com/ciclebyte/template_starter/api/v1/template_variables"
	dao "github.com/ciclebyte/template_starter/internal/dao"
	model "github.com/ciclebyte/template_starter/internal/model"
	do "github.com/ciclebyte/template_starter/internal/model/do"
	service "github.com/ciclebyte/template_starter/internal/service"
	liberr "github.com/ciclebyte/template_starter/library/liberr"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

func init() {
	service.RegisterTemplateVariables(New())
}

func New() *sTemplateVariables {
	return &sTemplateVariables{}
}

type sTemplateVariables struct {
}

func (s sTemplateVariables) List(ctx context.Context, req *api.TemplateVariablesListReq) (templateVariablesList []*model.TemplateVariablesInfo, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		m := dao.TemplateVariables.Ctx(ctx)
		columns := dao.TemplateVariables.Columns()
		if req.TemplateId != "" {
			idInt64 := gconv.Int64(req.TemplateId)
			m = m.Where(columns.TemplateId+" = ?", idInt64)
		}
		if req.Name != "" {
			m = m.Where(fmt.Sprintf("%s like ?", columns.Name), "%"+req.Name+"%")
		}
		if req.VariableType != "" {
			m = m.Where(columns.VariableType+" = ?", req.VariableType)
		}
		if req.IsRequired != 0 {
			m = m.Where(columns.IsRequired+" = ?", req.IsRequired)
		}

		err = m.Order("sort asc, created_at desc").Scan(&templateVariablesList)
		liberr.ErrIsNil(ctx, err, "获取变量列表失败")
	})
	return
}

func (s sTemplateVariables) Add(ctx context.Context, req *api.TemplateVariablesAddReq) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// TODO 查询是否已经存在

		// add
		_, err = dao.TemplateVariables.Ctx(ctx).Insert(do.TemplateVariables{
			TemplateId:   req.TemplateId,   // 所属模板ID
			Name:         req.Name,         // 变量名称
			VariableType: req.VariableType, // 变量类型
			Description:  req.Description,  // 变量描述
			DefaultValue: req.DefaultValue, // 变量默认值
			IsRequired:   req.IsRequired,   // 是否为必填变量
			Sort:         req.Sort,         // 显示顺序
		})
		liberr.ErrIsNil(ctx, err, "新增分类失败")
	})
	return
}

func (s sTemplateVariables) Edit(ctx context.Context, req *api.TemplateVariablesEditReq) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		idInt64 := gconv.Int64(req.Id)
		_, err = s.GetById(ctx, idInt64)
		liberr.ErrIsNil(ctx, err, "获取分类失败")
		//TODO 根据名称等查询是否存在

		//编辑
		_, err = dao.TemplateVariables.Ctx(ctx).WherePri(req.Id).Update(do.TemplateVariables{
			Id:           req.Id,           // 变量ID，自增主键
			TemplateId:   req.TemplateId,   // 所属模板ID
			Name:         req.Name,         // 变量名称
			VariableType: req.VariableType, // 变量类型
			Description:  req.Description,  // 变量描述
			DefaultValue: req.DefaultValue, // 变量默认值
			IsRequired:   req.IsRequired,   // 是否为必填变量
			Sort:         req.Sort,         // 显示顺序
		})
		liberr.ErrIsNil(ctx, err, "修改分类失败")
	})
	return
}

func (s sTemplateVariables) Delete(ctx context.Context, id int64) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = dao.TemplateVariables.Ctx(ctx).WherePri(id).Delete()
		liberr.ErrIsNil(ctx, err, "删除分类失败")
	})
	return
}

func (s sTemplateVariables) BatchDelete(ctx context.Context, ids []int64) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = dao.TemplateVariables.Ctx(ctx).Where(dao.TemplateVariables.Columns().Id+" in(?)", ids).Delete()
		liberr.ErrIsNil(ctx, err, "批量删除分类失败")
	})
	return
}

func (s sTemplateVariables) GetById(ctx context.Context, id int64) (res *model.TemplateVariablesInfo, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		err = dao.TemplateVariables.Ctx(ctx).Where(fmt.Sprintf("%s=?", dao.TemplateVariables.Columns().Id), id).Scan(&res)
		liberr.ErrIsNil(ctx, err, "获取分类失败")
	})
	return
}

