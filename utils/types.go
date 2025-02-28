package utils

type HTTPError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
