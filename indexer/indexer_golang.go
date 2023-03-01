package indexer

import (
	"context"

	"github.com/opensibyl/UnitSqueezor/log"
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
	for _, each := range functionWithPaths {
		// all the errors from tag will be ignored
		_, _ = apiClient.TagApi.ApiV1TagFuncPost(ctx).Payload(openapi.ServiceTagUpload{
			RepoId:    &repo,
			RevHash:   &rev,
			Signature: each.Signature,
			Tag:       &tagCase,
		}).Execute()

		// tag all, and all their calls
		_ = i.TagCaseInfluence(apiClient, each.GetSignature(), ctx)
	}

	return nil
}

func (i *GoIndexer) TagCaseInfluence(apiClient *openapi.APIClient, signature string, ctx context.Context) error {
	// if batch id changed, will recalc
	tagInfluence := i.config.GetInfluenceTag()

	repo := i.config.RepoInfo.Name
	rev := i.config.RepoInfo.CommitId

	// tag itself
	_, _ = apiClient.TagApi.ApiV1TagFuncPost(ctx).Payload(openapi.ServiceTagUpload{
		RepoId:    &repo,
		RevHash:   &rev,
		Signature: &signature,
		Tag:       &tagInfluence,
	}).Execute()
	log.Log.Infof("tag influence: %v", signature)

	functionContext, _, _ := apiClient.SignatureQueryApi.
		ApiV1SignatureFuncctxGet(ctx).
		Repo(repo).
		Rev(rev).
		Signature(signature).
		Execute()
	for _, each := range functionContext.Calls {
		_ = i.TagCaseInfluence(apiClient, each, ctx)
	}
	return nil
}
