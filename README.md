# gobuild

> work in progress, tui version of emacs compilation-mode. (?)
> [Changelog](README#Changelog)

I use helix (sometimes vi, emacs), i tried using wezterm regex to do this for
me.. was a bit of a chore, so i tried building something that was useful for me
:D

Built in Go as an excuse to learn the language and try out
[bubbletea](https://github.com/charmbracelet/bubbletea), which has a really cool
name.

### Usage

```bash
# pipe
go build . 2>&1 | gobuild

# file
gobuild somefile somefile2

# run a command directly
gobuild -c "make"

# interactive (no args)
gobuild
```

### Features

- pipe, file, and `-c`/`--cmd` input modes
- interactive mode — type a command to run while the TUI is open
- open matched line in editor (`$EDITOR`, currently defaults to `wez-hx`)
- colors via lipgloss (hardcoded for now)
- navigate between errors/matches

### Modes

| Mode         | Description                            |
| ------------ | -------------------------------------- |
| Results      | browse parsed output, navigate matches |
| Interactive  | type a command to run                  |
| Loading      | debating if needed or if i can ditch   |
| Passsthrough | temporary idea for warning only output |

### Keybinds

| Key            | Action                         | Mode       |
| -------------- | ------------------------------ | ----------- |
| `n` / `down`   | next match                     | Results     |
| `N` / `up`     | previous match                 | Results     |
| `enter`        | open in editor / enter command | All         |
| `:`            | switch to interactive mode     | Results     |
| `esc`          | back to results                | Interactive |
| `q` / `ctrl+c` | quit                           | Results     |
| `ctrl+c`       | quit                           | Interactive |

---

## Changelog

> expect breaking/crazy changes.

### [Unreleased] — in progress

- trying to get PIPE command feedback so maybe i can update it?
- offer a config file for colors, editor, etc.
- some flags to switch been stuff

---

### 0.0.4

- interactive mode with command input box (`:` to open, `esc` to close)
- dynamic result parsing after command runs
- command execution with simple editor integration (hard coded)
- open matched line directly in `$EDITOR` (`enter`)

---

### 0.0.3

- show all lines
- navigation still jumps only between matches
- severity grading (error / warning / note / hint)

---

### 0.0.2

- scrolling TUI with lipgloss styles
- matched lines rendered as `file:(line:col) Severity(msg)`
- `n`/`N` navigation between matches

---

### 0.0.1

- simple line scanner with regex support (odin, generic, grep, python)
