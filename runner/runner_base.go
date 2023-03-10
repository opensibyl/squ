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
	result := make(map[string][]string)
	err := baseRunner.fillRelatedCases(ctx, targetSignature, result, make([]string, 0))
	if err != nil {
		return nil, err
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
			log.Log.Infof("found related case: %v", functionWithSignature)
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

func (baseRunner *BaseRunner) fillRelatedCases(ctx context.Context, targetSignature string, endpointRecord map[string][]string, paths []string) error {
	cur, err := baseRunner.getFuncBySignature(ctx, targetSignature)
	if err != nil {
		return err
	}

	// endpoint, return
	reversedCalls := cur.ReverseCalls
	newPaths := append(paths, targetSignature)
	if len(reversedCalls) == 0 {
		log.Log.Infof("reach endpoint: %v %v", targetSignature, newPaths)
		endpointRecord[targetSignature] = []string{targetSignature}
		return nil
	}

	// performance issue here
	if len(reversedCalls) > 10 {
		log.Log.Infof("function referenced %d times, give up", len(reversedCalls))
		return nil
	}

	// else, on path, continue searching
	curEndpoints := make([]string, 0)
	for _, each := range reversedCalls {
		// self call
		if each == targetSignature {
			continue
		}
		if slices.Contains(newPaths, each) {
			continue
		}
		err := baseRunner.fillRelatedCases(ctx, each, endpointRecord, newPaths)
		if err != nil {
			return err
		}

		if item, ok := endpointRecord[each]; ok {
			curEndpoints = append(curEndpoints, item...)
		}
	}
	endpointRecord[targetSignature] = curEndpoints
	return nil
}
