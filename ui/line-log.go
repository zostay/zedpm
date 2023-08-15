package ui

import (
	"strings"
	"sync"
	"time"

	"github.com/zostay/go-std/set"

	"github.com/zostay/zedpm/pkg/log"
)

type widgetLogLine struct {
	action  string
	message string
	flags   set.Set[string]
	tick    int
	outcome string

	lock    sync.Mutex
	spinner chan struct{}
}

func newLogLine(msg string) *widgetLogLine {
	return &widgetLogLine{
		message: msg,
		flags:   set.New[string](),
	}
}

// SetActionKey assigns the action key to use when referring to this line.
func (l *widgetLogLine) SetActionKey(action string) {
	l.action = action
}

// ActionKey returns the action key assigned to this line.
func (l *widgetLogLine) ActionKey() string {
	return l.action
}

// startSpinner start a goroutine that automatically ticks the spinner.
func (l *widgetLogLine) startSpinner() {
	l.spinner = make(chan struct{})

	go func() {
		for {
			<-time.After(1 * time.Second)
			_, ok := <-l.spinner
			if !ok {
				return
			}
			l.IncTick()
		}
	}()
}

// stopSpinner stops the goroutine that automatically ticks the spinner.
func (l *widgetLogLine) stopSpinner() {
	close(l.spinner)
}

// AddFlags will add a flag to the widgetLogLine.
func (l *widgetLogLine) AddFlags(flags ...string) {
	for _, flag := range flags {
		l.flags.Insert(flag)
		if flag == "autospin" {
			l.startSpinner()
		}
	}
}

// RemoveFlags will remove a flag from the widgetLogLine.
func (l *widgetLogLine) RemoveFlags(flags ...string) {
	for _, flag := range flags {
		l.flags.Delete(flag)
		if flag == "autospin" {
			l.stopSpinner()
		}
	}
}

// SetOutcome will set the final outcome of this widgetLogLine. Setting this will
// automatically clear the "spin" or "autospin" flag as well if the given
// outcome is considered a final state.
func (l *widgetLogLine) SetOutcome(outcome string) {
	l.outcome = outcome
	if log.Outcome(outcome).IsFinal() {
		l.flags.Delete("spin")
		l.flags.Delete("autospin")
	}
}

// Outcome will return the final outcome set for this widgetLogLine.
func (l *widgetLogLine) Outcome() string {
	return l.outcome
}

// IncTick increments this tick count for this widgetLogLine.
func (l *widgetLogLine) IncTick() {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.tick++
}

// Tick returns the tick count for this widgetLogLine.
func (l *widgetLogLine) Tick() int {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.tick
}

// var spinner = []string{"-", "\\", "|", "/"}
var spinner = []string{" Λ̊ ", " Λ∘", " ̥Λ ", "∘Λ "}

// spin displays the spinner for the given tick.
func spin(tick int) string {
	return spinner[tick%4]
}

var (
	outcomeSymbols = map[string]string{
		"okay":  "👌",
		"pass":  "✅",
		"fail":  "❌",
		"skip":  "⏭",
		"retry": "🔁",
	}
)

// String will render the line, including the outcome, spin tick, etc.
func (l *widgetLogLine) String() string {
	out := &strings.Builder{}

	if l.flags.Contains("spin") || l.flags.Contains("autospin") {
		out.WriteByte('[')
		out.WriteString(spin(l.tick))
		out.WriteByte(']')
	} else {
		out.WriteString("     ")
	}

	out.WriteByte(' ')

	if l.outcome != "" {
		if outcomeSymbols[l.outcome] != "" {
			out.WriteByte('[')
			out.WriteString(outcomeSymbols[l.outcome])
			out.WriteByte(']')
			out.WriteByte(' ')
		} else {
			out.WriteString(l.outcome[:4])
		}
	} else {
		out.WriteString("    ")
	}
	out.WriteByte(' ')

	out.WriteString(l.message)

	return out.String()
}
