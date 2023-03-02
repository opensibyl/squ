package UnitSqueezer

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/opensibyl/UnitSqueezor/object"
	openapi "github.com/opensibyl/sibyl-go-client"
)

type Executor interface {
	GetRelatedCases(ctx context.Context, diffFuncMap DiffFuncMap) ([]*openapi.ObjectFunctionWithSignature, error)
	Execute(cases []*openapi.ObjectFunctionWithSignature, ctx context.Context) error
}

type BaseExecutor struct {
	config    *object.SharedConfig
	apiClient *openapi.APIClient
}

func (baseExecutor *BaseExecutor) GetRelatedCases(ctx context.Context, diffFuncMap DiffFuncMap) ([]*openapi.ObjectFunctionWithSignature, error) {
	retMap := make(map[string]interface{})
	for _, eachFuncList := range diffFuncMap {
		for _, eachFunction := range eachFuncList {
			for _, eachRelatedCase := range eachFunction.ReachBy {
				retMap[eachRelatedCase] = nil
			}
		}
	}
	r := make([]*openapi.ObjectFunctionWithSignature, 0, len(retMap))
	for k := range retMap {
		cur, _, err := baseExecutor.apiClient.
			SignatureQueryApi.
			ApiV1SignatureFuncGet(ctx).
			Repo(baseExecutor.config.RepoInfo.Name).
			Rev(baseExecutor.config.RepoInfo.CommitId).
			Signature(k).
			Execute()
		if err != nil {
			return nil, err
		}
		r = append(r, cur)
	}
	return r, nil
}

func NewGoExecutor(conf *object.SharedConfig) (Executor, error) {
	apiClient, err := conf.NewSibylClient()
	if err != nil {
		return nil, err
	}
	return &GoExecutor{
		&BaseExecutor{
			conf,
			apiClient,
		},
	}, nil
}

type GoExecutor struct {
	*BaseExecutor
}

func (g *GoExecutor) Execute(cases []*openapi.ObjectFunctionWithSignature, ctx context.Context) error {
	execCmdList := make([]string, 0, len(cases))
	for _, each := range cases {
		execCmdList = append(execCmdList, fmt.Sprintf("^%s$", each.GetName()))
	}
	caseRegex := strings.Join(execCmdList, "|")
	goTestCmd := exec.CommandContext(ctx, "go", "test", "--run", caseRegex)
	goTestCmd.Dir = g.config.SrcDir
	goTestCmd.Stdout = os.Stdout
	goTestCmd.Stderr = os.Stderr
	err := goTestCmd.Run()
	if err != nil {
		return err
	}
	return nil
}
