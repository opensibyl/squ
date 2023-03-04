package indexer

import (
	"context"
	"fmt"
	"sync"

	"github.com/opensibyl/UnitSqueezor/object"
)

type Indexer interface {
	UploadSrc(ctx context.Context) error
	// TagCases different framework should have different rules
	TagCases(ctx context.Context) error
	TagCaseInfluence(caseSignature string, taggedSet *sync.Map, signature string, ctx context.Context) error
}

func GetIndexer(indexerType object.IndexerType, config *object.SharedConfig) (Indexer, error) {
	repoInfo, err := object.GetRepoInfoFromDir(config.SrcDir)
	config.RepoInfo = repoInfo
	if err != nil {
		return nil, err
	}

	client, _ := config.NewSibylClient()
	if err != nil {
		return nil, err
	}

	baseIndexer := &BaseIndexer{
		config:    config,
		apiClient: client,
	}
	switch indexerType {
	case object.IndexerGolang:
		return &GoIndexer{
			baseIndexer,
		}, nil
	case object.IndexerJavaJUnit:
		return &JavaJunitIndexer{
			baseIndexer,
		}, nil
	}
	return nil, fmt.Errorf("no indexer named: %v", indexerType)
}
