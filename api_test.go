package squ

import (
	"testing"

	"github.com/opensibyl/squ/object"
)

func TestGolang(t *testing.T) {
	conf := object.DefaultConfig()
	conf.Dry = true
	conf.GraphOutput = "abc.txt"
	MainFlow(conf)
}

func TestMaven(t *testing.T) {
	t.Skip()
	conf := object.DefaultConfig()
	conf.Dry = true
	conf.SrcDir = "../jacoco"
	conf.Before = "HEAD~5"
	conf.IndexerType = object.IndexerJavaJUnit
	conf.RunnerType = object.RunnerMaven
	conf.JsonOutput = "./output.json"
	MainFlow(conf)
}
