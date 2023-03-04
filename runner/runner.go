package runner

import (
	"context"
	"fmt"

	"github.com/opensibyl/UnitSqueezor/object"
	openapi "github.com/opensibyl/sibyl-go-client"
)

type Runner interface {
	GetRelatedCases(ctx context.Context, diffFuncMap object.DiffFuncMap) ([]*openapi.ObjectFunctionWithSignature, error)
	Run(cases []*openapi.ObjectFunctionWithSignature, ctx context.Context) error
}

func GetRunner(runnerType object.RunnerType, conf object.SharedConfig) (Runner, error) {
	switch runnerType {
	case object.RunnerGo:
		return NewGolangRunner(&conf)
	}
	return nil, fmt.Errorf("no runner type named: %v", runnerType)
}
