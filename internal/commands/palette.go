package commands

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/elpdev/bubbleplate/internal/theme"
)

type PaletteModel struct {
	registry *Registry
	query    string
	selected int
	executed *Command
}

func NewPaletteModel(registry *Registry) PaletteModel {
	return PaletteModel{registry: registry}
}

func (m PaletteModel) Update(msg tea.Msg) (PaletteModel, tea.Cmd) {
	m.executed = nil
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up", "ctrl+p":
			if m.selected > 0 {
				m.selected--
			}
			return m, nil
		case "down", "ctrl+n":
			if m.selected < len(m.matches())-1 {
				m.selected++
			}
			return m, nil
		case "enter":
			matches := m.matches()
			if len(matches) == 0 {
				return m, nil
			}
			command := matches[m.selected]
			m.executed = &command
			return m, nil
		case "backspace", "ctrl+h":
			if len(m.query) > 0 {
				m.query = m.query[:len(m.query)-1]
				m.selected = 0
			}
			return m, nil
		case "space":
			m.query += " "
			m.selected = 0
			return m, nil
		}
		if len(msg.String()) == 1 {
			m.query += msg.String()
			m.selected = 0
			return m, nil
		}
	}
	if m.selected >= len(m.matches()) {
		m.selected = 0
	}
	return m, nil
}

func (m PaletteModel) View(t theme.Theme) string {
	matches := m.matches()
	var b strings.Builder
	b.WriteString(t.Title.Render("Command Palette"))
	b.WriteString("\n")
	query := m.query
	if query == "" {
		query = t.Muted.Render("type a command...")
	}
	b.WriteString("> " + query)
	b.WriteString("\n\n")

	if len(matches) == 0 {
		b.WriteString(t.Muted.Render("No commands found"))
	} else {
		for i, command := range matches {
			line := fmt.Sprintf("%-18s %s", command.Title, command.Description)
			if i == m.selected {
				line = t.Selected.Render(line)
			} else {
				line = t.Text.Render(line)
			}
			b.WriteString(line + "\n")
		}
	}

	return t.Modal.Width(62).Render(b.String())
}

func (m *PaletteModel) Reset() {
	m.query = ""
	m.selected = 0
	m.executed = nil
}

func (m PaletteModel) ExecutedCommand() *Command { return m.executed }

func (m PaletteModel) matches() []Command { return m.registry.Filter(m.query) }
