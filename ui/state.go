package ui

import (
	"strings"
	"unicode/utf8"
)

const (
	boundary           = "⏷ Active Tasks ⏷ ──── ⏶ Logs ⏶"
	headerLineLine     = `-`
	headerLineEllipsis = `…`
	headerWidth        = 70
)

type WidgetID int
type Widget []string
type State struct {
	term    *Terminal
	serial  WidgetID
	widgets map[WidgetID]Widget
	order   []WidgetID
}

func NewWidget(size int) Widget {
	return make(Widget, 0, size)
}

func NewState(term *Terminal, capacity int) *State {
	s := &State{
		term:    term,
		widgets: make(map[WidgetID]Widget, capacity),
		order:   make([]WidgetID, 0, capacity),
	}
	s.writeBoundary()
	return s
}

func makeHeader(line string) string {
	if utf8.RuneCountInString(line) > headerWidth-10 {
		line = string([]rune(line)[:headerWidth-11]) + headerLineEllipsis
	}
	trailerLen := headerWidth - 5 - utf8.RuneCountInString(line)
	return strings.Repeat(headerLineLine, 4) + " " +
		line + " " + strings.Repeat(headerLineLine, trailerLen)
}

func (s *State) writeBoundary() {
	s.term.WriteLine(makeHeader(boundary))
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

	for _, o := range s.order {
		p := s.widgets[o]
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

func (s *State) AddWidget(widget Widget) WidgetID {
	oldToBoundary := s.MovementsToBoundary()
	s.serial++
	s.widgets[s.serial] = widget
	s.order = append(s.order, s.serial)
	s.resizeAndDraw(oldToBoundary)
	return s.serial
}

func (s *State) DeleteWidget(n WidgetID) {
	oldToBoundary := s.MovementsToBoundary()
	for i, o := range s.order {
		if o == n {
			copy(s.order[i:], s.order[i+1:])
			s.order = s.order[:len(s.order)-1]
		}
	}
	delete(s.widgets, n)
	s.resizeAndDraw(oldToBoundary)
}

func (s *State) Title(n WidgetID) string {
	return s.widgets[n][0]
}

func (s *State) SetTitle(n WidgetID, line string) {
	line = makeHeader(line)
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

func (s *State) Log(n WidgetID, line string) {
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

func (s *State) Set(n WidgetID, m int, line string) {
	if len(s.widgets[n]) < m+1 {
		s.widgets[n] = append(s.widgets[n], strings.Repeat("", m+1-len(s.widgets[n])))
	}
	s.widgets[n][m] = line
	s.redraw("")
}

func (s *State) Close() {
	for len(s.widgets) > 0 {
		s.DeleteWidget(s.order[0])
	}
	s.term.MoveUp(1)
	s.term.ClearLine()
}
