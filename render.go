package main

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	lg "charm.land/lipgloss/v2"
)

type model struct {
	lines  []Data
	cursor int
	view   int
	width  int
	height int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		{
			m.width = msg.Width
			m.height = msg.Height
		}
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "n", "down":
			if m.cursor < len(m.lines)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}
		case "N", "up":
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.lines) - 1
			}
		}
	}

	visibleArea := m.height - 2
	if visibleArea > 0 {
		if m.cursor < m.view {
			m.view = m.cursor
		}
		if m.cursor >= m.view+visibleArea {
			m.view = m.cursor - visibleArea + 1
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	var sb strings.Builder
	var view tea.View
	view.AltScreen = true
	view.MouseMode = tea.MouseModeCellMotion

	if m.height == 0 || m.width == 0 {
		sb.WriteString("Loading ...")
		view.SetContent(sb.String())
		return view
	}

	visibleArea := m.height - 2

	drawLine := 0
	end := m.view + visibleArea
	if end > len(m.lines) {
		end = len(m.lines)
	}

	// Styles
	hudStyle := lg.NewStyle().
		Background(lg.Color("#202020")).
		Foreground(lg.Color("#ebdbb2")).
		Width(m.width)
	errStyle := lg.NewStyle().
		Foreground(lg.Color("#E32636"))
	warnStyle := lg.NewStyle().
		Foreground(lg.Color("#FFBF00"))
	commonStyle := lg.NewStyle().
		Foreground(lg.Color("#DEB887"))

	err := 0
	warn := 0
	normal := 0
	
	for i := m.view; i < end; i++ {
		currLine := m.lines[i]

		prefix := " "
		if m.cursor == i {
			prefix = "> "
		}

		var style lg.Style
		switch currLine.Sev {
		case None, Note, Info, Hint:
			style = commonStyle
			normal++
		case Warning:
			style = warnStyle
			warn++
		case Error:
			style = errStyle
			err++
		}

		sb.WriteString(fmt.Sprintf(
			"%s%s(%d:%d): %s\n",
			prefix, currLine.File,
			currLine.Lin, currLine.Col,
			style.Render(fmt.Sprintf("%s", currLine.Msg)),
		))

		drawLine++
	}

	for i := drawLine; i < visibleArea; i++ {
		sb.WriteString("\n")
	}

	hudText := fmt.Sprintf(" ☰ %s  %s  %s  [n/N: navigate | q: quit]",
		commonStyle.Render(fmt.Sprintf("%d info", normal)),
		errStyle.Render(fmt.Sprintf("%d errors", err)),
		warnStyle.Render(fmt.Sprintf("%d warnings", warn)),
	)

	hud := hudStyle.Render(hudText)

	sb.WriteString(hud)
	view.SetContent(sb.String())

	return view
}
