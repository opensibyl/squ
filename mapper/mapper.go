package mapper

import (
	"context"

	openapi "github.com/opensibyl/sibyl-go-client"
	"github.com/opensibyl/squ/indexer"
	"github.com/opensibyl/squ/object"
)

type Mapper interface {
	SetIndexer(indexer.Indexer)
	GetRelatedCaseSignatures(ctx context.Context, targetSignature string) (map[string]interface{}, error)
	Diff2Cases(ctx context.Context, diffMap object.DiffFuncMap) ([]*openapi.ObjectFunctionWithSignature, error)
}

func NewMapper() Mapper {
	return &BaseMapper{}
}
