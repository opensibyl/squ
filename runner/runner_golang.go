package runner

import (
	"fmt"
	"strings"

	"github.com/opensibyl/UnitSqueezor/object"
	openapi "github.com/opensibyl/sibyl-go-client"
)

type GoRunner struct {
	*BaseRunner
}

func NewGolangRunner(conf *object.SharedConfig) (Runner, error) {
	apiClient, err := conf.NewSibylClient()
	if err != nil {
		return nil, err
	}
	return &GoRunner{
		&BaseRunner{
			conf,
			apiClient,
		},
	}, nil
}

func (g *GoRunner) GetRunCommand(cases []*openapi.ObjectFunctionWithSignature) []string {
	// go test --run TestABC|TestDEF
	execCmdList := make([]string, 0, len(cases))
	for _, each := range cases {
		execCmdList = append(execCmdList, fmt.Sprintf("^%s$", each.GetName()))
	}
	caseRegex := strings.Join(execCmdList, "|")
	return []string{"go", "test", "--run", caseRegex, "-v"}
}
