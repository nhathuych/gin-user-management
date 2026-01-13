package util

type APIResponse[T any] struct {
	Status     string      `json:"status"`
	Message    string      `json:"message,omitempty"`
	Data       T           `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}
