package template_files

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
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

		// 动态生成文件路径
		filePath := s.generateFilePath(ctx, req.TemplateId, req.ParentId, req.FileName)

		// add
		_, err = dao.TemplateFiles.Ctx(ctx).Insert(do.TemplateFiles{
			TemplateId:  req.TemplateId,  // 所属模板ID
			FilePath:    filePath,        // 文件路径（相对路径）
			FileName:    req.FileName,    // 文件名
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

		// 动态生成文件路径
		filePath := s.generateFilePath(ctx, req.TemplateId, req.ParentId, req.FileName)

		//编辑
		_, err = dao.TemplateFiles.Ctx(ctx).WherePri(req.Id).Update(do.TemplateFiles{
			Id:          req.Id,          // 文件ID，自增主键
			TemplateId:  req.TemplateId,  // 所属模板ID
			FilePath:    filePath,        // 文件路径（相对路径）
			FileName:    req.FileName,    // 文件名
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

func (s sTemplateFiles) Rename(ctx context.Context, req *api.TemplateFilesRenameReq) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// 获取当前文件信息
		fileInfo, err := s.GetById(ctx, gconv.Int64(req.Id))
		liberr.ErrIsNil(ctx, err, "获取模板文件失败")

		// 检查新文件名是否为空
		if req.FileName == "" {
			liberr.ErrIsNil(ctx, fmt.Errorf("新文件名不能为空"), "新文件名不能为空")
		}

		// 检查新文件名是否与当前文件名相同
		if fileInfo.FileName == req.FileName {
			liberr.ErrIsNil(ctx, fmt.Errorf("新文件名与当前文件名相同"), "新文件名与当前文件名相同")
		}

		// 检查同级目录下是否存在同名文件
		count, err := dao.TemplateFiles.Ctx(ctx).Where(
			"template_id = ? AND parent_id = ? AND file_name = ? AND id != ?",
			fileInfo.TemplateId, fileInfo.ParentId, req.FileName, req.Id,
		).Count()
		liberr.ErrIsNil(ctx, err, "检查文件名冲突失败")
		if count > 0 {
			liberr.ErrIsNil(ctx, fmt.Errorf("同级目录下已存在同名文件"), "同级目录下已存在同名文件")
		}

		// 生成新的文件路径
		newFilePath := s.generateFilePath(ctx, fileInfo.TemplateId, fileInfo.ParentId, req.FileName)

		// 如果是目录，需要更新所有子文件的路径
		if fileInfo.IsDirectory == 1 {
			err = s.updateChildrenPaths(ctx, fileInfo.Id, fileInfo.FilePath, newFilePath)
			liberr.ErrIsNil(ctx, err, "更新子文件路径失败")
		}

		// 更新文件名和路径
		_, err = dao.TemplateFiles.Ctx(ctx).WherePri(req.Id).Update(do.TemplateFiles{
			FileName: req.FileName,
			FilePath: newFilePath,
		})
		liberr.ErrIsNil(ctx, err, "重命名文件失败")
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
		Fields("id, file_path, file_name, is_directory, parent_id, file_size, md5").
		Scan(&files)
	if err != nil {
		return
	}
	// 构建树
	idMap := make(map[int64]*api.FileTreeNode)
	for _, f := range files {
		node := &api.FileTreeNode{
			Id: f.Id, FilePath: f.FilePath, FileName: f.FileName, IsDirectory: f.IsDirectory,
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

	// 获取目录名
	dirName := filepath.Base(filePath)
	if dirName == "" || dirName == "." {
		// 如果是根目录或空路径，使用一个默认名称
		dirName = "root"
	}

	fmt.Printf("创建目录记录: filePath=%s, dirName=%s, parentId=%d\n", filePath, dirName, parentId)

	// 创建目录记录
	_, err = dao.TemplateFiles.Ctx(ctx).Insert(do.TemplateFiles{
		TemplateId:  templateId,
		FilePath:    filePath,
		FileName:    dirName,
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

	// 获取文件名
	fileName := filepath.Base(filePath)
	if fileName == "" {
		// 如果文件名为空，使用一个默认名称
		fileName = "untitled"
	}

	fmt.Printf("创建文件记录: filePath=%s, fileName=%s, parentId=%d, contentLength=%d\n", filePath, fileName, parentId, len(content))

	// 计算MD5
	md5 := gmd5.MustEncrypt(content)

	// 创建文件记录
	_, err = dao.TemplateFiles.Ctx(ctx).Insert(do.TemplateFiles{
		TemplateId:  templateId,
		FilePath:    filePath,
		FileName:    fileName,
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

// generateFilePath 根据父目录ID和文件名生成完整的文件路径
func (s *sTemplateFiles) generateFilePath(ctx context.Context, templateId interface{}, parentId int, fileName string) string {
	if parentId == 0 {
		// 根目录下的文件
		return fileName
	}

	// 查找父目录的路径
	var parent entity.TemplateFiles
	err := dao.TemplateFiles.Ctx(ctx).Where("id = ? AND template_id = ? AND is_directory = 1", parentId, templateId).Scan(&parent)
	if err != nil {
		// 找不到父目录，返回文件名
		return fileName
	}

	// 拼接父目录路径和文件名
	if parent.FilePath == "" {
		return fileName
	}
	return filepath.Join(parent.FilePath, fileName)
}

// updateChildrenPaths 更新目录下所有子文件的路径
func (s *sTemplateFiles) updateChildrenPaths(ctx context.Context, parentId int64, oldParentPath, newParentPath string) error {
	// 获取所有子文件
	var children []*entity.TemplateFiles
	err := dao.TemplateFiles.Ctx(ctx).Where("parent_id = ?", parentId).Scan(&children)
	if err != nil {
		return err
	}

	for _, child := range children {
		// 计算新的文件路径
		var newChildPath string
		if newParentPath == "" {
			newChildPath = child.FileName
		} else {
			newChildPath = filepath.Join(newParentPath, child.FileName)
		}

		// 更新子文件的路径
		_, err = dao.TemplateFiles.Ctx(ctx).WherePri(child.Id).Update(do.TemplateFiles{
			FilePath: newChildPath,
		})
		if err != nil {
			return err
		}

		// 如果子文件也是目录，递归更新其子文件
		if child.IsDirectory == 1 {
			err = s.updateChildrenPaths(ctx, child.Id, child.FilePath, newChildPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// UploadCode 上传代码文件
func (s *sTemplateFiles) UploadCode(ctx context.Context, req *api.TemplateFilesUploadCodeReq) (res *api.TemplateFilesUploadCodeRes, err error) {
	res = &api.TemplateFilesUploadCodeRes{}

	err = g.Try(ctx, func(ctx context.Context) {
		// 获取上传的文件
		file := g.RequestFromCtx(ctx).GetUploadFile("codeFile")
		if file == nil {
			liberr.ErrIsNil(ctx, fmt.Errorf("未找到上传的代码文件"), "请选择代码文件")
		}

		// 检查文件大小（1MB = 1024 * 1024 字节）
		if file.Size > 1024*1024 {
			liberr.ErrIsNil(ctx, fmt.Errorf("文件大小超过1MB限制"), "代码文件大小不能超过1MB")
		}

		// 获取文件名
		fileName := file.Filename
		if fileName == "" {
			liberr.ErrIsNil(ctx, fmt.Errorf("文件名不能为空"), "文件名不能为空")
		}

		// 读取文件内容
		src, err := file.Open()
		liberr.ErrIsNil(ctx, err, "打开上传文件失败")
		defer src.Close()
		fileContent, err := io.ReadAll(src)
		liberr.ErrIsNil(ctx, err, "读取文件内容失败")

		// 判断是否为文本文件
		isTextFile := s.isTextFile(fileName, fileContent)

		// 如果是文本文件，尝试解码为UTF-8字符串
		var content string
		if isTextFile {
			content = string(fileContent)
		}

		// 获取父目录ID
		parentId := int64(0)
		if req.ParentId != nil {
			parentId = gconv.Int64(req.ParentId)
		}

		// 生成文件路径
		filePath := s.generateFilePath(ctx, req.TemplateId, int(parentId), fileName)

		// 检查同级目录下是否存在同名文件
		count, err := dao.TemplateFiles.Ctx(ctx).Where(
			"template_id = ? AND parent_id = ? AND file_name = ?",
			req.TemplateId, parentId, fileName,
		).Count()
		liberr.ErrIsNil(ctx, err, "检查文件名冲突失败")
		if count > 0 {
			liberr.ErrIsNil(ctx, fmt.Errorf("同级目录下已存在同名文件"), "同级目录下已存在同名文件")
		}

		// 计算MD5
		md5 := gmd5.MustEncrypt(fileContent)

		// 创建文件记录
		_, err = dao.TemplateFiles.Ctx(ctx).Insert(do.TemplateFiles{
			TemplateId:  req.TemplateId,
			FilePath:    filePath,
			FileName:    fileName,
			FileContent: content,
			FileSize:    len(fileContent),
			IsDirectory: 0,
			Md5:         md5,
			Sort:        0,
			ParentId:    parentId,
		})
		liberr.ErrIsNil(ctx, err, "保存文件记录失败")

		// 设置返回结果
		res.FileName = fileName
		res.FileSize = len(fileContent)
		res.IsTextFile = isTextFile
		res.FileContent = content
		res.Message = "代码文件上传成功"
	})

	return
}

// isTextFile 判断文件是否为文本文件（仅检查前4KB内容，不再根据扩展名判断）
func (s sTemplateFiles) isTextFile(filename string, content []byte) bool {
	// 检查是否为文本文件
	textFileExtensions := []string{
		".txt", ".md", ".json", ".xml", ".yaml", ".yml", ".toml", ".ini", ".cfg", ".conf",
		".go", ".java", ".py", ".js", ".ts", ".vue", ".html", ".css", ".scss", ".less",
		".php", ".rb", ".rs", ".cpp", ".c", ".h", ".cs", ".swift", ".kt", ".scala",
		".sql", ".sh", ".bat", ".ps1", ".dockerfile", ".gitignore", ".env",
	}

	// 检查文件扩展名
	for _, ext := range textFileExtensions {
		if strings.HasSuffix(strings.ToLower(filename), ext) {
			return true
		}
	}

	// 检查文件内容是否包含null字节（二进制文件特征）
	if bytes.Contains(content, []byte{0}) {
		return false
	}

	// 检查是否包含大量可打印字符
	printableCount := 0
	totalCount := len(content)
	for _, b := range content {
		if b >= 32 && b <= 126 || b == 9 || b == 10 || b == 13 {
			printableCount++
		}
	}

	// 如果可打印字符占比超过90%，认为是文本文件
	return float64(printableCount)/float64(totalCount) > 0.9
}

// 获取模板函数映射
func (s sTemplateFiles) getTemplateFuncs() template.FuncMap {
	funcs := sprig.FuncMap()

	// 添加自定义函数
	funcs["formatGoPackage"] = func(packageName string) string {
		return strings.ReplaceAll(packageName, "-", "_")
	}

	funcs["formatJavaPackage"] = func(packageName string) string {
		return strings.ReplaceAll(packageName, "-", ".")
	}

	funcs["toSnakeCase"] = func(str string) string {
		// 转换为snake_case
		return strings.ToLower(strings.ReplaceAll(str, " ", "_"))
	}

	funcs["toCamelCase"] = func(str string) string {
		// 转换为camelCase
		words := strings.Fields(str)
		if len(words) == 0 {
			return ""
		}
		result := strings.ToLower(words[0])
		for i := 1; i < len(words); i++ {
			result += strings.Title(strings.ToLower(words[i]))
		}
		return result
	}

	funcs["toPascalCase"] = func(str string) string {
		// 转换为PascalCase
		words := strings.Fields(str)
		result := ""
		for _, word := range words {
			result += strings.Title(strings.ToLower(word))
		}
		return result
	}

	funcs["toKebabCase"] = func(str string) string {
		// 转换为kebab-case
		return strings.ToLower(strings.ReplaceAll(str, " ", "-"))
	}

	funcs["indent"] = func(spaces int, text string) string {
		// 缩进文本
		indent := strings.Repeat(" ", spaces)
		lines := strings.Split(text, "\n")
		for i, line := range lines {
			if line != "" {
				lines[i] = indent + line
			}
		}
		return strings.Join(lines, "\n")
	}

	funcs["wrapText"] = func(width int, text string) string {
		// 文本换行
		words := strings.Fields(text)
		if len(words) == 0 {
			return ""
		}

		var lines []string
		currentLine := words[0]

		for i := 1; i < len(words); i++ {
			if len(currentLine)+len(words[i])+1 <= width {
				currentLine += " " + words[i]
			} else {
				lines = append(lines, currentLine)
				currentLine = words[i]
			}
		}
		lines = append(lines, currentLine)

		return strings.Join(lines, "\n")
	}

	return funcs
}

// 渲染模板文件
func (s sTemplateFiles) Render(ctx context.Context, req *api.TemplateFilesRenderReq) (res *api.TemplateFilesRenderRes, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		// 1. 获取文件内容
		fileContent, err := s.GetFileContent(ctx, gconv.Int64(req.FileId))
		liberr.ErrIsNil(ctx, err, "获取文件内容失败")

		// 2. 获取文件信息
		fileInfo, err := s.GetById(ctx, gconv.Int64(req.FileId))
		liberr.ErrIsNil(ctx, err, "获取文件信息失败")

		// 3. 创建模板
		tmpl, err := template.New("template").Funcs(s.getTemplateFuncs()).Parse(fileContent)
		liberr.ErrIsNil(ctx, err, "解析模板失败")

		// 4. 渲染模板
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, req.TestVariables)
		liberr.ErrIsNil(ctx, err, "模板渲染失败")

		res = &api.TemplateFilesRenderRes{
			FileId:      gconv.Int64(req.FileId),
			FileName:    fileInfo.FileName,
			FileContent: buf.String(),
			Variables:   req.TestVariables,
		}
	})
	return
}
