package indexer

import (
	"context"

	"github.com/xxjwxc/gowp/workpool"
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
	wp := workpool.New(4)
	for _, each := range functionWithPaths {
		wp.Do(func() error {
			err := j.TagCase(each.GetSignature(), ctx)
			if err != nil {
				return err
			}
			return nil
		})
	}
	err = wp.Wait()
	if err != nil {
		return err
	}
	return nil
}
