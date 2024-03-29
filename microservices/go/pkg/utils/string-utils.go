package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// StringIsBlank checks if a string is empty or only has whitespaces.
func StringIsBlank(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// StringIsNotBlank checks if a string is not empty and has non-whitespace
// character(s).
func StringIsNotBlank(s string) bool {
	return !StringIsBlank(s)
}

func StringFirstNChars(s string, n int) string {
	i := 0
	for j := range s {
		if i == n {
			return s[:j]
		}
		i++
	}
	return s
}

func Int64ToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}

func StrToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func AnyToString(val any) string {
	if val == nil {
		return ""
	}
	if str, ok := val.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", val)
}
