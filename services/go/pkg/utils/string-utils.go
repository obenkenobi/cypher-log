package utils

import "strings"

// StringIsBlank checks if a string is empty or only has whitespaces.
func StringIsBlank(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// StringIsNotBlank checks if a string is not empty and has non-whitespace
// character(s).
func StringIsNotBlank(s string) bool {
	return !StringIsBlank(s)
}
