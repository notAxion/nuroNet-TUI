package constants

import (
	"reflect"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

type KeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Enter key.Binding
	Save  key.Binding
	Back  key.Binding
	Tab   key.Binding
	Help  key.Binding
	Quit  key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.Enter, k.Save, k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	tp := reflect.TypeOf(k)
	vl := reflect.ValueOf(k)
	bindings := make([][]key.Binding, tp.NumField())
	for i := 0; i < tp.NumField(); i++ {
		// bindings[i] = tp.Field(i).
		if _, isOkay := vl.Field(i).Interface().(key.Binding); !isOkay {
			continue
		}
		binding := vl.Field(i).Interface().(key.Binding)
		bindings[i] = []key.Binding{binding}
	}
	// return [][]key.Binding{
	// 	{k.Up}, {k.Down}, {k.Left}, {k.Right},
	// 	{k.Tab}, {k.Enter}, {k.Save}, {k.Help}, {k.Quit},
	// }
	return bindings
}

func New() help.Model {
	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "236",
		Dark:  "248",
	})

	descStyle := lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "239",
		Dark:  "242",
	})

	sepStyle := lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "#DDDADA",
		Dark:  "#3C3C3C",
	})

	return help.Model{
		ShortSeparator: " • ",
		FullSeparator:  " • ",
		Ellipsis:       "…",
		Styles: help.Styles{
			ShortKey:       keyStyle,
			ShortDesc:      descStyle,
			ShortSeparator: sepStyle,
			Ellipsis:       sepStyle.Copy(),
			FullKey:        keyStyle.Copy(),
			FullDesc:       descStyle.Copy(),
			FullSeparator:  sepStyle.Copy(),
		},
	}
}

var Keys = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "right"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Save: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("Ctrl+s", "save Network"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "more Help"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "Next"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("q", "quit"),
	),
}
