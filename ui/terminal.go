package ui

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"golang.org/x/sys/unix"
)

// TruncateString will take a string and shorten it so that it's visible length
// matches the given value. It will also make the last visible part of the
// string the value given as the ellipsis (which may be empty).
func TruncateString(
	line string,
	width int,
	ellipsis string,
) string {
	ellipsisWidth := StringWidth(ellipsis)
	truncateTo := width - ellipsisWidth

	truncLine := &strings.Builder{}
	cs := NewCharacters(line)
	soFar := 0
	for cs.Next() {
		thisWidth := cs.Width()
		if thisWidth+soFar > truncateTo {
			truncLine.WriteString(ellipsis)
			return truncLine.String()
		}

		truncLine.WriteString(cs.String())
		soFar += thisWidth
	}

	return truncLine.String()
}

// Terminal is a very simple tool for writing to and manipulating the terminal.
type Terminal struct {
	ty       *os.File
	istty    bool
	h, w     uint16
	ellipsis string
	lock     sync.RWMutex
}

// NewTerminal creates a Terminal object for the given terminal and returns a
// pointer to that object.
func NewTerminal(tty *os.File) *Terminal {
	t := &Terminal{ty: tty}
	t.detectTTY()
	if t.istty {
		go t.sigwinch()
	}
	return t
}

// SetEllipsis sets a string to add at the end of a truncated widgetLogLine. The default
// is to have no such string and just terminate the widgetLogLine at the end of the
// on-screen widgetLogLine.
func (t *Terminal) SetEllipsis(ellipsis string) {
	t.ellipsis = ellipsis
}

func (t *Terminal) detectTTY() {
	t.lock.Lock()
	defer t.lock.Unlock()
	winsize, err := unix.IoctlGetWinsize(int(t.ty.Fd()), unix.TIOCGWINSZ)
	if err == nil {
		t.h = winsize.Row
		t.w = winsize.Col
		t.istty = true
	}
}

func (t *Terminal) sigwinch() {
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGWINCH)
	for range signals {
		t.detectTTY()
	}
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

// Println will write a single widgetLogLine to the screen. This will blank any existing
// data on the current widgetLogLine before writing and will move the cursor down one
// widgetLogLine afterward.
func (t *Terminal) Println(line string) {
	t.ClearLine()
	_, _ = fmt.Fprintln(t.ty, line)
}

// WriteLine will write a single widgetLogLine to the screen. If the given widgetLogLine contains
// a newline, it will be replaced by U+2424 (SYMBOL FOR NEWLINE) on screen. This
// will also truncate the widgetLogLine so it is not longer than the terminal width. This
// will blank any existing data on the current widgetLogLine before writing and will move
// the cursor down one widgetLogLine afterward.
func (t *Terminal) WriteLine(line string) {
	t.ClearLine()
	line = strings.ReplaceAll(line, "\n", "\u2424")
	line = TruncateString(line, t.Width(), t.ellipsis)
	_, _ = fmt.Fprintln(t.ty, line)
}

// IsTTY returns true if the screen is a TTY. If IsTTY is false, the only
// feature that will do anything is WriteLine. All other methods become no-ops.
// In fact, even WriteLine does less: it will no longer clear lines.
func (t *Terminal) IsTTY() bool {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.istty
}

// Height returns the number of rows on the current TTY. Returns 0 if this is
// not a TTY.
func (t *Terminal) Height() int {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return int(t.h)
}

// Width returns the number of coumns in the current TTY. Returns 0 if this is
// not a TTY.
func (t *Terminal) Width() int {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return int(t.w)
}
