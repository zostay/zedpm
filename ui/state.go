package ui

import (
	"strings"
	"unicode/utf8"
)

const (
	boundary           = "⏷ Progress ⏷ ──── ⏶ Logs ⏶"
	headerLineLine     = `-`
	headerLineEllipsis = `…`
	defaultHeaderWidth = 78
)

type WidgetID int

type State struct {
	term    *Terminal
	serial  WidgetID
	widgets map[WidgetID]*widget
	order   []WidgetID
	width   int
}

func NewState(term *Terminal, capacity int) *State {
	s := &State{
		term:    term,
		widgets: make(map[WidgetID]*widget, capacity),
		order:   make([]WidgetID, 0, capacity),
		width:   defaultHeaderWidth,
	}
	s.writeBoundary()
	return s
}

func (s *State) makeHeader(line string) string {
	if utf8.RuneCountInString(line) > s.width-10 {
		line = string([]rune(line)[:s.width-11]) + headerLineEllipsis
	}
	trailerLen := s.width - 5 - utf8.RuneCountInString(line)
	return strings.Repeat(headerLineLine, 4) + " " +
		line + " " + strings.Repeat(headerLineLine, trailerLen)
}

func (s *State) writeBoundary() {
	s.term.WriteLine(s.makeHeader(boundary))
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
		w := s.widgets[o]
		w.Draw(s.term)
	}
}

func (s *State) redraw(line string) {
	toBoundary := s.MovementsToBoundary()
	s.term.MoveUp(toBoundary)
	s.draw(line)
}

func (s *State) AddWidget(n int) WidgetID {
	oldToBoundary := s.MovementsToBoundary()
	s.serial++
	s.widgets[s.serial] = newWidget(n)
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
	return s.widgets[n].Title()
}

func (s *State) SetTitle(n WidgetID, line string) {
	line = s.makeHeader(line)
	s.widgets[n].SetTitle(line)
}

func (s *State) MovementsToBoundary() int {
	l := 1
	for _, w := range s.widgets {
		l += w.Size()
	}
	return l
}

func (s *State) Redraw() {
	s.redraw("")
}

func (s *State) Log(n WidgetID, line string) {
	s.widgets[n].Log(line)
	s.redraw(line)
}

func (s *State) Set(n WidgetID, m int, line string) {
	s.widgets[n].Set(m, line)
	s.redraw("")
}

func (s *State) Close() {
	for len(s.widgets) > 0 {
		s.DeleteWidget(s.order[0])
	}
	s.term.MoveUp(1)
	s.term.ClearLine()
}
