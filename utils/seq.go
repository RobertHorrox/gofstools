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
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/cockroachdb/errors"

	"github.com/RobertHorrox/gofstools/internal/utils"
)

const (
	PatternSplitLengthBase = 2
	PatternSplitLengthIncr = 3
	RuneLength             = 1
	RuneBitSize            = 32
)

var ErrInvalidPattern = errors.New("invalid pattern format")

func decodeParameter(param string, isRune bool) (int32, error) {
	var err error

	var parsedParameter int32

	if isRune {
		parsedParameter, _ = utf8.DecodeRuneInString(param)
		if parsedParameter == utf8.RuneError {
			err = errors.Wrapf(ErrInvalidPattern, "Cannot decode rune %s", param)
		}
	} else {
		param64, parseErr := strconv.ParseInt(param, 10, RuneBitSize)
		parsedParameter = int32(param64)
		err = parseErr
	}

	return parsedParameter, err
}

func expandIntSeq(start int32, end int32, incr int32) []int32 {
	// If incr is negative then i will decrement, check i >= endInt
	// If incr is positive then i will increment, check i <= endInt
	var retVal []int32
	for i := start; (incr < 0 && i >= end) || (incr > 0 && i <= end); i += incr {
		retVal = append(retVal, i)
	}

	return retVal
}

func parsePattern(pattern string) (string, string, int32, error) {
	elements := strings.Split(pattern, "..")
	start := elements[0]
	end := elements[1]

	var incr int32 = 1

	var retErr error

	switch len(elements) {
	case PatternSplitLengthBase: // Base case
		break
	case PatternSplitLengthIncr: // Case with increment detected
		incr64, err := strconv.ParseInt(elements[2], 10, RuneBitSize)
		if err != nil {
			retErr = errors.Wrap(err, "Error Parsing Increment")
		}

		incr = int32(incr64) // We can downsize to int32 since RuneBitSize is 32 bits
	default: // Invalid pattern
		retErr = ErrInvalidPattern
	}

	return start, end, incr, retErr
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
	isRuneSeq := len(start) == RuneLength && len(end) == RuneLength

	// Decoded Params
	startDecoded, err := decodeParameter(start, isRuneSeq)
	if err != nil {
		return nil, err
	}

	endDecoded, err := decodeParameter(end, isRuneSeq)
	if err != nil {
		return nil, err
	}

	// Get Sequence
	intSeq := expandIntSeq(startDecoded, endDecoded, incr)

	fmtString := ""
	if isRuneSeq { // Rune Seq, just cast the ints to characters in sprintf
		fmtString = "%c"
	} else { // If not a rune seq then it's an Integer seq with a possible padding
		// Check for zero padding
		paddingLength := utils.IntMax(getZeroPadLength(start), getZeroPadLength(end))
		fmtString = "%0" + strconv.Itoa(paddingLength) + "d"
	}

	outputSeq := make([]string, len(intSeq))
	for i, v := range intSeq {
		outputSeq[i] = fmt.Sprintf(fmtString, v)
	}

	return outputSeq, nil
}
