package templates

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	api "github.com/ciclebyte/template_starter/api/v1/templates"
	template_expose_api "github.com/ciclebyte/template_starter/api/v1/template_expose"
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
			Icon:         req.Icon,         // 模板图标名称
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
			Icon:         req.Icon,         // 模板图标名称
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

		// 1. 自定义变量功能已移除，返回空数组
		res.CustomVariables = []interface{}{}

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
			TotalCustomVariables:  0, // 自定义变量功能已移除
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

// 分析模板变量
func (s sTemplates) AnalyzeVariables(ctx context.Context, templateId int64) (res *api.TemplatesAnalyzeVariablesRes, err error) {
	res = &api.TemplatesAnalyzeVariablesRes{
		DetectedVariables: []*api.DetectedVariable{},
		MissingVariables:  []*api.DetectedVariable{},
		UnusedVariables:   []string{},
		ConflictVariables: []*api.DetectedVariable{},
	}

	err = g.Try(ctx, func(ctx context.Context) {
		// 1. 获取模板文件列表
		var files []*model.TemplateFilesInfo
		err = dao.TemplateFiles.Ctx(ctx).Where(dao.TemplateFiles.Columns().TemplateId, templateId).Scan(&files)
		liberr.ErrIsNil(ctx, err, "获取模板文件失败")

		// 2. 获取当前变量定义
		var templateExposeInfo *template_expose_api.TemplateExposeInfo
		err = dao.TemplateExposeFields.Ctx(ctx).Where(dao.TemplateExposeFields.Columns().TemplateId, templateId).Scan(&templateExposeInfo)
		if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
			liberr.ErrIsNil(ctx, err, "获取模板暴露信息失败")
		}

		// 解析已定义的变量
		definedVars := make(map[string]interface{})
		if templateExposeInfo != nil && templateExposeInfo.FieldSchemaJson != "" {
			err = gconv.Scan(templateExposeInfo.FieldSchemaJson, &definedVars)
			if err != nil {
				g.Log().Warning(ctx, "解析变量定义失败:", err)
				definedVars = make(map[string]interface{})
			}
		}

		// 3. 分析模板文件中的变量使用
		detectedVars := make(map[string]*api.DetectedVariable)
		res.AnalyzedFileCount = len(files)

		g.Log().Info(ctx, "开始分析模板文件，共", len(files), "个文件")
		
		for _, file := range files {
			g.Log().Debug(ctx, "处理文件:", file.FileName, "IsDirectory:", file.IsDirectory, "FilePath:", file.FilePath)
			
			// 跳过目录
			if file.IsDirectory == 1 {
				g.Log().Debug(ctx, "跳过目录:", file.FileName)
				continue
			}
			
			// 读取数据库中存储的文件内容
			fileContent, err := dao.TemplateFiles.Ctx(ctx).Where("id = ?", file.Id).Fields("file_content").Value()
			if err != nil {
				g.Log().Warning(ctx, "读取文件内容失败:", file.FileName, err)
				continue
			}
			
			content := fileContent.String()
			g.Log().Debug(ctx, "文件内容长度:", len(content), "文件:", file.FileName)
			if len(content) > 0 {
				g.Log().Debug(ctx, "文件内容预览:", file.FileName, "=>", content[:min(len(content), 100)])
			}
			
			if content == "" {
				g.Log().Debug(ctx, "文件内容为空:", file.FileName)
				continue
			}

			// 使用正则表达式提取所有变量
			variables := s.extractVariablesFromContent(content)
			g.Log().Info(ctx, "文件", file.FileName, "提取到", len(variables), "个变量")
			
			// 调试：输出检测结果
			if len(variables) > 0 {
				g.Log().Debug(ctx, "文件", file.FileName, "检测到变量:", len(variables))
				for _, v := range variables {
					g.Log().Debug(ctx, "  变量:", v.Name, "类型:", v.Type, "上下文:", v.Contexts)
				}
			}
			
			for _, varInfo := range variables {
				if detectedVars[varInfo.Name] == nil {
					detectedVars[varInfo.Name] = &api.DetectedVariable{
						Name:        varInfo.Name,
						Type:        varInfo.Type,
						Files:       []string{},
						Contexts:    []string{},
						Suggestions: varInfo.Suggestions,
					}
				}
				
				// 添加文件信息
				if !s.containsString(detectedVars[varInfo.Name].Files, file.FileName) {
					detectedVars[varInfo.Name].Files = append(detectedVars[varInfo.Name].Files, file.FileName)
				}
				
				// 添加上下文信息
				for _, context := range varInfo.Contexts {
					if !s.containsString(detectedVars[varInfo.Name].Contexts, context) {
						detectedVars[varInfo.Name].Contexts = append(detectedVars[varInfo.Name].Contexts, context)
					}
				}
			}
		}

		// 4. 分类变量
		for _, detectedVar := range detectedVars {
			res.DetectedVariables = append(res.DetectedVariables, detectedVar)
			
			// 检查是否在已定义变量中
			if _, exists := definedVars[detectedVar.Name]; !exists {
				res.MissingVariables = append(res.MissingVariables, detectedVar)
			} else {
				// 检查类型冲突
				if definedVar, ok := definedVars[detectedVar.Name].(map[string]interface{}); ok {
					if definedType, exists := definedVar["type"]; exists {
						if definedType != detectedVar.Type && detectedVar.Type != "unknown" {
							conflictVar := *detectedVar
							conflictVar.Suggestions = fmt.Sprintf("已定义类型: %s, 检测到的使用模式: %s", definedType, detectedVar.Type)
							res.ConflictVariables = append(res.ConflictVariables, &conflictVar)
						}
					}
				}
			}
		}

		// 5. 查找未使用的变量
		for varName := range definedVars {
			if detectedVars[varName] == nil {
				res.UnusedVariables = append(res.UnusedVariables, varName)
			}
		}

		res.TotalVariableCount = len(detectedVars)
	})
	
	return
}

// 辅助函数：获取两个数的最小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 从内容中提取变量信息
func (s sTemplates) extractVariablesFromContent(content string) []*api.DetectedVariable {
	var variables []*api.DetectedVariable
	
	// 先检查内容中是否包含模板语法
	if !strings.Contains(content, "{{") {
		fmt.Printf("内容中不包含模板语法: %s\n", content[:min(len(content), 50)])
		return variables
	}
	
	fmt.Printf("开始分析内容，长度: %d\n", len(content))
	
	// 正则表达式匹配不同的模板语法
	patterns := map[string]string{
		"simple":    `\{\{\.([a-zA-Z_][a-zA-Z0-9_]*)\}\}`,               // {{.varName}}
		"range":     `\{\{\s*range\s+\.([a-zA-Z_][a-zA-Z0-9_]*)\s*\}\}`, // {{range .items}}
		"if":        `\{\{\s*if\s+\.([a-zA-Z_][a-zA-Z0-9_]*)\s*\}\}`,    // {{if .condition}}
		"with":      `\{\{\s*with\s+\.([a-zA-Z_][a-zA-Z0-9_]*)\s*\}\}`,  // {{with .value}}
		"index":     `\{\{\.([a-zA-Z_][a-zA-Z0-9_]*)\[`,                 // {{.array[index]}}
		"nested":    `\{\{\.([a-zA-Z_][a-zA-Z0-9_]*)\.`,                 // {{.object.field}}
		"function":  `\{\{\s*\w+\s+\.([a-zA-Z_][a-zA-Z0-9_]*)`,          // {{len .items}}
	}

	varMap := make(map[string]*api.DetectedVariable)

	for patternType, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(content, -1)
		
		// 调试：输出匹配结果
		if len(matches) > 0 {
			fmt.Printf("模式 %s 匹配到 %d 个结果\n", patternType, len(matches))
		}
		
		for _, match := range matches {
			if len(match) > 1 {
				varName := match[1]
				context := strings.TrimSpace(match[0])
				
				fmt.Printf("  发现变量: %s, 上下文: %s\n", varName, context)
				
				if varMap[varName] == nil {
					varMap[varName] = &api.DetectedVariable{
						Name:        varName,
						Type:        s.inferVariableType(patternType, context),
						Files:       []string{},
						Contexts:    []string{},
						Suggestions: s.generateTypeSuggestion(patternType, context),
					}
				}
				
				// 添加上下文
				if !s.containsString(varMap[varName].Contexts, context) {
					varMap[varName].Contexts = append(varMap[varName].Contexts, context)
				}
			}
		}
	}

	// 转换为切片
	for _, v := range varMap {
		variables = append(variables, v)
	}
	
	return variables
}

// 推断变量类型
func (s sTemplates) inferVariableType(patternType, context string) string {
	switch patternType {
	case "range":
		return "array" // range 通常用于数组或切片
	case "if":
		return "boolean" // if 通常用于布尔值
	case "index":
		return "array" // 索引访问通常是数组
	case "nested":
		return "object" // 嵌套访问通常是对象
	case "function":
		if strings.Contains(context, "len ") {
			return "array" // len 函数通常用于数组
		}
		return "unknown"
	default:
		return "string" // 默认为字符串
	}
}

// 生成类型建议
func (s sTemplates) generateTypeSuggestion(patternType, context string) string {
	switch patternType {
	case "range":
		return "建议使用 array 或 object_arr 类型，用于循环遍历"
	case "if":
		return "建议使用 boolean 类型，用于条件判断"
	case "index":
		return "建议使用 array 类型，支持索引访问"
	case "nested":
		return "建议使用 object 类型，支持嵌套属性访问"
	case "function":
		if strings.Contains(context, "len ") {
			return "建议使用 array 类型，用于长度计算"
		}
		return "根据使用上下文确定合适的类型"
	default:
		return "建议使用 string 类型，或根据实际用途选择其他类型"
	}
}
