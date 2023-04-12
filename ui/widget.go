package ui

import "strings"

// widget is an internal object used by State to track the state of individual
// widgets.
type widget struct {
	lines []string
	title bool
}

// newWidget creates a new widget object with the given height.
func newWidget(size int) *widget {
	return &widget{
		lines: make([]string, 0, size),
		title: false,
	}
}

// Size returns the height of the widget.
func (w *widget) Size() int {
	return cap(w.lines)
}

// Log will add the given data to the end of the widget, shifting all existing
// lines upward, if necessary. It will preserve the title line if SetTitle has
// been called.
func (w *widget) Log(line string) {
	switch {
	case len(w.lines) < cap(w.lines):
		w.lines = append(w.lines, line)

	case len(w.lines) == 1:
		w.lines[0] = line

	default:
		start := 0
		if w.title {
			start = 1
		}
		for i := start; i < len(w.lines)-1; i++ {
			w.lines[i] = w.lines[i+1]
		}
		w.lines[len(w.lines)-1] = line
	}
}

// Set will replace the line on the given row.
func (w *widget) Set(n int, line string) {
	if len(w.lines) < n+1 {
		w.lines = append(w.lines, strings.Repeat("", n+1-len(w.lines)))
	}
	w.lines[n] = line
}

// SetTitle will mark this widget as having a title line and set the first row
// to the given line.
func (w *widget) SetTitle(line string) {
	w.Set(0, line)
	w.title = true
}

// Title will return the title line, if one has been set. If no title has been
// set it will return an empty string.
func (w *widget) Title() string {
	if w.title {
		return w.lines[0]
	}
	return ""
}

// Draw will draw the widget to the screen, assuming the cursor is already in
// the correct place for doing so.
func (w *widget) Draw(term *Terminal) {
	for _, l := range w.lines {
		term.WriteLine(l)
	}

	for i := 0; i < cap(w.lines)-len(w.lines); i++ {
		term.WriteLine("")
	}
}
