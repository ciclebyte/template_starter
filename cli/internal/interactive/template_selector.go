package interactive

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ciclebyte/template_starter/cli/internal/client"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

// TemplateItem 模板项
type TemplateItem struct {
	Template client.Template
}

func (i TemplateItem) FilterValue() string { 
	return i.Template.Name + " " + i.Template.Description + " " + i.Template.Category
}

func (i TemplateItem) Title() string       { return i.Template.Name }
func (i TemplateItem) Description() string { 
	return fmt.Sprintf("%s • %s", i.Template.Category, i.Template.Description)
}

// TemplateSelectorModel TUI模型
type TemplateSelectorModel struct {
	list        list.Model
	choice      *client.Template
	quitting    bool
	searchInput textinput.Model
	searching   bool
	client      *client.Client
	allItems    []list.Item
}

// TemplateSelectedMsg 模板选择消息
type TemplateSelectedMsg struct {
	Template *client.Template
}

// NewTemplateSelectorModel 创建新的模板选择器模型
func NewTemplateSelectorModel(apiClient *client.Client) *TemplateSelectorModel {
	// 创建搜索输入框
	searchInput := textinput.New()
	searchInput.Placeholder = "搜索模板..."
	searchInput.CharLimit = 50
	searchInput.Width = 50

	// 创建列表
	l := list.New([]list.Item{}, itemDelegate{}, 80, 14)
	l.Title = "选择模板"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false) // 我们使用自定义搜索
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return &TemplateSelectorModel{
		list:        l,
		searchInput: searchInput,
		client:      apiClient,
		allItems:    []list.Item{},
	}
}

// itemDelegate 列表项委托
type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	if i, ok := item.(TemplateItem); ok {
		str := fmt.Sprintf("%s", i.Title())
		
		isSelected := index == m.Index()
		if isSelected {
			str = selectedItemStyle.Render("> " + str)
		} else {
			str = itemStyle.Render(str)
		}
		
		fmt.Fprint(w, str)
	}
}

// Init 初始化
func (m *TemplateSelectorModel) Init() tea.Cmd {
	return tea.Batch(
		m.loadTemplates(),
		textinput.Blink,
	)
}

// loadTemplates 加载模板列表
func (m *TemplateSelectorModel) loadTemplates() tea.Cmd {
	return func() tea.Msg {
		// TODO: 实际调用API获取模板列表
		// 这里先返回模拟数据
		templates := []client.Template{
			{
				ID:          "1",
				Name:        "go-web",
				Description: "Go Web应用模板",
				Category:    "backend",
			},
			{
				ID:          "2", 
				Name:        "vue3-admin",
				Description: "Vue3后台管理系统",
				Category:    "frontend",
			},
			{
				ID:          "3",
				Name:        "microservice",
				Description: "微服务架构模板",
				Category:    "backend",
			},
			{
				ID:          "4",
				Name:        "react-app",
				Description: "React应用模板",
				Category:    "frontend",
			},
		}
		
		return templatesLoadedMsg{templates: templates}
	}
}

// templatesLoadedMsg 模板加载完成消息
type templatesLoadedMsg struct {
	templates []client.Template
}

// Update 更新
func (m *TemplateSelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch {
		case msg.String() == "ctrl+c", msg.String() == "q":
			m.quitting = true
			return m, tea.Quit

		case msg.String() == "/":
			// 进入搜索模式
			m.searching = true
			m.searchInput.Focus()
			return m, nil

		case msg.String() == "esc":
			if m.searching {
				// 退出搜索模式
				m.searching = false
				m.searchInput.Blur()
				m.searchInput.SetValue("")
				m.filterTemplates("")
				return m, nil
			}

		case msg.String() == "enter":
			if m.searching {
				// 执行搜索
				m.searching = false
				m.searchInput.Blur()
				m.filterTemplates(m.searchInput.Value())
				return m, nil
			} else {
				// 选择模板
				i, ok := m.list.SelectedItem().(TemplateItem)
				if ok {
					m.choice = &i.Template
				}
				return m, tea.Quit
			}
		}

	case templatesLoadedMsg:
		// 加载模板列表
		items := make([]list.Item, len(msg.templates))
		for i, tmpl := range msg.templates {
			items[i] = TemplateItem{Template: tmpl}
		}
		m.allItems = items
		m.list.SetItems(items)
		return m, nil
	}

	// 处理搜索输入
	if m.searching {
		var cmd tea.Cmd
		m.searchInput, cmd = m.searchInput.Update(msg)
		return m, cmd
	}

	// 处理列表更新
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// filterTemplates 过滤模板
func (m *TemplateSelectorModel) filterTemplates(query string) {
	if query == "" {
		m.list.SetItems(m.allItems)
		return
	}

	query = strings.ToLower(query)
	var filtered []list.Item
	
	for _, item := range m.allItems {
		if templateItem, ok := item.(TemplateItem); ok {
			searchText := strings.ToLower(templateItem.FilterValue())
			if strings.Contains(searchText, query) {
				filtered = append(filtered, item)
			}
		}
	}
	
	m.list.SetItems(filtered)
}

// View 视图
func (m *TemplateSelectorModel) View() string {
	if m.quitting {
		return quitTextStyle.Render("已取消模板选择。")
	}

	var searchView string
	if m.searching {
		searchView = fmt.Sprintf("\n搜索: %s\n按 Enter 确认搜索，按 Esc 取消\n", m.searchInput.View())
	} else {
		searchView = "\n按 '/' 开始搜索，按 Enter 选择，按 'q' 退出\n"
	}

	return "\n" + m.list.View() + searchView
}

// GetSelectedTemplate 获取选中的模板
func (m *TemplateSelectorModel) GetSelectedTemplate() *client.Template {
	return m.choice
}

// SelectTemplate 启动模板选择器
func SelectTemplate(apiClient *client.Client) (*client.Template, error) {
	model := NewTemplateSelectorModel(apiClient)
	
	p := tea.NewProgram(model, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("运行模板选择器失败: %w", err)
	}

	if m, ok := finalModel.(*TemplateSelectorModel); ok {
		return m.GetSelectedTemplate(), nil
	}

	return nil, fmt.Errorf("无效的模型类型")
}