// Copyright 2022 Robert Horrox
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils_test

import (
	"testing"

	"github.com/RobertHorrox/gofstools/internal/test"

	"github.com/RobertHorrox/gofstools/utils"

	"github.com/stretchr/testify/assert"
)

type ParamTest struct {
	param       string
	IsRune      bool
	Expected    int32
	ExpectedErr error
}

func (p ParamTest) Name() string {
	name := p.param
	if p.IsRune {
		name += "_isRune"
	}

	return name
}

type PatternTest struct {
	pattern string
	start   string
	end     string
	incr    int32
	err     error
}

func (p PatternTest) Name() string {
	return p.pattern
}

func TestDecodeParam(t *testing.T) {
	t.Parallel()

	paramTests := []ParamTest{
		{"A", true, 65, nil},
		{"b", true, 98, nil},
		{"1", true, 49, nil},
		{"1", false, 1, nil},
		{"10", false, 10, nil},
		{"101", false, 101, nil},
		{"0005", false, 5, nil},
	}

	test.RunSubTests(t, paramTests, func(t *testing.T, param ParamTest) {
		t.Helper()
		decdedParam, err := utils.InternalDecodeParameter(param.param, param.IsRune)
		assert.Equal(t, param.ExpectedErr, err)
		assert.Equal(t, param.Expected, decdedParam)
	})
}

func TestExpandIntSeq(t *testing.T) {
	t.Parallel()

	IntSeqTests := []test.SequanceTest[int32]{
		{Start: 0, End: 5, Incr: 1, Expected: []int32{0, 1, 2, 3, 4, 5}},
		{Start: 5, End: 0, Incr: -1, Expected: []int32{5, 4, 3, 2, 1, 0}},
		{Start: 2, End: 10, Incr: 2, Expected: []int32{2, 4, 6, 8, 10}},
		{Start: 1, End: 10, Incr: 2, Expected: []int32{1, 3, 5, 7, 9}},
		{Start: 9, End: 0, Incr: -3, Expected: []int32{9, 6, 3, 0}},
		{Start: 8, End: 1, Incr: -3, Expected: []int32{8, 5, 2}},
		{Start: -5, End: 0, Incr: 1, Expected: []int32{-5, -4, -3, -2, -1, 0}},
		{Start: -6, End: 6, Incr: 3, Expected: []int32{-6, -3, 0, 3, 6}},
		{Start: 6, End: -6, Incr: -3, Expected: []int32{6, 3, 0, -3, -6}},
	}

	test.RunSubTests(t, IntSeqTests, func(t *testing.T, s test.SequanceTest[int32]) {
		t.Helper()
		intSeq := utils.InternalExpandIntSeq(s.Start, s.End, s.Incr)
		assert.Equal(t, s.Expected, intSeq)
	})
}

func TestParsePatern(t *testing.T) {
	t.Parallel()

	patterns := []PatternTest{
		{pattern: "1..5", start: "1", end: "5", incr: 1, err: nil},
		{pattern: "-1..5", start: "-1", end: "5", incr: 1, err: nil},
		{pattern: "-10..-5", start: "-10", end: "-5", incr: 1, err: nil},
		{pattern: "0001..5", start: "0001", end: "5", incr: 1, err: nil},
		{pattern: "1..0005", start: "1", end: "0005", incr: 1, err: nil},
		{pattern: "00001..005", start: "00001", end: "005", incr: 1, err: nil},
		{pattern: "a..e", start: "a", end: "e", incr: 1, err: nil},
		{pattern: "1..5..1", start: "1", end: "5", incr: 1, err: nil},
		{pattern: "1..5..3", start: "1", end: "5", incr: 3, err: nil},
		{pattern: "5..1..-1", start: "5", end: "1", incr: -1, err: nil},
		{pattern: "5..1..-3", start: "5", end: "1", incr: -3, err: nil},
		{pattern: "a..e..16", start: "a", end: "e", incr: 16, err: nil},
	}

	test.RunSubTests(t, patterns, func(t *testing.T, patternTest PatternTest) {
		t.Helper()
		start, end, incr, err := utils.InternalParsePattern(patternTest.pattern)
		assert.Equal(t, patternTest.start, start)
		assert.Equal(t, patternTest.end, end)
		assert.Equal(t, patternTest.incr, incr)
		assert.Equal(t, patternTest.err, err)
	})
}

func TestSeqExpansion(t *testing.T) {
	t.Parallel()

	sequences := []test.ExpansionTest{
		{Pattern: "1..5", Expected: []string{"1", "2", "3", "4", "5"}},
		{Pattern: "a..e", Expected: []string{"a", "b", "c", "d", "e"}},
		{Pattern: "!..&", Expected: []string{"!", "\"", "#", "$", "%", "&"}},
		{Pattern: "-..3", Expected: []string{"-", ".", "/", "0", "1", "2", "3"}},
		{Pattern: "Z..a", Expected: []string{"Z", "[", "\\", "]", "^", "_", "`", "a"}},
		{Pattern: "z..~", Expected: []string{"z", "{", "|", "}", "~"}},
		{Pattern: "9..12", Expected: []string{"9", "10", "11", "12"}},
		{Pattern: "0009..12", Expected: []string{"0009", "0010", "0011", "0012"}},
		{Pattern: "9..0012", Expected: []string{"0009", "0010", "0011", "0012"}},
		{Pattern: "0..8..2", Expected: []string{"0", "2", "4", "6", "8"}},
		{Pattern: "8..0..-2", Expected: []string{"8", "6", "4", "2", "0"}},
		{Pattern: "0..14..3", Expected: []string{"0", "3", "6", "9", "12"}},
		{Pattern: "a..z..3", Expected: []string{"a", "d", "g", "j", "m", "p", "s", "v", "y"}},
	}

	test.RunSubTests(t, sequences, func(t *testing.T, expansionTest test.ExpansionTest) {
		t.Helper()

		expand, err := utils.ExpandSeq(expansionTest.Pattern)
		assert.Nil(t, err)
		assert.Equal(t, expansionTest.Expected, expand)
	})
}
