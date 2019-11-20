package unmarkedhelper

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestFileSystem(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer, "a")
}
