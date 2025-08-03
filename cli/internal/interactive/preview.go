package interactive

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ciclebyte/template_starter/cli/internal/client"
)

type previewModel struct {
	files       []client.RenderedFile
	fileTree    []treeItem
	list        list.Model
	viewport    viewport.Model
	width       int
	height      int
	currentFile *client.RenderedFile
	projectName string
	confirmed   bool
	cancelled   bool
}

type treeItem struct {
	title    string
	desc     string
	file     *client.RenderedFile
	isDir    bool
	level    int
}

func (i treeItem) Title() string       { return i.title }
func (i treeItem) Description() string { return i.desc }
func (i treeItem) FilterValue() string { return i.title }

// PreviewFilesTUI 显示文件预览界面（TUI版本）
func PreviewFilesTUI(files []client.RenderedFile, projectName string) (bool, error) {
	// 构建文件树
	tree := buildFileTree(files)
	
	// 创建文件列表
	items := make([]list.Item, len(tree))
	for i, item := range tree {
		items[i] = item
	}

	// 创建列表组件
	l := list.New(items, newItemDelegate(), 40, 20)
	l.Title = "📁 文件预览"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	// 创建视口组件
	vp := viewport.New(80, 20)
	vp.SetContent("请选择左侧文件查看内容...")
	
	m := previewModel{
		files:       files,
		fileTree:    tree,
		list:        l,
		viewport:    vp,
		projectName: projectName,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	result, err := p.Run()
	if err != nil {
		return false, err
	}

	finalModel := result.(previewModel)
	return finalModel.confirmed, nil
}

func (m previewModel) Init() tea.Cmd {
	return nil
}

func (m previewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		// 调整布局
		listWidth := m.width / 3
		contentWidth := m.width - listWidth - 4
		
		m.list.SetWidth(listWidth)
		m.list.SetHeight(m.height - 6)
		
		m.viewport.Width = contentWidth
		m.viewport.Height = m.height - 6

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("q", "ctrl+c"))):
			m.cancelled = true
			return m, tea.Quit
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter", "y"))):
			m.confirmed = true
			return m, tea.Quit
		case key.Matches(msg, key.NewBinding(key.WithKeys("n"))):
			m.cancelled = true
			return m, tea.Quit
		case key.Matches(msg, key.NewBinding(key.WithKeys("tab"))):
			// 切换焦点（暂时不实现）
		case key.Matches(msg, key.NewBinding(key.WithKeys("r"))):
			// 刷新预览
			if m.currentFile != nil {
				m.viewport.SetContent(formatFileContent(*m.currentFile))
			}
		}
	}

	// 更新列表
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	
	// 检查选中的文件是否改变
	if selected := m.list.SelectedItem(); selected != nil {
		if item, ok := selected.(treeItem); ok && item.file != nil {
			if m.currentFile == nil || m.currentFile != item.file {
				m.currentFile = item.file
				m.viewport.SetContent(formatFileContent(*item.file))
			}
		}
	}

	// 更新视口
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m previewModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "正在初始化..."
	}

	// 简化布局，避免复杂边框
	listWidth := 40
	if m.width > 120 {
		listWidth = m.width / 3
	}
	contentWidth := m.width - listWidth - 4

	// 使用简单的边框样式，兼容Windows Terminal
	listBorder := strings.Repeat("-", listWidth)
	contentBorder := strings.Repeat("-", contentWidth)

	// 构建头部
	header := fmt.Sprintf("Project Preview: %s", m.projectName)
	
	// 构建文件列表部分
	listTitle := "File List"
	listContent := m.list.View()
	
	// 构建内容部分
	contentTitle := "File Content"
	contentView := m.viewport.View()
	
	// 组装左侧面板
	leftPanel := fmt.Sprintf("%s\n%s\n%s", listTitle, listBorder, listContent)
	
	// 组装右侧面板
	rightPanel := fmt.Sprintf("%s\n%s\n%s", contentTitle, contentBorder, contentView)
	
	// 使用简单的水平布局
	var body strings.Builder
	leftLines := strings.Split(leftPanel, "\n")
	rightLines := strings.Split(rightPanel, "\n")
	
	maxLines := len(leftLines)
	if len(rightLines) > maxLines {
		maxLines = len(rightLines)
	}
	
	for i := 0; i < maxLines && i < m.height-4; i++ {
		leftLine := ""
		rightLine := ""
		
		if i < len(leftLines) {
			leftLine = leftLines[i]
		}
		if i < len(rightLines) {
			rightLine = rightLines[i]
		}
		
		// 确保左侧宽度固定
		if len(leftLine) > listWidth {
			leftLine = leftLine[:listWidth-3] + "..."
		} else {
			leftLine += strings.Repeat(" ", listWidth-len(leftLine))
		}
		
		body.WriteString(fmt.Sprintf("%s  %s\n", leftLine, rightLine))
	}
	
	// 构建底部帮助
	help := "Up/Down: Navigate | Enter/y: Confirm | n: Cancel | r: Refresh | q: Quit"
	helpBorder := strings.Repeat("-", len(help))
	
	return fmt.Sprintf("%s\n%s\n\n%s\n%s\n%s",
		header,
		strings.Repeat("=", len(header)),
		body.String(),
		helpBorder,
		help,
	)
}

// buildFileTree 构建文件树结构
func buildFileTree(files []client.RenderedFile) []treeItem {
	var tree []treeItem
	
	// 按路径排序
	sort.Slice(files, func(i, j int) bool {
		return files[i].Path < files[j].Path
	})
	
	for i, file := range files {
		level := strings.Count(file.Path, "/")
		
		var desc string
		if file.IsDirectory {
			desc = "[DIR]"
		} else {
			desc = fmt.Sprintf("[FILE] %d bytes", len(file.Content))
		}
		
		// 使用简单的缩进和ASCII字符
		title := file.Path
		if level > 0 {
			title = strings.Repeat("  ", level) + "|- " + getFileName(file.Path)
		}
		
		tree = append(tree, treeItem{
			title: title,
			desc:  desc,
			file:  &files[i], // 正确的引用
			isDir: file.IsDirectory,
			level: level,
		})
	}
	
	return tree
}

// getFileName 从路径中提取文件名
func getFileName(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return path
}

// formatFileContent 格式化文件内容用于显示
func formatFileContent(file client.RenderedFile) string {
	if file.IsDirectory {
		return fmt.Sprintf("Directory: %s\n\nThis is a directory with no file content.", file.Path)
	}
	
	if file.Content == "" {
		return fmt.Sprintf("File: %s\n\nThis is an empty file.", file.Path)
	}
	
	header := fmt.Sprintf("File: %s\n%s\n\n", file.Path, strings.Repeat("-", 50))
	return header + file.Content
}

// newItemDelegate 创建列表项委托
func newItemDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	
	d.Styles.SelectedTitle = selectedItemStyle
	d.Styles.SelectedDesc = selectedItemStyle.Copy().Foreground(lipgloss.Color("244"))
	
	return d
}

// 样式定义已移至 styles.go 文件