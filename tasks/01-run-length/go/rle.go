package rle

import "fmt"

// Encode returns the run‑length encoding of UTF‑8 string s.
//
// "AAB" → "A2B1"
func Encode(s string) string {

	if len(s) == 0 {
		return ""
	}

	result := ""
	count := 0
	runes := []rune(s) // to handle special characters
	currentChar := runes[0]

	for i := 0; i < len(runes); i++ {
		if runes[i] == currentChar {
			count++
		} else {
			result += string(currentChar) + fmt.Sprintf("%d", count)
			currentChar = runes[i]
			count = 1
		}
	}

	result += string(currentChar) + fmt.Sprintf("%d", count)

	return result
}
