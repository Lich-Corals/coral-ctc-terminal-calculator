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

// Testing with "2 * (3 + (2 + 6) * 2) + sqrt(9) + 200 ** 0"
//
// Priorities:
// A. powers & factorials
// B. multiplication and division
// C. addition and substriction
//
// Processing from left to right

package main

import (
	"os"
	"regexp"
	"slices"
	"strings"
)

var (
	tokenSplitPoints     []string = []string{"(", ")"}
	functionNames        []string = []string{"sqrt"}
	numberRegex                   = regexp.MustCompile(`\d+`)
	factorialNumberRegex          = regexp.MustCompile(`\d+!`)
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
	var tokens = make([][]token, 2)
	var usedTokens = 1
	var otherTokens = 0
	var run bool = true
	tokens[usedTokens] = getUnprocessedTokens(strings.Split(arg, " "))
	for run {
		if usedTokens == 0 {
			usedTokens = 1
			otherTokens = 0
		} else {
			usedTokens = 0
			otherTokens = 1
		}
		tokens[usedTokens] = []token{}
		for _, tok := range tokens[usedTokens] {
			if len(tok.unprocessed) > 0 {
				for _, tP := range tok.unprocessed {
					processTokenPart(tP)
				}
			}
		}
	}
	return tokens[otherTokens]
}

// Process an unprocessed token part
func processTokenPart(part string) token {
	var s = false
	for _, sP := range tokenSplitPoints {
		if strings.Contains(part, sP) {
			s = true
			// TODO: Extract sub-tokens for data in brackets
		}
	}
	if !s {
		var tT, tP = getTokenTypeAndPriority(part)
		return token{token: tT, priority: tP, content: []string{part}, subTokens: nil, unprocessed: nil}
	}
	return emptyToken
}

// convert a list of strings into unprocessed tokens
func getUnprocessedTokens(data []string) []token {
	var tks []token
	for _, d := range data {
		tks = append(tks, token{token: unknownTokenType, priority: pX, content: nil, subTokens: nil, unprocessed: []string{d}})
	}
	return tks
}

// Get the type and priority of an isolated token
func getTokenTypeAndPriority(content string) (tokenType, tokenPriority) {
	if len(numberRegex.FindAllString(content, -1)) > 0 {
		return number, pX
	} else if len(factorialNumberRegex.FindAllString(content, -1)) > 0 {
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
		case "(":
			return openDelim, pX
		case ")":
			return closeDelim, pX
		}

		if slices.Contains(functionNames, content) {
			return functionName, pX
		}
	}
	return unknownTokenType, pX
}

// Get sub token of tokens like "(6" or "sqrt(9)"
func getSubTokens(token string) []string {
	tokens := strings.Split(token, "(")
	if len(tokens[0]) == 0 {

	}
	return tokens
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
