package template_files

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	api "github.com/ciclebyte/template_starter/api/v1/template_files"
	dao "github.com/ciclebyte/template_starter/internal/dao"
	model "github.com/ciclebyte/template_starter/internal/model"
	do "github.com/ciclebyte/template_starter/internal/model/do"
	"github.com/ciclebyte/template_starter/internal/model/entity"
	service "github.com/ciclebyte/template_starter/internal/service"
	liberr "github.com/ciclebyte/template_starter/library/liberr"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
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

func (s *sTemplateFiles) UploadZip(ctx context.Context, templateId int64) (successCount int, failedFiles []string, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// 获取上传的文件
		file := g.RequestFromCtx(ctx).GetUploadFile("zipFile")
		if file == nil {
			liberr.ErrIsNil(ctx, fmt.Errorf("未找到上传的ZIP文件"), "请选择ZIP文件")
		}

		// 打印调试信息
		fmt.Printf("上传文件信息: 文件名=%s, 大小=%d, Content-Type=%s\n",
			file.Filename, file.Size, file.Header.Get("Content-Type"))

		// 检查文件大小（1MB = 1024 * 1024 字节）
		if file.Size > 1024*1024 {
			liberr.ErrIsNil(ctx, fmt.Errorf("文件大小超过1MB限制"), "ZIP文件大小不能超过1MB")
		}

		// 检查文件类型
		if file.Header.Get("Content-Type") != "application/zip" &&
			file.Header.Get("Content-Type") != "application/x-zip-compressed" &&
			!strings.HasSuffix(file.Filename, ".zip") {
			liberr.ErrIsNil(ctx, fmt.Errorf("文件类型不是ZIP"), "请上传ZIP格式的文件")
		}

		// 创建临时目录
		tempDir := g.Cfg().MustGet(ctx, "server.tempPath", "").String()
		if tempDir == "" {
			// 如果没有配置，使用系统临时目录
			tempDir = os.TempDir()
		}
		// 确保临时目录存在
		if err := os.MkdirAll(tempDir, 0755); err != nil {
			liberr.ErrIsNil(ctx, err, "创建临时目录失败")
		}

		fmt.Printf("临时目录: %s\n", tempDir)

		// 使用安全的文件名
		safeFileName := strings.ReplaceAll(file.Filename, "/", "_")
		safeFileName = strings.ReplaceAll(safeFileName, "\\", "_")
		// 添加时间戳避免文件名冲突
		timestamp := gtime.Timestamp()
		uniqueFileName := fmt.Sprintf("%d_%s", timestamp, safeFileName)
		extractPath := filepath.Join(tempDir, "zip_extract_"+gconv.String(templateId)+"_"+gconv.String(timestamp))

		// 创建解压目录
		err = os.MkdirAll(extractPath, 0755)
		liberr.ErrIsNil(ctx, err, "创建临时目录失败")

		// 保存上传的文件到临时目录
		zipPath := filepath.Join(tempDir, uniqueFileName)
		fmt.Printf("保存文件到: %s\n", zipPath)

		// 检查保存前的文件信息
		fmt.Printf("保存前文件大小: %d 字节\n", file.Size)

		// 使用标准库保存文件，而不是 GoFrame 的 file.Save
		src, err := file.Open()
		if err != nil {
			liberr.ErrIsNil(ctx, fmt.Errorf("无法打开上传文件: %v", err), "文件保存失败")
		}
		defer src.Close()

		dst, err := os.Create(zipPath)
		if err != nil {
			liberr.ErrIsNil(ctx, fmt.Errorf("无法创建目标文件: %v", err), "文件保存失败")
		}
		defer dst.Close()

		written, err := io.Copy(dst, src)
		if err != nil {
			liberr.ErrIsNil(ctx, fmt.Errorf("复制文件失败: %v", err), "文件保存失败")
		}

		fmt.Printf("实际写入字节数: %d\n", written)

		// 确保文件写入完成
		dst.Sync()

		// 验证保存的文件是否存在且可读
		if _, err := os.Stat(zipPath); os.IsNotExist(err) {
			liberr.ErrIsNil(ctx, fmt.Errorf("保存的ZIP文件不存在: %s", zipPath), "文件保存失败")
		}

		// 检查保存后的文件信息
		savedFileInfo, err := os.Stat(zipPath)
		if err != nil {
			liberr.ErrIsNil(ctx, fmt.Errorf("无法获取保存文件信息: %v", err), "文件保存失败")
		}
		fmt.Printf("保存后文件大小: %d 字节\n", savedFileInfo.Size())

		if savedFileInfo.Size() == 0 {
			liberr.ErrIsNil(ctx, fmt.Errorf("保存的文件为空"), "文件保存失败")
		}

		fmt.Printf("开始解压文件: %s 到 %s\n", zipPath, extractPath)
		// 解压ZIP文件
		err = s.unzipFile(zipPath, extractPath)
		liberr.ErrIsNil(ctx, err, "解压ZIP文件失败")

		// 递归处理文件
		successCount, failedFiles = s.processExtractedFiles(ctx, templateId, extractPath, "")

		// 清理临时文件
		os.RemoveAll(extractPath)
		os.Remove(zipPath)
	})
	return
}

func (s *sTemplateFiles) unzipFile(zipPath, extractPath string) error {
	fmt.Printf("开始解压文件: %s\n", zipPath)

	// 检查文件是否存在
	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		return fmt.Errorf("ZIP文件不存在: %s", zipPath)
	}

	// 获取文件信息
	fileInfo, err := os.Stat(zipPath)
	if err != nil {
		return fmt.Errorf("无法获取文件信息: %v", err)
	}

	fmt.Printf("文件大小: %d 字节\n", fileInfo.Size())

	// 检查文件大小
	if fileInfo.Size() == 0 {
		return fmt.Errorf("ZIP文件为空")
	}

	// 检查ZIP文件头
	file, err := os.Open(zipPath)
	if err != nil {
		return fmt.Errorf("无法打开ZIP文件: %v", err)
	}

	header := make([]byte, 4)
	_, err = file.Read(header)
	if err != nil {
		file.Close()
		return fmt.Errorf("无法读取文件头: %v", err)
	}

	fmt.Printf("文件头: %02X %02X %02X %02X\n", header[0], header[1], header[2], header[3])

	// ZIP文件头应该是 PK\x03\x04
	if header[0] != 0x50 || header[1] != 0x4B || header[2] != 0x03 || header[3] != 0x04 {
		file.Close()
		return fmt.Errorf("不是有效的ZIP文件，文件头: %02X %02X %02X %02X", header[0], header[1], header[2], header[3])
	}

	// 关闭文件，然后重新打开用于解压
	file.Close()

	fmt.Printf("ZIP文件头验证通过，开始解压...\n")

	// 重新打开文件用于解压
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("无法打开ZIP文件: %v", err)
	}
	defer reader.Close()

	fmt.Printf("ZIP文件包含 %d 个文件/目录\n", len(reader.File))

	for _, zipFile := range reader.File {
		fmt.Printf("解压文件: %s\n", zipFile.Name)
		filePath := filepath.Join(extractPath, zipFile.Name)

		if zipFile.FileInfo().IsDir() {
			os.MkdirAll(filePath, zipFile.Mode())
			continue
		}

		if err = os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return fmt.Errorf("创建目录失败: %v", err)
		}

		outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zipFile.Mode())
		if err != nil {
			return fmt.Errorf("创建文件失败: %v", err)
		}

		rc, err := zipFile.Open()
		if err != nil {
			outFile.Close()
			return fmt.Errorf("打开ZIP内文件失败: %v", err)
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return fmt.Errorf("复制文件内容失败: %v", err)
		}
	}

	fmt.Printf("解压完成\n")
	return nil
}

func (s *sTemplateFiles) processExtractedFiles(ctx context.Context, templateId int64, basePath, relativePath string) (successCount int, failedFiles []string) {
	// 获取当前目录下的所有文件和文件夹
	currentPath := filepath.Join(basePath, relativePath)
	files, err := os.ReadDir(currentPath)
	if err != nil {
		failedFiles = append(failedFiles, relativePath+": "+err.Error())
		return
	}

	for _, file := range files {
		filePath := filepath.Join(relativePath, file.Name())
		if relativePath == "" {
			filePath = file.Name()
		}

		if file.IsDir() {
			// 创建目录记录
			err = s.createDirectoryRecord(ctx, templateId, filePath)
			if err != nil {
				failedFiles = append(failedFiles, filePath+": "+err.Error())
			} else {
				successCount++
			}

			// 递归处理子目录
			subSuccess, subFailed := s.processExtractedFiles(ctx, templateId, basePath, filePath)
			successCount += subSuccess
			failedFiles = append(failedFiles, subFailed...)
		} else {
			// 读取文件内容
			fullPath := filepath.Join(basePath, filePath)
			content, err := os.ReadFile(fullPath)
			if err != nil {
				failedFiles = append(failedFiles, filePath+": 读取文件失败")
				continue
			}

			// 创建文件记录
			err = s.createFileRecord(ctx, templateId, filePath, string(content))
			if err != nil {
				failedFiles = append(failedFiles, filePath+": "+err.Error())
			} else {
				successCount++
			}
		}
	}
	return
}

func (s *sTemplateFiles) createDirectoryRecord(ctx context.Context, templateId int64, filePath string) error {
	// 检查目录是否已存在
	count, err := dao.TemplateFiles.Ctx(ctx).Where("template_id = ? AND file_path = ? AND is_directory = 1", templateId, filePath).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // 目录已存在，跳过
	}

	// 获取父目录ID
	parentId := s.getParentId(ctx, templateId, filePath)

	// 创建目录记录
	_, err = dao.TemplateFiles.Ctx(ctx).Insert(do.TemplateFiles{
		TemplateId:  templateId,
		FilePath:    filePath,
		FileContent: "",
		FileSize:    0,
		IsDirectory: 1,
		Md5:         "",
		Sort:        0,
		ParentId:    parentId,
	})
	return err
}

func (s *sTemplateFiles) createFileRecord(ctx context.Context, templateId int64, filePath, content string) error {
	// 检查文件是否已存在
	count, err := dao.TemplateFiles.Ctx(ctx).Where("template_id = ? AND file_path = ? AND is_directory = 0", templateId, filePath).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // 文件已存在，跳过
	}

	// 获取父目录ID
	parentId := s.getParentId(ctx, templateId, filePath)

	// 计算MD5
	md5 := gmd5.MustEncrypt(content)

	// 创建文件记录
	_, err = dao.TemplateFiles.Ctx(ctx).Insert(do.TemplateFiles{
		TemplateId:  templateId,
		FilePath:    filePath,
		FileContent: content,
		FileSize:    len(content),
		IsDirectory: 0,
		Md5:         md5,
		Sort:        0,
		ParentId:    parentId,
	})
	return err
}

func (s *sTemplateFiles) getParentId(ctx context.Context, templateId int64, filePath string) int64 {
	// 获取父目录路径
	parentPath := filepath.Dir(filePath)
	if parentPath == "." || parentPath == "" {
		return 0 // 根目录
	}

	// 查找父目录ID
	var parent entity.TemplateFiles
	err := dao.TemplateFiles.Ctx(ctx).Where("template_id = ? AND file_path = ? AND is_directory = 1", templateId, parentPath).Scan(&parent)
	if err != nil {
		return 0 // 找不到父目录，设为根目录
	}
	return parent.Id
}
