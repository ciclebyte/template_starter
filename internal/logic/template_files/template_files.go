package template_files

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
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
	"github.com/gogf/gf/v2/errors/gerror"
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
		// 检查文件是否存在
		_, err = s.GetById(ctx, gconv.Int64(req.Id))
		liberr.ErrIsNil(ctx, err, "获取模板文件失败")

		// 只更新文件内容和相关字段
		_, err = dao.TemplateFiles.Ctx(ctx).WherePri(req.Id).Update(do.TemplateFiles{
			FileContent: req.FileContent,                         // 文件内容
			FileSize:    len(req.FileContent),                    // 重新计算文件大小
			Md5:         gmd5.MustEncryptString(req.FileContent), // 重新计算MD5
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

	// 合并内置函数
	builtinFuncs := service.BuiltinFunctions().GetTemplateFuncMap()
	for name, fn := range builtinFuncs {
		funcs[name] = fn
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
		err = tmpl.Execute(&buf, req.Variables)
		liberr.ErrIsNil(ctx, err, "模板渲染失败")

		res = &api.TemplateFilesRenderRes{
			FileId:      gconv.Int64(req.FileId),
			FileName:    fileInfo.FileName,
			FileContent: buf.String(),
			Variables:   req.Variables,
		}
	})
	return
}

// RenderFileTree 渲染整个文件树
// renderTemplateFiles 通用模板文件渲染函数
func (s sTemplateFiles) renderTemplateFiles(ctx context.Context, templateId int64, variables map[string]interface{}) ([]*api.RenderFileInfo, error) {
	// 1. 获取模板下的所有文件
	var files []*entity.TemplateFiles
	err := dao.TemplateFiles.Ctx(ctx).Where("template_id = ?", templateId).Scan(&files)
	if err != nil {
		return nil, err
	}

	// 2. 渲染并重建文件树
	renderedFiles := s.renderAndRebuildTree(files, variables)

	return renderedFiles, nil
}

func (s sTemplateFiles) RenderFileTree(ctx context.Context, req *api.TemplateFilesRenderFileTreeReq) (res *api.TemplateFilesRenderFileTreeRes, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		templateId := gconv.Int64(req.TemplateId)

		res = &api.TemplateFilesRenderFileTreeRes{
			TemplateId: templateId,
			Tree:       []*api.RenderFileInfo{},
			Variables:  req.Variables,
			TotalFiles: 0,
			TotalSize:  0,
		}

		// 使用通用渲染函数
		renderedFiles, err := s.renderTemplateFiles(ctx, templateId, req.Variables)
		liberr.ErrIsNil(ctx, err, "渲染模板文件失败")

		// 3. 构建树形结构
		res.Tree = s.buildTree(renderedFiles)

		// 4. 计算统计信息
		for _, file := range renderedFiles {
			res.TotalFiles++
			res.TotalSize += int64(file.FileSize)
		}
	})
	return
}

// renderAndRebuildTree 渲染文件并重建树结构
func (s sTemplateFiles) renderAndRebuildTree(files []*entity.TemplateFiles, variables map[string]interface{}) []*api.RenderFileInfo {
	fmt.Printf("=== 开始渲染和重建树形结构 ===\n")
	fmt.Printf("总文件数: %d\n", len(files))
	fmt.Printf("变量数据: %+v\n", variables)

	var result []*api.RenderFileInfo
	var nextId int64 = 10000
	pathToNode := make(map[string]*api.RenderFileInfo)
	originalPathToFinalPath := make(map[string]string) // 记录原始路径到最终路径的映射

	// 第一步：处理所有文件，渲染并创建基础节点
	fmt.Printf("\n=== 第一步：处理所有文件 ===\n")
	for i, file := range files {
		fmt.Printf("\n--- 处理文件 %d/%d ---\n", i+1, len(files))
		fmt.Printf("原始数据: ID=%d, 名称=%s, 路径=%s, 是否目录=%d, 父ID=%d\n",
			file.Id, file.FileName, file.FilePath, file.IsDirectory, file.ParentId)

		// 修复变量格式：将 {{/var}} 转换为 {{.var}}
		fileNameToRender := file.FileName
		filePathToRender := file.FilePath

		if strings.Contains(fileNameToRender, "{{/") {
			fileNameToRender = strings.ReplaceAll(fileNameToRender, "{{/", "{{.")
			fmt.Printf("  修复文件名变量格式: %s\n", fileNameToRender)
		}
		if strings.Contains(filePathToRender, "{{/") {
			filePathToRender = strings.ReplaceAll(filePathToRender, "{{/", "{{.")
			fmt.Printf("  修复文件路径变量格式: %s\n", filePathToRender)
		}

		// 渲染文件名和路径
		renderedName := file.FileName
		renderedPath := file.FilePath

		if fileNameToRender != "" {
			if tmpl, err := template.New("fileName").Funcs(s.getTemplateFuncs()).Parse(fileNameToRender); err == nil {
				var buf bytes.Buffer
				if tmpl.Execute(&buf, variables) == nil {
					renderedName = buf.String()
					fmt.Printf("  文件名渲染: %s -> %s\n", fileNameToRender, renderedName)
				} else {
					fmt.Printf("  文件名渲染失败: %v\n", err)
				}
			} else {
				fmt.Printf("  文件名模板解析失败: %v\n", err)
			}
		}

		if filePathToRender != "" {
			if tmpl, err := template.New("filePath").Funcs(s.getTemplateFuncs()).Parse(filePathToRender); err == nil {
				var buf bytes.Buffer
				if tmpl.Execute(&buf, variables) == nil {
					renderedPath = buf.String()
					fmt.Printf("  文件路径渲染: %s -> %s\n", filePathToRender, renderedPath)
				} else {
					fmt.Printf("  文件路径渲染失败: %v\n", err)
				}
			} else {
				fmt.Printf("  文件路径模板解析失败: %v\n", err)
			}
		}

		// 渲染文件内容
		renderedContent := file.FileContent
		if file.IsDirectory == 0 && file.FileContent != "" {
			contentToRender := file.FileContent
			if strings.Contains(contentToRender, "{{/") {
				contentToRender = strings.ReplaceAll(contentToRender, "{{/", "{{.")
			}

			if tmpl, err := template.New("fileContent").Funcs(s.getTemplateFuncs()).Parse(contentToRender); err == nil {
				var buf bytes.Buffer
				if tmpl.Execute(&buf, variables) == nil {
					renderedContent = buf.String()
				}
			}
		}

		// 检查是否需要分割目录名
		fmt.Printf("  检查是否需要分割目录名: 是否目录=%d, 渲染后名称=%s, 包含点号=%t\n",
			file.IsDirectory, renderedName, s.containsDots(renderedName))

		if file.IsDirectory == 1 && s.containsDots(renderedName) {
			// 需要分割的目录，创建多级目录结构
			fmt.Printf("  >>> 检测到需要分割的目录名: %s\n", renderedName)

			// 获取父路径
			parentPath := s.getParentPath(renderedPath)
			fmt.Printf("    父路径: %s\n", parentPath)
			currentPath := parentPath
			var parentId int64 = 0

			// 如果父路径存在，找到父节点
			if parentPath != "" {
				if parentNode, exists := pathToNode[parentPath]; exists {
					parentId = parentNode.Id
					currentPath = parentPath
					fmt.Printf("    找到父节点: ID=%d, 名称=%s\n", parentNode.Id, parentNode.FileName)
				} else {
					fmt.Printf("    父路径存在但找不到父节点: %s\n", parentPath)
				}
			} else {
				fmt.Printf("    没有父路径\n")
			}

			// 分割目录名并创建中间节点
			parts := strings.Split(renderedName, ".")
			fmt.Printf("    分割结果: %v\n", parts)

			for i, part := range parts {
				fmt.Printf("    处理第 %d 部分: %s\n", i+1, part)

				if currentPath == "" {
					currentPath = part
				} else {
					currentPath = currentPath + "/" + part
				}
				fmt.Printf("    当前路径: %s\n", currentPath)

				// 检查路径是否已存在
				if existingNode, exists := pathToNode[currentPath]; exists {
					fmt.Printf("    路径已存在，跳过: %s (ID: %d)\n", currentPath, existingNode.Id)
					parentId = existingNode.Id
					continue
				}

				// 创建新节点
				newNode := &api.RenderFileInfo{
					Id:          nextId,
					FilePath:    currentPath,
					FileName:    part,
					FileContent: "",
					FileSize:    0,
					IsDirectory: 1,
					ParentId:    int(parentId),
				}

				pathToNode[currentPath] = newNode
				result = append(result, newNode)
				fmt.Printf("    >>> 创建中间目录: %s (路径: %s, 父ID: %d, 新ID: %d)\n", part, currentPath, newNode.ParentId, nextId)
				parentId = nextId
				nextId++
			}

			// 记录最终路径，用于后续文件路径修正
			finalPath := currentPath
			originalPathToFinalPath[renderedPath] = finalPath
			fmt.Printf("    >>> 记录路径映射: %s -> %s\n", renderedPath, finalPath)

		} else {
			// 普通文件或不需要分割的目录
			// 检查是否需要修正路径（如果父目录被分割了）
			finalPath := renderedPath

			// 检查路径中是否包含需要替换的部分
			for originalPath, mappedPath := range originalPathToFinalPath {
				if strings.Contains(renderedPath, originalPath) {
					finalPath = strings.Replace(renderedPath, originalPath, mappedPath, 1)
					fmt.Printf("  >>> 修正文件路径: %s -> %s (替换: %s -> %s)\n", renderedPath, finalPath, originalPath, mappedPath)
					break
				}
			}

			// 确保目录属性正确
			isDirectory := file.IsDirectory
			if isDirectory == 1 {
				fmt.Printf("  >>> 创建目录节点: %s (路径: %s)\n", renderedName, finalPath)
			} else {
				fmt.Printf("  >>> 创建文件节点: %s (路径: %s)\n", renderedName, finalPath)
			}

			newNode := &api.RenderFileInfo{
				Id:          nextId,
				FilePath:    finalPath,
				FileName:    renderedName,
				FileContent: renderedContent,
				FileSize:    len(renderedContent),
				IsDirectory: isDirectory,
				ParentId:    0, // 先设为0，后面修正
			}

			pathToNode[finalPath] = newNode
			result = append(result, newNode)
			fmt.Printf("  >>> 普通节点创建完成: ID=%d, 是否目录=%d\n", nextId, isDirectory)
			nextId++
		}
	}

	// 第二步：根据路径重新设置父ID
	fmt.Printf("\n=== 第二步：重新设置父ID ===\n")
	fmt.Printf("当前节点总数: %d\n", len(result))

	for i, node := range result {
		fmt.Printf("\n--- 处理节点 %d/%d ---\n", i+1, len(result))
		fmt.Printf("节点: ID=%d, 名称=%s, 路径=%s, 当前父ID=%d\n", node.Id, node.FileName, node.FilePath, node.ParentId)

		parentPath := s.getParentPath(node.FilePath)
		fmt.Printf("计算父路径: %s\n", parentPath)

		if parentPath != "" {
			// 尝试多种路径格式匹配
			var parentNode *api.RenderFileInfo
			var found bool

			// 尝试直接匹配
			if parentNode, found = pathToNode[parentPath]; found {
				node.ParentId = int(parentNode.Id)
				fmt.Printf(">>> 设置父子关系: %s -> %s (ID: %d -> %d)\n", node.FileName, parentNode.FileName, node.Id, parentNode.Id)
			} else {
				// 尝试统一分隔符后匹配
				normalizedParentPath := strings.ReplaceAll(parentPath, "\\", "/")
				if parentNode, found = pathToNode[normalizedParentPath]; found {
					node.ParentId = int(parentNode.Id)
					fmt.Printf(">>> 设置父子关系(统一分隔符): %s -> %s (ID: %d -> %d)\n", node.FileName, parentNode.FileName, node.Id, parentNode.Id)
				} else {
					// 尝试反向分隔符匹配
					reverseParentPath := strings.ReplaceAll(parentPath, "/", "\\")
					if parentNode, found = pathToNode[reverseParentPath]; found {
						node.ParentId = int(parentNode.Id)
						fmt.Printf(">>> 设置父子关系(反向分隔符): %s -> %s (ID: %d -> %d)\n", node.FileName, parentNode.FileName, node.Id, parentNode.Id)
					} else {
						// 如果找不到父节点，设为根节点
						node.ParentId = 0
						fmt.Printf(">>> 找不到父节点，设为根节点: %s (父路径: %s, 尝试了: %s, %s)\n", node.FileName, parentPath, normalizedParentPath, reverseParentPath)
					}
				}
			}
		} else {
			// 没有父路径，设为根节点
			node.ParentId = 0
			fmt.Printf(">>> 没有父路径，设为根节点: %s\n", node.FileName)
		}
	}

	// 调试信息
	fmt.Printf("\n=== 最终结果 ===\n")
	fmt.Printf("渲染和重建完成，总文件数: %d\n", len(result))

	fmt.Printf("\n所有节点列表:\n")
	for i, file := range result {
		fmt.Printf("%d. ID=%d, 名称=%s, 路径=%s, 父ID=%d, 是否目录=%d\n",
			i+1, file.Id, file.FileName, file.FilePath, file.ParentId, file.IsDirectory)
	}

	fmt.Printf("\n路径映射表:\n")
	for path, node := range pathToNode {
		fmt.Printf("路径: %s -> 节点: ID=%d, 名称=%s\n", path, node.Id, node.FileName)
	}

	fmt.Printf("\n=== 渲染和重建树形结构完成 ===\n")
	return result
}

// DownloadZip 下载ZIP包
func (s sTemplateFiles) DownloadZip(ctx context.Context, req *api.TemplateFilesDownloadZipReq) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		templateId := gconv.Int64(req.TemplateId)

		// 1. 获取模板信息
		var templateInfo entity.Templates
		err := dao.Templates.Ctx(ctx).Where("id = ?", templateId).Scan(&templateInfo)
		liberr.ErrIsNil(ctx, err, "获取模板信息失败")

		// 2. 使用通用渲染函数获取渲染后的文件
		renderedFiles, err := s.renderTemplateFiles(ctx, templateId, req.Variables)
		liberr.ErrIsNil(ctx, err, "渲染模板文件失败")

		// 3. 确定ZIP文件名
		zipFileName := req.FileName
		if zipFileName == "" {
			zipFileName = templateInfo.Name
		}
		if !strings.HasSuffix(zipFileName, ".zip") {
			zipFileName += ".zip"
		}

		// 4. 创建临时目录
		tempDir := g.Cfg().MustGet(ctx, "server.tempPath", "").String()
		if tempDir == "" {
			tempDir = os.TempDir()
		}
		if err := os.MkdirAll(tempDir, 0755); err != nil {
			liberr.ErrIsNil(ctx, err, "创建临时目录失败")
		}

		// 5. 创建ZIP文件
		timestamp := gtime.Timestamp()
		zipPath := filepath.Join(tempDir, fmt.Sprintf("template_%d_%d.zip", templateId, timestamp))

		zipFile, err := os.Create(zipPath)
		liberr.ErrIsNil(ctx, err, "创建ZIP文件失败")
		defer zipFile.Close()

		zipWriter := zip.NewWriter(zipFile)
		defer zipWriter.Close()

		// 6. 将渲染后的文件添加到ZIP中
		fileCount := 0
		for _, file := range renderedFiles {
			if file.IsDirectory == 1 {
				// 跳过目录，只添加文件
				continue
			}

			// 统一路径分隔符为正斜杠
			zipPath := strings.ReplaceAll(file.FilePath, "\\", "/")

			// 创建ZIP文件条目
			zipEntry, err := zipWriter.Create(zipPath)
			if err != nil {
				fmt.Printf("创建ZIP条目失败: %s, 错误: %v\n", zipPath, err)
				continue
			}

			// 写入文件内容
			_, err = zipEntry.Write([]byte(file.FileContent))
			if err != nil {
				fmt.Printf("写入ZIP文件内容失败: %s, 错误: %v\n", zipPath, err)
				continue
			}

			fileCount++
			fmt.Printf("成功添加文件到ZIP: %s\n", zipPath)
		}

		fmt.Printf("ZIP文件创建完成，共添加 %d 个文件\n", fileCount)

		// 7. 关闭ZIP写入器
		zipWriter.Close()

		// 8. 设置响应头并返回ZIP文件
		response := g.RequestFromCtx(ctx).Response
		response.Header().Set("Content-Type", "application/zip")
		response.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", zipFileName))
		response.Header().Set("Content-Transfer-Encoding", "binary")

		// 9. 读取ZIP文件并写入响应
		zipData, err := os.ReadFile(zipPath)
		liberr.ErrIsNil(ctx, err, "读取ZIP文件失败")

		response.Write(zipData)

		// 10. 清理临时文件
		os.Remove(zipPath)
	})
	return
}

// buildTree 构建树形结构
func (s sTemplateFiles) buildTree(files []*api.RenderFileInfo) []*api.RenderFileInfo {
	// 创建ID到节点的映射
	idMap := make(map[int64]*api.RenderFileInfo)
	var rootNodes []*api.RenderFileInfo

	// 初始化所有节点
	for _, file := range files {
		file.Children = []*api.RenderFileInfo{}
		idMap[file.Id] = file
	}

	// 构建父子关系
	for _, file := range files {
		if file.ParentId == 0 {
			// 根节点
			rootNodes = append(rootNodes, file)
		} else {
			// 子节点，添加到父节点的children中
			parentId := int64(file.ParentId)
			if parent, exists := idMap[parentId]; exists {
				parent.Children = append(parent.Children, file)
			} else {
				// 如果找不到父节点，作为根节点处理
				fmt.Printf("警告: 找不到父节点 ID=%d，将 %s 作为根节点\n", parentId, file.FileName)
				rootNodes = append(rootNodes, file)
			}
		}
	}

	// 调试信息
	fmt.Printf("构建树形结构完成，根节点数量: %d\n", len(rootNodes))
	for i, root := range rootNodes {
		fmt.Printf("根节点 %d: ID=%d, 名称=%s, 子节点数量=%d\n", i, root.Id, root.FileName, len(root.Children))
	}

	// 对每个节点的children进行排序
	var sortChildren func(nodes []*api.RenderFileInfo)
	sortChildren = func(nodes []*api.RenderFileInfo) {
		for _, node := range nodes {
			if len(node.Children) > 0 {
				// 排序：目录在前，文件在后，同类型按名称排序
				sort.Slice(node.Children, func(i, j int) bool {
					if node.Children[i].IsDirectory != node.Children[j].IsDirectory {
						return node.Children[i].IsDirectory > node.Children[j].IsDirectory
					}
					return node.Children[i].FileName < node.Children[j].FileName
				})
				sortChildren(node.Children)
			}
		}
	}
	sortChildren(rootNodes)

	return rootNodes
}

// containsTemplateVariables 检查字符串是否包含模板变量
func (s sTemplateFiles) containsTemplateVariables(str string) bool {
	return strings.Contains(str, "{{") && strings.Contains(str, "}}")
}

// containsDots 检查字符串是否包含点号
func (s sTemplateFiles) containsDots(str string) bool {
	return strings.Contains(str, ".")
}

// splitPath 分割路径
func (s sTemplateFiles) splitPath(path string) []string {
	// 统一路径分隔符
	path = strings.ReplaceAll(path, "\\", "/")

	// 分割路径
	parts := strings.Split(path, "/")

	// 过滤空字符串并进一步分割包含点的部分
	var result []string
	for _, part := range parts {
		if part != "" {
			// 如果部分包含点，进一步分割
			if strings.Contains(part, ".") {
				dotParts := strings.Split(part, ".")
				for _, dotPart := range dotParts {
					if dotPart != "" {
						result = append(result, dotPart)
					}
				}
			} else {
				result = append(result, part)
			}
		}
	}

	return result
}

// getParentPath 获取父路径
func (s sTemplateFiles) getParentPath(path string) string {
	// 统一路径分隔符
	path = strings.ReplaceAll(path, "\\", "/")

	// 找到最后一个斜杠的位置
	lastSlash := strings.LastIndex(path, "/")
	if lastSlash == -1 {
		// 没有斜杠，说明是根路径
		return ""
	}

	// 返回父路径
	return path[:lastSlash]
}

func (s *sTemplateFiles) Move(ctx context.Context, req *api.TemplateFilesMoveReq) (err error) {
	// 1. 验证参数
	fileId := gconv.Int64(req.Id)
	newParentId := gconv.Int64(req.NewParentId)

	if fileId <= 0 {
		return gerror.New("文件ID无效")
	}

	// 2. 获取要移动的文件信息
	fileInfo, err := dao.TemplateFiles.Ctx(ctx).Where(dao.TemplateFiles.Columns().Id, fileId).One()
	if err != nil {
		return gerror.Wrap(err, "获取文件信息失败")
	}
	if fileInfo.IsEmpty() {
		return gerror.New("文件不存在")
	}

	// 3. 如果新父目录ID不为0，验证新父目录是否存在且为目录
	if newParentId != 0 {
		parentInfo, err := dao.TemplateFiles.Ctx(ctx).Where(dao.TemplateFiles.Columns().Id, newParentId).One()
		if err != nil {
			return gerror.Wrap(err, "获取父目录信息失败")
		}
		if parentInfo.IsEmpty() {
			return gerror.New("目标父目录不存在")
		}
		if parentInfo["is_directory"].Int() != 1 {
			return gerror.New("目标必须是目录")
		}

		// 验证不能移动到自己的子目录
		if err := s.validateNotMoveToChild(ctx, fileId, newParentId); err != nil {
			return err
		}
	}

	// 4. 更新文件的父目录ID
	_, err = dao.TemplateFiles.Ctx(ctx).
		Where(dao.TemplateFiles.Columns().Id, fileId).
		Update(g.Map{
			dao.TemplateFiles.Columns().ParentId:  newParentId,
			dao.TemplateFiles.Columns().UpdatedAt: gtime.Now(),
		})

	if err != nil {
		return gerror.Wrap(err, "移动文件失败")
	}

	// 5. 重新生成文件路径
	newFilePath, err := s.buildFilePath(ctx, fileId, gconv.Int64(fileInfo["template_id"]))
	if err != nil {
		return gerror.Wrap(err, "生成新文件路径失败")
	}

	// 6. 更新文件路径
	_, err = dao.TemplateFiles.Ctx(ctx).
		Where(dao.TemplateFiles.Columns().Id, fileId).
		Update(g.Map{
			dao.TemplateFiles.Columns().FilePath:  newFilePath,
			dao.TemplateFiles.Columns().UpdatedAt: gtime.Now(),
		})

	if err != nil {
		return gerror.Wrap(err, "更新文件路径失败")
	}

	// 7. 如果移动的是目录，需要更新其所有子项的文件路径
	if fileInfo["is_directory"].Int() == 1 {
		oldFilePath := fileInfo["file_path"].String()
		if err := s.updateChildrenPaths(ctx, fileId, oldFilePath, newFilePath); err != nil {
			return gerror.Wrap(err, "更新子项路径失败")
		}
	}

	return nil
}

// 验证不能移动到自己的子目录
func (s *sTemplateFiles) validateNotMoveToChild(ctx context.Context, fileId, newParentId int64) error {
	// 递归查找所有子目录
	var checkParent func(int64) error
	checkParent = func(parentId int64) error {
		if parentId == 0 {
			return nil
		}
		if parentId == fileId {
			return gerror.New("不能移动到自己的子目录")
		}

		// 获取父目录信息
		parentInfo, err := dao.TemplateFiles.Ctx(ctx).Where(dao.TemplateFiles.Columns().Id, parentId).One()
		if err != nil {
			return err
		}
		if parentInfo.IsEmpty() {
			return nil
		}

		// 递归检查上级目录
		return checkParent(parentInfo["parent_id"].Int64())
	}

	return checkParent(newParentId)
}

// 更新单个文件的路径
func (s *sTemplateFiles) updateFilePath(ctx context.Context, fileId, templateId int64) error {
	// 构建新的文件路径
	filePath, err := s.buildFilePath(ctx, fileId, templateId)
	if err != nil {
		return err
	}

	// 更新文件路径
	_, err = dao.TemplateFiles.Ctx(ctx).
		Where(dao.TemplateFiles.Columns().Id, fileId).
		Update(g.Map{
			dao.TemplateFiles.Columns().FilePath:  filePath,
			dao.TemplateFiles.Columns().UpdatedAt: gtime.Now(),
		})

	return err
}


// 构建文件路径
func (s *sTemplateFiles) buildFilePath(ctx context.Context, fileId, templateId int64) (string, error) {
	// 获取文件信息
	fileInfo, err := dao.TemplateFiles.Ctx(ctx).Where(dao.TemplateFiles.Columns().Id, fileId).One()
	if err != nil {
		return "", err
	}
	if fileInfo.IsEmpty() {
		return "", gerror.New("文件不存在")
	}

	fileName := fileInfo["file_name"].String()
	parentId := fileInfo["parent_id"].Int64()

	// 如果没有父目录，返回文件名
	if parentId == 0 {
		return fileName, nil
	}

	// 递归构建父目录路径
	parentPath, err := s.buildFilePath(ctx, parentId, templateId)
	if err != nil {
		return "", err
	}

	// 组合路径
	if parentPath == "" {
		return fileName, nil
	}

	return parentPath + "/" + fileName, nil
}
