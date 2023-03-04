package indexer

import (
	"context"
	"sync"

	"github.com/opensibyl/UnitSqueezor/object"
	openapi "github.com/opensibyl/sibyl-go-client"
)

type JavaJunitIndexer struct {
	*BaseIndexer
}

func (j *JavaJunitIndexer) TagCases(ctx context.Context) error {
	repo := j.config.RepoInfo.RepoId
	rev := j.config.RepoInfo.RevHash

	functionWithPaths, _, _ := j.apiClient.RegexQueryApi.
		ApiV1RegexFuncGet(ctx).
		Repo(repo).
		Rev(rev).
		Field("extras.annotations").
		Regex(".*@Test.*").
		Execute()

	// case is case, will not change
	tagCase := object.TagCase
	// tag cases
	var taggedMap sync.Map
	for _, eachCaseMethod := range functionWithPaths {
		// all the errors from tag will be ignored
		_, _ = j.apiClient.TagApi.ApiV1TagFuncPost(ctx).Payload(openapi.ServiceTagUpload{
			RepoId:    &repo,
			RevHash:   &rev,
			Signature: eachCaseMethod.Signature,
			Tag:       &tagCase,
		}).Execute()

		// tag all, and all their calls
		go j.TagCaseInfluence(eachCaseMethod.GetSignature(), &taggedMap, eachCaseMethod.GetSignature(), ctx)
	}

	return nil
}
