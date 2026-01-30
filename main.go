package main

import (
	"fmt"
	"io"
	"os"

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

	questionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#5A5A5A")).
			Padding(0, 1).
			MarginBottom(1)
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
	list        list.Model
	textInput   textinput.Model
	choice      string
	quitting    bool
	currentQuiz *quiz.Quiz
	quizItems   []list.Item
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
				m.currentQuiz = nil
				m.list.SetItems([]list.Item{
					item("Quiz"),
					item("Exit"),
				})
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
					} else if selectedItem == "Quiz" {
						m.choice = selectedItem
						// Get quiz data and set up quiz items
						quizData := quiz.GetQuiz()
						if quizData != nil {
							m.currentQuiz = quizData
							var quizItems []list.Item
							for _, option := range quizData.Options {
								quizItems = append(quizItems, item(fmt.Sprintf("%s", option)))
							}
							quizItems = append(quizItems, item("Back to Menu"))
							m.list.SetItems(quizItems)
						}
					}
				}
			} else if m.choice == "Quiz" {
				// Handle quiz option selection
				i, ok := m.list.SelectedItem().(item)
				if ok {
					selectedItem := string(i)
					if selectedItem == "Back to Menu" {
						// Return to main menu
						m.choice = ""
						m.currentQuiz = nil
						m.list.SetItems([]list.Item{
							item("Quiz"),
							item("Exit"),
						})
					} else {
						// Handle quiz answer selection here
						fmt.Printf("Selected answer: %s\n", selectedItem)
					}
				}
			}
			return m, nil
		}
	}

	// Update the list
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	var header string
	if m.choice == "Quiz" && m.currentQuiz != nil {
		header = questionStyle.Render(m.currentQuiz.Question)
	} else {
		header = titleStyle.Render("Wapuugotchi CLI")
	}

	return "\n" + header + "\n\n" + m.list.View()
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
