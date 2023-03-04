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
	Before   string    `json:"before"`
	After    string    `json:"after"`
}

func DefaultConfig() SharedConfig {
	return SharedConfig{
		".",
		nil,
		"http://127.0.0.1:9876",
		int(time.Now().UnixMicro()),
		"HEAD~1",
		"HEAD",
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

func (conf *SharedConfig) GetReachTag() string {
	return fmt.Sprintf("%s%d", TagPrefixReach, conf.BatchId)
}
