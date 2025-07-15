package models

type HTTPRequest struct{}

type HTTPResponse struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
	Error   bool   `json:"error"`
}

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
