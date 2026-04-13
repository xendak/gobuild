package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"log"
	
	"charm.land/bubbles/v2/key"
)

type Config struct {
	Editor string
	Block  bool

	WindowBg string

	// Accent: prompt symbol, borders, cursor.
	Accent string

	// Severity text colors.
	ErrorFg    string
	WarnFg     string
	NoteFg     string
	NormalFg   string
	FileFg     string
	LocationFg string

	// Bottom Status
	HudBg string
	HudFg string

	// Floating input overlay — inner text field.
	InputFg          string
	InputBorderColor string
	InputLabelFg     string
	InputCursorColor string

	// Outer rounded panel wrapping the input field.
	InputPanelBg     string
	InputPanelBorder string

	// Keys
	Quit       []string
	QuitPrompt []string
	NextMatch  []string
	PrevMatch  []string
	OpenPrompt []string
	OpenEditor []string
	EscPrompt  []string
	SubmitCmd  []string
}

type KeyBindings struct {
	Quit       key.Binding
	QuitPrompt key.Binding
	NextMatch  key.Binding
	PrevMatch  key.Binding
	OpenPrompt key.Binding
	OpenEditor key.Binding
	EscPrompt  key.Binding
	SubmitCmd  key.Binding
}

func (k KeyBindings) ShortHelp() []key.Binding {
	return []key.Binding{k.NextMatch, k.PrevMatch, k.Quit}
}

func (k KeyBindings) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.NextMatch, k.PrevMatch, k.OpenEditor},
		{k.OpenPrompt, k.EscPrompt, k.Quit},
	}
}

func BuildKeyMap(cfg Config) KeyBindings {
	return KeyBindings{
		Quit:       key.NewBinding(key.WithKeys(cfg.Quit...), key.WithHelp(cfg.Quit[0], "quit")),
		QuitPrompt: key.NewBinding(key.WithKeys(cfg.QuitPrompt...), key.WithHelp(cfg.QuitPrompt[0], "quit")),
		NextMatch:  key.NewBinding(key.WithKeys(cfg.NextMatch...), key.WithHelp(cfg.NextMatch[0], "next match")),
		PrevMatch:  key.NewBinding(key.WithKeys(cfg.PrevMatch...), key.WithHelp(cfg.PrevMatch[0], "prev match")),
		OpenPrompt: key.NewBinding(key.WithKeys(cfg.OpenPrompt...), key.WithHelp(cfg.OpenPrompt[0], "open prompt")),
		OpenEditor: key.NewBinding(key.WithKeys(cfg.OpenEditor...), key.WithHelp(cfg.OpenEditor[0], "open in editor")),
		EscPrompt:  key.NewBinding(key.WithKeys(cfg.EscPrompt...), key.WithHelp(cfg.EscPrompt[0], "close prompt")),
		SubmitCmd:  key.NewBinding(key.WithKeys(cfg.SubmitCmd...), key.WithHelp(cfg.SubmitCmd[0], "run command")),
	}
}

func DefaultConfig() Config {
	return Config{
		Editor: "wez-hx",
		Block:  false,

		WindowBg:         "2",
		Accent:           "12",
		ErrorFg:          "1",
		WarnFg:           "3",
		NoteFg:           "5",
		NormalFg:         "7",
		FileFg:           "4",
		LocationFg:       "2",
		HudBg:            "0",
		HudFg:            "7",
		InputFg:          "7",
		InputBorderColor: "12",
		InputLabelFg:     "7",
		InputCursorColor: "12",
		InputPanelBg:     "",
		InputPanelBorder: "12",

		Quit:       []string{"q", "ctrl+c"},
		QuitPrompt: []string{"ctrl+c"},
		NextMatch:  []string{"n", "down"},
		PrevMatch:  []string{"N", "up"},
		OpenPrompt: []string{":"},
		OpenEditor: []string{"enter"},
		EscPrompt:  []string{"esc"},
		SubmitCmd:  []string{"enter"},
	}
}

func LoadConfig() (Config, error) {
	cfg := DefaultConfig()

	path, err := configPath()
	if err != nil {
		return cfg, nil
	}

	data, err := os.Open(path)
	if err != nil {
		// NOTE(xendak): lets not check if no config exists.. we have sane defaults i hope?
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}
	defer data.Close()

	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		line := strings.ReplaceAll(scanner.Text(), " ", "")

		// # == Comment
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		before, after, valid := strings.Cut(line, ":")
		if !valid || after == "" {
			continue
		}
		log.Printf("%s -> %s", before, after)

		// NOTE(xendak): do i need trim here?
		key := strings.ToLower(before)
		value := after


		switch key {
		case "editor":
			cfg.Editor = value

		case "windowbg":
			cfg.WindowBg = value
		case "accent":
			cfg.Accent = value
		case "errorfg":
			cfg.ErrorFg = value
		case "warnfg":
			cfg.WarnFg = value
		case "notefg":
			cfg.NoteFg = value
		case "normalfg":
			cfg.NormalFg = value
		case "filefg":
			cfg.FileFg = value
		case "locationfg":
			cfg.LocationFg = value
		case "hudbg":
			cfg.HudBg = value
		case "hudfg":
			cfg.HudFg = value
		case "inputfg":
			cfg.InputFg = value
		case "inputbordercolor":
			cfg.InputBorderColor = value
		case "inputlabelfg":
			cfg.InputLabelFg = value
		case "inputcursorcolor":
			cfg.InputCursorColor = value
		case "inputpanelbg":
			cfg.InputPanelBg = value
		case "inputpanelborder":
			cfg.InputPanelBorder = value

		case "quit":
			cfg.Quit = strings.Split(value, ",")
		case "quitprompt":
			cfg.QuitPrompt = strings.Split(value, ",")
		case "nextmatch":
			cfg.NextMatch = strings.Split(value, ",")
		case "prevmatch":
			cfg.PrevMatch = strings.Split(value, ",")
		case "openprompt":
			cfg.OpenPrompt = strings.Split(value, ",")
		case "openeditor":
			cfg.OpenEditor = strings.Split(value, ",")
		case "escprompt":
			cfg.EscPrompt = strings.Split(value, ",")
		case "submitcmd":
			cfg.SubmitCmd = strings.Split(value, ",")
		}
	}

	return cfg, scanner.Err()
}

func configPath() (string, error) {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "gobuild", "config"), nil
}
