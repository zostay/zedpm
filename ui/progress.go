package ui

import (
	"os"

	"github.com/zostay/go-std/generic"
)

const defaultWidgetCount = 20

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

const progressWidget = 0

type phase struct {
	status    phaseStatus
	name      string
	operation string
}

type Progress struct {
	state        *State
	phases       []phase
	widgets      map[int]WidgetID
	currentPhase int
}

func NewProgress(tty *os.File) *Progress {
	term := NewTerminal(tty)
	state := NewState(term, defaultWidgetCount)
	return &Progress{
		state: state,
	}
}

func (p *Progress) SetPhases(phases []string) {
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

func (p *Progress) UpdateProgress() {
	var stop int
	phaseEnd := len(p.phases) - 1
	switch phaseEnd {
	case p.currentPhase, p.currentPhase + 1:
		stop = phaseEnd
	default:
		stop = p.currentPhase + 2
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
		p.state.Set(pw, i-start, status+" "+ph.name+op)
	}
}
