package runner

import (
	"context"
	"fmt"

	openapi "github.com/opensibyl/sibyl-go-client"
	"github.com/opensibyl/squ/object"
)

type BaseRunnerPart interface {
	Run(command []string, ctx context.Context) error
}

type Runner interface {
	BaseRunnerPart
	GetRunCommand(cases []*openapi.ObjectFunctionWithSignature) []string
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
