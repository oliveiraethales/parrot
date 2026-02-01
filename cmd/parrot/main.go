package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/oliveiraethales/parrot"
)

func main() {
	singlechecker.Main(parrot.Analyzer)
}
