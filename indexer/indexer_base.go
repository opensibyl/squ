package indexer

import (
	"context"
	"sync"

	"github.com/opensibyl/UnitSqueezor/object"
	openapi "github.com/opensibyl/sibyl-go-client"
	"github.com/opensibyl/sibyl2/cmd/sibyl/subs/upload"
	"golang.org/x/exp/slices"
)

var cache sync.Map
var l sync.Mutex

type CaseTagCache = map[string]*map[string]interface{}

type BaseIndexer struct {
	config       *object.SharedConfig
	apiClient    *openapi.APIClient
	tagCache     CaseTagCache
	specialCases []string
}

func (baseIndexer *BaseIndexer) GetTagCache() CaseTagCache {
	return baseIndexer.tagCache
}

func (baseIndexer *BaseIndexer) GetSpecialCases() []string {
	return baseIndexer.specialCases
}

func (baseIndexer *BaseIndexer) UploadSrc(_ context.Context) error {
	conf := upload.DefaultConfig()
	conf.Src = baseIndexer.config.SrcDir
	conf.Url = baseIndexer.config.SibylUrl
	conf.RepoId = baseIndexer.config.RepoInfo.RepoId
	conf.RevHash = baseIndexer.config.RepoInfo.RevHash

	// todo: use ctx
	err := upload.ExecWithConfig(conf)
	if err != nil {
		return err
	}
	return nil
}

func (baseIndexer *BaseIndexer) getFuncBySignature(ctx context.Context, signature string) (*openapi.Sibyl2FunctionContextSlim, error) {
	if item, ok := cache.Load(signature); ok {
		return item.(*openapi.Sibyl2FunctionContextSlim), nil
	} else {
		ret, _, err := baseIndexer.apiClient.SignatureQueryApi.
			ApiV1SignatureFuncctxGet(ctx).
			Repo(baseIndexer.config.RepoInfo.RepoId).
			Rev(baseIndexer.config.RepoInfo.RevHash).
			Signature(signature).
			Execute()
		if err != nil {
			return nil, err
		}
		cache.Store(signature, ret)
		return ret, nil
	}
}

func (baseIndexer *BaseIndexer) TagCase(caseSignature string, ctx context.Context) error {
	// if batch id changed, will recalc
	tagCase := object.TagCase
	repo := baseIndexer.config.RepoInfo.RepoId
	rev := baseIndexer.config.RepoInfo.RevHash
	_, err := baseIndexer.apiClient.TagApi.ApiV1TagFuncPost(ctx).Payload(openapi.ServiceTagUpload{
		RepoId:    &repo,
		RevHash:   &rev,
		Signature: &caseSignature,
		Tag:       &tagCase,
	}).Execute()
	if err != nil {
		return err
	}

	// query and store to cache
	_, err = baseIndexer.fillTagCache(caseSignature, caseSignature, make([]string, 0), ctx)
	if err != nil {
		return err
	}
	return nil
}

func (baseIndexer *BaseIndexer) fillTagCache(caseSignature string, curSignature string, paths []string, ctx context.Context) (bool, error) {
	functionContextSlim, err := baseIndexer.getFuncBySignature(ctx, curSignature)
	if err != nil {
		return false, err
	}
	calls := functionContextSlim.Calls
	// end
	if len(calls) == 0 {
		return true, nil
	}
	// avoid too large scale
	if item, ok := baseIndexer.tagCache[caseSignature]; ok {
		if len(*item) > 128 {
			baseIndexer.specialCases = append(baseIndexer.specialCases, caseSignature)
			return false, nil
		}
	}

	// loop
	if slices.Contains(paths, curSignature) {
		return true, nil
	}

	// accept this node
	l.Lock()
	m := baseIndexer.tagCache[caseSignature]
	if m == nil {
		newM := make(map[string]interface{})
		newM[curSignature] = nil
		baseIndexer.tagCache[caseSignature] = &newM
	} else {
		(*m)[curSignature] = nil
	}
	l.Unlock()

	// go on
	newPaths := append(paths, curSignature)
	for _, each := range calls {
		needContinue, err := baseIndexer.fillTagCache(caseSignature, each, newPaths, ctx)
		if err != nil {
			return false, err
		}
		if !needContinue {
			return false, nil
		}
	}
	return true, nil
}

func (baseIndexer *BaseIndexer) TagCases(_ context.Context) error {
	// TODO implement me
	panic("implement me")
}
