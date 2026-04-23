package modal

import (
	"charm.land/lipgloss/v2"
	"github.com/elpdev/bubbleplate/internal/theme"
)

func Overlay(base, content string, width, height int, _ theme.Theme) string {
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content, lipgloss.WithWhitespaceChars(" "))
}
