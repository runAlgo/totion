package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	newFileInput           textinput.Model
	noteTextArea           textarea.Model
	createFileInputVisible bool
	noteTextAreaVisible    bool
	currentFile            *os.File
	list                   list.Model
	showingList            bool
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

var (
	vaultDir    string
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	docStyle    = lipgloss.NewStyle().Margin(1, 2)
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error getting home directory", err)
	}

	vaultDir = fmt.Sprintf("%s/.totion", homeDir)
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-8)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			if m.createFileInputVisible {
				m.createFileInputVisible = false
			}
			if m.currentFile != nil {
				m.noteTextArea.SetValue("")
				m.currentFile = nil
			}
			if m.showingList {
				if m.list.FilterState() == list.Filtering {
					break
				}
				m.showingList = false
			}
			return m, nil
		case "ctrl+n":
			m.createFileInputVisible = true
			return m, nil
		case "ctrl+l":
			noteList := listFiles()
			m.list.SetItems(noteList)
			m.showingList = true
			return m, nil
		case "ctrl+d":
			if m.showingList {
				selectedItem, ok := m.list.SelectedItem().(item)
				if !ok {
					break
				}
				filePath := fmt.Sprintf("%s/%s", vaultDir, selectedItem.title)
				err := os.Remove(filePath)
				if err != nil {
					log.Printf("Failed to delete file: %v", err)
				} else {
					// refresh list after deletion
					m.list.SetItems(listFiles())
				}
			}
		case "ctrl+s":
			// textarea value -> write it in the file descriptor and close it
			if m.currentFile == nil {
				break // safety from privent error when press "ctrl+s": if current file is not opened
			}

			if err := m.currentFile.Truncate(0); err != nil {
				fmt.Println("can not save the file :(")
				return m, nil
			}

			if _, err := m.currentFile.Seek(0, 0); err != nil {
				fmt.Println("can not save the file :(")
				return m, nil
			}

			if _, err := m.currentFile.WriteString(m.noteTextArea.Value()); err != nil {
				fmt.Println("can not save the file :(")
				return m, nil
			}

			if err := m.currentFile.Close(); err != nil {
				fmt.Println("can not close the file.")
			}

			m.currentFile = nil
			m.noteTextArea.SetValue("")

			return m, nil
		case "enter":
			if m.currentFile != nil {
				break
			}

			if m.showingList {
				item, ok := m.list.SelectedItem().(item)
				if ok {
					filepath := fmt.Sprintf("%s/%s", vaultDir, item.title)
					content, err := os.ReadFile(filepath)
					if err != nil {
						log.Printf("failed to reading file: %v", err)
						return m, nil
					}
					m.noteTextArea.SetValue(string(content))

					f, err := os.OpenFile(filepath, os.O_RDWR, 0644)
					if err != nil {
						log.Printf("Error reading file: %v", err)
					}
					m.currentFile = f
					m.showingList = false
				}
				return m, nil
			}
			// todo: create file
			fileName := m.newFileInput.Value()
			if fileName != "" {
				filepath := fmt.Sprintf("%s/%s.md", vaultDir, fileName)

				// it checks is there already have same file exists on our system
				// if it exists it's stop and do not create same file again
				if _, err := os.Stat(filepath); err == nil {
					return m, nil
				}

				// If user input file not not exists, now create a file
				f, err := os.Create(filepath)
				if err != nil {
					log.Fatalf("%v", err)
				}
				m.currentFile = f
				m.createFileInputVisible = false
				m.newFileInput.SetValue("\n")

			}
			return m, nil

		}
	}

	if m.createFileInputVisible {
		m.newFileInput, cmd = m.newFileInput.Update(msg)
	}

	if m.currentFile != nil {
		m.noteTextArea, cmd = m.noteTextArea.Update(msg)
	}

	if m.showingList {
		m.list, cmd = m.list.Update(msg)
	}

	return m, cmd
}

// Return the updated model to the Bubble Tea runtime for processing.
// Note that we're not returning a command.

func (m model) View() string {
	var style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("16")).
		Background(lipgloss.Color("205")).
		PaddingTop(1).
		PaddingLeft(1).
		PaddingBottom(1).
		Width(50)

	var style2 = lipgloss.NewStyle().
		Bold(true)

		// fmt.Println(style.Render("Hello, kitty"))
	Welcome := style.Render("Welcome to Totion 🧠")
	help := style2.Render("Ctrl+N: new file • Ctrl+L: list • Esc: back • Ctrl+S: save • Ctrl+Q: quit")
	view := ""
	if m.createFileInputVisible {
		view = m.newFileInput.View()
	}

	if m.currentFile != nil {
		view = m.noteTextArea.View()
	}

	if m.showingList {
		view = m.list.View()
	}
	return fmt.Sprintf("\n%s\n\n%s\n\n%s", Welcome, view, help)
}

func initializeModel() model {

	err := os.MkdirAll(vaultDir, 0750)
	if err != nil {
		log.Fatal(err)
	}

	// initialise new file input
	ti := textinput.New()
	ti.Placeholder = "What would you like to call it?"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50
	ti.Cursor.Style = cursorStyle
	ti.PromptStyle = cursorStyle
	ti.TextStyle = cursorStyle

	// textarea

	ta := textarea.New()
	ta.Placeholder = "Write your note here..."
	ta.Focus()
	ta.ShowLineNumbers = false

	// list
	noteList := listFiles()

	finalList := list.New(noteList, list.NewDefaultDelegate(), 0, 0)
	finalList.Title = "All notes 📘"
	finalList.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("16")).
		Background(lipgloss.Color("254")).
		Padding(0, 1)

	return model{
		newFileInput:           ti,
		createFileInputVisible: false,
		noteTextArea:           ta,
		list:                   finalList,
	}
}

func main() {
	p := tea.NewProgram(initializeModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func listFiles() []list.Item {
	items := make([]list.Item, 0)

	entries, err := os.ReadDir(vaultDir)
	if err != nil {
		log.Fatal("Error reading notes")
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				continue
			}

			modTime := info.ModTime().Format("2006-01-02 15:04")

			items = append(items, item{
				title: entry.Name(),
				desc:  fmt.Sprintf("Modified: %s", modTime),
			})
		}

	}
	return items
}
