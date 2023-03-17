package runner

import (
	"fmt"
	"strings"

	openapi "github.com/opensibyl/sibyl-go-client"
	"github.com/opensibyl/squ/object"
)

type PytestRunner struct {
	*BaseRunner
}

func (p *PytestRunner) GetRunCommand(cases []*openapi.ObjectFunctionWithSignature) string {
	// py.test tests_directory/foo.py tests_directory/bar.py -k 'test_001 or test_some_other_test'
	if len(cases) == 0 {
		return "--ignore-glob=\"*\""
	}

	execCmdList := make([]string, 0, len(cases))
	for _, each := range cases {
		execCmdList = append(execCmdList, fmt.Sprintf("%s", each.GetName()))
	}
	caseRegex := strings.Join(execCmdList, " and ")
	finalCaseStr := fmt.Sprintf("-k=\"%s\"", caseRegex)
	return finalCaseStr
}

func NewPytestRunner(conf *object.SharedConfig) (Runner, error) {
	apiClient, err := conf.NewSibylClient()
	if err != nil {
		return nil, err
	}
	return &PytestRunner{
		&BaseRunner{
			conf,
			apiClient,
		},
	}, nil
}
