package templates

import (
	"context"
	"fmt"
	"strings"

	api "github.com/ciclebyte/template_starter/api/v1/templates"
	dao "github.com/ciclebyte/template_starter/internal/dao"
	model "github.com/ciclebyte/template_starter/internal/model"
	do "github.com/ciclebyte/template_starter/internal/model/do"
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
			Name:         req.Name,         // 模板名称
			Description:  req.Description,  // 模板详细描述
			Introduction: req.Introduction, // 模板详细介绍，支持Markdown格式
			CategoryId:   req.CategoryId,   // 所属分类ID
			IsFeatured:   req.IsFeatured,   // 是否推荐模板
			Logo:         req.Logo,         // 模板logo图片URL
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
			Id:           req.Id,           // 模板ID，自增主键
			Name:         req.Name,         // 模板名称
			Description:  req.Description,  // 模板详细描述
			Introduction: req.Introduction, // 模板详细介绍，支持Markdown格式
			CategoryId:   req.CategoryId,   // 所属分类ID
			IsFeatured:   req.IsFeatured,   // 是否推荐模板
			Logo:         req.Logo,         // 模板logo图片URL
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

		// 检查是否找到数据
		if res == nil {
			liberr.ErrIsNil(ctx, gerror.Newf("模板ID %d 不存在", id), "模板不存在")
		}

		// 查询并填充Languages字段
		var langs []model.TemplateLanguagesInfo
		err = dao.TemplateLanguages.Ctx(ctx).Where("template_id = ?", id).Scan(&langs)
		if err == nil {
			res.Languages = langs
		}
	})
	return
}

func (s sTemplates) GetVariables(ctx context.Context, templateId int64) (res *api.TemplatesVariablesRes, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		res = &api.TemplatesVariablesRes{}

		// 1. 获取用户自定义变量
		var customVariables []*model.TemplateVariablesInfo
		err = dao.TemplateVariables.Ctx(ctx).Where("template_id = ?", templateId).Scan(&customVariables)
		liberr.ErrIsNil(ctx, err, "获取自定义变量失败")
		res.CustomVariables = customVariables

		// 2. 获取模板文件树
		var fileTree []*model.TemplateFilesInfo
		err = dao.TemplateFiles.Ctx(ctx).Where("template_id = ?", templateId).Scan(&fileTree)
		liberr.ErrIsNil(ctx, err, "获取模板文件失败")

		// 3. 解析模板内容中的内置变量和函数
		builtinVars := make(map[string]*api.BuiltinVariableInfo)
		templateFuncs := make(map[string]*api.TemplateFunctionInfo)
		fileSet := make(map[string]bool)

		// 内置变量定义
		builtinVarDefs := map[string]struct {
			label       string
			description string
		}{
			"ProjectName": {"项目名称", "项目的基本名称"},
			"Author":      {"作者", "项目作者信息"},
			"PackageName": {"包名", "项目的包结构名称"},
		}

		// 模板函数定义
		funcDefs := map[string]struct {
			label       string
			description string
		}{
			"now":          {"当前时间", "返回当前时间"},
			"date":         {"格式化日期", "按指定格式返回当前日期"},
			"datetime":     {"格式化时间", "按指定格式返回当前时间"},
			"timestamp":    {"时间戳", "返回Unix时间戳"},
			"year":         {"当前年份", "返回当前年份"},
			"month":        {"当前月份", "返回当前月份"},
			"day":          {"当前日期", "返回当前日期"},
			"lower":        {"转小写", "将变量转换为小写"},
			"upper":        {"转大写", "将变量转换为大写"},
			"title":        {"首字母大写", "将变量首字母大写"},
			"camelcase":    {"驼峰命名", "转换为驼峰命名格式"},
			"snakecase":    {"下划线命名", "转换为下划线命名格式"},
			"kebabcase":    {"短横线命名", "转换为短横线命名格式"},
			"trim":         {"去除空格", "去除变量首尾空格"},
			"trunc":        {"截断字符串", "截断字符串到指定长度"},
			"randInt":      {"随机整数", "生成随机整数"},
			"randAlpha":    {"随机字母", "生成随机字母"},
			"randAlphaNum": {"随机字母数字", "生成随机字母数字"},
			"randNumeric":  {"随机数字", "生成随机数字"},
			"uuid":         {"UUID", "生成UUID"},
			"default":      {"默认值", "如果变量为空则使用默认值"},
			"if":           {"条件判断", "条件判断语句"},
			"eq":           {"相等判断", "判断两个变量是否相等"},
			"ne":           {"不等判断", "判断两个变量是否不相等"},
		}

		// 遍历文件树，解析内容
		for _, file := range fileTree {
			if file.IsDirectory == 1 {
				continue // 跳过目录
			}

			fileSet[file.FileName] = true

			// 获取文件内容
			fileContent, err := dao.TemplateFiles.Ctx(ctx).Where("id = ?", file.Id).Fields("file_content").Value()
			if err != nil {
				continue // 跳过无法读取的文件
			}

			content := fileContent.String()

			// 解析内置变量 {{.变量名}}
			for varName, def := range builtinVarDefs {
				if s.containsVariable(content, varName) {
					if builtinVars[varName] == nil {
						builtinVars[varName] = &api.BuiltinVariableInfo{
							Name:        varName,
							Label:       def.label,
							Description: def.description,
							UsageCount:  0,
							Files:       []string{},
						}
					}
					builtinVars[varName].UsageCount++
					if !s.containsString(builtinVars[varName].Files, file.FileName) {
						builtinVars[varName].Files = append(builtinVars[varName].Files, file.FileName)
					}
				}
			}

			// 解析模板函数 {{函数名}}
			for funcName, def := range funcDefs {
				if s.containsFunction(content, funcName) {
					if templateFuncs[funcName] == nil {
						templateFuncs[funcName] = &api.TemplateFunctionInfo{
							Name:        funcName,
							Label:       def.label,
							Description: def.description,
							UsageCount:  0,
							Files:       []string{},
						}
					}
					templateFuncs[funcName].UsageCount++
					if !s.containsString(templateFuncs[funcName].Files, file.FileName) {
						templateFuncs[funcName].Files = append(templateFuncs[funcName].Files, file.FileName)
					}
				}
			}
		}

		// 转换为切片
		for _, v := range builtinVars {
			res.BuiltinVariables = append(res.BuiltinVariables, v)
		}
		for _, f := range templateFuncs {
			res.TemplateFunctions = append(res.TemplateFunctions, f)
		}

		// 4. 统计信息
		res.Statistics = &api.VariableStatistics{
			TotalCustomVariables:  len(customVariables),
			TotalBuiltinVariables: len(builtinVars),
			TotalFunctions:        len(templateFuncs),
			TotalFiles:            len(fileSet),
		}
	})
	return
}

// 辅助方法：检查内容是否包含变量
func (s sTemplates) containsVariable(content, varName string) bool {
	// 检查 {{.变量名}} 格式
	pattern := fmt.Sprintf("{{.%s}}", varName)
	return strings.Contains(content, pattern)
}

// 辅助方法：检查内容是否包含函数
func (s sTemplates) containsFunction(content, funcName string) bool {
	// 检查 {{函数名}} 格式
	pattern := fmt.Sprintf("{{%s", funcName)
	return strings.Contains(content, pattern)
}

// 辅助方法：检查切片是否包含字符串
func (s sTemplates) containsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
