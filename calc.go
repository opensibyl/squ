package UnitSqueezer

import (
	"context"
	"os/exec"
	"strconv"
	"strings"

	"github.com/opensibyl/UnitSqueezor/object"
	openapi "github.com/opensibyl/sibyl-go-client"
	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/ext"
	"golang.org/x/exp/slices"
)

type DiffMap = map[string][]int
type DiffFuncMap = map[string][]*FunctionWithState

type DiffExtractor interface {
	ExtractDiffMap(_ context.Context) (DiffMap, error)
	ExtractDiffMethods(ctx context.Context) (DiffFuncMap, error)
}

func NewDiffExtractor(config *object.SharedConfig) (DiffExtractor, error) {
	apiClient, err := config.NewSibylClient()
	if err != nil {
		return nil, err
	}
	return &GitExtractor{
		config,
		apiClient,
	}, nil
}

type GitExtractor struct {
	config    *object.SharedConfig
	apiClient *openapi.APIClient
}

func (g *GitExtractor) ExtractDiffMap(_ context.Context) (DiffMap, error) {
	gitDiffCmd := exec.Command("git", "diff", "HEAD~1", "HEAD")
	gitDiffCmd.Dir = g.config.SrcDir
	patchRaw, err := gitDiffCmd.CombinedOutput()
	if err != nil {
		core.Log.Errorf("git cmd error: %s", patchRaw)
		panic(err)
	}

	return ext.Unified2Affected(patchRaw)
}

func (g *GitExtractor) ExtractDiffMethods(ctx context.Context) (map[string][]*FunctionWithState, error) {
	diffMap, err := g.ExtractDiffMap(ctx)
	if err != nil {
		return nil, err
	}

	// method level diff, and influence
	influencedMethods := make(map[string][]*FunctionWithState)
	for eachFile, eachLineList := range diffMap {
		eachLineStrList := make([]string, 0, len(eachLineList))
		for _, eachLine := range eachLineList {
			eachLineStrList = append(eachLineStrList, strconv.Itoa(eachLine))
		}
		functionWithSignatures, _, err := g.apiClient.BasicQueryApi.
			ApiV1FuncGet(ctx).
			Repo(g.config.RepoInfo.Name).
			Rev(g.config.RepoInfo.CommitId).
			File(eachFile).
			Lines(strings.Join(eachLineStrList, ",")).
			Execute()
		PanicIfErr(err)

		for _, eachFunc := range functionWithSignatures {
			eachFuncWithState := &FunctionWithState{
				ObjectFunctionWithSignature: &eachFunc,
				Reachable:                   false,
				ReachBy:                     make([]string, 0),
			}
			influencedMethods[eachFile] = append(influencedMethods[eachFile], eachFuncWithState)
		}
	}

	// reachable?
	influenceTag := g.config.GetReachTag()
	for _, methods := range influencedMethods {
		for _, eachMethod := range methods {
			if slices.Contains(eachMethod.Tags, influenceTag) {
				eachMethod.Reachable = true
				// reach by whom
				for _, eachTag := range eachMethod.Tags {
					if !strings.Contains(eachTag, object.TagPrefixReachBy) {
						continue
					}
					caseSignature := strings.TrimPrefix(eachTag, object.TagPrefixReachBy)
					eachMethod.ReachBy = append(eachMethod.ReachBy, caseSignature)
				}
			} else {
				eachMethod.Reachable = false
			}
		}
	}
	return influencedMethods, nil
}
