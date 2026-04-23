package theme

import "charm.land/lipgloss/v2"

func Default() Theme {
	primary := lipgloss.Color("#7D56F4")
	muted := lipgloss.Color("#6B7280")
	border := lipgloss.Color("#D1D5DB")
	return Theme{
		Name:     "Default",
		Text:     lipgloss.NewStyle().Foreground(lipgloss.Color("#111827")),
		Muted:    lipgloss.NewStyle().Foreground(muted),
		Title:    lipgloss.NewStyle().Bold(true).Foreground(primary),
		Selected: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF")).Background(primary).Padding(0, 1),
		Header:   lipgloss.NewStyle().Foreground(lipgloss.Color("#111827")).Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(border).Padding(0, 1),
		Sidebar:  lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, true, false, false).BorderForeground(border).Padding(1, 1),
		Main:     lipgloss.NewStyle().Padding(1, 2),
		Footer:   lipgloss.NewStyle().Foreground(muted).Border(lipgloss.NormalBorder(), true, false, false, false).BorderForeground(border).Padding(0, 1),
		Modal:    lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(primary).Padding(1, 2).Background(lipgloss.Color("#FFFFFF")),
		Border:   lipgloss.NewStyle().Foreground(border),
		Info:     lipgloss.NewStyle().Foreground(lipgloss.Color("#2563EB")),
		Warn:     lipgloss.NewStyle().Foreground(lipgloss.Color("#B45309")),
	}
}

func MutedDark() Theme {
	primary := lipgloss.Color("#A78BFA")
	muted := lipgloss.Color("#9CA3AF")
	border := lipgloss.Color("#374151")
	bg := lipgloss.Color("#111827")
	return Theme{
		Name:     "Muted Dark",
		Text:     lipgloss.NewStyle().Foreground(lipgloss.Color("#E5E7EB")).Background(bg),
		Muted:    lipgloss.NewStyle().Foreground(muted).Background(bg),
		Title:    lipgloss.NewStyle().Bold(true).Foreground(primary).Background(bg),
		Selected: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#111827")).Background(primary).Padding(0, 1),
		Header:   lipgloss.NewStyle().Foreground(lipgloss.Color("#E5E7EB")).Background(bg).Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(border).Padding(0, 1),
		Sidebar:  lipgloss.NewStyle().Background(bg).Border(lipgloss.NormalBorder(), false, true, false, false).BorderForeground(border).Padding(1, 1),
		Main:     lipgloss.NewStyle().Background(bg).Padding(1, 2),
		Footer:   lipgloss.NewStyle().Foreground(muted).Background(bg).Border(lipgloss.NormalBorder(), true, false, false, false).BorderForeground(border).Padding(0, 1),
		Modal:    lipgloss.NewStyle().Foreground(lipgloss.Color("#E5E7EB")).Background(lipgloss.Color("#1F2937")).Border(lipgloss.RoundedBorder()).BorderForeground(primary).Padding(1, 2),
		Border:   lipgloss.NewStyle().Foreground(border).Background(bg),
		Info:     lipgloss.NewStyle().Foreground(lipgloss.Color("#60A5FA")).Background(bg),
		Warn:     lipgloss.NewStyle().Foreground(lipgloss.Color("#FBBF24")).Background(bg),
	}
}
