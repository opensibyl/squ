package indexer

import (
	"context"

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
	conf.RepoId = baseIndexer.config.RepoInfo.Name
	conf.RevHash = baseIndexer.config.RepoInfo.CommitId
	upload.ExecWithConfig(conf)
	return nil
}

func (baseIndexer *BaseIndexer) TagCaseInfluence(caseSignature string, signature string, ctx context.Context) error {
	// if batch id changed, will recalc
	tagReach := baseIndexer.config.GetReachTag()
	tagReachBy := object.TagPrefixReachBy + caseSignature

	repo := baseIndexer.config.RepoInfo.Name
	rev := baseIndexer.config.RepoInfo.CommitId

	// tag itself
	_, _ = baseIndexer.apiClient.TagApi.ApiV1TagFuncPost(ctx).Payload(openapi.ServiceTagUpload{
		RepoId:    &repo,
		RevHash:   &rev,
		Signature: &signature,
		Tag:       &tagReach,
	}).Execute()
	_, _ = baseIndexer.apiClient.TagApi.ApiV1TagFuncPost(ctx).Payload(openapi.ServiceTagUpload{
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
		_ = baseIndexer.TagCaseInfluence(caseSignature, each, ctx)
	}
	return nil
}
