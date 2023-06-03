package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSmprintf(t *testing.T) {
	str := Smprintf("Test #%[answer]d", "answer", 42)
	assert.Equal(t, "Test #42", str)

	str = Smprintf("Test #%[answer]d", "answer")
	assert.Equal(t, "Test #%!d(<nil>)", str)

	str = Smprintf("Test #%[answer]d", []byte("answer"), 42)
	assert.Equal(t, "Test #%!d(<nil>)", str)

	str = Smprintf("%[a]d + %[b]d = %[c]d", "a", 1, "b", 2, "c", 3)
	assert.Equal(t, "1 + 2 = 3", str)

	str = Smprintf("%[c]d + %[b]d = %[a]d", "a", 1, "b", 2, "c", 3)
	assert.Equal(t, "3 + 2 = 1", str)

	str = Smprintf("%[c]d + %[c]d = %[c]d", "a", 1, "b", 2, "c", 3)
	assert.Equal(t, "3 + 3 = 3", str)

	str = Smprintf("%[c]d + %[b]d = %[c]d", "a", 1, "b", 2, "c", 3)
	assert.Equal(t, "3 + 2 = 3", str)

	str = Smprintf("%[c]d + %[d]d = %[c]d", "a", 1, "b", 2, "c", 3)
	assert.Equal(t, "3 + %!d(<nil>) = 3", str)
}
