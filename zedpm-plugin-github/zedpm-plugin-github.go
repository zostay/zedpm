// Package main is the command that runs the zedpm-plugin-github plugin.
package main

import (
	"github.com/zostay/zedpm/plugin/metal"
	"github.com/zostay/zedpm/zedpm-plugin-github/githubImpl"
)

func main() {
	metal.RunPlugin(&githubImpl.Plugin{})
}
