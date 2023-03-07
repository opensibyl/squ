package indexer

import (
	"context"
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

	// tag cases
	for _, eachCaseMethod := range functionWithPaths {
		err := i.TagCase(eachCaseMethod.GetSignature(), ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
