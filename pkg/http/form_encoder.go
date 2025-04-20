package http

import (
	"encoding/json"
	"github.com/ajg/form"
	"io"
	"strings"
)

type formEncoder struct {
}

func (f *formEncoder) Encode(v interface{}) (io.Reader, error) {
	values, err := form.EncodeToValues(v)
	if err != nil {
		return nil, err
	}
	return strings.NewReader(values.Encode()), nil
}

func (f *formEncoder) Decode(b io.ReadCloser, v interface{}) error {
	return json.NewDecoder(b).Decode(v)
}
