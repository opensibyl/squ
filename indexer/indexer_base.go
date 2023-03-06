package indexer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/dominikbraun/graph"
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
}

func (baseIndexer *BaseIndexer) UploadSrc(_ context.Context) error {
	conf := upload.DefaultConfig()
	conf.Src = baseIndexer.config.SrcDir
	conf.Url = baseIndexer.config.SibylUrl
	conf.RepoId = baseIndexer.config.RepoInfo.RepoId
	conf.RevHash = baseIndexer.config.RepoInfo.RevHash
	err := upload.ExecWithConfig(conf)
	if err != nil {
		return err
	}
	baseIndexer.graphCache = conf.BizContext.GraphCache
	if err != nil {
		return fmt.Errorf("no graph cache found after upload")
	}

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
	return nil
}

func (baseIndexer *BaseIndexer) TagCaseInfluence(caseSignature string, taggedSet *sync.Map, signature string, ctx context.Context) error {
	// if batch id changed, will recalc
	tagReach := baseIndexer.config.GetReachTag()
	tagReachBy := object.TagPrefixReachBy + caseSignature

	repo := baseIndexer.config.RepoInfo.RepoId
	rev := baseIndexer.config.RepoInfo.RevHash

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

	// tag
	for k := range tagMap {
		go func() {
			baseIndexer.apiClient.TagApi.ApiV1TagFuncPost(ctx).Payload(openapi.ServiceTagUpload{
				RepoId:    &repo,
				RevHash:   &rev,
				Signature: &k,
				Tag:       &tagReach,
			}).Execute()
			baseIndexer.apiClient.TagApi.ApiV1TagFuncPost(ctx).Payload(openapi.ServiceTagUpload{
				RepoId:    &repo,
				RevHash:   &rev,
				Signature: &k,
				Tag:       &tagReachBy,
			}).Execute()
		}()
	}
	return nil
}
