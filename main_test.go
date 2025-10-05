package main

import (
	"testing"
)

// q is the query; e is the expected value
type test struct {
	q string
	e float64
}

// Test a few valid calculations
func TestValid(t *testing.T) {
	tests := []test{{q: "1 + 1 - 2", e: 0.0}, {q: "1 / 2", e: 0.5}, {q: "5 * (5 + 5)", e: 50.0}, {q: "2 * (3 + (2 + 6.1) * 4) + ( 2 // 9 ! ) + 200 ** 0 + -5 * (1 + 2) ! + 5 % 2", e: 645.1952191045343}, {q: "8 log 10", e: 0.9030899869919434}}
	for _, te := range tests {
		got := GetSum(GetTokens(te.q))
		if got != te.e {
			t.Errorf("Got %f; wanted %f", got, te.e)
		}
	}
}
