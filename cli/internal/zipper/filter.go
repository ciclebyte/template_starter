package zipper

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// FileFilter 文件过滤器
type FileFilter struct {
	config           *Config
	excludePatterns  []*regexp.Regexp
	includePatterns  []*regexp.Regexp
	binaryExtensions map[string]bool
	excludeDirs      map[string]bool
	excludeFiles     map[string]bool
}

// NewFileFilter 创建文件过滤器
func NewFileFilter(config *Config) (*FileFilter, error) {
	filter := &FileFilter{
		config:          config,
		excludePatterns: make([]*regexp.Regexp, 0),
		includePatterns: make([]*regexp.Regexp, 0),
		binaryExtensions: make(map[string]bool),
		excludeDirs:     make(map[string]bool),
		excludeFiles:    make(map[string]bool),
	}

	// 初始化默认规则
	filter.initDefaultRules()

	// 加载配置文件
	if config.ConfigFile != "" {
		if err := filter.loadConfigFile(config.ConfigFile); err != nil {
			return nil, err
		}
	}

	// 编译命令行指定的模式
	if err := filter.compilePatterns(); err != nil {
		return nil, err
	}

	return filter, nil
}

// initDefaultRules 初始化默认规则
func (f *FileFilter) initDefaultRules() {
	// 默认排除的目录
	defaultExcludeDirs := []string{
		".git", ".svn", ".hg", ".bzr",           // 版本控制
		"node_modules", "bower_components",       // Node.js
		"vendor", "vendors",                      // 依赖包
		".vscode", ".idea", ".vs",               // IDE配置
		"__pycache__", ".pytest_cache",          // Python
		"target", "build", "dist", "out",        // 编译输出
		".gradle", ".maven",                     // Java构建工具
		"bin", "obj",                            // .NET
		".next", ".nuxt",                        // 前端框架
		"coverage", ".coverage",                 // 测试覆盖率
		"logs", "log",                           // 日志目录
		"tmp", "temp", ".tmp", ".temp",          // 临时目录
		".sass-cache", ".less-cache",            // CSS预处理器缓存
	}

	for _, dir := range defaultExcludeDirs {
		f.excludeDirs[dir] = true
	}

	// 默认排除的文件模式
	defaultExcludeFiles := []string{
		".DS_Store", "Thumbs.db", "desktop.ini", // 系统文件
		"*.log", "*.tmp", "*.temp",              // 临时和日志文件
		"*.swp", "*.swo", "*~",                  // 编辑器临时文件
		".env", ".env.local", ".env.*.local",    // 环境变量文件
		"*.pid", "*.lock",                       // 锁文件
		"core", "core.*",                        // 核心转储文件
		"*.orig", "*.rej",                       // 合并冲突文件
	}

	for _, pattern := range defaultExcludeFiles {
		f.excludeFiles[pattern] = true
	}

	// 常见的二进制文件扩展名
	binaryExts := []string{
		// 可执行文件
		".exe", ".dll", ".so", ".dylib", ".bin", ".app",
		// 归档文件
		".zip", ".tar", ".gz", ".rar", ".7z", ".bz2", ".xz",
		// 图片文件
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".ico", ".svg",
		// 音频文件
		".mp3", ".wav", ".flac", ".aac", ".ogg", ".m4a",
		// 视频文件
		".mp4", ".avi", ".mov", ".wmv", ".flv", ".mkv", ".webm",
		// 字体文件
		".ttf", ".otf", ".woff", ".woff2", ".eot",
		// 文档文件
		".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx",
		// 数据库文件
		".db", ".sqlite", ".sqlite3", ".mdb",
		// Java
		".class", ".jar", ".war", ".ear",
		// .NET
		".pdb", ".idb",
		// Python
		".pyc", ".pyo", ".pyd",
		// Node.js
		".node",
	}

	for _, ext := range binaryExts {
		f.binaryExtensions[ext] = true
	}
}

// loadConfigFile 加载配置文件
func (f *FileFilter) loadConfigFile(configPath string) error {
	file, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("打开配置文件失败: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentSection := ""

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 检查是否是节标题
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = strings.Trim(line, "[]")
			continue
		}

		// 根据节处理配置
		switch currentSection {
		case "exclude_dirs":
			f.excludeDirs[line] = true
		case "exclude_files":
			f.excludeFiles[line] = true
		case "binary_extensions":
			if !strings.HasPrefix(line, ".") {
				line = "." + line
			}
			f.binaryExtensions[line] = true
		}
	}

	return scanner.Err()
}

// compilePatterns 编译正则表达式模式
func (f *FileFilter) compilePatterns() error {
	// 编译排除模式
	for _, pattern := range f.config.ExcludePatterns {
		regex, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("编译排除模式失败 '%s': %w", pattern, err)
		}
		f.excludePatterns = append(f.excludePatterns, regex)
	}

	// 编译包含模式
	for _, pattern := range f.config.IncludePatterns {
		regex, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("编译包含模式失败 '%s': %w", pattern, err)
		}
		f.includePatterns = append(f.includePatterns, regex)
	}

	return nil
}

// ShouldExclude 判断是否应该排除文件或目录
func (f *FileFilter) ShouldExclude(path string, info os.FileInfo) bool {
	fileName := info.Name()
	
	// 强制包含模式优先级最高
	if f.matchesIncludePatterns(path) {
		return false
	}

	// 检查隐藏文件
	if !f.config.IncludeHidden && f.isHidden(fileName) {
		return true
	}

	// 检查目录
	if info.IsDir() {
		return f.shouldExcludeDir(fileName)
	}

	// 检查文件
	return f.shouldExcludeFile(path, fileName, info)
}

// isHidden 判断是否是隐藏文件
func (f *FileFilter) isHidden(fileName string) bool {
	return strings.HasPrefix(fileName, ".")
}

// shouldExcludeDir 判断是否排除目录
func (f *FileFilter) shouldExcludeDir(dirName string) bool {
	// 检查默认排除目录
	if f.excludeDirs[dirName] {
		return true
	}

	// 检查自定义排除模式
	return f.matchesExcludePatterns(dirName)
}

// shouldExcludeFile 判断是否排除文件
func (f *FileFilter) shouldExcludeFile(path, fileName string, info os.FileInfo) bool {
	// 检查文件名模式
	for pattern := range f.excludeFiles {
		if matched, _ := filepath.Match(pattern, fileName); matched {
			return true
		}
	}

	// 检查自定义排除模式
	if f.matchesExcludePatterns(path) {
		return true
	}

	// 检查二进制文件
	if !f.config.IncludeBinary && f.isBinaryFile(path, fileName, info) {
		return true
	}

	return false
}

// isBinaryFile 判断是否是二进制文件
func (f *FileFilter) isBinaryFile(path, fileName string, info os.FileInfo) bool {
	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(fileName))
	if f.binaryExtensions[ext] {
		return true
	}

	// 检查文件内容（简单启发式）
	if info.Size() > 0 && info.Size() < 8192 { // 只检查小文件
		return f.isBinaryContent(path)
	}

	return false
}

// isBinaryContent 通过内容判断是否是二进制文件
func (f *FileFilter) isBinaryContent(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	// 读取前512字节
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return false
	}

	// 检查是否包含null字节（简单的二进制检测）
	for i := 0; i < n; i++ {
		if buffer[i] == 0 {
			return true
		}
	}

	return false
}

// matchesExcludePatterns 检查是否匹配排除模式
func (f *FileFilter) matchesExcludePatterns(path string) bool {
	for _, pattern := range f.excludePatterns {
		if pattern.MatchString(path) {
			return true
		}
	}
	return false
}

// matchesIncludePatterns 检查是否匹配包含模式
func (f *FileFilter) matchesIncludePatterns(path string) bool {
	// 如果没有定义包含模式，则不强制包含
	if len(f.includePatterns) == 0 {
		return false
	}

	for _, pattern := range f.includePatterns {
		if pattern.MatchString(path) {
			return true
		}
	}
	return false
}