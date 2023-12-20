package main

import (
	"github.com/zostay/zedpm/plugin-golangci/golangciImpl"
	"github.com/zostay/zedpm/plugin/metal"
)

func main() {
	metal.RunPlugin(&golangciImpl.Plugin{})
}
