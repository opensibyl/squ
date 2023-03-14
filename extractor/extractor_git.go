package extractor

import (
	"context"
	"os/exec"
	"strconv"
	"strings"

	openapi "github.com/opensibyl/sibyl-go-client"
	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/ext"
	"github.com/opensibyl/squ/log"
	"github.com/opensibyl/squ/object"
)

type GitExtractor struct {
	config    *object.SharedConfig
	apiClient *openapi.APIClient
}

func (g *GitExtractor) ExtractDiffMap(_ context.Context) (object.DiffMap, error) {
	gitDiffCmd := exec.Command("git", "diff", g.config.Before, g.config.After)
	gitDiffCmd.Dir = g.config.SrcDir
	patchRaw, err := gitDiffCmd.CombinedOutput()
	if err != nil {
		core.Log.Errorf("git cmd error: %s", patchRaw)
		panic(err)
	}

	return ext.Unified2Affected(patchRaw)
}

func (g *GitExtractor) ExtractDiffMethods(ctx context.Context) (map[string][]*object.FunctionWithState, error) {
	diffMap, err := g.ExtractDiffMap(ctx)
	if err != nil {
		return nil, err
	}

	// method level diff, and influence
	influencedMethods := make(map[string][]*object.FunctionWithState)
	for eachFile, eachLineList := range diffMap {
		eachLineStrList := make([]string, 0, len(eachLineList))
		for _, eachLine := range eachLineList {
			eachLineStrList = append(eachLineStrList, strconv.Itoa(eachLine))
		}
		functionWithSignatures, _, err := g.apiClient.BasicQueryApi.
			ApiV1FuncGet(ctx).
			Repo(g.config.RepoInfo.RepoId).
			Rev(g.config.RepoInfo.RevHash).
			File(eachFile).
			Lines(strings.Join(eachLineStrList, ",")).
			Execute()
		if err != nil {
			return nil, err
		}
		log.Log.Infof("%s %v => functions %d", eachFile, eachLineList, len(functionWithSignatures))
		for _, eachFunc := range functionWithSignatures {
			eachFuncWithState := &object.FunctionWithState{
				ObjectFunctionWithSignature: eachFunc,
				Reachable:                   false,
				ReachBy:                     make([]string, 0),
			}
			influencedMethods[eachFile] = append(influencedMethods[eachFile], eachFuncWithState)
		}
	}
	return influencedMethods, nil
}
