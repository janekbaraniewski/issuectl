package main

import (
	cli "github.com/janekbaraniewski/issuectl/cmd/issuectl"
)

var BuildVersion string

func main() {
	cli.Execute(BuildVersion)
}
