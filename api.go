// Copyright 2023 williamfzc (https://github.com/williamfzc), opensibyl team (https://github.com/opensibyl)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package squ

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/dominikbraun/graph/draw"
	"github.com/goccy/go-graphviz"
	openapi "github.com/opensibyl/sibyl-go-client"
	"github.com/opensibyl/sibyl2/pkg/server"
	object2 "github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/opensibyl/squ/extractor"
	"github.com/opensibyl/squ/indexer"
	"github.com/opensibyl/squ/log"
	"github.com/opensibyl/squ/mapper"
	"github.com/opensibyl/squ/object"
	"github.com/opensibyl/squ/runner"
)

/*
index (preparation)
1. upload all the files to sibyl2
2. search and tag all the test methods
3. calc and tag all the test methods influencing scope

extract (extract data from workspace)
1. calc diff between current and previous
2. find methods influenced by diff

mapper (mapping between cases and diff)
1. search related cases

runner (can be implemented by different languages)
1. build test commands for different frameworks
2. call cmd
*/

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func MainFlow(conf object.SharedConfig) {
	// init config
	absSrcDir, err := filepath.Abs(conf.SrcDir)
	PanicIfErr(err)
	conf.SrcDir = absSrcDir

	rootContext := context.Background()
	rootContext, cancel := context.WithCancel(rootContext)
	defer cancel()

	rootStartTime := time.Now()

	// 0. start sibyl2 backend
	if conf.LocalSibyl() {
		log.Log.Infof("using local sibyl, starting ...")
		portNum, err := strconv.Atoi(conf.GetSibylPort())
		PanicIfErr(err)
		go func() {
			config := object2.DefaultExecuteConfig()
			// for performance
			config.BindingConfigPart.DbType = object2.DriverTypeInMemory
			config.EnableLog = false
			config.Port = portNum
			err := server.Execute(config, rootContext)
			PanicIfErr(err)
		}()
		log.Log.Infof("sibyl2 backend ready")
	} else {
		log.Log.Infof("using remote sibyl, skip starting server")
	}

	// 1. index
	// todo: if using remote sibyl server, pull data directly
	curIndexer, err := indexer.GetIndexer(conf.IndexerType, &conf)
	err = curIndexer.UploadSrc(rootContext)
	PanicIfErr(err)
	err = curIndexer.TagCases(rootContext)
	PanicIfErr(err)
	log.Log.Infof("indexer ready")

	// 2. extract
	// line level diff
	calcContext, cancel := context.WithCancel(rootContext)
	defer cancel()
	curExtractor, err := extractor.NewDiffExtractor(&conf)
	PanicIfErr(err)
	diffMap, err := curExtractor.ExtractDiffMethods(calcContext)
	PanicIfErr(err)
	log.Log.Infof("diff calc ready: %v", len(diffMap))

	// 3. mapper
	curMapper := mapper.NewMapper()
	curMapper.SetIndexer(curIndexer)
	casesToRun, err := curMapper.Diff2Cases(rootContext, diffMap)
	PanicIfErr(err)
	log.Log.Infof("case analyzer done, before: %d, after: %d", len(curIndexer.GetCaseSet()), len(casesToRun))

	if conf.JsonOutput != "" {
		o := &Output{}
		o.DiffMap = diffMap
		o.Cases = casesToRun
		diffFuncOutputBytes, err := json.Marshal(o)
		PanicIfErr(err)
		err = os.WriteFile(conf.JsonOutput, diffFuncOutputBytes, os.ModePerm)
		PanicIfErr(err)
	}
	if conf.GraphOutput != "" {
		buf := bytes.NewBuffer([]byte{})
		err = draw.DOT(curIndexer.GetSibylCache().CallGraph.Graph, buf)
		PanicIfErr(err)

		// todo: node attr should be applied here
		g, err := graphviz.ParseBytes(buf.Bytes())
		PanicIfErr(err)
		gviz := graphviz.New()
		err = gviz.RenderFilename(g, graphviz.Format(graphviz.DOT), conf.GraphOutput)
		PanicIfErr(err)
	}

	// 4. runner
	log.Log.Infof("runner scope")
	if len(casesToRun) == 0 {
		log.Log.Infof("no cases need to run")
		return
	}
	curRunner, err := runner.GetRunner(conf.RunnerType, conf)
	PanicIfErr(err)

	prepareTotalCost := time.Since(rootStartTime)
	log.Log.Infof("prepare stage finished, total cost: %d ms", prepareTotalCost.Milliseconds())
	cmd := curRunner.GetRunCommand(casesToRun)
	log.Log.Infof("runner cmd: %v", cmd)
	if !conf.Dry {
		log.Log.Infof("start running cases: %v", len(casesToRun))
		_ = curRunner.Run(cmd, rootContext)
	}
}

type Output struct {
	DiffMap object.DiffFuncMap                     `json:"diff"`
	Cases   []*openapi.ObjectFunctionWithSignature `json:"cases"`
}
