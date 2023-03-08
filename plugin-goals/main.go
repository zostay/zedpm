// Package main is the command that runs the zedpm-plugin-goals plugin.
package main

import (
	"github.com/zostay/zedpm/plugin-goals/goalsImpl"
	"github.com/zostay/zedpm/plugin/metal"
)

func main() {
	metal.RunPlugin(&goalsImpl.Plugin{})
}
