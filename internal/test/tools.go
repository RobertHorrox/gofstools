package test

import (
	"fmt"
	"testing"
)

type ExpansionTest struct {
	Pattern  string
	Expected []string
}

func (e ExpansionTest) Name() string {
	return e.Pattern
}

type SequanceTest[T comparable] struct {
	Start    T
	End      T
	Incr     T
	Expected []T
}

func (s SequanceTest[T]) Name() string {
	return fmt.Sprintf("%v_%v_%v", s.Start, s.End, s.Incr)
}

type Named interface {
	Name() string
}

func RunSubTests[T Named](t *testing.T, tests []T, testFunc func(*testing.T, T)) {
	t.Helper()

	for _, test := range tests {
		local := test
		t.Run(test.Name(), func(subt *testing.T) {
			subt.Parallel()
			testFunc(subt, local)
		})
	}
}
