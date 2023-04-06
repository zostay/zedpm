package ui

import "strings"

type widget struct {
	lines []string
	title bool
}

func newWidget(size int) *widget {
	return &widget{
		lines: make([]string, 0, size),
		title: false,
	}
}

func (w *widget) Size() int {
	return cap(w.lines)
}

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

func (w *widget) Set(n int, line string) {
	if len(w.lines) < n+1 {
		w.lines = append(w.lines, strings.Repeat("", n+1-len(w.lines)))
	}
	w.lines[n] = line
}

func (w *widget) SetTitle(line string) {
	w.Set(0, line)
	w.title = true
}

func (w *widget) Title() string {
	if w.title {
		return w.lines[0]
	}
	return ""
}

func (w *widget) Draw(term *Terminal) {
	for _, l := range w.lines {
		term.WriteLine(l)
	}

	for i := 0; i < cap(w.lines)-len(w.lines); i++ {
		term.WriteLine("")
	}
}
