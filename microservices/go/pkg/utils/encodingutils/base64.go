package encodingutils

import "encoding/base64"

func DecodeBase64String(src string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(src)
}

func EncodeBase64String(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}
