package a

import "testing"

func TestFoo(t *testing.T) {
	foo(t)
	bar(t)
	baz(0)
}

func foo(t *testing.T) { // want `unmarked test helper`
	bar(t)
}

func bar(t *testing.T) {
	t.Helper()
	t.Errorf("bar %d", int64(0))
	qux(t)
}

func baz(_ int) {
	_ = []int(nil)
}

func qux(s *testing.T) { // want `unmarked test helper`
	s.Error("qux")
}
