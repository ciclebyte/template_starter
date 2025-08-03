package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ciclebyte/template_starter/cli/internal/client"
)

// Generator 项目生成器
type Generator struct {
	OutputDir string
	Force     bool
}

// NewGenerator 创建新的生成器实例
func NewGenerator(outputDir string, force bool) *Generator {
	return &Generator{
		OutputDir: outputDir,
		Force:     force,
	}
}

// GenerateProject 生成项目
func (g *Generator) GenerateProject(projectName string, renderedFiles []client.RenderedFile) error {
	projectDir := filepath.Join(g.OutputDir, projectName)
	
	// 检查目录是否存在
	if _, err := os.Stat(projectDir); err == nil && !g.Force {
		return fmt.Errorf("目录 %s 已存在，使用 --force 强制覆盖", projectDir)
	}
	
	// 创建项目目录
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("创建项目目录失败: %w", err)
	}
	
	// 生成文件
	for _, file := range renderedFiles {
		if err := g.writeFile(projectDir, file); err != nil {
			return fmt.Errorf("写入文件 %s 失败: %w", file.Path, err)
		}
	}
	
	fmt.Printf("项目 %s 创建成功！\n", projectName)
	return nil
}

// writeFile 写入单个文件
func (g *Generator) writeFile(projectDir string, file client.RenderedFile) error {
	fullPath := filepath.Join(projectDir, file.Path)
	
	if file.IsDirectory {
		// 创建目录
		return os.MkdirAll(fullPath, 0755)
	} else {
		// 创建文件
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建目录失败: %w", err)
		}
		
		return os.WriteFile(fullPath, []byte(file.Content), 0644)
	}
}

