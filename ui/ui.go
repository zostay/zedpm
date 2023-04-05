package ui

import (
	"fmt"
	"os"
	"strings"
)

type Terminal struct {
	ty *os.File
}

func NewTerminal(tty *os.File) *Terminal {
	return &Terminal{tty}
}

func (t *Terminal) MoveUp(n int) {
	_, _ = fmt.Fprintf(t.ty, "\x1b[%dA", n)
}

func (t *Terminal) ClearLine() {
	_, _ = fmt.Fprint(t.ty, "\x1b[2K")
}

func (t *Terminal) ClearLines(n int) {
	_, _ = fmt.Fprint(t.ty, strings.Repeat("\x1b[2K\n", n))
}

func (t *Terminal) AddLines(n int) {
	_, _ = fmt.Fprint(t.ty, strings.Repeat("\n", n))
}

func (t *Terminal) WriteLine(line string) {
	t.ClearLine()
	strings.ReplaceAll(line, "\n", "\u2424")
	_, _ = fmt.Fprintln(t.ty, line)
}
