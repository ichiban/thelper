package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/ichiban/thelper"
)

func main() {
	singlechecker.Main(thelper.Analyzer)
}
