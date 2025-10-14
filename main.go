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

// Priorities:
// (X. Numbers, factorials, groups)
// A. Powers, roots, absolute-function and sine, etc.
// B. Multiplication, division, nCr, nPr and modulo
// C. Addition, subtraction, logarithms
//
// Processing from left to right

package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

// The current package version
const pkgVersion = "0.5.0"

// Global variables
var (
	currentInputMode             = noInuputMode
	calculationSuccess           = true
	lastAnswer         []float64 = nil
)

var (
	numberRegex                 = regexp.MustCompile(`^-{0,1}\d+(?:\.\d+){0,1}$`)
	decimalNumberRegex          = regexp.MustCompile(`^-{0,1}\d+\.\d+$`)
	constants          []string = []string{"pi", "tau", "e", "g", "phi", "c", "ans"}
)

type tokenType int8
type tokenPriority int8
type inputMode int8
type neededArgumetns int8

// Input mode mostly determines whether there is an argument; if there are none, continuous mode will be used.
const (
	noInuputMode inputMode = iota
	single
	continuous
)

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
	modulo
	logarithm
	constant
	nPr
	nCr
	sine
	cosine
	tangent
	sineDegrees
	cosineDegrees
	tangentDegrees
	absoluteFunction
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

// Which which arguments are needed for an operation
const (
	noArgs neededArgumetns = iota
	left
	right
	leftRight
)

const (
	ansiRed   = "\033[31m"
	ansiReset = "\033[0m"
	ansiBlue  = "\033[34m"
)

type token struct {
	token      tokenType
	priority   tokenPriority
	content    []string
	subTokens  []token
	sum        float64
	neededArgs neededArgumetns
}

// Get tokens from input string
func GetTokens(arg string) []token {
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
		userError("Unmatched ')'")
	} else if brS > 0 {
		userError("Unmatched '('")
	}

	var grouped = makeGroups(ungrouped)

	return grouped
}

// calculate
func GetSum(tokens []token) float64 {
	var summed = []token{}
	for _, tok := range tokens {
		switch tok.token {
		case number:
			val, err := strconv.ParseFloat(tok.content[0], 64)
			if err != nil {
				userError(fmt.Sprint("Cannot convert to number: ", tok.content[0]))
			}
			tok.sum = val
		case openDelim:
			tok.sum = GetSum(tok.subTokens) // Calculates parenthesized-contents recursively
		}
		summed = append(summed, tok)
	}
	// Faculties are converted into numbers before any other calculations start
	var withFactorials = []token{}
	for i, tok := range summed {
		if tok.token == factorial {
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
			var ru = true
			var skip = false
			for i, tok := range groups[this] {
				if tok.priority == pri && ru {
					var a = tok
					var b = tok
					var missingArg = false
					switch tok.neededArgs {
					case leftRight:
						if i-1 >= 0 && i+1 < len(groups[this]) {
							a = groups[this][i-1]
							b = groups[this][i+1]
						} else {
							missingArg = true
						}
					case left:
						if i-1 >= 0 {
							a = groups[this][i-1]
						} else {
							missingArg = true
						}
					case right:
						if i+1 < len(groups[this]) {
							b = groups[this][i+1]
						} else {
							missingArg = true
						}
					}
					if missingArg {
						userError(fmt.Sprintf("Operator at index %v missing argument(s): %s\nNote: The first index is 0", i, tok.content[0]))
						groups[other] = append(groups[other], tok) // Those add dummy arguments to avoid panics in the following code;
						groups[other] = append(groups[other], tok) // The result will not be visible anyway.
						run = false
					}
					tok.priority = pX
					switch tok.token {
					case power:
						tok.sum = math.Pow(a.sum, b.sum)
					case root:
						if b.sum < 0 && int(a.sum)%2 == 0 {
							userError(fmt.Sprint("Negative numbers do not have roots of even numbers: ", b.sum, " // ", a.sum))
						}
						if a.sum == 0 {
							userError(fmt.Sprint("Can't get the 0th root: ", a.sum, " // ", b.sum))
						}
						tok.sum = math.Pow(absolute(b.sum), 1.0/a.sum)
						if b.sum < 0 {
							tok.sum = -tok.sum
						}
					case multiplication:
						tok.sum = a.sum * b.sum
					case division:
						if b.sum == 0.0 {
							if a.sum == 0.0 {
								println(ansiRed, fmt.Sprint("You can't divide by 0: ", a.sum, " / ", b.sum), ansiBlue, "\nNever gonna give you up!\nNever gonna let you down\n...", ansiReset)
								os.Exit(69)
							} else {
								userError(fmt.Sprint("You can't divide by 0: ", a.sum, " / ", b.sum))
							}
						}
						tok.sum = a.sum / b.sum
					case modulo:
						if len(decimalNumberRegex.FindAllString(a.content[0], -1)) != 0.0 || len(decimalNumberRegex.FindAllString(b.content[0], -1)) != 0 {
							userError(fmt.Sprint("Cannot perform modulo on float values: ", a.sum, " % ", b.sum))
						}
						tok.sum = float64(int(a.sum) % int(b.sum))
					case logarithm:
						if a.sum == 0 || b.sum == 0 {
							userError(fmt.Sprint("Logarithm with zero as base or x: ", a.sum, " log ", b.sum))
						}
						tok.sum = math.Log(a.sum) / math.Log(b.sum)
					case absoluteFunction:
						tok.sum = absolute(b.sum)
					case nCr:
						tok.sum = calcNCr(a.sum, b.sum)
					case nPr:
						tok.sum = calcNPr(a.sum, b.sum)
					case addition:
						tok.sum = a.sum + b.sum
					case subtraction:
						tok.sum = a.sum - b.sum
					case sine:
						tok.sum = math.Sin(b.sum)
					case cosine:
						tok.sum = math.Cos(b.sum)
					case tangent:
						tok.sum = math.Tan(b.sum)
					case sineDegrees:
						tok.sum = math.Sin(degreesToRadians(b.sum))
					case cosineDegrees:
						tok.sum = math.Cos(degreesToRadians(b.sum))
					case tangentDegrees:
						tok.sum = math.Tan(degreesToRadians(b.sum))
					}
					tok.token = number
					if !missingArg {
						switch tok.neededArgs {
						case leftRight:
							groups[other] = groups[other][:len(groups[other])-1]        // Pop a
							groups[other] = append(groups[other], tok)                  // Insert self at a's place
							groups[other] = append(groups[other], groups[other][i:]...) // Add everything after b
						case right:
							groups[other] = append(groups[other], tok)                    // Add self
							groups[other] = append(groups[other], groups[other][i+1:]...) // Add everything after b
						case left:
							groups[other] = groups[other][:len(groups[other])-1]        // Pop a
							groups[other] = append(groups[other], tok)                  // Add self at a's place
							groups[other] = append(groups[other], groups[other][i:]...) // Add everything after self
						}
					}
					skip = true
					ru = false
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
			if ru {
				run = false
			}
		}
	}
	if len(groups[other]) > 1 {
		userError(fmt.Sprint("Too many calculation results: ", groups[other][0].sum, " ", groups[other][1].sum, "\nMaybe you forgot an operator?"))
	}
	return groups[other][0].sum
}

// Find the factorial of a number
func fact(n float64) float64 {
	if n != float64(int64(n)) {
		userError(fmt.Sprint("Factorial numbers can't be based on decimal numbers: ", fmt.Sprintf("%f", n), " !"))
		return 0.0
	}
	if n < 0 {
		userError(fmt.Sprint("You can't get the factorial of a negative number: ", fmt.Sprintf("%f", n), " !"))
		return 0.0
	}
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
		userError(fmt.Sprint("One token may never contain both '(' and ')': ", part))
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
		var tT, tP, tA = getTokenProperties(part)
		if tT == constant {
			tT = number
			nP := 0.0
			switch strings.Replace(part, "-", "", 1) {
			case "pi":
				nP = math.Pi
			case "tau":
				nP = math.Pi * 2
			case "e":
				nP = math.E
			case "phi":
				nP = math.Phi
			case "g":
				nP = 9.8066500 // https://oeis.org/A072915
			case "c":
				nP = 299792458.0 // https://oeis.org/A003678
			case "ans":
				if lastAnswer != nil {
					nP = lastAnswer[0]
				} else {
					userError("Can't use `ans` without previous answer!")
				}
			case "0":
				nP = 0
			case "1":
				nP = 1
			default:
				userError(fmt.Sprint("Constant not defined: ", part, "\nThis looks like a bug; please report."))
			}
			if strings.Contains(part, "-") {
				nP = -nP
			}
			part = strconv.FormatFloat(nP, 'f', -1, 64)
		}
		tokens = append(tokens, token{token: tT, priority: tP, content: []string{part}, subTokens: nil, neededArgs: tA})
	}
	return tokens
}

// Get the type and priority of an isolated token
func getTokenProperties(content string) (tokenType, tokenPriority, neededArgumetns) {
	if len(numberRegex.FindAllString(content, -1)) == 1 {
		return number, pX, noArgs
	} else if slices.Contains(constants, content) {
		return constant, pX, noArgs
	} else {
		switch content {
		case "+":
			return addition, pC, leftRight
		case "-":
			return subtraction, pC, leftRight
		case "*":
			return multiplication, pB, leftRight
		case "/":
			return division, pB, leftRight
		case "**":
			return power, pA, leftRight
		case "//":
			return root, pA, leftRight
		case "!":
			return factorial, pA, left
		case "%":
			return modulo, pB, leftRight
		case "log":
			return logarithm, pC, leftRight
		case "nCr":
			return nCr, pB, leftRight
		case "nPr":
			return nPr, pB, leftRight
		case "sin":
			return sine, pA, right
		case "cos":
			return cosine, pA, right
		case "tan":
			return tangent, pA, right
		case "dsin":
			return sineDegrees, pA, right
		case "dcos":
			return cosineDegrees, pA, right
		case "dtan":
			return tangentDegrees, pA, right
		case "abs":
			return absoluteFunction, pA, right
		case "(":
			return openDelim, pX, noArgs
		case ")":
			return closeDelim, pX, noArgs
		}
	}
	userError(fmt.Sprint("Unknown token: ", content))
	return unknownTokenType, pX, noArgs
}

// Convert degrees to radians
func degreesToRadians(x float64) float64 {
	return x * (math.Pi / 180)
}

// make a number absolute
func absolute(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// nCr-combinations
func calcNCr(n float64, r float64) float64 {
	return fact(n) / (fact(r) * fact(n-r))
}

// nPr-permutations
func calcNPr(n float64, r float64) float64 {
	return fact(n) / fact(n-r)
}

// Show an error to the user
func userError(content string) {
	if calculationSuccess {
		fmt.Printf("%v%s%v\n", ansiRed, content, ansiReset)
		println(ansiBlue, "If you believe this is a bug, please open an issue: https://github.com/Lich-Corals/coral-ctc-terminal-calculator/issues", ansiReset)

	}
	switch currentInputMode {
	case single:
		os.Exit(1)
	case continuous:
		calculationSuccess = false
	}
}

// Show a licence notice to the user
func showLicence() {
	println("\nCoral-CTC-Terminal-Calculator  Copyright (C) 2025  Linus Tibert\nThis program comes with ABSOLUTELY NO WARRANTY.\nThis is free software, and you are welcome to redistribute it\nunder certain conditions. You can view the licence here:\nhttps://github.com/Lich-Corals/coral-ctc-terminal-calculator/\n", ansiBlue, "\nFor questions and issues, please head to the GitHub repository above.", ansiReset)

}

// Tell the user to get help somewhere else
func showHelp() {
	println("Please take a look at the git repository for detailed instructions on how to use this program:\nhttps://github.com/Lich-Corals/coral-ctc-terminal-calculator/")
}

// Add a -x for every constant x
func updateConstants() {
	var newConstants []string
	for _, c := range constants {
		newConstants = append(newConstants, fmt.Sprintf("-%s", c))
		newConstants = append(newConstants, c)
	}
	constants = newConstants
}

func main() {
	updateConstants()
	var terminalArguments = os.Args
	var tokens []token
	var sum float64
	if len(terminalArguments) == 1 {
		currentInputMode = continuous
	} else if len(terminalArguments) > 2 {
		userError("Too many arguments!")
	} else {
		currentInputMode = single
	}
	switch currentInputMode {
	case single:
		switch terminalArguments[1] {
		case "--licence", "--license":
			showLicence()
			os.Exit(0)
		case "--version":
			println(ansiBlue, pkgVersion, ansiReset)
			os.Exit(0)
		case "--help", "-h", "help":
			showHelp()
			os.Exit(0)
		}
		tokens = GetTokens(terminalArguments[len(terminalArguments)-1])
		sum = GetSum(tokens)
		fmt.Println(strconv.FormatFloat(sum, 'f', -1, 64))
	case continuous:
		showLicence()
		cursorRune := "█"
		commandHistory := []string{}
		historyPos := 0
		for true {
			buf := bytes.NewBufferString("")
			cursorPos := 0
			fmt.Printf("> %s\n", cursorRune)
			line := ""
			keyboard.Listen(func(key keys.Key) (stop bool, err error) {
				fmt.Printf("\033[1A\033[K") // Clears the last line
				switch key.Code {
				case keys.Enter:
					line = buf.String()
					if len(commandHistory) == 0 || commandHistory[len(commandHistory)-1] != buf.String() {
						commandHistory = append(commandHistory, buf.String())
					}
					historyPos = len(commandHistory)
					fmt.Printf("> %s\n", buf)
					return true, nil
				case keys.Backspace: // Remove last character
					if cursorPos > 0 {
						n := bytes.NewBuffer(buf.Bytes()[:cursorPos-1])
						n.Write(buf.Bytes()[cursorPos:])
						buf = n
						cursorPos -= 1
					}
				case keys.Delete: // Remove next character
					if cursorPos < buf.Len() {
						n := bytes.NewBuffer(buf.Bytes()[0:cursorPos])
						a := bytes.NewBuffer(buf.Bytes()[cursorPos+1:])
						n.Write(a.Bytes())
						buf = n
					}
				case keys.Left:
					if cursorPos > 0 {
						cursorPos -= 1
					}
				case keys.Right:
					if cursorPos < buf.Len() {
						cursorPos += 1
					}
				case keys.Up:
					if historyPos > 0 {
						historyPos -= 1
						buf = bytes.NewBufferString(commandHistory[historyPos])
					} else if len(commandHistory) > 0 {
						buf = bytes.NewBufferString(commandHistory[0])
					}
					cursorPos = buf.Len()
				case keys.Down:
					if historyPos+1 < len(commandHistory) {
						historyPos += 1
						buf = bytes.NewBufferString(commandHistory[historyPos])
					} else if historyPos+1 == len(commandHistory) {
						historyPos += 1
						buf = bytes.NewBufferString("")
					}
					cursorPos = buf.Len()
				case keys.RuneKey, keys.Space: // Insert a character
					n := bytes.NewBuffer(buf.Bytes()[:cursorPos])
					a := bytes.NewBuffer(slices.Clone(buf.Bytes()[cursorPos:]))
					n.WriteByte(byte(key.Runes[0]))
					n.Write(a.Bytes())
					cursorPos += 1
					buf = n
				}
				fmt.Printf("> %s%s%s\n", string(buf.Bytes()[0:cursorPos]), cursorRune, string(buf.Bytes()[cursorPos:])) // Prints the line with cursor
				return false, nil
			})
			skipCommand := false
			switch line {
			case ":q", "exit", "exit()", "fuck":
				os.Exit(0)
			case "help":
				showHelp()
				skipCommand = true
			}
			if !skipCommand {
				tokens := GetTokens(line)
				sum := GetSum(tokens)
				if calculationSuccess {
					fmt.Println(strconv.FormatFloat(sum, 'f', -1, 64))
					lastAnswer = []float64{sum}
				} else {
					calculationSuccess = true
				}
			}

		}
	}
}
