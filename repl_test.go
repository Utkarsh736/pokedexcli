package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		// Edge cases - add these too!
		{
			input:    "HELLO",
			expected: []string{"hello"},
		},
		{
			input:    "  ",
			expected: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		
		// Check length first
		if len(actual) != len(c.expected) {
			t.Errorf("cleanInput(%q) returned slice of length %d, expected %d",
				c.input, len(actual), len(c.expected))
			continue // Skip word-by-word check if lengths don't match
		}
		
		// Check each word
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("cleanInput(%q) word %d = %q, expected %q",
					c.input, i, word, expectedWord)
			}
		}
	}
}

