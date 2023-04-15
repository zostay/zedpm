package main

import (
	"os"

	"github.com/zostay/zedpm/cmd"
)

func main() {
	exitStatus := cmd.Execute()
	os.Exit(exitStatus)
}
