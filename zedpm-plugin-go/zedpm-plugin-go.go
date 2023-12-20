package main

import (
	"github.com/zostay/zedpm/plugin-go/goImpl"
	"github.com/zostay/zedpm/plugin/metal"
)

func main() {
	metal.RunPlugin(&goImpl.Plugin{})
}
