package object

import (
	"net/url"
	"time"

	openapi "github.com/opensibyl/sibyl-go-client"
)

type SharedConfig struct {
	SrcDir      string      `json:"srcDir"`
	RepoInfo    *RepoInfo   `json:"repoInfo"`
	SibylUrl    string      `json:"sibylUrl"`
	BatchId     int         `json:"batchId"`
	Before      string      `json:"before"`
	After       string      `json:"after"`
	JsonOutput  string      `json:"jsonOutput"`
	Dry         bool        `json:"dry"`
	IndexerType IndexerType `json:"indexerType"`
	RunnerType  RunnerType  `json:"runnerType"`
	ScanLimit   int         `json:"scanLimit"`
}

func DefaultConfig() SharedConfig {
	return SharedConfig{
		".",
		nil,
		"http://127.0.0.1:9876",
		int(time.Now().UnixMicro()),
		"HEAD~1",
		"HEAD",
		"",
		false,
		IndexerGolang,
		RunnerGolang,
		128,
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
