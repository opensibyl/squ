package runner

import (
	"context"
	"sync"

	"github.com/opensibyl/UnitSqueezor/log"
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
	var result sync.Map
	err := baseRunner.fillRelatedCases(ctx, targetSignature, &result, make([]string, 0))
	if err != nil {
		return nil, err
	}

	ret := make([]*openapi.ObjectFunctionWithSignature, 0)
	result.Range(func(key any, _ any) bool {
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

func (baseRunner *BaseRunner) getFuncBySignature(ctx context.Context, signature string) (*openapi.Sibyl2FunctionContextSlim, error) {
	if item, ok := cache.Load(signature); ok {
		return item.(*openapi.Sibyl2FunctionContextSlim), nil
	} else {
		ret, _, err := baseRunner.apiClient.SignatureQueryApi.
			ApiV1SignatureFuncctxGet(ctx).
			Repo(baseRunner.config.RepoInfo.RepoId).
			Rev(baseRunner.config.RepoInfo.RevHash).
			Signature(signature).
			Execute()
		if err != nil {
			return nil, err
		}
		cache.Store(signature, ret)
		return ret, nil
	}
}

func (baseRunner *BaseRunner) fillRelatedCases(ctx context.Context, targetSignature string, result *sync.Map, l []string) error {
	cur, err := baseRunner.getFuncBySignature(ctx, targetSignature)
	if err != nil {
		return err
	}
	// endpoint, store and return
	rcalls := cur.ReverseCalls
	if len(rcalls) == 0 {
		log.Log.Infof("reach endpoint: %v", targetSignature)
		result.Store(targetSignature, nil)
		return nil
	}
	// continue searching
	for _, each := range rcalls {
		// self call
		if each == targetSignature {
			continue
		}
		// loop call
		if slices.Contains(l, each) {
			continue
		}
		err := baseRunner.fillRelatedCases(ctx, each, result, append(l, each))
		if err != nil {
			return err
		}
	}
	return nil
}
