package kfka

type Retry[T any] struct {
	Value T   `json:"value"`
	Tries int `json:"tries"`
}
