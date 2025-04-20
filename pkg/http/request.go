package http

import "time"

type Request struct {
	Url         string
	Body        interface{}         // the request body, should be json serializable
	QueryParams map[string][]string // the request query url params
	Headers     map[string][]string // to set any custom headers, if any
	PathParams  map[string]string
	Timeout     time.Duration
}
