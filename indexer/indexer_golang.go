package indexer

import (
	"context"

	"github.com/opensibyl/UnitSqueezor/object"
	openapi "github.com/opensibyl/sibyl-go-client"
	"github.com/opensibyl/sibyl2/cmd/sibyl/subs/upload"
)

type GoIndexer struct {
	config *object.SharedConfig
}

func (i *GoIndexer) UploadSrc(_ context.Context) error {
	conf := upload.DefaultConfig()
	conf.Src = i.config.SrcDir
	conf.Url = i.config.SibylUrl
	conf.RepoId = i.config.RepoInfo.Name
	conf.RevHash = i.config.RepoInfo.CommitId
	upload.ExecWithConfig(conf)
	return nil
}

func (i *GoIndexer) TagCases(apiClient *openapi.APIClient, ctx context.Context) error {
	repo := i.config.RepoInfo.Name
	rev := i.config.RepoInfo.CommitId

	functionWithPaths, _, _ := apiClient.RegexQueryApi.
		ApiV1RegexFuncGet(ctx).
		Repo(repo).
		Rev(rev).
		Field("name").
		Regex("^Test.*").
		Execute()

	// case is case, will not change
	tagCase := object.TagCase
	// tag cases
	for _, eachCaseMethod := range functionWithPaths {
		// all the errors from tag will be ignored
		_, _ = apiClient.TagApi.ApiV1TagFuncPost(ctx).Payload(openapi.ServiceTagUpload{
			RepoId:    &repo,
			RevHash:   &rev,
			Signature: eachCaseMethod.Signature,
			Tag:       &tagCase,
		}).Execute()

		// tag all, and all their calls
		_ = i.TagCaseInfluence(apiClient, eachCaseMethod.GetSignature(), eachCaseMethod.GetSignature(), ctx)
	}

	return nil
}

func (i *GoIndexer) TagCaseInfluence(apiClient *openapi.APIClient, caseSignature string, signature string, ctx context.Context) error {
	// if batch id changed, will recalc
	tagReach := i.config.GetReachTag()
	tagReachBy := object.TagPrefixReachBy + caseSignature

	repo := i.config.RepoInfo.Name
	rev := i.config.RepoInfo.CommitId

	// tag itself
	_, _ = apiClient.TagApi.ApiV1TagFuncPost(ctx).Payload(openapi.ServiceTagUpload{
		RepoId:    &repo,
		RevHash:   &rev,
		Signature: &signature,
		Tag:       &tagReach,
	}).Execute()
	_, _ = apiClient.TagApi.ApiV1TagFuncPost(ctx).Payload(openapi.ServiceTagUpload{
		RepoId:    &repo,
		RevHash:   &rev,
		Signature: &signature,
		Tag:       &tagReachBy,
	}).Execute()

	functionContext, _, _ := apiClient.SignatureQueryApi.
		ApiV1SignatureFuncctxGet(ctx).
		Repo(repo).
		Rev(rev).
		Signature(signature).
		Execute()
	for _, each := range functionContext.Calls {
		_ = i.TagCaseInfluence(apiClient, caseSignature, each, ctx)
	}
	return nil
}
