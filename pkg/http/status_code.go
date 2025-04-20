package http

import "net/http"

type StatusCode int

func (s StatusCode) Code() int {
	return int(s)
}

// String returns a text string for the HTTP status code.
// If mapping not found, then returns http.StatusText
func (s StatusCode) String() string {
	statusText, exists := StatusCodeToStringMap[s]
	if !exists {
		return http.StatusText(s.Code())
	}

	return statusText
}

// Is2xx method returns true if HTTP status `code >= 200 and <= 299` otherwise false.
func (s StatusCode) Is2xx() bool {
	return s.Code() > 199 && s.Code() < 300
}

// Is3xx method returns true if HTTP status `code >= 300 and <= 399` otherwise false.
func (s StatusCode) Is3xx() bool {
	return s.Code() > 299 && s.Code() < 400
}

// Is4xx method returns true if HTTP status `code >= 400 and <= 499` otherwise false.
func (s StatusCode) Is4xx() bool {
	return s.Code() > 399 && s.Code() < 500
}

// Is5xx method returns true if HTTP status `code >= 500 and <= 599` otherwise false.
func (s StatusCode) Is5xx() bool {
	return s.Code() > 499 && s.Code() < 600
}

const (
	StatusOK      StatusCode = 200
	StatusCreated StatusCode = 201

	StatusMovedPermanently StatusCode = 301
	StatusFound            StatusCode = 302

	StatusBadRequest          StatusCode = 400
	StatusUnauthorized        StatusCode = 401
	StatusPaymentRequired     StatusCode = 402
	StatusForbidden           StatusCode = 403
	StatusNotFound            StatusCode = 404
	StatusRequestTimeout      StatusCode = 408
	StatusUnprocessableEntity StatusCode = 422
	StatusTooManyRequests     StatusCode = 429

	StatusInternalServerError StatusCode = 500
	StatusNotImplemented      StatusCode = 501
	StatusBadGateway          StatusCode = 502
)

var StatusCodeToStringMap = map[StatusCode]string{
	StatusBadRequest:          "Invalid Request",
	StatusInternalServerError: "Something went wrong",
	StatusUnauthorized:        "You don't have access to this action",
	StatusOK:                  "Success",
	StatusForbidden:           "Validation failed",
	StatusUnprocessableEntity: "Invalid Request",
	StatusRequestTimeout:      "Request Timeout",
}

type APIMethod string

func (s APIMethod) String() string {
	return string(s)
}

const (
	APIGet     APIMethod = "GET"
	APIPost    APIMethod = "POST"
	APIPut     APIMethod = "PUT"
	APIDelete  APIMethod = "DELETE"
	APIPatch   APIMethod = "PATCH"
	APIHead    APIMethod = "HEAD"
	APIOptions APIMethod = "OPTIONS"
)
