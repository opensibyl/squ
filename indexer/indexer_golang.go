package indexer

import (
	"context"
	"sync"

	"github.com/opensibyl/UnitSqueezor/object"
	openapi "github.com/opensibyl/sibyl-go-client"
)

type GoIndexer struct {
	*BaseIndexer
}

func (i *GoIndexer) TagCases(ctx context.Context) error {
	repo := i.config.RepoInfo.RepoId
	rev := i.config.RepoInfo.RevHash

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
	var taggedMap sync.Map
	for _, eachCaseMethod := range functionWithPaths {
		// all the errors from tag will be ignored
		_, _ = i.apiClient.TagApi.ApiV1TagFuncPost(ctx).Payload(openapi.ServiceTagUpload{
			RepoId:    &repo,
			RevHash:   &rev,
			Signature: eachCaseMethod.Signature,
			Tag:       &tagCase,
		}).Execute()

		// tag all, and all their calls
		go i.TagCaseInfluence(eachCaseMethod.GetSignature(), &taggedMap, eachCaseMethod.GetSignature(), ctx)
	}

	return nil
}
