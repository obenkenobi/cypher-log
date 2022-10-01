package utils

func IdentityFunc[T any](v T) T {
	return v
}

func CastToAny[T any](v T) any {
	return v
}
