package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cloudboy-jh/annotr/internal/config"
)

type modelStep int

const (
	modelStepDetecting modelStep = iota
	modelStepSelectProvider
	modelStepSelectModel
	modelStepDone
)

type ModelSelectModel struct {
	step             modelStep
	ollamaFound      bool
	ollamaModels     []config.OllamaModel
	providers        []string
	selectedIdx      int
	selectedProvider string
	selectedModel    string
	config           *config.Config
	err              error
	quitting         bool
}

func NewModelSelectModel(cfg *config.Config) ModelSelectModel {
	return ModelSelectModel{
		step:             modelStepDetecting,
		providers:        []string{"Ollama (local)", "Claude (Anthropic)", "OpenAI", "Groq"},
		config:           cfg,
		selectedProvider: cfg.DefaultProvider,
		selectedModel:    cfg.DefaultModel,
	}
}

type modelDetectMsg struct {
	found  bool
	models []config.OllamaModel
}

func detectOllamaForModel() tea.Msg {
	found, models, _ := config.DetectOllama()
	return modelDetectMsg{found: found, models: models}
}

func (m ModelSelectModel) Init() tea.Cmd {
	return detectOllamaForModel
}

func (m ModelSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			return m.handleEnter()
		case "up", "k":
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}
		case "down", "j":
			m.selectedIdx++
			m.selectedIdx = m.clampSelection()
		}

	case modelDetectMsg:
		m.ollamaFound = msg.found
		m.ollamaModels = msg.models
		m.step = modelStepSelectProvider
		return m, nil
	}

	return m, nil
}

func (m ModelSelectModel) handleEnter() (tea.Model, tea.Cmd) {
	switch m.step {
	case modelStepSelectProvider:
		providerMap := map[int]string{0: "ollama", 1: "anthropic", 2: "openai", 3: "groq"}
		m.selectedProvider = providerMap[m.selectedIdx]

		if m.selectedProvider == "ollama" && !m.ollamaFound {
			return m, nil
		}

		if m.selectedProvider != "ollama" {
			if m.config.APIKeys[m.selectedProvider] == "" {
				m.err = fmt.Errorf("no API key configured for %s. Run 'annotr init' to configure", m.selectedProvider)
				m.step = modelStepDone
				return m, tea.Quit
			}
		}

		m.step = modelStepSelectModel
		m.selectedIdx = 0
		return m, nil

	case modelStepSelectModel:
		if m.selectedProvider == "ollama" {
			if m.selectedIdx < len(m.ollamaModels) {
				m.selectedModel = m.ollamaModels[m.selectedIdx].Name
			}
		} else {
			models := getModelsForProvider(m.selectedProvider)
			if m.selectedIdx < len(models) {
				m.selectedModel = models[m.selectedIdx].ID
			}
		}

		m.config.DefaultProvider = m.selectedProvider
		m.config.DefaultModel = m.selectedModel
		if err := m.config.Save(); err != nil {
			m.err = err
		}
		m.step = modelStepDone
		return m, tea.Quit
	}

	return m, nil
}

func (m ModelSelectModel) clampSelection() int {
	var max int
	switch m.step {
	case modelStepSelectProvider:
		max = len(m.providers) - 1
	case modelStepSelectModel:
		if m.selectedProvider == "ollama" {
			max = len(m.ollamaModels) - 1
		} else {
			max = len(getModelsForProvider(m.selectedProvider)) - 1
		}
	default:
		max = 0
	}
	if max < 0 {
		max = 0
	}
	if m.selectedIdx > max {
		return max
	}
	return m.selectedIdx
}

func (m ModelSelectModel) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	b.WriteString(TitleStyle.Render("annotr model") + "\n\n")

	switch m.step {
	case modelStepDetecting:
		b.WriteString("Checking for Ollama...\n")

	case modelStepSelectProvider:
		b.WriteString(SubtitleStyle.Render("Select Provider") + "\n\n")
		b.WriteString(DimStyle.Render(fmt.Sprintf("Current: %s / %s", m.config.DefaultProvider, m.config.DefaultModel)) + "\n\n")

		for i, p := range m.providers {
			disabled := false
			if i == 0 && !m.ollamaFound {
				disabled = true
			}
			if i > 0 {
				providerKey := map[int]string{1: "anthropic", 2: "openai", 3: "groq"}[i]
				if m.config.APIKeys[providerKey] == "" {
					disabled = true
				}
			}

			if i == m.selectedIdx {
				if disabled {
					b.WriteString(SelectedBullet() + " " + DimStyle.Render(p+" (not configured)") + "\n")
				} else {
					b.WriteString(SelectedBullet() + " " + SelectedStyle.Render(p) + "\n")
				}
			} else {
				if disabled {
					b.WriteString(Bullet() + " " + DimStyle.Render(p+" (not configured)") + "\n")
				} else {
					b.WriteString(Bullet() + " " + UnselectedStyle.Render(p) + "\n")
				}
			}
		}

	case modelStepSelectModel:
		b.WriteString(SubtitleStyle.Render("Select Model") + "\n\n")
		b.WriteString(DimStyle.Render(fmt.Sprintf("Provider: %s", m.selectedProvider)) + "\n\n")

		var models []struct{ id, name string }
		if m.selectedProvider == "ollama" {
			for _, om := range m.ollamaModels {
				models = append(models, struct{ id, name string }{om.Name, om.Name})
			}
		} else {
			for _, pm := range getModelsForProvider(m.selectedProvider) {
				models = append(models, struct{ id, name string }{pm.ID, pm.Name})
			}
		}

		for i, model := range models {
			if i == m.selectedIdx {
				b.WriteString(SelectedBullet() + " " + SelectedStyle.Render(model.name) + "\n")
			} else {
				b.WriteString(Bullet() + " " + UnselectedStyle.Render(model.name) + "\n")
			}
		}

	case modelStepDone:
		if m.err != nil {
			b.WriteString(Cross() + " " + ErrorStyle.Render("Error: "+m.err.Error()) + "\n")
		} else {
			b.WriteString(Checkmark() + " Model updated!\n\n")
			b.WriteString("Provider: " + m.selectedProvider + "\n")
			b.WriteString("Model: " + m.selectedModel + "\n")
		}
	}

	b.WriteString("\n" + HelpStyle.Render("↑/↓ navigate • enter select • q quit"))

	return BoxStyle.Render(b.String())
}

func (m ModelSelectModel) Config() *config.Config {
	return m.config
}

func (m ModelSelectModel) Error() error {
	return m.err
}
