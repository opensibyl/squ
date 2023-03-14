package log

import (
	"path/filepath"

	"github.com/opensibyl/squ/object"
	"go.uber.org/zap"
)

var Log *zap.SugaredLogger

func InitLogger(config object.SharedConfig) {
	conf := zap.NewProductionConfig()
	if config.DebugMode {
		conf.OutputPaths = []string{
			filepath.Join(config.SrcDir, "squ-debug.log"),
		}
	} else {
		// else, quiet
		conf.Level.SetLevel(zap.ErrorLevel)
	}
	logger, _ := conf.Build()
	defer logger.Sync() // flushes buffer, if any
	Log = logger.Sugar()
}

func init() {
	InitLogger(object.DefaultConfig())
}
