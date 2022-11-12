package main

import (
	"log"
	"weights/nets"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

type model struct {
	net   nets.Network
	input textarea.Model
}

func newModel() model {
	m := model{
		net:   *nets.NewNetwork(),
		input: textarea.New(),
	}
	return m
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top, m.input.View(), m.net.ShowChange())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
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
