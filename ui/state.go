package ui

import "fmt"

type Widget []string
type State []Widget

func NewWidget(size int) Widget {
	return make(Widget, 0, size)
}

func NewState(capacity int) State {
	s := make(State, 0, capacity)
	WriteBoundary()
	return s
}

func resizeAndDraw(prev, s State) {
	var (
		newToBoundary = s.MovementsToBoundary()
		oldToBoundary = prev.MovementsToBoundary()
	)

	if newToBoundary > oldToBoundary {
		AddLines(newToBoundary - oldToBoundary)
	} else if newToBoundary < oldToBoundary {
		MoveUp(oldToBoundary - newToBoundary)
		ClearLines(oldToBoundary - newToBoundary)
		MoveUp(oldToBoundary - newToBoundary)
	}

	MoveUp(newToBoundary)

	draw(s, "")
}

func draw(s State, logLine string) {
	if logLine != "" {
		ClearLine()
		fmt.Println(logLine)
	}

	WriteBoundary()

	for _, p := range s {
		for _, l := range p {
			ClearLine()
			fmt.Println(l)
		}

		for i := 0; i < cap(p)-len(p); i++ {
			ClearLine()
			fmt.Println("")
		}
	}
}

func redraw(s State, line string) {
	toBoundary := s.MovementsToBoundary()
	MoveUp(toBoundary)
	draw(s, line)
}

func (s State) AddWidget(widget Widget) State {
	newState := append(s, widget)
	resizeAndDraw(s, newState)
	return newState
}

func (s State) DeleteWidget(n int) State {
	copy(s[n:], s[n+1:])
	newState := s[:len(s)-1]
	resizeAndDraw(s, newState)
	return newState
}

func (s State) Title(n int) string {
	return s[n][0]
}

func (s State) SetTitle(n int, line string) {
	s[n][0] = line
}

func (s State) MovementsToBoundary() int {
	l := 1
	for _, p := range s {
		l += cap(p)
	}
	return l
}

func (s State) Redraw() {
	redraw(s, "")
}

func (s State) Log(n int, line string) {
	switch {
	case len(s[n]) < cap(s[n]):
		s[n] = append(s[n], line)

	case len(s[n]) == 1:
		s[n][0] = line

	default:
		for i := 1; i < len(s[n])-1; i++ {
			s[n][i] = s[n][i+1]
		}
		s[n][len(s[n])-1] = line
	}

	redraw(s, line)
}

func (s State) Close() {
	for len(s) > 0 {
		s = s.DeleteWidget(0)
	}
	MoveUp(1)
	ClearLine()
}
