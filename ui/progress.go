package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/zostay/go-std/generic"

	"github.com/zostay/zedpm/pkg/log"
)

const (
	defaultWidgetCount = 20
	standardWidgetSize = 4
	compactWidgetSize  = 1
)

const progressWidget = ""

type phaseStatus int

const (
	phaseUnknown  phaseStatus = iota
	phasePending              // red
	phaseWorking              // yellow
	phaseComplete             // green
)

var (
	defaultIcon  = purpleCircle
	phaseIconMap = map[phaseStatus]statusIcon{
		phaseUnknown:  purpleCircle,
		phasePending:  redCircle,
		phaseWorking:  yellowCircle,
		phaseComplete: greenCircle,
	}
)

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
	} else {
		p.compact = (p.state.term.Height()-9)/taskCount < standardWidgetSize
	}

	p.UpdateProgress()
}

func (p *Progress) RegisterTask(name, title string) {
	if p.state == nil {
		panic("cannot call RegisterTask before SetPhases")
	}

	w := p.state.AddWidget(p.TaskWidgetSize())
	p.state.SetTitle(w, title)
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
		var icon statusIcon
		var hasIcon bool
		if icon, hasIcon = phaseIconMap[ph.status]; !hasIcon {
			icon = defaultIcon
		}
		p.state.SetStatus(pw, i-start, ph.name, icon, ph.operation)
	}
}

func (p *Progress) TaskWidgetSize() int {
	if p.compact {
		return compactWidgetSize
	}
	return standardWidgetSize
}

// Log will process a log widgetLogLine, typically by writing it to the screen. If
// widgets are present, it will add the widgetLogLine to the appropriate widget. It will
// also handle a number of special @<name> fields:
//
//	@task        - name of the task to log this with
//	@operation   - the operation the task is performing
//	@action      - an action key to identify some persistent state
//	@actionFlags - flags to modify how the log is displayed (e.g., "spin" and
//	               "autospin")
//	@outcome     - outcome is the final outcome of an action
//	@tick        - tick will cause a widgetLogLine with an associated "spin" flag to move
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

	var (
		task, op, action string
		outcome          log.Outcome
		skip             = false
		tick             = false
		actionFlags      []string
	)
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
			case "@actionFlags":
				var isStringSlice bool
				if actionFlags, isStringSlice = args[i+1].([]string); isStringSlice {
					skip = true
				}
			case "@action":
				if outcome == "" {
					skip = true
				}
				action = fmt.Sprintf("%v", args[i+1])
			case "@outcome":
				outcome = log.Outcome(fmt.Sprintf("%v", args[i+1]))
				skip = false
			case "@tick":
				tick = true
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

	if !skip {
		if task != "" {
			p.taskLog(task, op, fmt.Sprintf("%s %s", level, message), action)
		}

		p.state.Log(line)
	}

	if task != "" {
		if outcome != "" {
			p.taskOutcome(task, action, outcome)
		}
		if actionFlags != nil {
			p.taskAddFlags(task, action, actionFlags)
		}
		if tick {
			p.taskTick(task, action)
		}
	}
}

// taskLog will log the given message widgetLogLine to the widget for taskName and update
// the operation to the given op.
func (p *Progress) taskLog(taskName string, op string, line string, action string) {
	w := p.widgets[taskName]
	title := p.state.Title(w)
	if p.compact {
		p.state.Set(w, 0, title+": "+line)
		p.state.SetActionKey(w, 0, action)
		p.state.LogWidget(-1, line)
	} else {
		p.state.LogWidget(w, line)
		p.state.SetActionKey(w, -1, action)
	}

	p.phases[p.currentPhase].operation = op

	p.UpdateProgress()
}

// taskOutcome will update the outcome of the given action for the given task.
func (p *Progress) taskOutcome(
	taskName string,
	action string,
	outcome log.Outcome,
) {
	w := p.widgets[taskName]
	p.state.SetOutcome(w, action, outcome)
}

// taskAddFlags will add the given flags to the given action for the given task.
func (p *Progress) taskAddFlags(taskName string, action string, flags []string) {
	w := p.widgets[taskName]
	p.state.AddFlags(w, action, flags)
}

// taskTick will increment the tick for the given action for the given task.
func (p *Progress) taskTick(taskName string, action string) {
	w := p.widgets[taskName]
	p.state.IncTick(w, action)
}

// Close should always be called when finished with the progress widgetLogLine. When
// state is used, this will close all widgets and move the cursor to the final
// widgetLogLine.
func (p *Progress) Close() {
	if p.state != nil {
		p.state.Close()
		p.state = nil
	}
}
