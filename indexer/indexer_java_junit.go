package indexer

import (
	"context"
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

	// tag cases
	for _, eachCaseMethod := range functionWithPaths {
		// tag all, and all their calls
		j.TagCaseInfluence(eachCaseMethod.GetSignature(), ctx)
	}

	return nil
}
