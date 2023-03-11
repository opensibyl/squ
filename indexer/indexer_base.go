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
var tagCacheLock sync.Mutex

type CaseTagCache = map[string]*map[string]interface{}

type BaseIndexer struct {
	config      *object.SharedConfig
	apiClient   *openapi.APIClient
	tagCache    CaseTagCache
	giveUpCases *sync.Map
}

func (baseIndexer *BaseIndexer) GetTagCache() CaseTagCache {
	return baseIndexer.tagCache
}

func (baseIndexer *BaseIndexer) GetGiveUpCases() *sync.Map {
	return baseIndexer.giveUpCases
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
	err = baseIndexer.fillTagCache(caseSignature, caseSignature, make([]string, 0), ctx)
	if err != nil {
		return err
	}
	return nil
}

func (baseIndexer *BaseIndexer) fillTagCache(caseSignature string, curSignature string, paths []string, ctx context.Context) error {
	// already give up?
	if _, alreadyGiveUp := baseIndexer.giveUpCases.Load(caseSignature); alreadyGiveUp {
		return nil
	}

	functionContextSlim, err := baseIndexer.getFuncBySignature(ctx, curSignature)
	if err != nil {
		return err
	}
	calls := functionContextSlim.Calls
	// endpoint
	if len(calls) == 0 {
		return nil
	}

	// avoid too large scale
	tagCacheLock.Lock()
	if item, ok := baseIndexer.GetTagCache()[caseSignature]; ok {
		if len(*item) > baseIndexer.config.ScanLimit {
			baseIndexer.GetGiveUpCases().Store(caseSignature, nil)
			tagCacheLock.Unlock()
			return nil
		}
	}
	tagCacheLock.Unlock()

	// loop
	if slices.Contains(paths, curSignature) {
		return nil
	}

	// accept this node
	tagCacheLock.Lock()
	m := baseIndexer.GetTagCache()[caseSignature]
	if m == nil {
		newM := make(map[string]interface{})
		newM[curSignature] = nil
		baseIndexer.GetTagCache()[caseSignature] = &newM
	} else {
		(*m)[curSignature] = nil
	}
	tagCacheLock.Unlock()

	// go on
	newPaths := append(paths, curSignature)
	for _, each := range calls {
		err := baseIndexer.fillTagCache(caseSignature, each, newPaths, ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (baseIndexer *BaseIndexer) TagCases(_ context.Context) error {
	// TODO implement me
	panic("implement me")
}
