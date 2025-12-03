package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/johnhorton/annotr/internal/config"
)

type step int

const (
	stepDetecting step = iota
	stepSelectProvider
	stepSelectOllamaModel
	stepEnterAPIKey
	stepSelectCloudModel
	stepSelectStyle
	stepConfirm
	stepDone
)

type InitModel struct {
	step           step
	ollamaFound    bool
	ollamaModels   []config.OllamaModel
	providers      []string
	selectedIdx    int
	apiKeyInput    textinput.Model
	selectedProvider string
	selectedModel    string
	selectedStyle    string
	config         *config.Config
	err            error
	quitting       bool
}

func NewInitModel() InitModel {
	ti := textinput.New()
	ti.Placeholder = "sk-ant-..."
	ti.Focus()
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '*'

	return InitModel{
		step:        stepDetecting,
		providers:   []string{"Claude (Anthropic)", "OpenAI", "Groq"},
		apiKeyInput: ti,
		config:      config.DefaultConfig(),
	}
}

type detectMsg struct {
	found  bool
	models []config.OllamaModel
}

func detectOllama() tea.Msg {
	found, models, _ := config.DetectOllama()
	return detectMsg{found: found, models: models}
}

func (m InitModel) Init() tea.Cmd {
	return detectOllama
}

func (m InitModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	case detectMsg:
		m.ollamaFound = msg.found
		m.ollamaModels = msg.models
		if m.ollamaFound && len(m.ollamaModels) > 0 {
			m.step = stepSelectOllamaModel
		} else {
			m.step = stepSelectProvider
		}
		return m, nil
	}

	if m.step == stepEnterAPIKey {
		var cmd tea.Cmd
		m.apiKeyInput, cmd = m.apiKeyInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m InitModel) handleEnter() (tea.Model, tea.Cmd) {
	switch m.step {
	case stepSelectProvider:
		providerMap := map[int]string{0: "anthropic", 1: "openai", 2: "groq"}
		m.selectedProvider = providerMap[m.selectedIdx]
		m.step = stepEnterAPIKey
		m.apiKeyInput.Placeholder = getPlaceholder(m.selectedProvider)
		m.selectedIdx = 0
		return m, textinput.Blink

	case stepSelectOllamaModel:
		if m.selectedIdx < len(m.ollamaModels) {
			m.selectedModel = m.ollamaModels[m.selectedIdx].Name
			m.selectedProvider = "ollama"
			m.step = stepSelectStyle
			m.selectedIdx = 0
		} else {
			m.step = stepSelectProvider
			m.selectedIdx = 0
		}
		return m, nil

	case stepEnterAPIKey:
		key := m.apiKeyInput.Value()
		if config.ValidateAPIKey(m.selectedProvider, key) {
			m.config.APIKeys[m.selectedProvider] = key
			m.step = stepSelectCloudModel
			m.selectedIdx = 0
		}
		return m, nil

	case stepSelectCloudModel:
		models := getModelsForProvider(m.selectedProvider)
		if m.selectedIdx < len(models) {
			m.selectedModel = models[m.selectedIdx].ID
			m.step = stepSelectStyle
			m.selectedIdx = 0
		}
		return m, nil

	case stepSelectStyle:
		styles := []string{"line", "block", "docstring"}
		m.selectedStyle = styles[m.selectedIdx]
		m.step = stepConfirm
		return m, nil

	case stepConfirm:
		m.config.DefaultProvider = m.selectedProvider
		m.config.DefaultModel = m.selectedModel
		m.config.CommentStyle = m.selectedStyle
		if err := m.config.Save(); err != nil {
			m.err = err
		}
		m.step = stepDone
		return m, tea.Quit
	}

	return m, nil
}

func (m InitModel) clampSelection() int {
	var max int
	switch m.step {
	case stepSelectProvider:
		max = len(m.providers) - 1
	case stepSelectOllamaModel:
		max = len(m.ollamaModels)
	case stepSelectCloudModel:
		max = len(getModelsForProvider(m.selectedProvider)) - 1
	case stepSelectStyle:
		max = 2
	default:
		max = 0
	}
	if m.selectedIdx > max {
		return max
	}
	return m.selectedIdx
}

func (m InitModel) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	b.WriteString(TitleStyle.Render("annotr init") + "\n\n")

	switch m.step {
	case stepDetecting:
		b.WriteString("Checking for Ollama...\n")

	case stepSelectProvider:
		b.WriteString(SubtitleStyle.Render("Select LLM Provider") + "\n\n")
		for i, p := range m.providers {
			if i == m.selectedIdx {
				b.WriteString(SelectedBullet() + " " + SelectedStyle.Render(p) + "\n")
			} else {
				b.WriteString(Bullet() + " " + UnselectedStyle.Render(p) + "\n")
			}
		}
		b.WriteString("\n" + DimStyle.Render("[Install Ollama for local/free option]") + "\n")

	case stepSelectOllamaModel:
		b.WriteString(Checkmark() + " Ollama found at localhost:11434\n")
		b.WriteString(Checkmark() + fmt.Sprintf(" %d models available\n\n", len(m.ollamaModels)))
		b.WriteString(SubtitleStyle.Render("Select Model") + "\n\n")
		for i, model := range m.ollamaModels {
			if i == m.selectedIdx {
				b.WriteString(SelectedBullet() + " " + SelectedStyle.Render(model.Name) + "\n")
			} else {
				b.WriteString(Bullet() + " " + UnselectedStyle.Render(model.Name) + "\n")
			}
		}
		if m.selectedIdx == len(m.ollamaModels) {
			b.WriteString(SelectedBullet() + " " + SelectedStyle.Render("[Use API key instead]") + "\n")
		} else {
			b.WriteString(Bullet() + " " + DimStyle.Render("[Use API key instead]") + "\n")
		}

	case stepEnterAPIKey:
		b.WriteString(SubtitleStyle.Render(fmt.Sprintf("Configure %s", strings.Title(m.selectedProvider))) + "\n\n")
		b.WriteString("API Key: " + m.apiKeyInput.View() + "\n")

	case stepSelectCloudModel:
		b.WriteString(Checkmark() + " API key validated\n\n")
		b.WriteString(SubtitleStyle.Render("Select Model") + "\n\n")
		models := getModelsForProvider(m.selectedProvider)
		for i, model := range models {
			if i == m.selectedIdx {
				b.WriteString(SelectedBullet() + " " + SelectedStyle.Render(model.Name) + "\n")
			} else {
				b.WriteString(Bullet() + " " + UnselectedStyle.Render(model.Name) + "\n")
			}
		}

	case stepSelectStyle:
		b.WriteString(SubtitleStyle.Render("Comment Style") + "\n\n")
		styles := []struct{ name, desc string }{
			{"Line comments", "// comment or # comment"},
			{"Block comments", "/* comment */ or \"\"\" comment \"\"\""},
			{"Language-specific", "Auto-detect best format per language"},
		}
		for i, s := range styles {
			if i == m.selectedIdx {
				b.WriteString(SelectedBullet() + " " + SelectedStyle.Render(s.name) + "\n")
				b.WriteString("    " + DimStyle.Render(s.desc) + "\n")
			} else {
				b.WriteString(Bullet() + " " + UnselectedStyle.Render(s.name) + "\n")
				b.WriteString("    " + DimStyle.Render(s.desc) + "\n")
			}
		}

	case stepConfirm:
		b.WriteString(SubtitleStyle.Render("Configuration Complete") + "\n\n")
		b.WriteString(Checkmark() + " Provider: " + m.selectedProvider + "\n")
		b.WriteString(Checkmark() + " Model: " + m.selectedModel + "\n")
		b.WriteString(Checkmark() + " Style: " + m.selectedStyle + "\n")
		b.WriteString("\n" + DimStyle.Render("Press Enter to save") + "\n")

	case stepDone:
		if m.err != nil {
			b.WriteString(Cross() + " " + ErrorStyle.Render("Error: "+m.err.Error()) + "\n")
		} else {
			b.WriteString(Checkmark() + " Saved to: ~/.annotr/config.json\n\n")
			b.WriteString("Ready to use! Try:\n")
			b.WriteString(InputStyle.Render("  $ annotr file.go") + "\n")
		}
	}

	b.WriteString("\n" + HelpStyle.Render("↑/↓ navigate • enter select • q quit"))

	return BoxStyle.Render(b.String())
}

func (m InitModel) Config() *config.Config {
	return m.config
}

func (m InitModel) Error() error {
	return m.err
}

func getPlaceholder(provider string) string {
	switch provider {
	case "anthropic":
		return "sk-ant-..."
	case "openai":
		return "sk-..."
	case "groq":
		return "gsk_..."
	default:
		return ""
	}
}

func getModelsForProvider(provider string) []config.Model {
	manifest, _ := config.LoadModelsManifest()
	if manifest == nil {
		manifest = config.DefaultModelsManifest()
	}
	if p, ok := manifest.Providers[provider]; ok {
		return p.Models
	}
	return nil
}
