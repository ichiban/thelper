package thelper

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestFileSystem(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer, "a")
	analysistest.Run(t, testdata, Analyzer, "b")
}
