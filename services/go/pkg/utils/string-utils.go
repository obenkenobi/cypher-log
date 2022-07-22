package utils

import "strings"

func StringIsBlank(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func StringIsNotBlank(s string) bool {
	return !StringIsBlank(s)
}
