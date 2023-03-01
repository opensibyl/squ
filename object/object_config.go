package object

import (
	"fmt"
	"net/url"
	"time"

	openapi "github.com/opensibyl/sibyl-go-client"
)

type SharedConfig struct {
	SrcDir   string    `json:"srcDir"`
	RepoInfo *RepoInfo `json:"repoInfo"`
	SibylUrl string    `json:"sibylUrl"`
	BatchId  int       `json:"batchId"`
}

func DefaultConfig() SharedConfig {
	return SharedConfig{
		".",
		nil,
		"http://127.0.0.1:9876",
		int(time.Now().UnixMicro()),
	}
}

func (conf *SharedConfig) NewSibylClient() (*openapi.APIClient, error) {
	parsed, err := url.Parse(conf.SibylUrl)
	if err != nil {
		return nil, err
	}

	configuration := openapi.NewConfiguration()
	configuration.Scheme = parsed.Scheme
	configuration.Host = parsed.Host
	return openapi.NewAPIClient(configuration), nil
}

func (conf *SharedConfig) GetInfluenceTag() string {
	return fmt.Sprintf("%s%d", TagPrefixInfluence, conf.BatchId)
}
