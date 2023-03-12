package indexer

import (
	"context"
	"os/signal"
	"syscall"
	"testing"

	"github.com/opensibyl/sibyl2/pkg/server"
	object2 "github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/opensibyl/squ/log"
	"github.com/opensibyl/squ/object"
	"github.com/stretchr/testify/assert"
)

func TestIndexerBase(t *testing.T) {
	ctx := context.Background()
	sibylContext, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		config := object2.DefaultExecuteConfig()
		// for performance
		config.BindingConfigPart.DbType = object2.DriverTypeInMemory
		config.EnableLog = false
		err := server.Execute(config, sibylContext)
		assert.Nil(t, err)
	}()
	defer stop()
	log.Log.Infof("sibyl2 backend ready")

	conf := object.DefaultConfig()
	conf.SrcDir = "../"
	curIndexer, err := GetIndexer(object.IndexerGolang, &conf)

	assert.Nil(t, err)
	err = curIndexer.UploadSrc(ctx)
	assert.Nil(t, err)
	err = curIndexer.TagCases(ctx)
	assert.Nil(t, err)
}
