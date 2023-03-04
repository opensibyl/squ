package runner

import (
	"context"

	"github.com/opensibyl/UnitSqueezor/log"
	"github.com/opensibyl/UnitSqueezor/object"
	openapi "github.com/opensibyl/sibyl-go-client"
)

type BaseRunner struct {
	config    *object.SharedConfig
	apiClient *openapi.APIClient
}

func (baseRunner *BaseRunner) GetRelatedCases(ctx context.Context, diffFuncMap object.DiffFuncMap) ([]*openapi.ObjectFunctionWithSignature, error) {
	retMap := make(map[string]interface{})
	for _, eachFuncList := range diffFuncMap {
		for _, eachFunction := range eachFuncList {
			if !eachFunction.Reachable {
				log.Log.Warnf("func %v can not be reached by cases", eachFunction.GetName())
			}
			for _, eachRelatedCase := range eachFunction.ReachBy {
				log.Log.Infof("func %v reached by %v", eachFunction.GetName(), eachRelatedCase)
				retMap[eachRelatedCase] = nil
			}
		}
	}
	r := make([]*openapi.ObjectFunctionWithSignature, 0, len(retMap))
	for k := range retMap {
		cur, _, err := baseRunner.apiClient.
			SignatureQueryApi.
			ApiV1SignatureFuncGet(ctx).
			Repo(baseRunner.config.RepoInfo.RepoId).
			Rev(baseRunner.config.RepoInfo.RevHash).
			Signature(k).
			Execute()
		if err != nil {
			return nil, err
		}
		r = append(r, cur)
	}
	return r, nil
}
