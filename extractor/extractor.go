package extractor

import (
	"context"

	"github.com/opensibyl/UnitSqueezor/object"
)

type DiffExtractor interface {
	ExtractDiffMap(_ context.Context) (object.DiffMap, error)
	ExtractDiffMethods(ctx context.Context) (object.DiffFuncMap, error)
}

func NewDiffExtractor(config *object.SharedConfig) (DiffExtractor, error) {
	apiClient, err := config.NewSibylClient()
	if err != nil {
		return nil, err
	}
	return &GitExtractor{
		config,
		apiClient,
	}, nil
}
