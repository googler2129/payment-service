package http

import "net/http"

type ResponseParams struct {
	Headers    http.Header
	StatusCode int
}
