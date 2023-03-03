package main

import (
	"github.com/zostay/zedpm/plugin-git/gitImpl"
	"github.com/zostay/zedpm/plugin/metal"
)

func main() {
	metal.RunPlugin(&gitImpl.Plugin{})
}
