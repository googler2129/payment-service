package http

import (
	"bytes"
	"encoding/json"
	"io"
)

type JsonEncoder struct {
}

func (m *JsonEncoder) Encode(v interface{}) (io.Reader, error) {
	s, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(s), nil
}

func (m *JsonEncoder) Decode(b io.ReadCloser, v interface{}) error {
	return json.NewDecoder(b).Decode(v)
}
