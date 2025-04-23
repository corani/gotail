package main

// A simple program demonstrating the text area component from the Bubbles
// component library.

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	var (
		lines int
		help  bool
	)

	flag.IntVar(&lines, "lines", 5, "the number of lines to print")
	flag.IntVar(&lines, "n", 5, "the number of lines to print")
	flag.BoolVar(&help, "help", false, "print help screen")
	flag.Parse()

	if help {
		flag.PrintDefaults()
		return
	}

	args := flag.Args()

	if len(args) == 0 {
		log.Fatal("no command provided")
	}

	reader, err := execCommand(args[0], args[1:]...)
	if err != nil {
		log.Fatal(err)
	}

	p := tea.NewProgram(initialModel(reader, lines))

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func execCommand(program string, args ...string) (io.ReadCloser, error) {
	cmd := exec.Command(program, args...)
	out, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout

	cmd.Start()

	go cmd.Wait()

	return out, err
}

type model struct {
	viewport viewport.Model
	messages []string
	height   int
	reader   *bufio.Reader
	closer   io.Closer
}

func initialModel(reader io.ReadCloser, height int) model {
	vp := viewport.New(30, height)

	return model{
		messages: []string{},
		viewport: vp,
		height:   height,
		reader:   bufio.NewReader(reader),
		closer:   reader,
	}
}

func (m model) Init() tea.Cmd {
	return m.tail
}

type tailMsg string

func (m model) tail() tea.Msg {
	line, err := m.reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	return tailMsg(line)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var vpCmd tea.Cmd

	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = min(m.height, msg.Height)

		if len(m.messages) > 0 {
			// Wrap content before setting it.
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "")))
		}
		m.viewport.GotoBottom()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.closer.Close()

			return m, tea.Quit
		}
	case tailMsg:
		m.messages = appendN(m.messages, string(msg), 100)
		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "")))
		m.viewport.GotoBottom()

		return m, m.tail
	}

	return m, vpCmd
}

// append appends a line to the list of lines, and if the list is longer
// than the given size, it removes the first line.
func appendN(list []string, item string, size int) []string {
	if len(list) >= size {
		list = list[1:]
	}

	return append(list, item)
}

func (m model) View() string {
	return m.viewport.View()
}
