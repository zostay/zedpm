// Package main is the command that runs the zedpm-plugin-github plugin.
package main

import (
	"github.com/zostay/zedpm/plugin-github/githubImpl"
	"github.com/zostay/zedpm/plugin/metal"
)

func main() {
	metal.RunPlugin(&githubImpl.Plugin{})
}
