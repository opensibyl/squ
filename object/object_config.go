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
	GraphOutput string      `json:"svgOutput"`
	Dry         bool        `json:"dry"`
	IndexerType IndexerType `json:"indexerType"`
	RunnerType  RunnerType  `json:"runnerType"`
	CmdTemplate string      `json:"cmdTemplate"`
}

func DefaultConfig() SharedConfig {
	return SharedConfig{
		".",
		nil,
		"http://127.0.0.1:9875",
		int(time.Now().UnixMicro()),
		"HEAD~1",
		"HEAD",
		"",
		"",
		false,
		IndexerGolang,
		RunnerGolang,
		"",
	}
}

func (conf *SharedConfig) NewSibylClient() (*openapi.APIClient, error) {
	parsed, err := conf.parseSibylUrl()
	if err != nil {
		return nil, err
	}

	configuration := openapi.NewConfiguration()
	configuration.Scheme = parsed.Scheme
	configuration.Host = parsed.Host
	return openapi.NewAPIClient(configuration), nil
}

func (conf *SharedConfig) LocalSibyl() bool {
	parsed, err := conf.parseSibylUrl()
	if err != nil {
		return false
	}
	hostName := parsed.Hostname()
	if hostName == "127.0.0.1" || hostName == "localhost" {
		return true
	}
	return false
}

func (conf *SharedConfig) GetSibylPort() string {
	parsed, err := conf.parseSibylUrl()
	if err != nil {
		return ""
	}
	return parsed.Port()
}

func (conf *SharedConfig) parseSibylUrl() (*url.URL, error) {
	return url.Parse(conf.SibylUrl)
}
