package main

import (
	"github.com/zostay/zedpm/plugin/metal"
	"github.com/zostay/zedpm/zedpm-plugin-golangci/golangciImpl"
)

func main() {
	metal.RunPlugin(&golangciImpl.Plugin{})
}
