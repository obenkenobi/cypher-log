package encodingutils

import "encoding/base64"

func DecodeBase64(src []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(src))
}

func EncodeBase64(src []byte) []byte {
	dst := base64.StdEncoding.EncodeToString(src)
	return []byte(dst)
}
