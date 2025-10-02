// Coral-Terminal-Calculator - Minimalist terminal calculator
// Copyright (C) 2025  Linus Tibert
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public Licence as published
// by the Free Software Foundation, either version 3 of the Licence, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public Licence for more details.
//
// You should have received a copy of the GNU Affero General Public Licence
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// Testing with "2 * (3 + (2 + 6) * 2) + ( 9! // 2 ) + 200 ** 0"
//
// Priorities:
// A. powers, factorials, roots
// B. multiplication and division
// C. addition and substriction
//
// Processing from left to right

package main

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
)

var (
	functionNames        []string = []string{}
	numberRegex                   = regexp.MustCompile(`^\d+$`)
	factorialNumberRegex          = regexp.MustCompile(`^\d+!$`)
)

var (
	emptyToken = token{token: unknownTokenType, priority: pX, content: nil, subTokens: nil, unprocessed: nil}
)

type tokenType int8
type tokenPriority int8

const (
	unknownTokenType tokenType = iota
	number
	addition
	substriction
	multiplication
	division
	power
	root
	factorialNumber
	openDelim
	closeDelim
	functionName
)

const (
	pX tokenPriority = iota
	pA
	pB
	pC
)

type token struct {
	token       tokenType
	priority    tokenPriority
	content     []string
	subTokens   []token
	unprocessed []string
}

// Get tokens from input string
func getTokens(arg string) []token {
	var ungrouped = []token{}

	for tok := range strings.SplitSeq(arg, " ") {
		for _, t := range processToken(tok) {
			ungrouped = append(ungrouped, t)
		}
	}

	for _, tex := range ungrouped {
		fmt.Println(tex.content, tex.token)
	}
	return ungrouped
}

// Split a string like 'sqrt(9' into sub-tokens
func getSubTokens(part string) []string {
	var retVal = []string{}
	var s = false
	if strings.Contains(part, ")") && strings.Contains(part, "(") {
		panic(fmt.Sprint("One token may never contain both '(' and ')': ", part))
	}
	if strings.Contains(part, "(") {
		s = true
		var pts = strings.Split(part, "(")
		retVal = append(retVal, "(")
		if len(pts[1]) > 0 {
			retVal = append(retVal, pts[1])
		}
	} else if strings.Contains(part, ")") {
		s = true
		var pts = strings.Split(part, ")")
		if len(pts[0]) > 0 {
			retVal = append(retVal, pts[0])
		}
		retVal = append(retVal, ")")
	}
	if !s {
		retVal = []string{part}
	}
	return retVal
}

// Process an unprocessed token part
func processToken(part string) []token {
	var parts = getSubTokens(part)
	var tokens = []token{}
	for _, part := range parts {
		var tT, tP = getTokenTypeAndPriority(part)
		tokens = append(tokens, token{token: tT, priority: tP, content: []string{part}, subTokens: nil, unprocessed: nil})
	}
	return tokens
}

// Get the type and priority of an isolated token
func getTokenTypeAndPriority(content string) (tokenType, tokenPriority) {
	if len(numberRegex.FindAllString(content, -1)) == 1 {
		return number, pX
	} else if len(factorialNumberRegex.FindAllString(content, -1)) == 1 {
		return factorialNumber, pA
	} else {
		switch content {
		case "+":
			return addition, pC
		case "-":
			return substriction, pC
		case "*":
			return multiplication, pB
		case "/":
			return division, pB
		case "**":
			return power, pA
		case "//":
			return root, pA
		case "(":
			return openDelim, pX
		case ")":
			return closeDelim, pX
		}

		if slices.Contains(functionNames, content) {
			return functionName, pX
		}
	}
	panic(fmt.Sprint("Unknown token: ", content))
	return unknownTokenType, pX
}

func main() {
	var terminalArguments = os.Args
	var tokens []token
	for i, v := range terminalArguments {
		if i == len(terminalArguments)-1 { // Use the last argument to get operations
			tokens = getTokens(v)
		}
	}
	println(tokens)
}
