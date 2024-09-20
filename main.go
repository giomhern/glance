package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type status int

const (
	todo status = iota
	inProgress
	done
)

/* CUSTOM ITEM */

type Task struct {
	status      status
	title       string
	description string
}

// implement the list.item interface

func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}
func (t Task) Description() string {
	return t.description
}

/* MAIN MODEL */

type Model struct {
	list list.Model
	err  error
}

func New() *Model {
	return &Model{}
}

// TODO: call this on tea.windowsize message
func (m *Model) initList(width int, height int) {
	m.list = list.New(
		[]list.Item{},
		list.NewDefaultDelegate(),
		width,
		height,
	)
	m.list.Title = "Todo List"
	m.list.SetItems([]list.Item{
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
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.initList(msg.Width, msg.Height)
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.list.View()
}

func main() {
	m := New()
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
