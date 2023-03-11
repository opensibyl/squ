package indexer

import (
	"context"
	"fmt"

	"github.com/opensibyl/UnitSqueezor/object"
	"github.com/opensibyl/sibyl2"
)

type Indexer interface {
	UploadSrc(ctx context.Context) error
	GetCaseSet() map[string]interface{}
	GetSibylCache() *sibyl2.FuncGraph
	GetVertexesWithSignature(s string) []string
	// TagCases different framework should have different rules
	TagCases(ctx context.Context) error
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
		config:        config,
		apiClient:     client,
		caseSet:       make(map[string]interface{}),
		vertexMapping: make(map[string]*map[string]interface{}),
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
