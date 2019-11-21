# thelper

[Go](https://golang.org/) [static analyzer](https://godoc.org/golang.org/x/tools/go/analysis) that reports where you forgot to call `t.Helper()`.

## As a command

First, install `thelper` command via `go get`.

```shellsession
$ go get github.com/ichiban/thelper/cmd/thelper
```

Then, run `thelper` command at your package directory.

```shellsession
$ thelper ./...
/Users/ichiban/src/thelper/testdata/src/a/a_test.go:11:1: unmarked test helper: call t.Helper()
/Users/ichiban/src/thelper/testdata/src/a/a_test.go:25:1: unmarked test helper: call s.Helper()
```

## As an `analysis.Analyzer`

First, install the package via `go get`.

```shellsession
$ go get github.com/ichiban/thelper
```

Then, include `thelper.Analyzer` in your checker.

```go
package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shift"

	"github.com/ichiban/thelper"
)

func main() {
	multichecker.Main(
		// other analyzers of your choice
		nilfunc.Analyzer,
		printf.Analyzer,
		shift.Analyzer,

		thelper.Analyzer,
	)
}
```

## What happens if I don't call `t.Helper()`?

`go test` shows FAILs with an unhelpful line number.

Let's take an example of a test with an unmarked test helper.
`TestFoo(t)` fails because `testChdir(t, "/this/directory/does/not/exist")` fails.

```go
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
func testChdir(t *testing.T, dir string) func() {
	old, err := os.Getwd()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := os.Chdir(dir); err != nil {
		t.Fatalf("err: %s", err) // b_test.go:22
	}

	return func() { os.Chdir(old) }
}
```

`b_test.go:22` points to `t.Fatalf("err: %s", err)` inside of `testChDir()` which isn't so helpful to understand why the test failed.

```shellsession
$ go test
--- FAIL: TestFoo (0.00s)
    b_test.go:22: err: chdir /this/directory/does/not/exist: no such file or directory
FAIL
exit status 1
FAIL    github.com/ichiban/thelper/testdata/src/b       0.006s
```

By marking `testChdir()`, we can get a meaningful line number.

```go
package b

import (
	"os"
	"testing"
)

func TestFoo(t *testing.T) {
	defer testChdir(t, "/this/directory/does/not/exist")() // b_test.go:9

	// ...
}

// https://speakerdeck.com/mitchellh/advanced-testing-with-go?slide=30
func testChdir(t *testing.T, dir string) func() {
	t.Helper()

	old, err := os.Getwd()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := os.Chdir(dir); err != nil {
		t.Fatalf("err: %s", err)
	}

	return func() { os.Chdir(old) }
}
```

Now, `b_test.go:9` points to `defer testChdir(t, "/this/directory/does/not/exist")()`.

```shellsession
$ go test
--- FAIL: TestFoo (0.00s)
    b_test.go:9: err: chdir /this/directory/does/not/exist: no such file or directory
FAIL
exit status 1
FAIL    github.com/ichiban/thelper/testdata/src/b       0.006s
```

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

This package is based on [ikawaha](https://github.com/ikawaha)'s idea and advices.
