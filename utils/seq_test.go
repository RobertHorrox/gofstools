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

func TestSeqExpansion(t *testing.T) {
	t.Parallel()

	sequences := []test.ExpansionTest{
		{
			Pattern: "1..5", Expected: []string{"1", "2", "3", "4", "5"},
		},
		{
			Pattern: "a..e", Expected: []string{"a", "b", "c", "d", "e"},
		},
		{
			Pattern: "!..&", Expected: []string{"!", "\"", "#", "$", "%", "&"},
		},
		{
			Pattern: "-..3", Expected: []string{"-", ".", "/", "0", "1", "2", "3"},
		},
		{
			Pattern: "Z..a", Expected: []string{"Z", "[", "\\", "]", "^", "_", "`", "a"},
		},
		{
			Pattern: "z..~", Expected: []string{"z", "{", "|", "}", "~"},
		},
		{
			Pattern: "9..12", Expected: []string{"9", "10", "11", "12"},
		},
		{
			Pattern: "0009..12", Expected: []string{"0009", "0010", "0011", "0012"},
		},
		{
			Pattern: "9..0012", Expected: []string{"0009", "0010", "0011", "0012"},
		},
		{
			Pattern: "0..8..2", Expected: []string{"0", "2", "4", "6", "8"},
		},
		{
			Pattern: "0..14..3", Expected: []string{"0", "3", "6", "9", "12"},
		},
		{
			Pattern: "a..z..3", Expected: []string{"a", "d", "g", "j", "m", "p", "s", "v", "y"},
		},
	}

	for _, seq := range sequences {
		localSeq := seq
		t.Run(localSeq.Pattern, func(subt *testing.T) {
			subt.Parallel()
			expand, err := utils.ExpandSeq(localSeq.Pattern)
			assert.Nil(subt, err)
			assert.Equal(subt, localSeq.Expected, expand)
		})
	}
}
