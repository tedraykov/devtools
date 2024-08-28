package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Tools interface {
  Run() error
  InstallZsh() error
  InstallGo() error
  InstallGcc() error
  InstallMake() error
  InstallUnzip() error
  InstallDocker() error
  InstallTmux() error
  InstallNode() error
  InstallPython() error
  InstallPoetry() error
  InstallNeovim() error
  InstallBitwarden() error
}

type item struct {
	title    string
	selected bool
}

type viewState int

const (
	osSelection viewState = iota
	toolSelection
)

type model struct {
	state      viewState
	osChoices  []string
	osCursor   int
	osSelected string
	tools      []item
	toolCursor int
}

func initialModel() model {
	return model{
		state:     osSelection,
		osChoices: []string{"Ubuntu", "MacOS"},
		tools: []item{
			{title: "All tools selected", selected: true},
			{title: "zsh", selected: true},
			{title: "make", selected: true},
			{title: "gcc", selected: true},
			{title: "unzip", selected: true},
			{title: "docker", selected: true},
      {title: "tmux", selected: true},
      {title: "go", selected: true},
      {title: "node", selected: true},
      {title: "python", selected: true},
      {title: "poetry", selected: true},
      {title: "neovim", selected: true},
      {title: "bitwarden", selected: true},
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case osSelection:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "l":
				if m.osCursor > 0 {
					m.osCursor--
				}
			case "down", "k":
				if m.osCursor < len(m.osChoices)-1 {
					m.osCursor++
				}
			case "enter":
				m.osSelected = m.osChoices[m.osCursor]
				m.state = toolSelection
			}
		case toolSelection:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "l":
				if m.toolCursor > 0 {
					m.toolCursor--
				}
			case "down", "k":
				if m.toolCursor < len(m.tools)-1 {
					m.toolCursor++
				}
			case " ":
				if m.toolCursor == 0 {
					allSelected := !m.tools[0].selected
					for i := range m.tools {
						m.tools[i].selected = allSelected
					}
				} else {
					m.tools[m.toolCursor].selected = !m.tools[m.toolCursor].selected
					allSelected := true
					anySelected := false
					for i := 1; i < len(m.tools); i++ {
						if m.tools[i].selected {
							anySelected = true
						} else {
							allSelected = false
						}
					}
					if allSelected {
						m.tools[0].selected = true
					} else if anySelected {
						m.tools[0].selected = true
					} else {
						m.tools[0].selected = false
					}
				}
			case "enter":
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	switch m.state {
	case osSelection:
		s := "Select an operating system:\n\n"
		for i, choice := range m.osChoices {
			cursor := " "
			if m.osCursor == i {
				cursor = ">"
			}
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}
		s += "\nPress enter to select, up/down to move\n"
		return s
	case toolSelection:
		s := fmt.Sprintf("Selected OS: %s\n\nSelect tools:\n\n", m.osSelected)
		for i, item := range m.tools {
			cursor := " "
			if m.toolCursor == i {
				cursor = ">"
			}
			checked := " "
			if i == 0 {
				allSelected := true
				anySelected := false
				for j := 1; j < len(m.tools); j++ {
					if m.tools[j].selected {
						anySelected = true
					} else {
						allSelected = false
					}
				}
				if allSelected {
					checked = "x"
				} else if anySelected {
					checked = "-"
				}
			} else if item.selected {
				checked = "x"
			}
			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, item.title)
		}
		s += "\nPress space to select/unselect, up/down to move, enter to submit\n"
		return s
	default:
		return "Unknown state"
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		return
	}

	if model, ok := m.(model); ok {
		fmt.Printf("Selected OS: %s\n", model.osSelected)
		fmt.Println("Selected tools:")
		var selected []string
		for i, item := range model.tools {
			if item.selected && i > 0 {
				selected = append(selected, item.title)
			}
		}
		if len(selected) == 0 {
			fmt.Println("No tools selected.")
      return
		}

    var tools Tools

    if model.osSelected == "Ubuntu" {
      tools = &UbuntuTools{tools: selected}
    }

    if model.osSelected == "MacOS" {
       tools = &MacOsTools{tools: selected}
    }

    if tools == nil {
      fmt.Println("Unknown OS")
      return
    }

    if err := tools.Run(); err != nil {
      fmt.Printf("Error installing and configuring tools: %v", err)
      return
    }
	}
}
