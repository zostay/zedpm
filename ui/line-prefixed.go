package ui

import "fmt"

const prefixWidth = 12

type prefixedLine struct {
	title string
	widgetLine
}

func newPrefixedLine(title string, line widgetLine) *prefixedLine {
	return &prefixedLine{
		title:      title,
		widgetLine: line,
	}
}

func (l *prefixedLine) Title() string {
	return l.title
}

func (l *prefixedLine) SetTitle(title string) {
	l.title = title
}

func (l *prefixedLine) String() string {
	f := fmt.Sprintf("%%%ds: %%s", prefixWidth)
	return fmt.Sprintf(f, l.title, l.widgetLine.String())
}
