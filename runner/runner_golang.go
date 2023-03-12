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
	finalCaseStr := fmt.Sprintf("--run=\"%s\"", caseRegex)

	if g.config.CmdTemplate != "" {
		s := fmt.Sprintf(g.config.CmdTemplate, finalCaseStr)
		return strings.Split(s, " ")
	}
	// default commands
	return []string{"go", "test", "./...", "-v", finalCaseStr}
}
