package mapper

import (
	"context"

	"github.com/opensibyl/UnitSqueezor/indexer"
	"github.com/opensibyl/UnitSqueezor/object"
	openapi "github.com/opensibyl/sibyl-go-client"
)

type Mapper interface {
	SetIndexer(indexer.Indexer)
	GetRelatedCaseSignatures(ctx context.Context, targetSignature string) (map[string]interface{}, error)
	Diff2Cases(ctx context.Context, diffMap object.DiffFuncMap) ([]*openapi.ObjectFunctionWithSignature, error)
}

func NewMapper() Mapper {
	return &BaseMapper{}
}
