package UnitSqueezer

import (
	"context"
	"path/filepath"
)

/*
index
1. upload all the files to sibyl2
2. search and tag all the test methods
3. calc and tag all the test methods influencing scope

calc
1. calc diff between current and previous
2. find methods influenced by diff
3. search related cases

execute (can be implemented by different languages
1. build test commands for different frameworks
2. call cmd
*/

type SharedConfig struct {
	SrcDir   string    `json:"srcDir"`
	RepoInfo *RepoInfo `json:"repoInfo"`
	SibylUrl string    `json:"sibylUrl"`
}

func DefaultConfig() SharedConfig {
	return SharedConfig{
		".",
		nil,
		"http://127.0.0.1:9876",
	}
}

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func MainFlow() {
	conf := DefaultConfig()
	absSrcDir, err := filepath.Abs(conf.SrcDir)
	PanicIfErr(err)
	conf.SrcDir = absSrcDir

	sharedContext := context.Background()

	indexer, err := NewIndexer(&conf)
	err = indexer.UploadSrc(sharedContext)
	PanicIfErr(err)
}
