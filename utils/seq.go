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

package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/RobertHorrox/gofstools/internal/utils"
)

const (
	PatternSplitLengthBase = 2
	PatternSplitLengthIncr = 3
	RuneLength             = 1
)

var (
	ErrParsingIncr    = errors.New("cannot decode increment")
	ErrInvalidPattern = errors.New("invalid pattern format")
)

func parsePattern(pattern string) (string, string, int, error) {
	elements := strings.Split(pattern, "..")
	start := elements[0]
	end := elements[1]
	incr := 1

	var err error

	switch len(elements) {
	case PatternSplitLengthBase: // Base case
		break
	case PatternSplitLengthIncr: // Case with increment detected
		incr, err = strconv.Atoi(elements[2])
		if err != nil {
			err = ErrParsingIncr
		}
	default: // Invalid pattern
		err = ErrInvalidPattern
	}

	return start, end, incr, err
}

func getZeroPadLength(value string) int {
	length := 0
	if value[0] == '0' {
		length = len(value)
	}

	return length
}

func ExpandSeq(pattern string) ([]string, error) {
	start, end, incr, err := parsePattern(pattern)
	if err != nil {
		return nil, err
	}

	// First check for runes
	var retVal []string

	if len(start) == RuneLength && len(end) == RuneLength {
		startRune, _ := utf8.DecodeRuneInString(start)
		endRune, _ := utf8.DecodeRuneInString(end)

		for i := startRune; i <= endRune; i += rune(incr) {
			retVal = append(retVal, string(i))
		}

		return retVal, nil
	}

	// Check for integers
	startInt, startErr := strconv.Atoi(start)
	endInt, endErr := strconv.Atoi(end)

	if startErr != nil || endErr != nil {
		return retVal, ErrInvalidPattern
	}

	// Check for zero padding
	paddingLength := utils.IntMax(getZeroPadLength(start), getZeroPadLength(end))
	fmtStr := "%0" + strconv.Itoa(paddingLength) + "d"

	for i := startInt; i <= endInt; i += incr {
		retVal = append(retVal, fmt.Sprintf(fmtStr, i))
	}

	return retVal, nil
}
