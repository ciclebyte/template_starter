package interactive

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

// 共用样式定义
var (
	titleStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Padding(0, 1)

	headerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("62")).
		Bold(true).
		Align(lipgloss.Center)

	selectedItemStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("62")).
		Bold(true)

	paginationStyle = list.DefaultStyles().PaginationStyle.
		PaddingLeft(4)

	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Align(lipgloss.Center)
)