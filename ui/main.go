package main

import (
	"fmt"
	"time"

	"github.com/gosuri/uilive"
)

// uilive works by going up N lines based on how many lines there are to be
// updated and then overwrites them. That's all it does.
//
// Literally, all it does is clear one line and move the cursor up and repeats
// that for the number of lines setup. It would be easier to implement this
// directly than to use uilive to do what I want.

func main() {
	tail := uilive.New()
	status := tail.Newline()
	for j := 0; j < 3; j++ {
		for i := 0; i < 10; i++ {
			tail.Start()
			logLine := fmt.Sprintf("testing %d", i)
			fmt.Fprintln(tail, logLine)
			fmt.Fprintf(status, "Status %d\n", j)
			time.Sleep(1 * time.Second)
			tail.Stop()
			for k := 0; k < 5; k++ {
				fmt.Println(logLine)
			}
		}
	}
}
