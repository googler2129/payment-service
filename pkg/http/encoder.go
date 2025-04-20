package http

import "io"

type Encoder interface {
	Encode(v interface{}) (io.Reader, error)
	Decode(b io.ReadCloser, v interface{}) error
}
