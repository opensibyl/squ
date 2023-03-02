package indexer

import (
	"context"
	"sync"

	"github.com/opensibyl/UnitSqueezor/object"
)

type Indexer interface {
	UploadSrc(ctx context.Context) error
	// TagCases different framework should have different rules
	TagCases(ctx context.Context) error
	TagCaseInfluence(caseSignature string, taggedSet *sync.Map, signature string, ctx context.Context) error
}

func NewIndexer(config *object.SharedConfig) (Indexer, error) {
	repoInfo, err := object.GetRepoInfoFromDir(config.SrcDir)
	config.RepoInfo = repoInfo
	if err != nil {
		return nil, err
	}

	client, _ := config.NewSibylClient()
	if err != nil {
		return nil, err
	}

	return &GoIndexer{
		&BaseIndexer{
			config:    config,
			apiClient: client,
		},
	}, nil
}
