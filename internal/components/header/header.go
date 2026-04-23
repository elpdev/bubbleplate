package header

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/elpdev/bubbleplate/internal/theme"
)

type Model struct {
	AppName     string
	ScreenTitle string
	Version     string
}

func View(m Model, width int, t theme.Theme) string {
	left := t.Title.Render(m.AppName)
	right := fmt.Sprintf("%s  %s", m.ScreenTitle, m.Version)
	content := lipgloss.JoinHorizontal(lipgloss.Center, left, lipgloss.PlaceHorizontal(max(0, width-lipgloss.Width(left)-lipgloss.Width(right)-4), lipgloss.Left, ""), right)
	return t.Header.Width(width).Height(3).Render(content)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
