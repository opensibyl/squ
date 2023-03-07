package indexer

import (
	"context"
	"errors"
	"fmt"

	"github.com/dominikbraun/graph"
	"github.com/opensibyl/UnitSqueezor/log"
	"github.com/opensibyl/UnitSqueezor/object"
	openapi "github.com/opensibyl/sibyl-go-client"
	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/cmd/sibyl/subs/upload"
)

type BaseIndexer struct {
	config                *object.SharedConfig
	apiClient             *openapi.APIClient
	graphCache            *sibyl2.FuncGraph
	graphSignatureMapping map[string][]string
	caseTagCache          object.CaseTagCache
}

func (baseIndexer *BaseIndexer) UploadSrc(_ context.Context) error {
	conf := upload.DefaultConfig()
	conf.Src = baseIndexer.config.SrcDir
	conf.Url = baseIndexer.config.SibylUrl
	conf.RepoId = baseIndexer.config.RepoInfo.RepoId
	conf.RevHash = baseIndexer.config.RepoInfo.RevHash

	// todo: use ctx
	err := upload.ExecWithConfig(conf)
	if err != nil {
		return err
	}

	// cache the graph
	baseIndexer.graphCache = conf.BizContext.GraphCache
	if err != nil {
		return fmt.Errorf("no graph cache found after upload")
	}
	// for mapping graph key and sibyl key
	baseIndexer.graphSignatureMapping = make(map[string][]string, 0)
	cg := baseIndexer.graphCache.CallGraph
	adjacencyMap, err := cg.AdjacencyMap()
	if err != nil {
		return err
	}
	for k := range adjacencyMap {
		functionWithPath, _ := cg.Vertex(k)
		signature := functionWithPath.GetSignature()
		baseIndexer.graphSignatureMapping[signature] = append(baseIndexer.graphSignatureMapping[signature], k)
	}
	// init
	baseIndexer.caseTagCache = make(object.CaseTagCache)
	return nil
}

func (baseIndexer *BaseIndexer) TagCaseInfluence(caseSignature string, _ context.Context) error {
	cases, ok := baseIndexer.graphSignatureMapping[caseSignature]
	if !ok {
		return errors.New("no case found: " + caseSignature)
	}
	g := baseIndexer.graphCache.CallGraph.Graph
	tagMap := make(map[string]interface{})
	for _, each := range cases {
		err := graph.BFS(g, each, func(k string) bool {
			functionWithPath, err := g.Vertex(k)
			if err != nil {
				return false
			}
			tagMap[functionWithPath.GetSignature()] = nil
			return false
		})
		if err != nil {
			return err
		}
	}

	log.Log.Infof("case %s influence: %d", caseSignature, len(tagMap))
	baseIndexer.caseTagCache[caseSignature] = tagMap
	return nil
}

func (baseIndexer *BaseIndexer) GetTagMap() *object.CaseTagCache {
	return &baseIndexer.caseTagCache
}
