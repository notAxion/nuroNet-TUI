package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/notAxion/nuroNet-TUI/constants"
	"github.com/notAxion/nuroNet-TUI/input"
	"github.com/notAxion/nuroNet-TUI/nets"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// margin   = lipgloss.NewStyle().Margin(0, 0, 0, 0)
	margin = lipgloss.NewStyle().Margin(1, 2, 0, 2)
	common = lipgloss.NewStyle().
		Margin(4, 2, 0, 0).
		Align(lipgloss.Center)
		// MarginTop(5).

	brd = lipgloss.Border{
		Left: lipgloss.ThickBorder().Left,
	}

	docStyle = common.Copy().
			Padding(1, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("238"))

	focusStyle = common.Copy().
			Padding(1, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("57"))

	modelStyle = common.Copy().
			Height(14). // just by counting lines
			Border(brd).
			BorderForeground(lipgloss.Color("57")).
		// Margin(4+3, 0, 3).
		PaddingLeft(2)

	helpStyle = lipgloss.NewStyle().
			MarginLeft(1).MarginBottom(1)
)

/*
var (
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	cursorLineStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("57")).
			Foreground(lipgloss.Color("230"))

	placeholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("238"))

	endOfBufferStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("235"))

	focusedPlaceholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("99"))

	focusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("238"))

	blurredBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.HiddenBorder())
)
*/

func main() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatalln(err)
	}
}

type Focus int

const (
	ws = iota
	bias
)

type model struct {
	net      nets.Network
	input    textinput.Model
	wsList   list.Model
	biasList list.Model
	focus    Focus
	keys     constants.KeyMap
	help     help.Model
}

func newModel() model {
	inp := textinput.New()
	inp.Prompt = "> "
	inp.Placeholder = ""
	inp.Width = 50

	m := model{
		net:   *nets.NewNetwork(),
		input: inp,
		focus: ws,
		keys:  constants.Keys,
		help:  constants.New(),
	}
	m.input.Validate = m.ValidateInput

	m.wsList = list.New(input.NewItems(m.net.WS), list.NewDefaultDelegate(), 0, 0)
	m.wsList.Title = "Wights"
	m.wsList.SetShowHelp(false)
	m.biasList = list.New(input.NewItems(m.net.Bias), list.NewDefaultDelegate(), 0, 0)
	m.biasList.Title = "Bias"
	m.biasList.SetShowHelp(false)
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	strs := []string{}
	strs = append(strs, m.wsStyle(m.wsList.View()))
	strs = append(strs, m.biasStyle(m.biasList.View()))
	strs = append(strs, modelStyle.Render(m.net.ShowChange()))
	template := lipgloss.JoinHorizontal(lipgloss.Left, strs...)
	template = lipgloss.JoinVertical(0, template, helpStyle.Render(m.help.View(m.keys)))
	// template += "\n" + helpStyle.Render(m.help.View(m.keys))
	if m.input.Focused() {
		// templete = margin.Render(templete + "\n" + m.input.View())
		return margin.Render(template + "\n" + m.input.View())
	}

	return margin.Render(template + "\n")
	// return templete
}

func (m model) wsStyle(s string) string {
	if m.focus == ws {
		return focusStyle.Render(s)
	} else if m.focus == bias {
		return docStyle.Render(s)
	} else {
		return ""
	}
}

func (m model) biasStyle(s string) string {
	if m.focus == ws {
		return docStyle.Render(s)
	} else if m.focus == bias {
		return focusStyle.Render(s)
	} else {
		return ""
	}
}

// func (m model) margin(s string) string {
// 	ms := lipgloss.NewStyle().MarginBottom(2)
// }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		docStyle.Width(msg.Width / 4)
		focusStyle.Width(msg.Width / 4)
		// modelStyle.Height(msg.Height - modelStyle.GetVerticalFrameSize() - margin.GetVerticalMargins())
		// vert := modelStyle.GetVerticalFrameSize() + margin.GetVerticalMargins()
		vert := msg.Height - margin.GetVerticalFrameSize() - modelStyle.GetHeight()
		modelStyle.Margin(common.GetMarginTop()-3+vert/2, 0, vert/2-1)
		// modelStyle.Margin(3+vert/2, 0, vert/2)

		// modelStyle.Width(msg.Width/2)
		h, v := docStyle.GetFrameSize()
		m.wsList.SetSize(msg.Width-h, msg.Height-v)
		h, v = focusStyle.GetFrameSize()
		m.biasList.SetSize(msg.Width-h, msg.Height-v)
		m.help.Width = msg.Width
	case tea.KeyMsg:
		if m.input.Focused() {
			switch {
			case key.Matches(msg, m.keys.Quit, m.keys.Back):
				m.input.SetValue("")
				m.input.Blur()
			case key.Matches(msg, m.keys.Enter):
				var value float32
				if ff, err := strconv.ParseFloat(m.input.Value(), 32); err == nil && m.input.Value() != "" {
					value = float32(ff)
					switch m.focus {
					case ws:
						m.net.WSChange(m.wsList.Index(), value)
					case bias:
						m.net.BiasChange(m.biasList.Index(), value)

					}
				}
				m.input.SetValue("")
				m.input.Blur()
			}
			m.input.CursorEnd()
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			switch {
			case key.Matches(msg, m.keys.Help):
				m.help.ShowAll = !m.help.ShowAll
			case key.Matches(msg, m.keys.Quit):
				return m, tea.Quit
			case key.Matches(msg, m.keys.Save):
				err := m.net.Save()
				if err != nil {
					m.net.LogError(err)
				}
			}
			switch m.focus {
			case ws:
				switch {
				case key.Matches(msg, m.keys.Tab): //, "h", "j":
					m.focus = bias
				case key.Matches(msg, m.keys.Enter):
					value := fmt.Sprint(m.net.WS[m.wsList.Index()])
					m.input.SetValue(value)
					m.input.Focus()
				default:
					m.wsList, cmd = m.wsList.Update(msg)
					cmds = append(cmds, cmd)
				}
			case bias:
				switch {
				case key.Matches(msg, m.keys.Tab): //, "h", "j":
					m.focus = ws
				case key.Matches(msg, m.keys.Enter):
					value := fmt.Sprint(m.net.Bias[m.biasList.Index()])
					m.input.SetValue(value)
					m.input.Focus()
				default:
					m.biasList, cmd = m.biasList.Update(msg)
					cmds = append(cmds, cmd)
				}
			}
			// switch msg.String() {
			// case "ctrl+c", "q":
			// 	return m, tea.Quit
			// case "tab": //, "h", "j":
			// 	if m.focus == ws {
			// 		m.focus = bias
			// 	} else {
			// 		m.focus = ws
			// 	}
			// case "enter":
			// 	value := ""
			// 	switch m.focus {
			// 	case ws:
			// 		value = fmt.Sprint(m.net.WS[m.wsList.Index()])
			// 	case bias:
			// 		value = fmt.Sprint(m.net.Bias[m.biasList.Index()])
			// 	}
			// 	m.input.SetValue(value)
			// 	m.input.Focus()
			// default:
			// 	switch m.focus {
			// 	case ws:
			// 		m.wsList, cmd = m.wsList.Update(msg)
			// 	case bias:
			// 		m.biasList, cmd = m.biasList.Update(msg)
			// 	}
			// 	cmds = append(cmds, cmd)
			// }
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) ValidateInput(s string) error {
	_, err := strconv.ParseFloat(s, 32)
	return err
}

/*
func getStuff(ws, bias []float32) (string, int, float32) {
	choice := ""
	index := 0
	fmt.Print("What to change: ")
	fmt.Scanf("%1s%d", &choice, &index)
	if choice == "w" {
		fmt.Println(ws[index])
	} else if choice == "b" {
		fmt.Println(bias[index])
	}

	fmt.Print("\nHow much Change: ")
	var input float32
	fmt.Scan(&input)

	return choice, index, input
}

// Tried different random numbers
// Have found this seed is a good starting point
// weights seed 1667659351435904
// bias seed 1667659351435945
// func randFlArrn(input []float32, min, max int, seed int64) {
func _(input []float32, min, max int, seed int64) {
	// seed := time.Now().UnixMicro()
	r := rand.New(rand.NewSource(seed))
	// fmt.Println("seed", seed)
	if max-min <= 0 {
		panic("min is bigger than max :[ ")
	}
	for i := range input {
		input[i] = float32(r.Intn(max-min)+min) - r.Float32()
	}
}

func main() {
	net := nets.NewNetwork()

	net.ShowChange()
	if err := net.Save(); err != nil {
		fmt.Println("Network Saving error: ", err)
	}
	for {
		choice, index, change := getStuff(net.WS, net.Bias)
		if choice == "w" {
			net.Change(index, net.WS[index]+change)
		} else if choice == "b" {
			net.Bias[index] += change
		}

		net.ShowChange()
		fmt.Println("")
		fmt.Printf("weights: %v\n", net)
		fmt.Printf("weights.Bias: %v\n", net.Bias)
		fmt.Println()

	}
}
*/
