package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"wapuugotchicli/quiz"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	itemStyle = lipgloss.NewStyle().PaddingLeft(4)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("170"))
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(strs ...string) string {
			if len(strs) > 0 {
				return selectedItemStyle.Render("> " + strs[0])
			}
			return selectedItemStyle.Render("> ")
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	list      list.Model
	textInput textinput.Model
	choice    string
	quitting  bool
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "q":
			// Only quit if we're in the main menu, not in a submenu
			if m.choice == "" {
				m.quitting = true
				return m, tea.Quit
			} else {
				// Return to main menu
				m.choice = ""
				return m, nil
			}
		case "enter":
			// Handle specific menu items
			if m.choice == "" {
				i, ok := m.list.SelectedItem().(item)
				if ok {
					selectedItem := string(i)
					if selectedItem == "Exit" {
						m.quitting = true
						return m, tea.Quit
					} else {
						m.choice = selectedItem
					}
				}
			} else {
				// Return to main menu when viewing a submenu
				m.choice = ""
			}
			return m, nil
		}
	}

	// Only update the list when we're in the main menu
	if m.choice == "" {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	switch m.choice {
	case "Quiz":
		return quiz.Display() + "\nPress enter to return to the menu."
	case "View Files":
		filesText := listFiles()
		return filesText + "\nPress enter to return to the menu."
	}

	if m.quitting {
		return "\nGoodbye!\n"
	}
	return "\n" + titleStyle.Render("Wapuugotchi CLI") + "\n\n" + m.list.View()
}

func listFiles() string {
	cmd := exec.Command("ls", "-l")
	output, err := cmd.Output()
	if err != nil {
		return "Error executing ls -l: " + err.Error()
	}

	return "Files in current directory:\n" + string(output)
}

func main() {
	items := []list.Item{
		item("Quiz"),
		item("Exit"),
	}

	const defaultWidth = 20
	const listHeight = 14

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.PaginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	l.Styles.HelpStyle = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	l.SetShowHelp(false)
	l.SetShowTitle(false)

	m := model{list: l}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
