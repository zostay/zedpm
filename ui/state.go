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

// WidgetID allows the caller to address individual widgets.
type WidgetID int

// NoWidget is the WidgetID to use when an operation (see State.Log) needs to be
// performed without referring to a widget.
const NoWidget WidgetID = -1

// State tracks the low-level state of output and writes that state to the
// terminal.
type State struct {
	term    *Terminal
	serial  WidgetID
	widgets map[WidgetID]*widget
	order   []WidgetID
	width   int
}

// NewState creates a new State object attached to the given terminal with room
// for capacity widgets, initially. Capacity will expand as needed.
//
// State manages the terminal in two basic sections separated by a boundary
// line. Above the boundary line is a log. Below is zero or more widgets which
// are drawn and redrawn with every change.
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

// makeHeader puts a piece of text inside a line to make it stand out on screen.
func (s *State) makeHeader(line string) string {
	if utf8.RuneCountInString(line) > s.width-10 {
		line = string([]rune(line)[:s.width-11]) + headerLineEllipsis
	}
	trailerLen := s.width - 5 - utf8.RuneCountInString(line)
	return strings.Repeat(headerLineLine, 4) + " " +
		line + " " + strings.Repeat(headerLineLine, trailerLen)
}

// writeBoundary will write the boundary line inside a header to the screen at
// the current cursor position.
func (s *State) writeBoundary() {
	s.term.WriteLine(s.makeHeader(boundary))
}

// resizeAndDraw will handle shrinkage or growth of the managed screen
// real-estate following changes to the widgets. The oldToBoundary value is the
// value of MovementsToBoundary for the state prior to resizing.
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

// draw will write the given log line, if any is given, and then will redraw the
// entire state to the screen. This call assumes that the cursor has already
// been moved to the location of the boundary line.
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

// redraw moves the cursor from the bottom of the terminal state to the boundary
// and then calls draw.
func (s *State) redraw(line string) {
	toBoundary := s.MovementsToBoundary()
	s.term.MoveUp(toBoundary)
	s.draw(line)
}

// AddWidget creates a new widget with the given row height and then resizes
// and redraws the State on the terminal. The returned WidgetID should be used
// to make any changes to this widget going forward.
func (s *State) AddWidget(n int) WidgetID {
	oldToBoundary := s.MovementsToBoundary()
	s.serial++
	s.widgets[s.serial] = newWidget(n)
	s.order = append(s.order, s.serial)
	s.resizeAndDraw(oldToBoundary)
	return s.serial
}

// DeleteWidget removes the widget with the given WidgetID. Then it resizes and
// redraws the State on the terminal.
func (s *State) DeleteWidget(id WidgetID) {
	oldToBoundary := s.MovementsToBoundary()
	for i, o := range s.order {
		if o == id {
			copy(s.order[i:], s.order[i+1:])
			s.order = s.order[:len(s.order)-1]
		}
	}
	delete(s.widgets, id)
	s.resizeAndDraw(oldToBoundary)
}

// Title will return the title string for the widget, if one has been set by
// SetTitle.
func (s *State) Title(id WidgetID) string {
	return s.widgets[id].Title()
}

// SetTitle sets the title line to use on a widget.
func (s *State) SetTitle(id WidgetID, line string) {
	line = s.makeHeader(line)
	s.widgets[id].SetTitle(line)
}

// MovementsToBoundary states how many cursor movements are required to move the
// cursor from the bottom of the on-screen State to the boundary line.
func (s *State) MovementsToBoundary() int {
	l := 1
	for _, w := range s.widgets {
		l += w.Size()
	}
	return l
}

// Redraw triggers a redraw. This should never be necessary to call directly.
func (s *State) Redraw() {
	s.redraw("")
}

// Log will write a log to the widget identified by the given WidgetID. If no
// such widget exists or the NoWidget constant is passed, the log is recorded
// above the boundary without writing anything to the state below.
func (s *State) Log(id WidgetID, line string) {
	if widget, widgetExists := s.widgets[id]; widgetExists {
		widget.Log(line)
	}
	s.redraw(line)
}

// Set will set the value of a specific widget line to the given value. If the
// widget given by the WidgetID does not exist, this method does nothing.
func (s *State) Set(id WidgetID, n int, line string) {
	if widget, widgetExists := s.widgets[id]; widgetExists {
		widget.Set(n, line)
		s.redraw("")
	}
}

// Close should always be called before program termination or when the State
// object is about to give up control of the terminal. This will erase all the
// widgets and move the cursor to where the boundary line was, just below the
// end of the log.
func (s *State) Close() {
	for len(s.widgets) > 0 {
		s.DeleteWidget(s.order[0])
	}
	s.term.MoveUp(1)
	s.term.ClearLine()
}
