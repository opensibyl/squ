package runner

import (
	"fmt"
	"strings"

	openapi "github.com/opensibyl/sibyl-go-client"
	"github.com/opensibyl/squ/object"
)

type GoRunner struct {
	*BaseRunner
}

func (g *GoRunner) GetRunCommand(cases []*openapi.ObjectFunctionWithSignature) string {
	// go test --run=TestABC|TestDEF
	if len(cases) == 0 {
		return "--run=\"^$\""
	}

	execCmdList := make([]string, 0, len(cases))
	for _, each := range cases {
		execCmdList = append(execCmdList, fmt.Sprintf("^%s$", each.GetName()))
	}
	caseRegex := strings.Join(execCmdList, "|")
	finalCaseStr := fmt.Sprintf("--run=\"%s\"", caseRegex)
	return finalCaseStr
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
