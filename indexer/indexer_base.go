package indexer

import (
	"context"

	openapi "github.com/opensibyl/sibyl-go-client"
	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/cmd/sibyl/subs/upload"
	object2 "github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/opensibyl/squ/log"
	"github.com/opensibyl/squ/object"
)

type BaseIndexer struct {
	config        *object.SharedConfig
	apiClient     *openapi.APIClient
	caseSet       map[string]interface{}
	sibylCache    upload.ExecuteCacheMap
	vertexMapping map[string]*map[string]interface{}
}

func (baseIndexer *BaseIndexer) GetSibylCache() *sibyl2.FuncGraph {
	var curCache *upload.ExecuteCache
	for _, v := range baseIndexer.sibylCache {
		// only take the first one
		curCache = v
		break
	}
	return curCache.AnalyzeGraph
}

func (baseIndexer *BaseIndexer) GetCaseSet() map[string]interface{} {
	return baseIndexer.caseSet
}

func (baseIndexer *BaseIndexer) GetVertexesWithSignature(s string) []string {
	m := baseIndexer.vertexMapping[s]
	ret := make([]string, 0, len(*m))
	for k := range *m {
		ret = append(ret, k)
	}
	return ret
}

func (baseIndexer *BaseIndexer) UploadSrc(_ context.Context) error {
	conf := upload.DefaultConfig()
	conf.Src = baseIndexer.config.SrcDir
	conf.Url = baseIndexer.config.SibylUrl
	conf.RepoId = baseIndexer.config.RepoInfo.RepoId
	conf.RevHash = baseIndexer.config.RepoInfo.RevHash

	// todo: use ctx
	sibylCache, err := upload.ExecCurRevWithConfig(conf.Src, &object2.WorkspaceConfig{
		RepoId:  conf.RepoId,
		RevHash: conf.RevHash,
	}, conf)
	if err != nil {
		return err
	}
	baseIndexer.sibylCache = sibylCache
	cg := baseIndexer.GetSibylCache().CallGraph
	adjacencyMap, err := cg.AdjacencyMap()
	if err != nil {
		return err
	}
	for k := range adjacencyMap {
		functionWithPath, _ := cg.Vertex(k)
		signature := functionWithPath.GetSignature()

		m := baseIndexer.vertexMapping[signature]
		if m == nil {
			newM := make(map[string]interface{})
			baseIndexer.vertexMapping[signature] = &newM
			m = &newM
		}
		(*m)[k] = nil
	}
	log.Log.Infof("indexer done")
	return nil
}

func (baseIndexer *BaseIndexer) TagCase(caseSignature string, ctx context.Context) error {
	tagCase := object.TagCase
	repo := baseIndexer.config.RepoInfo.RepoId
	rev := baseIndexer.config.RepoInfo.RevHash
	_, err := baseIndexer.apiClient.TagApi.ApiV1TagFuncPost(ctx).Payload(openapi.ServiceTagUpload{
		RepoId:    &repo,
		RevHash:   &rev,
		Signature: &caseSignature,
		Tag:       &tagCase,
	}).Execute()
	if err != nil {
		return err
	}

	baseIndexer.caseSet[caseSignature] = nil
	return nil
}

func (baseIndexer *BaseIndexer) GetFuncWithSignature(ctx context.Context, s string) (*openapi.ObjectFunctionWithSignature, error) {
	caseObject, _, err := baseIndexer.apiClient.SignatureQueryApi.
		ApiV1SignatureFuncGet(ctx).
		Repo(baseIndexer.config.RepoInfo.RepoId).
		Rev(baseIndexer.config.RepoInfo.RevHash).
		Signature(s).
		Execute()
	if err != nil {
		return nil, err
	}
	return caseObject, nil
}
