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

// Testing with "2 * (3 + (2 + 6) * 4) + ( 2 // 9! ) + 200 ** 0 + -5"
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
	"strings"
)

var (
	numberRegex          = regexp.MustCompile(`^-{0,1}\d+(?:,\d+){0,1}$`)
	factorialNumberRegex = regexp.MustCompile(`^-{0,1}\d+!$`)
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
)

const (
	pX tokenPriority = iota
	pA
	pB
	pC
)

type token struct {
	token     tokenType
	priority  tokenPriority
	content   []string
	subTokens []token
}

// Get tokens from input string
func getTokens(arg string) []token {
	var ungrouped = []token{}
	for tok := range strings.SplitSeq(arg, " ") {
		for _, t := range processToken(tok) {
			ungrouped = append(ungrouped, t)
		}
	}

	var brS = 0
	for _, tok := range ungrouped {
		switch tok.token {
		case openDelim:
			brS += 1
		case closeDelim:
			brS -= 1
		}
	}
	if brS < 0 {
		panic("Unmatched ')'")
	} else if brS > 0 {
		panic("Unmatched '('")
	}

	var grouped = makeGroups(ungrouped)

	fmt.Println(grouped)

	return ungrouped
}

// Move all contents of delimiters into the sub-tokens part of their opening delimiter
func makeGroups(tokens []token) []token {
	var grouped = []token{}
	var startIndex = 0
	for ii, tok := range tokens {
		if ii >= startIndex {
			if tok.token == openDelim {
				var deS = 1
				for i, t := range tokens[ii+1:] {
					switch t.token {
					case openDelim:
						deS += 1
					case closeDelim:
						deS -= 1
					}
					if deS == 0 {
						tok.subTokens = makeGroups(tokens[ii : i+ii+1])
						grouped = append(grouped, tok)
						startIndex = i + ii + 1
						break
					}
				}
			} else if tok.token != closeDelim {
				grouped = append(grouped, tok)
			}
		}
	}
	return grouped
}

// Split a string like '(6' into sub-tokens
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

// Get tokens from a string like "6", "(9" or "**"
func processToken(part string) []token {
	var parts = getSubTokens(part)
	var tokens = []token{}
	for _, part := range parts {
		var tT, tP = getTokenTypeAndPriority(part)
		tokens = append(tokens, token{token: tT, priority: tP, content: []string{part}, subTokens: nil})
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
	}
	panic(fmt.Sprint("Unknown token: ", content))
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
