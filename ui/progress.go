package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/zostay/go-std/generic"
)

const (
	defaultWidgetCount = 20
	standardWidgetSize = 4
	compactWidgetSize  = 1
)

type phaseStatus int

const (
	phaseUnknown  phaseStatus = iota
	phasePending              // red
	phaseWorking              // yellow
	phaseComplete             // green
)

const (
	redCircle    = "\U0001f534"
	yellowCircle = "\U0001f7e1"
	greenCircle  = "\U0001f7e2"
	purpleCircle = "\U0001f3e3"
)

const progressWidget = ""

type phase struct {
	status    phaseStatus
	name      string
	operation string
}

type Progress struct {
	term         *Terminal
	state        *State
	phases       []phase
	widgets      map[string]WidgetID
	currentPhase int
	compact      bool
	title        map[string]string
}

func NewProgress(tty *os.File) *Progress {
	return &Progress{
		term:         NewTerminal(tty),
		currentPhase: -1,
		widgets:      make(map[string]WidgetID, defaultWidgetCount),
	}
}

func (p *Progress) SetPhases(phases []string) {
	p.state = NewState(p.term, defaultWidgetCount)

	p.phases = make([]phase, len(phases))
	for i, name := range phases {
		p.phases[i].status = phasePending
		p.phases[i].name = name
	}

	if pw, hasProgressWidget := p.widgets[progressWidget]; hasProgressWidget {
		p.state.DeleteWidget(pw)
	}

	phaseCount := len(p.phases)
	if len(p.phases) > 4 {
		phaseCount = 4
	}

	p.widgets[progressWidget] = p.state.AddWidget(phaseCount)

	p.UpdateProgress()
}

func (p *Progress) HasPhases() bool {
	return p.state != nil
}

func (p *Progress) StartPhase(phase int, taskCount int) {
	if p.state == nil {
		panic("cannot call StartPhase before SetPhases")
	}

	for i := 0; i < phase; i++ {
		p.phases[i].operation = ""
		p.phases[i].status = phaseComplete
	}

	p.phases[phase].status = phaseWorking

	for wkey, widget := range p.widgets {
		if wkey == progressWidget {
			continue
		}

		p.state.DeleteWidget(widget)
		delete(p.widgets, wkey)
	}

	p.currentPhase = phase
	if taskCount == 0 {
		p.compact = true
		p.title = map[string]string{}
	} else {
		p.compact = (p.state.term.Height()-9)/taskCount < standardWidgetSize
		p.title = make(map[string]string, taskCount)
	}

	p.UpdateProgress()
}

func (p *Progress) RegisterTask(name, title string) {
	if p.state == nil {
		panic("cannot call RegisterTask before SetPhases")
	}

	p.title[name] = title
	w := p.state.AddWidget(p.TaskWidgetSize())
	p.widgets[name] = w
	if p.compact {
		p.state.Set(w, 0, title+": ...")
	} else {
		p.state.SetTitle(w, title)
	}

	p.UpdateProgress()
}

func (p *Progress) UpdateProgress() {
	if p.state == nil {
		return
	}

	var stop int
	phaseEnd := len(p.phases) - 1
	currentPhase := generic.Max(0, p.currentPhase)
	switch phaseEnd {
	case currentPhase, currentPhase + 1:
		stop = phaseEnd
	default:
		stop = generic.Min(phaseEnd, currentPhase+3)
	}
	start := generic.Max(0, stop-3)

	pw := p.widgets[progressWidget]
	for i := start; i <= stop; i++ {
		ph := p.phases[i]

		var status string
		switch ph.status {
		case phaseComplete:
			status = greenCircle
		case phaseWorking:
			status = yellowCircle
		case phasePending:
			status = redCircle
		case phaseUnknown:
			status = purpleCircle
		}

		op := ""
		if ph.operation != "" {
			op = " [" + ph.operation + "]"
		}
		p.state.Set(pw, i-start, " "+status+" "+ph.name+op)
	}
}

func (p *Progress) TaskWidgetSize() int {
	if p.compact {
		return compactWidgetSize
	}
	return standardWidgetSize
}

func (p *Progress) Log(
	name,
	level,
	message string,
	args ...any,
) {
	argsBlock := ""
	if len(args) > 0 {
		argStr := &strings.Builder{}
		key := ""
		for i, v := range args {
			if i%2 == 0 {
				key = fmt.Sprintf("%v", v)
			} else {
				if i >= 2 {
					_, _ = fmt.Fprint(argStr, " ")
				}
				_, _ = fmt.Fprintf(argStr, "%s=%v", key, v)
			}
		}

		if len(args)%2 == 1 {
			if len(args) > 1 {
				_, _ = fmt.Fprint(argStr, " ")
			}
			_, _ = fmt.Fprintf(argStr, "_=%s", key)
		}

		argsBlock = fmt.Sprintf(" [%s]", argStr.String())
	}

	var task, op string
	if p.state != nil {
		for i := 0; i < len(args); i += 2 {
			if i+1 >= len(args) {
				break
			}

			switch args[i] {
			case "@task":
				task = fmt.Sprintf("%v", args[i+1])
			case "@operation":
				op = fmt.Sprintf("%v", args[i+1])
			}

			if op != "" && task != "" {
				break
			}
		}
	}

	line := fmt.Sprintf("%12s / %-6s : %s%s", name, level, message, argsBlock)

	if p.state == nil {
		p.term.Println(line)
		return
	}

	if task != "" {
		p.taskLog(task, op, fmt.Sprintf("%s %s", level, message))
	}

	p.state.Log(line)
}

func (p *Progress) taskLog(taskName string, op string, line string) {
	title := p.title[taskName]
	w := p.widgets[taskName]
	if p.compact {
		p.state.Set(w, 0, title+": "+line)
		p.state.LogWidget(-1, line)
	} else {
		p.state.LogWidget(w, line)
	}

	p.phases[p.currentPhase].operation = op

	p.UpdateProgress()
}

func (p *Progress) Close() {
	if p.state != nil {
		p.state.Close()
		p.state = nil
	}
}
