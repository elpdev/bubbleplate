package generatorui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/elpdev/bubbleplate/internal/generator"
	"github.com/elpdev/bubbleplate/internal/theme"
)

type field int

const (
	fieldAppName field = iota
	fieldModulePath
	fieldOutputDir
	fieldDisplayName
	fieldDescription
	fieldDockerImage
	fieldForce
	fieldCreate
	fieldCount
)

type state int

const (
	stateEditing state = iota
	stateGenerating
	stateDone
)

type Model struct {
	width  int
	height int

	selected field
	values   map[field]string
	force    bool
	message  string
	state    state
	result   generator.Result
	config   generator.Config
	theme    theme.Theme
}

type generatedMsg struct {
	result generator.Result
	config generator.Config
	err    error
}

func New() Model {
	return Model{
		values: make(map[field]string),
		theme:  theme.Phosphor(),
	}
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg { return tea.RequestWindowSize() }
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case generatedMsg:
		if msg.err != nil {
			m.state = stateEditing
			m.message = msg.err.Error()
			return m, nil
		}
		m.state = stateDone
		m.result = msg.result
		m.config = msg.config
		m.message = ""
		return m, nil
	case tea.KeyPressMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		return m, tea.Quit
	}

	if m.state == stateGenerating {
		return m, nil
	}
	if m.state == stateDone {
		switch msg.String() {
		case "q", "enter":
			return m, tea.Quit
		case "b":
			m.state = stateEditing
			return m, nil
		}
		return m, nil
	}

	switch msg.String() {
	case "up", "shift+tab":
		m.previousField()
		return m, nil
	case "down", "tab":
		m.nextField()
		return m, nil
	case "enter":
		if m.selected == fieldCreate {
			return m.generate()
		}
		if m.selected == fieldForce {
			m.force = !m.force
			return m, nil
		}
		m.nextField()
		return m, nil
	case "ctrl+g":
		return m.generate()
	case "backspace", "ctrl+h":
		m.deleteLastRune()
		return m, nil
	case "space":
		m.appendText(" ")
		return m, nil
	}

	text := msg.String()
	if len([]rune(text)) == 1 {
		m.appendText(text)
		return m, nil
	}
	return m, nil
}

func (m *Model) nextField() {
	m.message = ""
	m.selected = (m.selected + 1) % fieldCount
}

func (m *Model) previousField() {
	m.message = ""
	if m.selected == 0 {
		m.selected = fieldCount - 1
		return
	}
	m.selected--
}

func (m *Model) appendText(text string) {
	m.message = ""
	if m.selected == fieldForce || m.selected == fieldCreate {
		return
	}
	m.values[m.selected] += text
}

func (m *Model) deleteLastRune() {
	m.message = ""
	if m.selected == fieldForce || m.selected == fieldCreate {
		return
	}
	value := []rune(m.values[m.selected])
	if len(value) == 0 {
		return
	}
	m.values[m.selected] = string(value[:len(value)-1])
}

func (m Model) generate() (tea.Model, tea.Cmd) {
	config, err := m.configFromValues()
	if err != nil {
		m.message = err.Error()
		return m, nil
	}
	m.state = stateGenerating
	m.message = "Creating project..."
	return m, func() tea.Msg {
		result, err := generator.Generate(config)
		return generatedMsg{result: result, config: config, err: err}
	}
}

func (m Model) configFromValues() (generator.Config, error) {
	appName := strings.TrimSpace(m.values[fieldAppName])
	if appName == "" {
		return generator.Config{}, fmt.Errorf("app name is required")
	}
	modulePath := strings.TrimSpace(m.values[fieldModulePath])
	if modulePath == "" {
		return generator.Config{}, fmt.Errorf("module path is required")
	}

	config := generator.NewConfig(appName)
	config.ModulePath = modulePath
	config.Force = m.force
	if value := strings.TrimSpace(m.values[fieldOutputDir]); value != "" {
		config.OutputDir = value
	}
	if value := strings.TrimSpace(m.values[fieldDisplayName]); value != "" {
		config.DisplayName = value
	}
	if value := strings.TrimSpace(m.values[fieldDescription]); value != "" {
		config.Description = value
	}
	if value := strings.TrimSpace(m.values[fieldDockerImage]); value != "" {
		config.DockerImage = value
	}
	return config, nil
}

func (m Model) View() tea.View {
	view := tea.NewView(m.render())
	view.BackgroundColor = m.theme.Background
	return view
}

func (m Model) render() string {
	if m.state == stateDone && m.result.Files > 0 {
		return m.renderDone()
	}
	return m.renderForm()
}

func (m Model) renderForm() string {
	var b strings.Builder
	b.WriteString(m.theme.Title.Render("Bubbleplate"))
	b.WriteString("\n")
	b.WriteString(m.theme.Muted.Render("Create a new Bubble Tea TUI project. Press ctrl+g to generate, esc to cancel."))
	b.WriteString("\n\n")

	for field := fieldAppName; field < fieldCount; field++ {
		b.WriteString(m.renderField(field))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(m.theme.Muted.Render("tab/down next  shift+tab/up previous  enter next/toggle/create  ctrl+g generate"))

	if m.message != "" {
		b.WriteString("\n\n")
		style := m.theme.Info
		if m.state == stateDone {
			style = m.theme.Warn
		}
		b.WriteString(style.Render(m.message))
	}

	return lipgloss.NewStyle().Padding(1, 2).Render(b.String())
}

func (m Model) renderField(field field) string {
	label := fieldLabel(field)
	value := m.displayValue(field)
	if field == m.selected && m.state == stateEditing {
		if field != fieldForce && field != fieldCreate {
			value += "█"
		}
		return m.theme.Selected.Render(fmt.Sprintf("%-14s %s", label, value))
	}
	return fmt.Sprintf("%-14s %s", label, m.theme.Text.Render(value))
}

func (m Model) displayValue(field field) string {
	if field == fieldForce {
		if m.force {
			return "yes"
		}
		return "no"
	}
	if field == fieldCreate {
		return "Create project"
	}

	value := m.values[field]
	if value != "" {
		return value
	}
	placeholder := m.placeholder(field)
	if placeholder == "" {
		return ""
	}
	return m.theme.Muted.Render(placeholder)
}

func (m Model) placeholder(field field) string {
	appName := strings.TrimSpace(m.values[fieldAppName])
	if appName == "" {
		switch field {
		case fieldAppName:
			return "hackernews"
		case fieldModulePath:
			return "github.com/acme/hackernews"
		}
		return ""
	}
	config := generator.NewConfig(appName)
	switch field {
	case fieldOutputDir:
		return config.OutputDir
	case fieldDisplayName:
		return config.DisplayName
	case fieldDescription:
		return config.Description
	case fieldDockerImage:
		return config.DockerImage
	}
	return ""
}

func (m Model) renderDone() string {
	lines := []string{
		m.theme.Title.Render("Project created"),
		"",
		fmt.Sprintf("Created %s (%d files)", m.result.OutputDir, m.result.Files),
		"",
		"Next steps:",
		fmt.Sprintf("  cd %s", m.result.OutputDir),
		"  go mod tidy",
		"  go test ./...",
		fmt.Sprintf("  go run ./cmd/%s", m.config.BinaryName),
		"",
		m.theme.Muted.Render("Press enter or q to quit."),
	}
	return lipgloss.NewStyle().Padding(1, 2).Render(strings.Join(lines, "\n"))
}

func fieldLabel(field field) string {
	switch field {
	case fieldAppName:
		return "App name"
	case fieldModulePath:
		return "Module path"
	case fieldOutputDir:
		return "Output dir"
	case fieldDisplayName:
		return "Display name"
	case fieldDescription:
		return "Description"
	case fieldDockerImage:
		return "Docker image"
	case fieldForce:
		return "Force"
	case fieldCreate:
		return "Create"
	default:
		return ""
	}
}
