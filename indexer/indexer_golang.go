package indexer

import (
	"context"

	"github.com/opensibyl/UnitSqueezor/object"
	openapi "github.com/opensibyl/sibyl-go-client"
)

type GoIndexer struct {
	*BaseIndexer
}

func (i *GoIndexer) TagCases(ctx context.Context) error {
	repo := i.config.RepoInfo.Name
	rev := i.config.RepoInfo.CommitId

	functionWithPaths, _, _ := i.apiClient.RegexQueryApi.
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
		_, _ = i.apiClient.TagApi.ApiV1TagFuncPost(ctx).Payload(openapi.ServiceTagUpload{
			RepoId:    &repo,
			RevHash:   &rev,
			Signature: eachCaseMethod.Signature,
			Tag:       &tagCase,
		}).Execute()

		// tag all, and all their calls
		_ = i.TagCaseInfluence(eachCaseMethod.GetSignature(), eachCaseMethod.GetSignature(), ctx)
	}

	return nil
}
