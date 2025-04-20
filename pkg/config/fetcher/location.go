package fetcher

import (
	"context"
	"fmt"

	"os"

	"github.com/mercor/payment-service/pkg/config/model"
)

type native struct {
	location string
}

func NewNativeFetcher(ctx context.Context, location string) (Fetcher, error) {
	nf := &native{
		location: location,
	}

	return nf, nil
}

func (nf *native) GetConfig(ctx context.Context) (*model.Config, error) {
	data, err := os.ReadFile(nf.location)
	if err != nil {
		return nil, fmt.Errorf("error while reading local file: %w", err)
	}

	configMap, err := ParseYAMLToConfigMap(string(data))
	if err != nil {
		return nil, err
	}

	return model.NewConfig(configMap), nil
}
