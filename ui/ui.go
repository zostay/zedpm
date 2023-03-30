package ui

import (
	"fmt"
	"strings"
)

func MoveUp(n int) {
	_, _ = fmt.Printf("\x1b[%dA", n)
}

func ClearLine() {
	_, _ = fmt.Print("\x1b[2K")
}

func ClearLines(n int) {
	_, _ = fmt.Print(strings.Repeat("\x1b[2K\n", n))
}

func AddLines(n int) {
	_, _ = fmt.Print(strings.Repeat("\n", n))
}

func WriteBoundary() {
	ClearLine()
	fmt.Println("---- ⏷ Active Tasks ⏷ ---- ⏶ Logs ⏶ ----")
}
