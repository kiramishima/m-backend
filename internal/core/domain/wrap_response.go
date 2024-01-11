package domain

type WrapResponse[T any] struct {
	Data T `json:"data"`
}
