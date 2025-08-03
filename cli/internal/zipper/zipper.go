package zipper

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Config 压缩器配置
type Config struct {
	SourcePath      string   // 源路径
	OutputFile      string   // 输出文件
	IncludeHidden   bool     // 包含隐藏文件
	IncludeBinary   bool     // 包含二进制文件
	ConfigFile      string   // 配置文件路径
	ExcludePatterns []string // 额外排除模式
	IncludePatterns []string // 强制包含模式
	Verbose         bool     // 详细输出
	DryRun          bool     // 预览模式
}

// Result 压缩结果
type Result struct {
	FileCount      int      // 文件数量
	DirCount       int      // 目录数量
	SkippedCount   int      // 跳过文件数量
	OriginalSize   int64    // 原始大小
	CompressedSize int64    // 压缩后大小
	SkippedFiles   []string // 跳过的文件列表
}

// Zipper 压缩器
type Zipper struct {
	config  *Config
	filter  *FileFilter
	result  *Result
	zipFile *os.File
	writer  *zip.Writer
}

// NewZipper 创建新的压缩器
func NewZipper(config *Config) (*Zipper, error) {
	// 创建文件过滤器
	filter, err := NewFileFilter(config)
	if err != nil {
		return nil, fmt.Errorf("创建文件过滤器失败: %w", err)
	}

	return &Zipper{
		config: config,
		filter: filter,
		result: &Result{
			SkippedFiles: make([]string, 0),
		},
	}, nil
}

// Zip 执行压缩
func (z *Zipper) Zip() (*Result, error) {
	// 检查源路径
	if _, err := os.Stat(z.config.SourcePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("源路径不存在: %s", z.config.SourcePath)
	}

	// 预览模式不创建实际文件
	if !z.config.DryRun {
		if err := z.createZipFile(); err != nil {
			return nil, err
		}
		defer z.closeZipFile()
	}

	// 遍历文件
	err := filepath.Walk(z.config.SourcePath, z.walkFunc)
	if err != nil {
		return nil, fmt.Errorf("遍历文件失败: %w", err)
	}

	// 获取压缩文件大小
	if !z.config.DryRun && z.zipFile != nil {
		if stat, err := z.zipFile.Stat(); err == nil {
			z.result.CompressedSize = stat.Size()
		}
	}

	return z.result, nil
}

// createZipFile 创建zip文件
func (z *Zipper) createZipFile() error {
	// 确保输出目录存在
	outputDir := filepath.Dir(z.config.OutputFile)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 创建zip文件
	file, err := os.Create(z.config.OutputFile)
	if err != nil {
		return fmt.Errorf("创建zip文件失败: %w", err)
	}

	z.zipFile = file
	z.writer = zip.NewWriter(file)
	return nil
}

// closeZipFile 关闭zip文件
func (z *Zipper) closeZipFile() {
	if z.writer != nil {
		z.writer.Close()
	}
	if z.zipFile != nil {
		z.zipFile.Close()
	}
}

// walkFunc 文件遍历函数
func (z *Zipper) walkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		if z.config.Verbose {
			fmt.Printf("跳过错误文件: %s (%v)\n", path, err)
		}
		z.result.SkippedCount++
		z.result.SkippedFiles = append(z.result.SkippedFiles, path+": "+err.Error())
		return nil // 继续处理其他文件
	}

	// 获取相对路径
	relPath, err := filepath.Rel(z.config.SourcePath, path)
	if err != nil {
		return err
	}

	// 跳过根目录
	if relPath == "." {
		return nil
	}

	// 检查是否应该排除
	if z.filter.ShouldExclude(path, info) {
		if z.config.Verbose {
			fmt.Printf("跳过: %s\n", relPath)
		}
		z.result.SkippedCount++
		z.result.SkippedFiles = append(z.result.SkippedFiles, relPath)
		
		// 如果是目录，跳过整个目录
		if info.IsDir() {
			return filepath.SkipDir
		}
		return nil
	}

	// 统计信息
	if info.IsDir() {
		z.result.DirCount++
	} else {
		z.result.FileCount++
		z.result.OriginalSize += info.Size()
	}

	if z.config.Verbose {
		if info.IsDir() {
			fmt.Printf("添加目录: %s\n", relPath)
		} else {
			fmt.Printf("添加文件: %s (%s)\n", relPath, formatSize(info.Size()))
		}
	}

	// 预览模式不写入文件
	if z.config.DryRun {
		return nil
	}

	// 添加到zip文件
	return z.addToZip(path, relPath, info)
}

// addToZip 添加文件到zip
func (z *Zipper) addToZip(filePath, zipPath string, info os.FileInfo) error {
	// 统一使用正斜杠作为路径分隔符
	zipPath = strings.ReplaceAll(zipPath, "\\", "/")

	// 创建zip文件头
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = zipPath
	header.Method = zip.Deflate

	// 如果是目录，确保路径以/结尾
	if info.IsDir() {
		if !strings.HasSuffix(header.Name, "/") {
			header.Name += "/"
		}
	}

	// 创建文件写入器
	writer, err := z.writer.CreateHeader(header)
	if err != nil {
		return err
	}

	// 如果是目录，不需要写入内容
	if info.IsDir() {
		return nil
	}

	// 打开源文件
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 复制文件内容
	_, err = io.Copy(writer, file)
	return err
}

// formatSize 格式化文件大小
func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}