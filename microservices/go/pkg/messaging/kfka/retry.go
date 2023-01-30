package kfka

type RetryDto[T any] struct {
	Value T   `json:"value"`
	Tries int `json:"tries"`
}

func CreateRetryDto[T any](value T, tries int) RetryDto[T] {
	return RetryDto[T]{
		Value: value,
		Tries: tries,
	}
}
