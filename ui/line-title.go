package ui

// titledLine is a widgetLine that shows a title in a heading line.
type titledLine struct {
	width  int
	title  string
	action string
}

// newTitledLine creates a new titledLine with the given title and screen width.
func newTitledLine(title string, width int) *titledLine {
	return &titledLine{
		title: title,
		width: width,
	}
}

// SetActionKey sets the action key for the line.
func (l *titledLine) SetActionKey(action string) {
	l.action = action
}

// ActionKey returns the action key for the line.
func (l *titledLine) ActionKey() string {
	return l.action
}

// Title returns the title of the line.
func (l *titledLine) Title() string {
	return l.title
}

// SetTitle sets the title of the line.
func (l *titledLine) SetTitle(title string) {
	l.title = title
}

// String returns the string representation of the line, showing the title.
func (l *titledLine) String() string {
	return makeHeader(l.title, l.width)
}
