package main

import (
	"github.com/zostay/zedpm/plugin-goals/goalsImpl"
	"github.com/zostay/zedpm/plugin/metal"
)

func main() {
	metal.RunPlugin(&goalsImpl.Plugin{})
}
