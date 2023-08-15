package ui

import "github.com/zostay/zedpm/pkg/log"

// widget is an internal object used by State to track the state of individual
// widgets.
type widget struct {
	lines []widgetLine
	title string
	width int
}

// newWidget creates a new widget object with the given height.
func newWidget(size int, width int) *widget {
	return &widget{
		lines: make([]widgetLine, 0, size),
		title: "",
		width: width,
	}
}

// Size returns the height of the widget.
func (w *widget) Size() int {
	return cap(w.lines)
}

func (w *widget) makeLogLine(msg string) widgetLine {
	var line widgetLine = newLogLine(msg)
	if cap(w.lines) == 1 && w.title != "" {
		line = newPrefixedLine(w.title, line)
	}
	return line
}

// Log will add the given data to the end of the widget, shifting all existing
// lines upward, if necessary. It will preserve the title widgetLine if SetTitle has
// been called.
func (w *widget) Log(msg string) {
	line := w.makeLogLine(msg)
	switch {
	case len(w.lines) < cap(w.lines):
		w.lines = append(w.lines, line)

	case len(w.lines) == 1:
		w.lines[0] = line

	default:
		start := 0
		if cap(w.lines) > 1 && w.title != "" {
			start = 1
		}
		copy(w.lines[start:], w.lines[start+1:])
		w.lines[len(w.lines)-1] = line
	}
}

// padLines inserts blanks before line n.
func (w *widget) padLines(n int) {
	for len(w.lines) < n+1 {
		w.lines = append(w.lines, w.makeLogLine(""))
	}
}

// Set will replace the widgetLine on the given row. This will grow the widget
// if the new position requires a larger widget to show.
func (w *widget) Set(n int, msg string) {
	w.padLines(n)
	w.lines[n] = w.makeLogLine(msg)
}

// SetStatus will replace the widgetLine on the given row with a status
// widgetLine.
func (w *widget) SetStatus(n int, name string, icon statusIcon, op string) {
	w.padLines(n)
	sl := newStatusLine(name)
	sl.SetIcon(icon)
	sl.SetOperation(op)
	w.lines[n] = sl
}

// SetTitle will mark this widget as having a title widgetLine and set the first row
// to the given widgetLine.
func (w *widget) SetTitle(msg string) {
	w.title = msg
	var line widgetLine
	if cap(w.lines) == 1 {
		line = newPrefixedLine(w.title, newLogLine(headerLineEllipsis))
	} else {
		line = newTitledLine(w.title, w.width)
	}

	if len(w.lines) == 0 {
		w.lines = append(w.lines, line)
	} else {
		w.lines[0] = line
	}
}

// Title will return the title string, if one has been set. If no title has been
// set it will return an empty string.
func (w *widget) Title() string {
	return w.title
}

// SetActionKey assigns the action key to use when referring to a given widget
// line or the last one if -1 is passed as the index.
func (w *widget) SetActionKey(n int, key string) {
	if n == -1 {
		n = len(w.lines) - 1
	}
	w.lines[n].SetActionKey(key)
}

// getActionIndex returns the line with the assigned action key. Returns -1 if
// no line has been assigned the given key.
func (w *widget) getActionIndex(key string) int {
	for i, line := range w.lines {
		if line.ActionKey() == key {
			return i
		}
	}
	return -1
}

// addNFlags will add the given flags to the widget on line n.
func (w *widget) addNFlags(n int, flags ...string) {
	if fl, isFlagged := w.lines[n].(widgetFlagLine); isFlagged {
		fl.AddFlags(flags...)
	}
}

// AddFlags will add the given flags to the widget line with the given action
// key.
func (w *widget) AddFlags(action string, flags ...string) {
	n := w.getActionIndex(action)
	if n < 0 {
		return
	}
	w.addNFlags(n, flags...)
}

// SetOutcome will set the outcome for the widget line with the given action
// key.
func (w *widget) SetOutcome(action string, outcome log.Outcome) {
	n := w.getActionIndex(action)
	if n < 0 {
		return
	}
	if ol, hasOutcome := w.lines[n].(widgetOutcomeLine); hasOutcome {
		ol.SetOutcome(outcome)
	}
}

// IncTick will increment the tick count for the widget line with the given
// action key.
func (w *widget) IncTick(action string) {
	n := w.getActionIndex(action)
	if n < 0 {
		return
	}
	if tl, isTicked := w.lines[n].(widgetTickLine); isTicked {
		tl.IncTick()
	}
}

// Draw will draw the widget to the screen, assuming the cursor is already in
// the correct place for doing so.
func (w *widget) Draw(term *Terminal) {
	for _, l := range w.lines {
		term.WriteLine(l.String())
	}

	for i := 0; i < cap(w.lines)-len(w.lines); i++ {
		term.WriteLine("")
	}
}
