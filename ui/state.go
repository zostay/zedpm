package ui

import (
	"strings"
)

type Widget []string
type State struct {
	term    *Terminal
	widgets []Widget
}

func NewWidget(size int) Widget {
	return make(Widget, 0, size)
}

func NewState(term *Terminal, capacity int) *State {
	w := make([]Widget, 0, capacity)
	s := &State{
		term:    term,
		widgets: w,
	}
	s.writeBoundary()
	return s
}

func (s *State) writeBoundary() {
	s.term.WriteLine("---- ⏷ Active Tasks ⏷ ---- ⏶ Logs ⏶ ----")
}

func (s *State) resizeAndDraw(oldToBoundary int) {
	newToBoundary := s.MovementsToBoundary()
	if newToBoundary > oldToBoundary {
		s.term.AddLines(newToBoundary - oldToBoundary)
	} else if newToBoundary < oldToBoundary {
		s.term.MoveUp(oldToBoundary - newToBoundary)
		s.term.ClearLines(oldToBoundary - newToBoundary)
		s.term.MoveUp(oldToBoundary - newToBoundary)
	}

	s.term.MoveUp(newToBoundary)

	s.draw("")
}

func (s *State) draw(logLine string) {
	if logLine != "" {
		s.term.WriteLine(logLine)
	}

	s.writeBoundary()

	for _, p := range s.widgets {
		for _, l := range p {
			s.term.WriteLine(l)
		}

		for i := 0; i < cap(p)-len(p); i++ {
			s.term.WriteLine("")
		}
	}
}

func (s *State) redraw(line string) {
	toBoundary := s.MovementsToBoundary()
	s.term.MoveUp(toBoundary)
	s.draw(line)
}

func (s *State) AddWidget(widget Widget) {
	oldToBoundary := s.MovementsToBoundary()
	s.widgets = append(s.widgets, widget)
	s.resizeAndDraw(oldToBoundary)
}

func (s *State) DeleteWidget(n int) {
	oldToBoundary := s.MovementsToBoundary()
	copy(s.widgets[n:], s.widgets[n+1:])
	s.widgets = s.widgets[:len(s.widgets)-1]
	s.resizeAndDraw(oldToBoundary)
}

func (s *State) Title(n int) string {
	return s.widgets[n][0]
}

func (s *State) SetTitle(n int, line string) {
	if len(s.widgets[n]) < 1 {
		s.widgets[n] = append(s.widgets[n], line)
		return
	}
	s.widgets[n][0] = line
}

func (s *State) MovementsToBoundary() int {
	l := 1
	for _, p := range s.widgets {
		l += cap(p)
	}
	return l
}

func (s *State) Redraw() {
	s.redraw("")
}

func (s *State) Log(n int, line string) {
	switch {
	case len(s.widgets[n]) < cap(s.widgets[n]):
		s.widgets[n] = append(s.widgets[n], line)

	case len(s.widgets[n]) == 1:
		s.widgets[n][0] = line

	default:
		for i := 1; i < len(s.widgets[n])-1; i++ {
			s.widgets[n][i] = s.widgets[n][i+1]
		}
		s.widgets[n][len(s.widgets[n])-1] = line
	}

	s.redraw(line)
}

func (s *State) Set(n, m int, line string) {
	if len(s.widgets[n]) < m+1 {
		s.widgets[n] = append(s.widgets[n], strings.Repeat("", m+1-len(s.widgets[n])))
	}
	s.widgets[n][m] = line
	s.redraw("")
}

func (s *State) Close() {
	for len(s.widgets) > 0 {
		s.DeleteWidget(0)
	}
	s.term.MoveUp(1)
	s.term.ClearLine()
}
