package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/ichiban/unmarkedhelper"
)

func main() {
	singlechecker.Main(unmarkedhelper.Analyzer)
}
