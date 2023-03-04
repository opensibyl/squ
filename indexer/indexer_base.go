package indexer

import (
	"context"
	"sync"

	"github.com/opensibyl/UnitSqueezor/object"
	openapi "github.com/opensibyl/sibyl-go-client"
	"github.com/opensibyl/sibyl2/cmd/sibyl/subs/upload"
)

type BaseIndexer struct {
	config    *object.SharedConfig
	apiClient *openapi.APIClient
}

func (baseIndexer *BaseIndexer) UploadSrc(_ context.Context) error {
	conf := upload.DefaultConfig()
	conf.Src = baseIndexer.config.SrcDir
	conf.Url = baseIndexer.config.SibylUrl
	conf.RepoId = baseIndexer.config.RepoInfo.RepoId
	conf.RevHash = baseIndexer.config.RepoInfo.RevHash
	upload.ExecWithConfig(conf)
	return nil
}

func (baseIndexer *BaseIndexer) TagCaseInfluence(caseSignature string, taggedSet *sync.Map, signature string, ctx context.Context) error {
	// if batch id changed, will recalc
	tagReach := baseIndexer.config.GetReachTag()
	tagReachBy := object.TagPrefixReachBy + caseSignature

	repo := baseIndexer.config.RepoInfo.RepoId
	rev := baseIndexer.config.RepoInfo.RevHash

	// tag itself
	if _, tagged := taggedSet.Load(signature); tagged {
		// stop
		return nil
	} else {
		taggedSet.Store(signature, nil)
	}

	go baseIndexer.apiClient.TagApi.ApiV1TagFuncPost(ctx).Payload(openapi.ServiceTagUpload{
		RepoId:    &repo,
		RevHash:   &rev,
		Signature: &signature,
		Tag:       &tagReach,
	}).Execute()
	go baseIndexer.apiClient.TagApi.ApiV1TagFuncPost(ctx).Payload(openapi.ServiceTagUpload{
		RepoId:    &repo,
		RevHash:   &rev,
		Signature: &signature,
		Tag:       &tagReachBy,
	}).Execute()

	functionContext, _, _ := baseIndexer.apiClient.SignatureQueryApi.
		ApiV1SignatureFuncctxGet(ctx).
		Repo(repo).
		Rev(rev).
		Signature(signature).
		Execute()
	for _, each := range functionContext.Calls {
		go baseIndexer.TagCaseInfluence(caseSignature, taggedSet, each, ctx)
	}
	return nil
}
