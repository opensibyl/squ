package UnitSqueezer

import (
	"context"
	"net/url"

	openapi "github.com/opensibyl/sibyl-go-client"
	"github.com/opensibyl/sibyl2/cmd/sibyl/subs/upload"
)

type GoIndexer struct {
	config *SharedConfig
}

func (i *GoIndexer) UploadSrc(_ context.Context) error {
	conf := upload.DefaultConfig()
	conf.Src = i.config.SrcDir
	conf.Url = i.config.SibylUrl
	conf.RepoId = i.config.RepoInfo.Name
	conf.RevHash = i.config.RepoInfo.CommitId
	upload.ExecWithConfig(conf)
	return nil
}

func (i *GoIndexer) TagCases(ctx context.Context) error {
	parsed, err := url.Parse(i.config.SibylUrl)
	if err != nil {
		return err
	}

	configuration := openapi.NewConfiguration()
	configuration.Scheme = parsed.Scheme
	configuration.Host = parsed.Host
	apiClient := openapi.NewAPIClient(configuration)

	functionWithPaths, _, _ := apiClient.RegexQueryApi.
		ApiV1RegexFuncctxGet(ctx).
		Repo(i.config.RepoInfo.Name).
		Rev(i.config.RepoInfo.CommitId).
		Field("name").
		Regex("^Test.*").
		Execute()

	for _, each := range functionWithPaths {
		Log.Infof("f: %v %v\n", *each.Name, each.Calls)
	}
	// tag all, and all their calls
	return nil
}
