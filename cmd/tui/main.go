package tui

import (
	"fmt"
	"strings"
	"time"

	"homevision/internal/houses/repository"
	"homevision/internal/houses/service"

	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	housesRepository = repository.NewHousesRepository()
	housesService    = service.NewHousesService(housesRepository)
)

type Model struct {
	choices          []string
	cursor           int
	selected         string
	spinner          spinner.Model
	loading          bool
	downloaded       bool
	downloadedImages []string
	paginator        paginator.Model
}

type tickMsg time.Time

func NewPaginator() paginator.Model {
	var items []string
	for i := 1; i < 101; i++ {
		text := fmt.Sprintf("Item %d", i)
		items = append(items, text)
	}

	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 10
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
	p.SetTotalPages(len(items))

	return p
}

func StartTUI() {
	m := Model{
		choices:          []string{"Download images", "List images"},
		spinner:          spinner.New(),
		paginator:        NewPaginator(),
		downloadedImages: []string{},
	}

	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "b":
			// Go back to the initial screen
			m.selected = ""
			m.loading = false
			m.downloaded = false

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter", " ":
			if m.selected == "" {
				m.selected = m.choices[m.cursor]
				if m.selected == "Download images" {
					m.loading = true
					m.downloaded = false
					return m, tea.Batch(spinner.Tick, performAction)
				} else if m.selected == "List images" {
					m.downloadedImages = housesService.GetDownloadedImages()
					return m, nil
				}
			}

		case "left":
			if m.selected == "List images" && len(m.downloadedImages) > 0 {
				m.paginator.PrevPage()
			}

		case "right":
			if m.selected == "List images" && len(m.downloadedImages) > 0 {
				m.paginator.NextPage()
			}
		}

	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case tickMsg:
		m.loading = false
		m.downloaded = true
		return m, nil

	case paginator.KeyMap:
		var cmd tea.Cmd
		m.paginator, cmd = m.paginator.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	if m.loading {
		return fmt.Sprintf("\n%s Downloading images...\n", m.spinner.View())
	}

	if m.downloaded {
		return "\n Download finished!\n \n b: back • q: quit\n"
	}

	if m.selected == "" {
		return mainView(m)
	}

	if m.selected == "List images" {
		return listView(m)
	}

	return ""
}

func mainView(m Model) string {
	s := "\n Choose an option:\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\n q: quit\n"

	return s
}

func listView(m Model) string {
	var b strings.Builder

	if len(m.downloadedImages) == 0 {
		b.WriteString("\n You have no images\n")
		b.WriteString("\n b: back • q: quit\n")
	} else {
		b.WriteString("\n Downloaded images:\n")
		m.paginator.PerPage = 10
		m.paginator.SetTotalPages(len(m.downloadedImages))
		start, end := m.paginator.GetSliceBounds(len(m.downloadedImages))
		for _, item := range m.downloadedImages[start:end] {
			b.WriteString("\n  • " + item + "\n")
		}
		b.WriteString("\n  " + m.paginator.View())
		b.WriteString("\n\n ←/→ page • b: back • q: quit\n")
	}
	return b.String()
}

func performAction() tea.Msg {
	// Simulate a delay for the API call
	housesService.DownloadImages()
	return tickMsg(time.Now())
}
