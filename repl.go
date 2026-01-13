package main

import (
	"strings"
)

func cleanInput(text string) []string {
	// Step 1: Convert to lowercase
	lowered := strings.ToLower(text)

	// Step 2: Trim leading/trailing whitespace
	trimmed := strings.TrimSpace(lowered)

	// Step 3: Split on whitespace
	words := strings.Fields(trimmed)

	return words
}
