package ui

import (
	"github.com/rivo/uniseg"
)

const esc byte = 0x1b
const ctrl byte = '['

// Characters is a scanner that returns character groups for a terminal widgetLogLine one
// at a time. Each group represents either some visible element (a grapheme or
// space or 0-width component of a unicode string) or a complete terminal escape
// sequence. Only a subset of CSI escapes are supported (i.e., only including
// those that may be used by the code of zedpm or plugins).
type Characters struct {
	str        string
	isEsc      bool
	seg        string
	start, end int
	width      int
}

// NewCharacters creates a new Characters iterator.
func NewCharacters(str string) *Characters {
	return &Characters{
		str: str,
	}
}

// Bytes returns a byte slice which corresponds to the current grapheme cluster
// or escape sequence.
func (c *Characters) Bytes() []byte {
	return []byte(c.seg)
}

// Positions returns the start and end locations of the current escape seqeuence
// or grapheme cluster in the string.
func (c *Characters) Positions() (int, int) {
	return c.start, c.end
}

// IsEscape returns true if the current cluster of bytes is an escape sequence.
func (c *Characters) IsEscape() bool {
	return c.isEsc
}

// Runes returns a rune slice which corresponds to the current grapheme cluster
// or escape sequence.
func (c *Characters) Runes() []rune {
	return []rune(c.seg)
}

// String returns a string which corresponse to the current grapheme cluster or
// terminal escape sequence.
func (c *Characters) String() string {
	return c.seg
}

// Width returns the visible width of the current grapheme or terminal sequence
// (terminal sequences are always treated as zero-width).
func (c *Characters) Width() int {
	return c.width
}

// Next moves to the next grapheme or terminal escape sequence. This returns
// false if no more graphemes or escape sequence remain.
func (c *Characters) Next() bool {
	if c.end >= len(c.str) {
		return false
	}

	if c.str[c.end] == esc {
		if c.scanEscapeSequence() {
			return true
		}
	}

	return c.scanGrapheme()
}

// scanGrapheme finds the next grapheme cluster from the string and sets up for
// that.
func (c *Characters) scanGrapheme() bool {
	seg, _, width, _ := uniseg.FirstGraphemeClusterInString(c.str[c.end:], -1)
	c.isEsc = false
	c.seg = seg
	c.start, c.end = c.end, c.end+len(seg)
	c.width = width
	return true
}

// scanEscapeSequence attempts to detect a CSI escape sequence and sets the
// state of Characters to ensure that the current segment is that sequence. It
// returns true if an escape sequence can be isolated or false otherwise.
func (c *Characters) scanEscapeSequence() bool {
	i := c.end

	// Look for ESC
	if i >= len(c.str) || c.str[i] != esc {
		return false
	}
	i++

	// Look for the control sequence introduction character
	if i >= len(c.str) || c.str[i] != ctrl {
		return false
	}
	i++

	// We have CSI. Now, look for numbers and semicolons.
	for {
		if i >= len(c.str) {
			return false
		}

		ch := c.str[i]
		if ch != ';' && !(ch >= '0' && ch <= '9') {
			break
		}
		i++
	}

	// If we didn't find at least one digit, that does not appear to be a
	// control sequence
	if i-c.end == 2 {
		return false
	}

	// See if the last character is a CSI termination char and, if so, get set
	// up and return true.
	ch := c.str[i]
	if ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' {
		c.isEsc = true
		c.start, c.end = c.end, i+1
		c.seg = c.str[c.start:c.end]
		c.width = 0
		return true
	}

	// Almost a CSI sequence, but not quite.
	return false
}

// StringWidth will calculate the width of the string, taking control sequences
// into account as well as graphemes and their sizes.
func StringWidth(str string) int {
	cs := NewCharacters(str)
	total := 0
	for cs.Next() {
		total += cs.Width()
	}
	return total
}
