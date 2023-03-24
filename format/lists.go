package format

import (
	"fmt"
	"strings"
)

// And connects the given values together in an English list using an "and" as a
// conjunction and an Oxford comma.
func And(values ...string) string {
	// TODO This should be replaced with proper localization.
	switch len(values) {
	case 0:
		return ""
	case 1:
		return values[0]
	case 2:
		return fmt.Sprintf("%s and %s", values[0], values[1])
	}

	// TODO this is expensive, though in the cases I use it, maybe it doesn't matter
	return fmt.Sprintf("%s, and %s",
		strings.Join(values[:len(values)-1], ", "),
		values[len(values)-1],
	)
}
