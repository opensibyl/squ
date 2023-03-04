package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
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

func (g *GoRunner) Run(cases []*openapi.ObjectFunctionWithSignature, ctx context.Context) error {
	// go test --run TestABC|TestDEF
	execCmdList := make([]string, 0, len(cases))
	for _, each := range cases {
		execCmdList = append(execCmdList, fmt.Sprintf("^%s$", each.GetName()))
	}
	caseRegex := strings.Join(execCmdList, "|")
	goTestCmd := exec.CommandContext(ctx, "go", "test", "--run", caseRegex, "-v")
	goTestCmd.Dir = g.config.SrcDir
	goTestCmd.Stdout = os.Stdout
	goTestCmd.Stderr = os.Stderr
	err := goTestCmd.Run()
	if err != nil {
		return err
	}
	return nil
}
