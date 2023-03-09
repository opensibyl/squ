package runner

import (
	"context"
	"sync"

	"github.com/opensibyl/UnitSqueezor/object"
	openapi "github.com/opensibyl/sibyl-go-client"
	"golang.org/x/exp/slices"
)

var cache sync.Map
var reachCache sync.Map

type BaseRunner struct {
	config    *object.SharedConfig
	apiClient *openapi.APIClient
}

func (baseRunner *BaseRunner) GetRelatedCases(ctx context.Context, targetSignature string) ([]*openapi.ObjectFunctionWithSignature, error) {
	result := make(map[string]interface{})
	resultList, err := baseRunner.fillRelatedCases(ctx, targetSignature, make([]string, 0))
	if err != nil {
		return nil, err
	}
	for _, each := range resultList {
		result[each] = nil
	}

	ret := make([]*openapi.ObjectFunctionWithSignature, 0)
	for eachEndpoint := range result {
		functionWithSignature, _, err := baseRunner.apiClient.SignatureQueryApi.
			ApiV1SignatureFuncGet(ctx).
			Repo(baseRunner.config.RepoInfo.RepoId).
			Rev(baseRunner.config.RepoInfo.RevHash).
			Signature(eachEndpoint).
			Execute()
		if err != nil {
			return nil, err
		}
		// it's a case
		if slices.Contains(functionWithSignature.GetTags(), object.TagCase) {
			ret = append(ret, functionWithSignature)
		}
	}

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

func (baseRunner *BaseRunner) fillRelatedCases(ctx context.Context, targetSignature string, l []string) ([]string, error) {
	cur, err := baseRunner.getFuncBySignature(ctx, targetSignature)
	if err != nil {
		return nil, err
	}
	// endpoint, store and return
	reversedCalls := cur.ReverseCalls

	// path, continue searching
	v, ok := reachCache.Load(targetSignature)
	if ok {
		return v.([]string), nil
	}

	endpoints := make([]string, 0)
	for _, each := range reversedCalls {
		// self call
		if each == targetSignature {
			continue
		}
		// loop call
		if slices.Contains(l, each) {
			continue
		}
		subEndpoints, err := baseRunner.fillRelatedCases(ctx, each, append(l, each))
		if err != nil {
			return nil, err
		}
		endpoints = append(endpoints, subEndpoints...)
	}
	// store to cache
	reachCache.Store(targetSignature, endpoints)
	return endpoints, nil
}
