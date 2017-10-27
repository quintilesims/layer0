package main

import (
	"flag"
)

var DryRun = flag.Bool("dry", false, "")

func init() {
	flag.Parse()
}
