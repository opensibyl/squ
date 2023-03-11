package runner

import (
	"context"
	"fmt"

	"github.com/opensibyl/UnitSqueezor/object"
	openapi "github.com/opensibyl/sibyl-go-client"
)

type Runner interface {
	GetRunCommand(cases []*openapi.ObjectFunctionWithSignature) []string
	Run(cases []*openapi.ObjectFunctionWithSignature, ctx context.Context) error
}

func GetRunner(runnerType object.RunnerType, conf object.SharedConfig) (Runner, error) {
	switch runnerType {
	case object.RunnerGolang:
		return NewGolangRunner(&conf)
	case object.RunnerMaven:
		return NewMavenRunner(&conf)
	}
	return nil, fmt.Errorf("no runner type named: %v", runnerType)
}
