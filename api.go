package UnitSqueezer

import (
	"context"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/opensibyl/UnitSqueezor/indexer"
	"github.com/opensibyl/UnitSqueezor/log"
	"github.com/opensibyl/UnitSqueezor/object"
	openapi "github.com/opensibyl/sibyl-go-client"
	"golang.org/x/exp/slices"
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
	conf := object.DefaultConfig()
	absSrcDir, err := filepath.Abs(conf.SrcDir)
	PanicIfErr(err)
	conf.SrcDir = absSrcDir
	apiClient, err := conf.NewSibylClient()
	PanicIfErr(err)

	// todo: start sibyl2 server

	// 1. upload and tag
	curIndexer, err := indexer.NewIndexer(&conf)
	err = curIndexer.UploadSrc(sharedContext)
	PanicIfErr(err)
	err = curIndexer.TagCases(apiClient, sharedContext)
	PanicIfErr(err)

	// 2. calc
	// line level diff
	extractor, err := NewDiffExtractor(&conf)
	PanicIfErr(err)
	diffMap, err := extractor.ExtractDiffMap()
	PanicIfErr(err)
	log.Log.Infof("diff map: %v", diffMap)

	// method level diff, and influence
	influenceTag := conf.GetReachTag()
	influencedMethods := make(map[string][]*FunctionWithState)
	for eachFile, eachLineList := range diffMap {
		eachLineStrList := make([]string, 0, len(eachLineList))
		for _, eachLine := range eachLineList {
			eachLineStrList = append(eachLineStrList, strconv.Itoa(eachLine))
		}
		functionWithSignatures, _, err := apiClient.BasicQueryApi.
			ApiV1FuncGet(sharedContext).
			Repo(conf.RepoInfo.Name).
			Rev(conf.RepoInfo.CommitId).
			File(eachFile).
			Lines(strings.Join(eachLineStrList, ",")).
			Execute()
		PanicIfErr(err)

		for _, eachFunc := range functionWithSignatures {
			eachFuncWithState := &FunctionWithState{
				ObjectFunctionWithSignature: &eachFunc,
				Reachable:                   false,
				ReachBy:                     make([]string, 0),
			}
			influencedMethods[eachFile] = append(influencedMethods[eachFile], eachFuncWithState)
		}
	}

	// reachable?
	for _, methods := range influencedMethods {
		for _, eachMethod := range methods {
			if slices.Contains(eachMethod.Tags, influenceTag) {
				log.Log.Infof("function reachable: %v", eachMethod.GetSignature())
				eachMethod.Reachable = true
			} else {
				log.Log.Infof("function changed but unreachable: %v", eachMethod.GetSignature())
				eachMethod.Reachable = false
			}
		}
	}
}

type FunctionWithState struct {
	*openapi.ObjectFunctionWithSignature

	Reachable bool
	ReachBy   []string
}
