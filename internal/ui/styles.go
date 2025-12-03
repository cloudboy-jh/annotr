package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#59fb2b")).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginBottom(1)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42"))

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	WarningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214"))

	SelectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#59fb2b")).
			Bold(true)

	UnselectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	DimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2)

	InputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#59fb2b"))

	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1)
)

func Checkmark() string {
	return SuccessStyle.Render("✓")
}

func Cross() string {
	return ErrorStyle.Render("✗")
}

func Bullet() string {
	return DimStyle.Render("○")
}

func SelectedBullet() string {
	return SelectedStyle.Render("●")
}
