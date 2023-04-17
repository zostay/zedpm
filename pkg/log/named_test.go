package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSmprintf(t *testing.T) {
	str := Smprintf("Test #%[answer]d", "answer", 42)
	assert.Equal(t, "Test #42", str)

	str = Smprintf("Test #%[answer]d", "answer")
	assert.Equal(t, "Test #%!d(BADINDEX)", str)

	str = Smprintf("Test #%[answer]d", []byte("answer"), 42)
	assert.Equal(t, "Test #%!d(BADINDEX)", str)

	str = Smprintf("%[a]d + %[b]d = %[c]d", "a", 2, "b", 2, "c", 4)
	assert.Equal(t, "2 + 2 = 4", str)
}
