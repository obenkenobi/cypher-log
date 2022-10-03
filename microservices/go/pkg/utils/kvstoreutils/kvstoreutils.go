package kvstoreutils

import "strings"

func CombineKeySections(keyStrings ...string) string {
	return strings.Join(keyStrings, "/")
}
