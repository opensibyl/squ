package runner

import (
	"context"

	"github.com/opensibyl/UnitSqueezor/object"
	openapi "github.com/opensibyl/sibyl-go-client"
)

type Runner interface {
	GetRelatedCases(ctx context.Context, diffFuncMap object.DiffFuncMap) ([]*openapi.ObjectFunctionWithSignature, error)
	Run(cases []*openapi.ObjectFunctionWithSignature, ctx context.Context) error
}
