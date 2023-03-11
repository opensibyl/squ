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

	functionWithPaths, _, err := j.apiClient.RegexQueryApi.
		ApiV1RegexFuncGet(ctx).
		Repo(repo).
		Rev(rev).
		Field("extras.annotations").
		Regex(".*@Test.*").
		Execute()
	if err != nil {
		return err
	}

	// tag cases
	for _, eachCaseMethod := range functionWithPaths {
		err := j.TagCase(eachCaseMethod.GetSignature(), ctx)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	return nil
}
