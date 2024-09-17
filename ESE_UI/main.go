package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const ENTRY_COUNT = 1522
const ESE_PATH = "../EveryStoryEverList.md"

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readChoicesFromFile(filename string) []string {
	// THIS CAN BE REFINED TO USE THE NUM OF THE LAST ENTRY
	choices := make([]string, ENTRY_COUNT)
	currentChoice := 0

	file, err := os.Open(filename)
	check(err)
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')

		if err == io.EOF {
			break
		} else {
			check(err)
		}

		choices[currentChoice] = line
		currentChoice++
	}

	return choices
}

func initialModel() model {
	return model{
		choices:  readChoicesFromFile(ESE_PATH),
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "Mark the comics you have read.\n\n"

	for i, choice := range m.choices {

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there has been an error: %v", err)
		os.Exit(1)
	}
}
