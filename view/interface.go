package view

import (
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/chroma/quick"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joaovrmoraes/http-go-cli/model"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	subtitleStyle = func() lipgloss.Style {
		return lipgloss.NewStyle().Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

type modelTea struct {
	content  string
	header   string
	subtitle string
	ready    bool
	viewport viewport.Model
}

func (m modelTea) Init() tea.Cmd {
	return nil
}

func (m modelTea) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m modelTea) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s\n%s", m.headerView(), m.subtitleView(), m.viewport.View(), m.footerView())
}

func (m modelTea) headerView() string {
	title := titleStyle.Render(m.header)
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m modelTea) subtitleView() string {
	return subtitleStyle.Render(m.subtitle)
}

func (m modelTea) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func formatHeaders(headers map[string][]string) string {
	var builder strings.Builder
	for key, values := range headers {
		builder.WriteString(fmt.Sprintf("%s: %v\n", key, values))
	}
	return builder.String()
}

func StartInterface(jsonResponse, header string, headers map[string][]string) {
	var coloredJSON strings.Builder
	err := quick.Highlight(&coloredJSON, jsonResponse, "json", "terminal", "monokai")
	if err != nil {
		fmt.Println("could not highlight JSON:", err)
		os.Exit(1)
	}

	var coloredHistory strings.Builder
	err = quick.Highlight(&coloredHistory, "Request History", "plaintext", "terminal", "monokai")
	if err != nil {
		fmt.Println("could not highlight history:", err)
		os.Exit(1)
	}

	subtitle := formatHeaders(headers)

	vp := viewport.New(80, 20)
	vp.SetContent(coloredJSON.String() + "\n\n" + coloredHistory.String() + "\n" + model.DisplayHistory())

	// content := fmt.Sprintf("%s\n%s", coloredHistory.String(), coloredJSON.String())

	p := tea.NewProgram(
		modelTea{header: header, subtitle: subtitle, viewport: vp, ready: true},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}
