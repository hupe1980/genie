package main

import (
	"github.com/hupe1980/genie/cmd"
)

var (
	version = "dev"
)

func main() {
	cmd.Execute(version)
}
