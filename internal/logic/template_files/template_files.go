package template_files

import (
	"context"
	"fmt"

	api "github.com/ciclebyte/template_starter/api/v1/template_files"
	dao "github.com/ciclebyte/template_starter/internal/dao"
	model "github.com/ciclebyte/template_starter/internal/model"
	do "github.com/ciclebyte/template_starter/internal/model/do"
	"github.com/ciclebyte/template_starter/internal/model/entity"
	service "github.com/ciclebyte/template_starter/internal/service"
	liberr "github.com/ciclebyte/template_starter/library/liberr"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

func init() {
	service.RegisterTemplateFiles(New())
}

func New() *sTemplateFiles {
	return &sTemplateFiles{}
}

type sTemplateFiles struct {
}

func (s sTemplateFiles) List(ctx context.Context, req *api.TemplateFilesListReq) (total interface{}, templateFilesList []*model.TemplateFilesInfo, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		m := dao.TemplateFiles.Ctx(ctx)
		columns := dao.TemplateFiles.Columns()
		if req.TemplateId != "" {
			m = m.Where(columns.TemplateId+" = ?", gconv.Int64(req.TemplateId))
		}

		total, err = m.Count()
		liberr.ErrIsNil(ctx, err, "获取模板文件列表失败")
		orderBy := req.OrderBy
		if orderBy == "" {
			orderBy = "created_at desc"
		}
		err = m.Page(req.PageNum, req.PageSize).Order(orderBy).Scan(&templateFilesList)
		liberr.ErrIsNil(ctx, err, "获取模板文件列表失败")
	})
	return
}

func (s sTemplateFiles) Add(ctx context.Context, req *api.TemplateFilesAddReq) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// TODO 查询是否已经存在

		// add
		_, err = dao.TemplateFiles.Ctx(ctx).Insert(do.TemplateFiles{
			TemplateId:  req.TemplateId,  // 所属模板ID
			FilePath:    req.FilePath,    // 文件路径（相对路径）
			FileContent: req.FileContent, // 文件内容
			FileSize:    req.FileSize,    // 文件大小（字节）
			IsDirectory: req.IsDirectory, // 是否为目录
			Md5:         req.Md5,         // md5
			Sort:        req.Sort,        // 排序
			ParentId:    req.ParentId,    // 父目录ID，如果是文件则指向所属目录
		})
		liberr.ErrIsNil(ctx, err, "新增模板文件失败")
	})
	return
}

func (s sTemplateFiles) Edit(ctx context.Context, req *api.TemplateFilesEditReq) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = s.GetById(ctx, gconv.Int64(req.Id))
		liberr.ErrIsNil(ctx, err, "获取模板文件失败")
		//TODO 根据名称等查询是否存在

		//编辑
		_, err = dao.TemplateFiles.Ctx(ctx).WherePri(req.Id).Update(do.TemplateFiles{
			Id:          req.Id,          // 文件ID，自增主键
			TemplateId:  req.TemplateId,  // 所属模板ID
			FilePath:    req.FilePath,    // 文件路径（相对路径）
			FileContent: req.FileContent, // 文件内容
			FileSize:    req.FileSize,    // 文件大小（字节）
			IsDirectory: req.IsDirectory, // 是否为目录
			Md5:         req.Md5,         // md5
			Sort:        req.Sort,        // 排序
			ParentId:    req.ParentId,    // 父目录ID，如果是文件则指向所属目录
		})
		liberr.ErrIsNil(ctx, err, "修改模板文件失败")
	})
	return
}

func (s sTemplateFiles) Delete(ctx context.Context, id int64) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = dao.TemplateFiles.Ctx(ctx).WherePri(id).Delete()
		liberr.ErrIsNil(ctx, err, "删除模板文件失败")
	})
	return
}

func (s sTemplateFiles) BatchDelete(ctx context.Context, ids []int64) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = dao.TemplateFiles.Ctx(ctx).Where(dao.TemplateFiles.Columns().Id+" in(?)", ids).Delete()
		liberr.ErrIsNil(ctx, err, "批量删除模板文件失败")
	})
	return
}

func (s sTemplateFiles) GetById(ctx context.Context, id int64) (res *model.TemplateFilesInfo, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		err = dao.TemplateFiles.Ctx(ctx).Where(fmt.Sprintf("%s=?", dao.TemplateFiles.Columns().Id), id).Scan(&res)
		liberr.ErrIsNil(ctx, err, "获取模板文件失败")
	})
	return
}

func (s *sTemplateFiles) FileTree(ctx context.Context, req *api.TemplatesFileTreeReq) (res *api.TemplatesFileTreeRes, err error) {
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

func (s *sTemplateFiles) GetFileContent(ctx context.Context, id int64) (fileContent string, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		var entityFile entity.TemplateFiles
		err = dao.TemplateFiles.Ctx(ctx).Where(dao.TemplateFiles.Columns().Id+"=?", id).Fields("file_content").Scan(&entityFile)
		liberr.ErrIsNil(ctx, err, "获取文件内容失败")
		fileContent = entityFile.FileContent
	})
	return
}
