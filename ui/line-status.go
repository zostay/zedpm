package ui

import "strings"

// statusLine is a widgetLine that shows the status of a task.
type statusLine struct {
	icon      statusIcon
	name      string
	action    string
	operation string
}

// newStatusLine creates a new statusLine with the given name.
func newStatusLine(name string) *statusLine {
	return &statusLine{
		name: name,
		icon: purpleCircle,
	}
}

// SetActionKey sets the action key for the line.
func (l *statusLine) SetActionKey(action string) {
	l.action = action
}

// ActionKey returns the action key for the line.
func (l *statusLine) ActionKey() string {
	return l.action
}

// String returns the string representation of the line, showing the icon, name,
// and operation.
func (l *statusLine) String() string {
	out := &strings.Builder{}
	out.WriteString(string(l.icon))
	out.WriteByte(' ')
	out.WriteString(l.name)

	if l.operation != "" {
		out.WriteString(" [")
		out.WriteString(l.operation)
		out.WriteByte(']')
	}

	return out.String()
}

// Title returns the name of the status line.
func (l *statusLine) Title() string {
	return l.name
}

// SetTitle sets the name of the status line.
func (l *statusLine) SetTitle(name string) {
	l.name = name
}

// Icon returns the icon for the status line.
func (l *statusLine) Icon() statusIcon {
	return l.icon
}

// SetIcon sets the icon for the status line.
func (l *statusLine) SetIcon(icon statusIcon) {
	l.icon = icon
}

// Operation returns the operation for the status line.
func (l *statusLine) Operation() string {
	return l.operation
}

// SetOperation sets the operation for the status line.
func (l *statusLine) SetOperation(operation string) {
	l.operation = operation
}
