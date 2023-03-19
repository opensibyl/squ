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
	// pytest test/test_api.py::test_analyse test/test_video.py::test_read_from_file
	// https://docs.pytest.org/en/7.1.x/how-to/usage.html#specifying-which-tests-to-run
	if len(cases) == 0 {
		return "--ignore-glob=\"*\""
	}

	execCmdList := make([]string, 0, len(cases))
	for _, each := range cases {
		// use relative path
		relPath := strings.TrimPrefix(each.GetPath(), p.config.SrcDir)
		eachStr := fmt.Sprintf("%s::%s", relPath, each.GetName())
		execCmdList = append(execCmdList, eachStr)
	}
	finalStr := strings.Join(execCmdList, " ")
	return finalStr
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
