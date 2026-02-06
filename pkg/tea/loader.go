package tea

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-isatty"
)

type TaskResult struct {
	Team    string
	Year    int
	Success bool
}

type model struct {
	totalTasks int
	completed  int
	quitting   bool
	taskCh     chan TaskResult
	lastTasks  []string
}

type taskMsg TaskResult

type Program struct {
	taskCh      chan TaskResult
	programDone chan struct{}
	simpleMode  bool
	program     *tea.Program
	firstTask   bool
}

func NewProgram(totalTasks int) *Program {
	p := &Program{
		taskCh:      make(chan TaskResult),
		programDone: make(chan struct{}),
		simpleMode:  !isatty.IsTerminal(os.Stderr.Fd()),
		firstTask:   true,
	}

	if !p.simpleMode {
		opts := []tea.ProgramOption{
			tea.WithOutput(os.Stderr),
			tea.WithInput(os.Stdin),
		}
		p.program = tea.NewProgram(model{
			totalTasks: totalTasks,
			taskCh:     p.taskCh,
			lastTasks:  []string{},
		}, opts...)

		go func() {
			if _, err := p.program.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			}
			close(p.programDone)
		}()
	}

	return p
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		result, ok := <-m.taskCh
		if !ok {
			return tea.QuitMsg{}
		}
		return taskMsg(result)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case taskMsg:
		m.completed++
		m.lastTasks = append(m.lastTasks, formatTask(TaskResult(msg)))
		if len(m.lastTasks) > 5 {
			m.lastTasks = m.lastTasks[1:]
		}
		if m.completed >= m.totalTasks {
			return m, tea.Quit
		}
		return m, func() tea.Msg {
			result, ok := <-m.taskCh
			if !ok {
				return tea.QuitMsg{}
			}
			return taskMsg(result)
		}
	case tea.QuitMsg:
		return m, tea.Quit
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	bar := createProgressBar(m.completed, m.totalTasks)

	var taskLines []string
	for _, task := range m.lastTasks {
		if strings.HasPrefix(task, "✅") {
			taskLines = append(taskLines, successStyle.Render(task))
		} else {
			taskLines = append(taskLines, skipStyle.Render(task))
		}
	}
	tasks := strings.Join(taskLines, "\n  ")

	return fmt.Sprintf(
		"\n%s\n\n%s [%s] %d/%d\n\n  %s\n",
		titleStyle.Render("🏈 Fantasy Football Fetcher"),
		progressStyle.Render("Progress"),
		bar,
		m.completed,
		m.totalTasks,
		tasks,
	)
}

func (p *Program) Start() {
	if p.simpleMode {
		fmt.Fprintln(os.Stderr, "\n  🏈 Fantasy Football Fetcher")
	}
}

func (p *Program) Update(result TaskResult) {
	if p.simpleMode {
		if p.firstTask {
			fmt.Fprint(os.Stderr, "\n\n")
			p.firstTask = false
		}
		fmt.Fprintf(os.Stderr, "  %s\n", formatTask(result))
	} else {
		p.taskCh <- result
	}
}

func (p *Program) Quit() {
	close(p.taskCh)
	if p.simpleMode {
		fmt.Fprintln(os.Stderr, "\n\n  Done!")
	} else {
		<-p.programDone
	}
}

func formatTask(tr TaskResult) string {
	status := "✅"
	if !tr.Success {
		status = "⏭️"
	}
	return fmt.Sprintf("%s %s %d", status, tr.Team, tr.Year)
}

func createProgressBar(completed, total int) string {
	if total == 0 {
		return strings.Repeat("░", 30)
	}
	percent := float64(completed) / float64(total) * 100
	barWidth := 30
	filled := int(percent / 100 * float64(barWidth))
	return strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
}

var (
	titleStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)
	progressStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	successStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("70"))
	skipStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
)
