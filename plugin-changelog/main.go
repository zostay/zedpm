package main

import (
	"github.com/zostay/zedpm/plugin-changelog/changelogImpl"
	"github.com/zostay/zedpm/plugin/metal"
)

func main() {
	metal.RunPlugin(&changelogImpl.Plugin{})
}
