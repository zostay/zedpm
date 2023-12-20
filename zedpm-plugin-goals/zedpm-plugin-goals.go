// Package main is the command that runs the zedpm-plugin-goals plugin.
package main

import (
	"github.com/zostay/zedpm/plugin/metal"
	"github.com/zostay/zedpm/zedpm-plugin-goals/goalsImpl"
)

func main() {
	metal.RunPlugin(&goalsImpl.Plugin{})
}
