package ui

import (
	"fmt"
	"os"
	"strings"
)

// Terminal is a very simple tool for writing to and manipulating the terminal.
type Terminal struct {
	ty *os.File
}

// NewTerminal creates a Terminal object for the given tty and returns a pointer
// to that object.
func NewTerminal(tty *os.File) *Terminal {
	return &Terminal{tty}
}

// MoveUp will move the cursor up n rows on-screen.
func (t *Terminal) MoveUp(n int) {
	_, _ = fmt.Fprintf(t.ty, "\x1b[%dA", n)
}

// ClearLine will clear the row the cursor is currently on.
func (t *Terminal) ClearLine() {
	_, _ = fmt.Fprint(t.ty, "\x1b[2K")
}

// ClearLines will clear n lines below the row the cursor is currently on.
func (t *Terminal) ClearLines(n int) {
	_, _ = fmt.Fprint(t.ty, strings.Repeat("\x1b[2K\n", n))
}

// AddLines will move the cursor down n lines.
func (t *Terminal) AddLines(n int) {
	_, _ = fmt.Fprint(t.ty, strings.Repeat("\n", n))
}

// WriteLine will write a single line to the screen. If the given line contains
// a newline, it will be replaced by U+2424 (SYMBOL FOR NEWLINE) on screen.
// This will blank any existing data on the current line before writing and will
// move the cursor down one line afterward.
func (t *Terminal) WriteLine(line string) {
	t.ClearLine()
	strings.ReplaceAll(line, "\n", "\u2424")
	_, _ = fmt.Fprintln(t.ty, line)
}
