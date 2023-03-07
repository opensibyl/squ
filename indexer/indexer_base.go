package indexer

import (
	"context"

	"github.com/opensibyl/UnitSqueezor/object"
	openapi "github.com/opensibyl/sibyl-go-client"
	"github.com/opensibyl/sibyl2/cmd/sibyl/subs/upload"
)

type BaseIndexer struct {
	config    *object.SharedConfig
	apiClient *openapi.APIClient
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
	return nil
}

func (baseIndexer *BaseIndexer) TagCase(caseSignature string, ctx context.Context) error {
	// if batch id changed, will recalc
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
	return nil
}
