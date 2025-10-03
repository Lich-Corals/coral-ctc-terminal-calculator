// Coral-CTC-Terminal-Calculator - Minimal terminal calculator
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

// Testing with "2 * (3 + (2 + 6.1) * 4) + ( 2 // 9 ! ) + 200 ** 0 + -5 * (1 + 2) !"
//
// Priorities:
// (X. Numbers, factorials, groups)
// A. Powers, roots
// B. Multiplication and division
// C. Addition and subtraction
//
// Processing from left to right

package main

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	numberRegex        = regexp.MustCompile(`^-{0,1}\d+(?:\.\d+){0,1}$`)
	decimalNumberRegex = regexp.MustCompile(`^-{0,1}\d+\.\d+$`)
)

type tokenType int8
type tokenPriority int8

const (
	unknownTokenType tokenType = iota
	number
	addition
	subtraction
	multiplication
	division
	power
	root
	factorial
	openDelim
	closeDelim
)

// Priorities for different types of tokens
const (
	pX tokenPriority = iota
	pA
	pB
	pC
)

const (
	ansiRed   = "\033[31m"
	ansiReset = "\033[0m"
	ansiBlue  = "\033[34m"
)

type token struct {
	token     tokenType
	priority  tokenPriority
	content   []string
	subTokens []token
	sum       float64
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
		quitWithError("Unmatched ')'")
	} else if brS > 0 {
		quitWithError("Unmatched '('")
	}

	var grouped = makeGroups(ungrouped)

	return grouped
}

// calculate
func getSum(tokens []token) float64 {
	var summed = []token{}
	for _, tok := range tokens {
		switch tok.token {
		case number:
			val, err := strconv.ParseFloat(tok.content[0], 64)
			if err != nil {
				quitWithError(fmt.Sprint("Cannot convert to number: ", tok.content[0]))
			}
			tok.sum = val
		case openDelim:
			tok.sum = getSum(tok.subTokens) // Calculates parenthesized-contents recursively
		}
		summed = append(summed, tok)
	}
	// Faculties are converted into numbers before any other calculations start
	var withFactorials = []token{}
	for i, tok := range summed {
		if tok.token == factorial {
			if len(decimalNumberRegex.FindAllString(summed[i-1].content[0], -1)) != 0 {
				quitWithError(fmt.Sprint("Factorial numbers can't be based on decimal numbers: ", summed[i-1].sum, " !"))
			}
			if summed[i-1].sum < 0 {
				quitWithError(fmt.Sprint("You can't get the factorial of a negative number: ", summed[i-1].sum, " !"))
			}
			tok.token = number
			tok.priority = pX
			tok.sum = fact(summed[i-1].sum)
			withFactorials = withFactorials[:i-1]
		}
		withFactorials = append(withFactorials, tok)
	}
	var priorities = []tokenPriority{pA, pB, pC}
	var groups = make([][]token, 2) // The following code is switching between two groups to prevent editing of the slice it is currently iterating over
	var this = 0
	var other = 1
	groups[this] = withFactorials
	for _, pri := range priorities {
		var run = true
		for run {
			groups[other] = []token{}
			var r = true
			var skip = false
			for i, tok := range groups[this] {
				if tok.priority == pri && r {
					var a = groups[this][i-1]
					var b = groups[this][i+1]
					tok.priority = pX
					switch tok.token {
					case power:
						tok.sum = math.Pow(a.sum, b.sum)
					case root:
						tok.sum = math.Pow(b.sum, 1.0/a.sum)
					case multiplication:
						tok.sum = a.sum * b.sum
					case division:
						if b.sum == 0.0 {
							if a.sum == 0.0 {
								println(ansiRed, fmt.Sprint("You can't divide by 0: ", a.sum, " / ", b.sum), ansiBlue, "\nNever gonna give you up!\nNever gonna let you down\n...", ansiReset)
								os.Exit(69)
							} else {
								quitWithError(fmt.Sprint("You can't divide by 0: ", a.sum, " / ", b.sum))
							}
						}
						tok.sum = a.sum / b.sum
					case addition:
						tok.sum = a.sum + b.sum
					case subtraction:
						tok.sum = a.sum - b.sum
					}
					tok.token = number
					groups[other] = groups[other][:len(groups[other])-1]
					groups[other] = append(groups[other], tok)
					groups[other] = append(groups[other], groups[other][i:]...)
					skip = true
					r = false
				} else {
					if skip {
						skip = false
					} else {
						groups[other] = append(groups[other], tok)
					}
				}
			}

			this = other
			if this == 1 {
				other = 0
			} else {
				other = 1
			}
			if r {
				run = false
			}
		}
	}
	return groups[other][0].sum
}

// Find the factorial of a number
func fact(n float64) float64 {
	if n == 1 || n == 0 {
		return 1.0
	}
	var fact = n * fact(n-1)
	return fact
}

// Move all contents of parentheses into the sub-tokens part of their opening delimiter
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
		quitWithError(fmt.Sprint("One token may never contain both '(' and ')': ", part))
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
	} else {
		switch content {
		case "+":
			return addition, pC
		case "-":
			return subtraction, pC
		case "*":
			return multiplication, pB
		case "/":
			return division, pB
		case "**":
			return power, pA
		case "//":
			return root, pA
		case "!":
			return factorial, pA
		case "(":
			return openDelim, pX
		case ")":
			return closeDelim, pX
		}
	}
	quitWithError(fmt.Sprint("Unknown token: ", content))
	return unknownTokenType, pX
}

// Show an error to the user
func quitWithError(content string) {
	println(ansiRed, content, ansiReset)
	println(ansiBlue, "If you believe this is a bug, please open an issue: https://github.com/Lich-Corals/coral-ctc-terminal-calculator/issues", ansiReset)
	os.Exit(1)
}

func main() {
	var terminalArguments = os.Args
	var tokens []token
	var sum float64
	if len(terminalArguments) < 2 {
		quitWithError("Missing an argument!")
	} else if len(terminalArguments) > 2 {
		quitWithError("Too many arguments!")
	}
	for i, v := range terminalArguments {
		if i == len(terminalArguments)-1 { // Use the last argument to get operations
			tokens = getTokens(v)
			sum = getSum(tokens)
		}
	}
	fmt.Println(sum)
}
