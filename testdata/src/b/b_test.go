package b

import (
	"os"
	"testing"
)

func TestFoo(t *testing.T) {
	defer testChdir(t, "/this/directory/does/not/exist")()

	// ...
}

// https://speakerdeck.com/mitchellh/advanced-testing-with-go?slide=30
func testChdir(t *testing.T, dir string) func() { // want `unmarked test helper`
	old, err := os.Getwd()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := os.Chdir(dir); err != nil {
		t.Fatalf("err: %s", err)
	}

	return func() { os.Chdir(old) }
}
