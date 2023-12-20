package main

import (
	"github.com/zostay/zedpm/plugin/metal"
	"github.com/zostay/zedpm/zedpm-plugin-go/goImpl"
)

func main() {
	metal.RunPlugin(&goImpl.Plugin{})
}
