// Package main is the program for running the zedpm-plugin-changelog plugin.
package main

import (
	"github.com/zostay/zedpm/plugin/metal"
	"github.com/zostay/zedpm/zedpm-plugin-git/gitImpl"
)

func main() {
	metal.RunPlugin(&gitImpl.Plugin{})
}
