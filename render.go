package main

import (
	"fmt"
	"log"
	"strings"

	tea "charm.land/bubbletea/v2"
	lg "charm.land/lipgloss/v2"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
)

// AppState NOTE(xendak): currently only using Interative and Results
type AppState int

const (
	Interactive = iota
	Loading
	Passthrough
	Results
)

type model struct {
	state  AppState
	input  textinput.Model
	msg    Message
	cursor int
	view   int
	width  int
	height int
	cmd    tea.Cmd
	cfg    Config
	keys   KeyBindings
}

func newModel(cfg Config) model {
	input := textinput.New()
	input.Focus()
	input.SetWidth(48)
	input.Prompt = "> "
	input.ShowSuggestions = true

	style := textinput.DefaultDarkStyles()

	colorAccent := lg.Color(cfg.Accent)

	style.Focused.Text = lg.NewStyle().Foreground(lg.Color(cfg.InputFg))
	style.Focused.Prompt = lg.NewStyle().Foreground(colorAccent)
	style.Cursor.Color = colorAccent

	input.SetStyles(style)

	keys := BuildKeyMap(cfg)

	return model{
		input: input,
		cfg:   cfg,
		keys:  keys,
	}
}

func (m model) Init() tea.Cmd {
	return m.cmd
}

// TODO(xendak): we don't need to check match if we actively use Passthrough/Results/Interative states.. eventually
// NOTE(xendak): i need to remember to avoid *model, violates bubbletea principles
func findNext(msg Message, cur int) int {
	if msg.match < 0 {
		return 0
	}
	cur = (cur + 1) % len(msg.Lines)
	for !(msg.Lines[cur].Match) {
		cur = (cur + 1) % len(msg.Lines)
	}
	return cur
}

func findPrev(msg Message, cur int) int {
	if msg.match < 0 {
		return 0
	}
	// NOTE(xendak): c = 0? c - 1 => -1, segfault, then we add maxCount to fix.
	cur = (cur - 1 + len(msg.Lines)) % len(msg.Lines)
	for !(msg.Lines[cur].Match) {
		cur = (cur - 1 + len(msg.Lines)) % len(msg.Lines)
	}
	return cur
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case cmdOut:
		// TODO(xendak): catalog the errors
		if msg.err != nil {
			return m, nil
		}

		m.msg = parseMsg(string(msg.out))
		m.cursor = max(0, m.msg.match)
		m.view = 0
		m.state = Results

	case tea.KeyPressMsg:
    log.Printf("KeyPress: %q", msg.String())
		switch m.state {
		case Interactive:
			switch {
			case key.Matches(msg, m.keys.QuitPrompt):
				return m, tea.Quit
			case key.Matches(msg, m.keys.EscPrompt):
				if len(m.msg.Lines) > 0 {
					m.state = Results
					m.input.Blur()
				} else {
					if len(m.input.Value()) < 1 {
						return m, tea.Quit
					}
				}
			case key.Matches(msg, m.keys.SubmitCmd):
				val := m.input.Value()
				if val == "" {
					val = "make"
				}

				m.state = Results
				m.input.Blur()

				save(val)

				return m, runCommand(val)
			}
		default:
			switch {
			case key.Matches(msg, m.keys.Quit):
				return m, tea.Quit
			// TODO(xendak): add Horizontal movement if we didn't/can't wrap
			case key.Matches(msg, m.keys.NextMatch):
				m.cursor = findNext(m.msg, m.cursor)
			case key.Matches(msg, m.keys.PrevMatch):
				m.cursor = findPrev(m.msg, m.cursor)
			case key.Matches(msg, m.keys.OpenPrompt):
				m.input.SetValue("")
				m.state = Interactive
				m.input.Focus()
				m.input.SetSuggestions(getSuggestions())

				return m, nil
			case key.Matches(msg, m.keys.OpenPrompt):
				curr := m.msg.Lines[m.cursor]
				var arg strings.Builder

				// TODO(xendak) remove the hardcode and offer configs
				editor := "wez-hx"
				fmt.Fprintf(&arg, "%s:%d:%d", curr.File, curr.Lin, curr.Col)

				if m.cfg.Block {
					return m, openEditor(editor, arg.String())
				}
				return m, openEditorAsync(editor, arg.String())
			}

		}
	}

	if m.state == Interactive {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd

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
	var view tea.View
	view.AltScreen = true
	view.MouseMode = tea.MouseModeCellMotion

	if m.height == 0 || m.width == 0 {
		view.SetContent("Loading ...")
		return view
	}

	results := m.renderResults()

	if m.state == Interactive {
		box := m.renderInput()
		boxW := 58
		boxH := 6
		x := (m.width - boxW) / 2
		y := (m.height - boxH) / 2

		comp := lg.NewCompositor(
			lg.NewLayer(results),
			lg.NewLayer(box).X(x).Y(y),
		)
		view.SetContent(comp.Render())

		if c := m.input.Cursor(); c != nil {
			c.X += x + 3
			c.Y += y + 3
			view.Cursor = c
		}
		return view
	}

	view.SetContent(results)
	return view
}

func (m model) renderInput() string {
	labelStyle := lg.NewStyle().
		Foreground(lg.Color(m.cfg.NormalFg)).
		Width(48).
		Bold(true)

	inputStyle := lg.NewStyle().
		Foreground(lg.Color(m.cfg.InputFg)).
		Padding(0, 1).
		Border(lg.RoundedBorder()).
		BorderForeground(lg.Color(m.cfg.InputBorderColor)).
		Width(48)

	outerStyle := lg.NewStyle().
		Border(lg.RoundedBorder()).
		BorderForeground(lg.Color(m.cfg.InputBorderColor)).
		Padding(0, 2).
		Width(56)

	inner := lg.JoinVertical(lg.Left,
		labelStyle.Render("Command (default: make)"),
		inputStyle.Render(m.input.View()),
	)

	return outerStyle.Render(inner)
}

func (m model) renderResults() string {
	var sb strings.Builder

	visibleArea := m.height - 2

	drawLine := 0
	end := min(m.view+visibleArea, len(m.msg.Lines))

	hudStyle := lg.NewStyle().
		Background(lg.Color(m.cfg.HudBg)).
		Foreground(lg.Color(m.cfg.HudFg))

	errStyle := lg.NewStyle().Foreground(lg.Color(m.cfg.ErrorFg))
	warnStyle := lg.NewStyle().Foreground(lg.Color(m.cfg.WarnFg))
	noteStyle := lg.NewStyle().Foreground(lg.Color(m.cfg.NoteFg))

	fileStyle := lg.NewStyle().Foreground(lg.Color(m.cfg.FileFg))
	locationStyle := lg.NewStyle().Foreground(lg.Color(m.cfg.LocationFg))

	err := 0
	warn := 0
	normal := 0

	for i := m.view; i < end; i++ {
		currLine := m.msg.Lines[i]

		prefix := " "
		if m.cursor == i {
			prefix = "> "
		}

		// var style lg.Style
		switch currLine.Sev {
		case None, Note, Info, Hint:
			// style = commonStyle
			normal++
		case Warning:
			// style = warnStyle
			warn++
		case Error:
			// style = errStyle
			err++
		}

		if currLine.Match {
			fmt.Fprintf(&sb, "%s%s(%s): %s\n",
				prefix,
				fileStyle.Render(currLine.File),
				locationStyle.Render(fmt.Sprintf("%d:%d", currLine.Lin, currLine.Col)),
				// style.Render(currLine.Sev.String()),
				currLine.Msg)
		} else {
			fmt.Fprintf(&sb, "%s\n", currLine.Raw)
		}

		drawLine++
	}

	for i := drawLine; i < visibleArea; i++ {
		sb.WriteString("\n")
	}

	n := m.keys.NextMatch.Help().Key
	p := m.keys.PrevMatch.Help().Key
	q := m.keys.Quit.Help().Key

	hudText := hudStyle.Render(" ☰ ") +
		noteStyle.Inherit(hudStyle).Render(fmt.Sprintf("%d info", normal)) +
		hudStyle.Render("  ") +
		errStyle.Inherit(hudStyle).Render(fmt.Sprintf("%d errors", err)) +
		hudStyle.Render("  ") +
		warnStyle.Inherit(hudStyle).Render(fmt.Sprintf("%d warnings", warn)) +
		hudStyle.Render(fmt.Sprintf("  [%s/%s: navigate | %s: quit]", n, p, q))

	hud := hudStyle.
		Width(m.width).
		MaxHeight(1).
		Render(hudText)

	// hud := lg.PlaceHorizontal(m.width, lg.Left, hudText, lg.WithWhitespaceStyle(hudStyle))

	sb.WriteString(hud)

	return sb.String()
}
