package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type status int

const divisor = 4

const (
	todo status = iota
	inProgress
	done
)

var models []tea.Model

const (
	board status = iota
	form
)

/* STYLING */
var (
	columnStyle  = lipgloss.NewStyle().Padding(1, 2)
	focusedStyle = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("62"))
	// helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("41"))
)

/* FORM MODEL */

type Form struct {
	title   textinput.Model
	desc    textarea.Model
	focused status
}

func NewForm(focused status) *Form {
	form := &Form{focused: focused}
	form.title = textinput.New()
	form.title.Focus()
	form.desc = textarea.New()

	return form

}

func (m Form) Init() tea.Cmd {
	return nil
}

func (m Form) CreateTask() tea.Msg {
	task := NewTask(m.focused, m.title.Value(), m.desc.Value())
	fmt.Println(task)
	return TaskCreatedMsg{Task: task}
}

type TaskCreatedMsg struct {
	Task Task
}

func (m Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.title.Focused() {
				m.title.Blur()
				m.desc.Focus()
				return m, textarea.Blink
			} else {
				return models[board], m.CreateTask
			}
		}
	}

	if m.title.Focused() {
		m.title, cmd = m.title.Update(msg)
	} else {
		m.desc, cmd = m.desc.Update(msg)
	}
	return m, cmd
}

func (m Form) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.title.View(), m.desc.View())
}

/* CUSTOM ITEM */

type Task struct {
	status      status
	title       string
	description string
}

// implement the list.item interface

func NewTask(status status, title string, description string) Task {
	return Task{status: status, title: title, description: description}
}

func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}
func (t Task) Description() string {
	return t.description
}

func (t *Task) Next() {
	if t.status == done {
		t.status = todo
	} else {
		t.status++
	}
}

/* MAIN MODEL */

type Model struct {
	focused status
	lists   []list.Model
	// err      error
	loaded   bool
	quitting bool
}

/* MODEL MANAGEMENT */
// var models []tea.Model

// const (
// 	model status = iota
// 	form
// )

func New() *Model {
	return &Model{}
}

// DEBUG: index issue
func (m *Model) MoveToNext() tea.Msg {
	selectItem := m.lists[m.focused].SelectedItem()
	selectedTask := selectItem.(Task)
	m.lists[selectedTask.status].RemoveItem(m.lists[m.focused].Index())
	selectedTask.Next()
	m.lists[selectedTask.status].InsertItem(len(m.lists[selectedTask.status].Items())-1, list.Item(selectedTask))
	return nil

}

// TODO: go to next and prev list
func (m *Model) Next() {
	if m.focused == done {
		m.focused = todo
	} else {
		m.focused++
	}

}

func (m *Model) Prev() {
	if m.focused == todo {
		m.focused = done
	} else {
		m.focused--
	}
}

func (m *Model) initLists(width int, height int) {
	defaultList := list.New(
		[]list.Item{},
		list.NewDefaultDelegate(),
		width/divisor,
		height/2,
	)
	defaultList.SetShowHelp(false)
	m.lists = []list.Model{defaultList, defaultList, defaultList}

	// init the todo's
	m.lists[todo].Title = "Todo List"
	m.lists[todo].SetItems([]list.Item{
		Task{
			status:      todo,
			title:       "Buy milk",
			description: "Strawberry milk",
		},
		Task{
			status:      todo,
			title:       "Eat sushi",
			description: "California roll & miso soup",
		},
		Task{
			status:      todo,
			title:       "Fold laundry",
			description: "or wear wrinkly t-shirts",
		},
	})

	// init in progress
	m.lists[inProgress].Title = "In Progress"
	m.lists[inProgress].SetItems([]list.Item{
		Task{
			status:      todo,
			title:       "Write code",
			description: "Don't worry, it's go",
		},
	})

	// init done
	m.lists[done].Title = "Done"
	m.lists[done].SetItems([]list.Item{
		Task{
			status:      todo,
			title:       "Stay cool",
			description: "Drink lots of water",
		},
	})
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			columnStyle.Width(msg.Width / divisor)
			focusedStyle.Width(msg.Width / divisor)
			columnStyle.Height(msg.Height - divisor)
			focusedStyle.Height(msg.Height - divisor)
			m.initLists(msg.Width, msg.Height)
			m.loaded = true
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "left", "h":
			m.Prev()
		case "right", "l":
			m.Next()
		case "enter":
			return m, m.MoveToNext
		case "n":
			form := NewForm(m.focused)
			return form, nil // Switch to form model
		}
	case TaskCreatedMsg:
		task := msg.Task
		m.lists[task.status].InsertItem(len(m.lists[task.status].Items()), task)
		return m, nil // Switch back to main model
	}
	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}
	if m.loaded {
		todoView := m.lists[todo].View()
		inProgressView := m.lists[inProgress].View()
		doneView := m.lists[done].View()
		switch m.focused {
		case inProgress:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				columnStyle.Render(todoView),
				focusedStyle.Render(inProgressView),
				columnStyle.Render(doneView),
			)
		case done:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				columnStyle.Render(todoView),
				columnStyle.Render(inProgressView),
				focusedStyle.Render(doneView),
			)
		default:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				focusedStyle.Render(todoView),
				columnStyle.Render(inProgressView),
				columnStyle.Render(doneView),
			)
		}
	} else {
		return "Loading..."
	}
}

func main() {
	models = []tea.Model{New(), NewForm(todo)}
	m := models[board]
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
