// Package main runs the zedpm-plugin-changelog plugin.
package main

import (
	"github.com/zostay/zedpm/plugin/metal"
	"github.com/zostay/zedpm/zedpm-plugin-changelog/changelogImpl"
)

func main() {
	metal.RunPlugin(&changelogImpl.Plugin{})
}
