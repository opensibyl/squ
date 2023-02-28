package UnitSqueezer

import (
	"context"
)

type Indexer interface {
	UploadSrc(ctx context.Context) error
	// TagCases different framework should have different rules
	TagCases(ctx context.Context) error
}

func NewIndexer(config *SharedConfig) (*GoIndexer, error) {
	repoInfo, err := GetRepoInfoFromDir(config.SrcDir)
	config.RepoInfo = repoInfo
	if err != nil {
		return nil, err
	}

	return &GoIndexer{
		config: config,
	}, nil
}
