// Package main is the program for running the zedpm-plugin-changelog plugin.
package main

import (
	"github.com/zostay/zedpm/plugin-git/gitImpl"
	"github.com/zostay/zedpm/plugin/metal"
)

func main() {
	metal.RunPlugin(&gitImpl.Plugin{})
}
