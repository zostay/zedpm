package main

import (
	"github.com/zostay/zedpm/plugin-github/githubImpl"
	"github.com/zostay/zedpm/plugin/metal"
)

func main() {
	metal.RunPlugin(&githubImpl.Plugin{})
}
