package runner

import (
	"context"

	"github.com/dominikbraun/graph"
	"github.com/opensibyl/UnitSqueezor/indexer"
	"github.com/opensibyl/UnitSqueezor/log"
	"github.com/opensibyl/UnitSqueezor/object"
	openapi "github.com/opensibyl/sibyl-go-client"
)

type BaseRunner struct {
	config    *object.SharedConfig
	apiClient *openapi.APIClient
}

func (baseRunner *BaseRunner) GetRelatedCases(_ context.Context, targetSignature string, indexer indexer.Indexer) (map[string]interface{}, error) {
	g := indexer.GetSibylCache().ReverseCallGraph
	vertexes := indexer.GetVertexesWithSignature(targetSignature)
	caseSet := indexer.GetCaseSet()
	matchedCases := make(map[string]interface{}, 0)
	for _, eachV := range vertexes {
		err := graph.BFS(g.Graph, eachV, func(k string) bool {
			functionWithPath, err := g.Graph.Vertex(k)
			if err != nil {
				return true
			}
			s := functionWithPath.GetSignature()

			if _, ok := caseSet[s]; ok {
				// reach
				matchedCases[s] = nil
			}
			return false
		})
		if err != nil {
			return nil, err
		}
	}
	log.Log.Infof("cases related to %v: %d", targetSignature, len(matchedCases))
	return matchedCases, nil
}

func (baseRunner *BaseRunner) Signature2Case(ctx context.Context, s string) (*openapi.ObjectFunctionWithSignature, error) {
	caseObject, _, err := baseRunner.apiClient.SignatureQueryApi.
		ApiV1SignatureFuncGet(ctx).
		Repo(baseRunner.config.RepoInfo.RepoId).
		Rev(baseRunner.config.RepoInfo.RevHash).
		Signature(s).
		Execute()
	if err != nil {
		return nil, err
	}
	return caseObject, nil
}

func (baseRunner *BaseRunner) Diff2Cases(ctx context.Context, diffMap object.DiffFuncMap, indexer indexer.Indexer) ([]*openapi.ObjectFunctionWithSignature, error) {
	casesToRunRaw := make(map[string]interface{})
	for fileName, eachFunctionList := range diffMap {
		log.Log.Infof("handle modified file: %s, functions: %d", fileName, len(eachFunctionList))
		for _, eachFunc := range eachFunctionList {
			cases, err := baseRunner.GetRelatedCases(ctx, eachFunc.GetSignature(), indexer)
			if err != nil {
				return nil, err
			}
			// merge
			for k := range cases {
				casesToRunRaw[k] = nil
			}
		}
	}
	casesToRun := make([]*openapi.ObjectFunctionWithSignature, 0)
	for eachCase := range casesToRunRaw {
		functionWithSignature, err := baseRunner.Signature2Case(ctx, eachCase)
		if err != nil {
			return nil, err
		}
		casesToRun = append(casesToRun, functionWithSignature)
	}
	return casesToRun, nil
}
