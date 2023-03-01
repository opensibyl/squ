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

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func MainFlow() {
	sharedContext := context.Background()

	// init config
	conf := DefaultConfig()
	absSrcDir, err := filepath.Abs(conf.SrcDir)
	PanicIfErr(err)
	conf.SrcDir = absSrcDir
	apiClient, err := conf.NewSibylClient()
	PanicIfErr(err)

	// todo: start sibyl2 server

	// 1. upload and tag
	indexer, err := NewIndexer(&conf)
	err = indexer.UploadSrc(sharedContext)
	PanicIfErr(err)
	err = indexer.TagCases(apiClient, sharedContext)
	PanicIfErr(err)
}
