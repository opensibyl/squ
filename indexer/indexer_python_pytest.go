package indexer

import "context"

type PythonPytestIndexer struct {
	*BaseIndexer
}

func (p *PythonPytestIndexer) TagCases(ctx context.Context) error {
	repo := p.config.RepoInfo.RepoId
	rev := p.config.RepoInfo.RevHash

	functionWithPaths, _, err := p.apiClient.RegexQueryApi.
		ApiV1RegexFuncGet(ctx).
		Repo(repo).
		Rev(rev).
		Field("name").
		Regex("^test_.*").
		Execute()
	if err != nil {
		return err
	}

	// tag cases
	for _, eachCaseMethod := range functionWithPaths {
		err := p.TagCase(eachCaseMethod.GetSignature(), ctx)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	return nil
}
