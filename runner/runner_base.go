package runner

import (
	"context"
	"sync"

	"github.com/opensibyl/UnitSqueezor/object"
	openapi "github.com/opensibyl/sibyl-go-client"
	"golang.org/x/exp/slices"
)

var cache sync.Map

type BaseRunner struct {
	config    *object.SharedConfig
	apiClient *openapi.APIClient
}

func (baseRunner *BaseRunner) GetRelatedCases(ctx context.Context, targetSignature string) ([]*openapi.ObjectFunctionWithSignature, error) {
	// todo: slow in big repo
	var endpoints sync.Map
	var walked sync.Map
	err := baseRunner.fillRelatedCases(ctx, targetSignature, &endpoints, &walked)
	if err != nil {
		return nil, err
	}

	ret := make([]*openapi.ObjectFunctionWithSignature, 0)
	endpoints.Range(func(key any, _ any) bool {
		// todo: multi request here
		functionWithSignature, _, err := baseRunner.apiClient.SignatureQueryApi.
			ApiV1SignatureFuncGet(ctx).
			Repo(baseRunner.config.RepoInfo.RepoId).
			Rev(baseRunner.config.RepoInfo.RevHash).
			Signature(targetSignature).
			Execute()
		if err != nil {
			return false
		}
		// it's a case
		if slices.Contains(functionWithSignature.GetTags(), object.TagCase) {
			ret = append(ret, functionWithSignature)
		}
		return true
	})

	return ret, nil
}

func (baseRunner *BaseRunner) fillRelatedCases(ctx context.Context, targetSignature string, m *sync.Map, walked *sync.Map) error {
	var rcalls []string
	if item, ok := cache.Load(targetSignature); ok {
		ctxSlim := item.([]string)
		rcalls = ctxSlim
	} else {
		item, _, err := baseRunner.apiClient.SignatureQueryApi.
			ApiV1SignatureFuncctxGet(ctx).
			Repo(baseRunner.config.RepoInfo.RepoId).
			Rev(baseRunner.config.RepoInfo.RevHash).
			Signature(targetSignature).
			Execute()
		if err != nil {
			return err
		}
		rcalls = item.GetReverseCalls()
		cache.Store(targetSignature, rcalls)
	}

	walked.Store(targetSignature, nil)
	if len(rcalls) == 0 {
		m.Store(targetSignature, nil)
		return nil
	}
	for _, each := range rcalls {
		if _, ok := walked.Load(each); ok {
			continue
		}
		go baseRunner.fillRelatedCases(ctx, each, m, walked)
	}
	return nil
}
