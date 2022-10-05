package businessobjects

type AppSecretBo struct {
	Kid string
	Key []byte
}

func NewAppSecretBo(kid string, key []byte) AppSecretBo {
	return AppSecretBo{Kid: kid, Key: key}
}
