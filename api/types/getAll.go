package types

type GetAllDTO[T any] struct {
	Count uint64 `json:"count"`
	Items []T    `json:"items"`
}
