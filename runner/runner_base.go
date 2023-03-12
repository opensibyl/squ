package runner

import (
	"context"
	"os"
	"os/exec"

	openapi "github.com/opensibyl/sibyl-go-client"
	"github.com/opensibyl/squ/object"
)

type BaseRunner struct {
	config    *object.SharedConfig
	apiClient *openapi.APIClient
}

func (b *BaseRunner) Run(command []string, ctx context.Context) error {
	realCmd := exec.CommandContext(ctx, command[0], command[1:]...)
	realCmd.Dir = b.config.SrcDir
	realCmd.Stdout = os.Stdout
	realCmd.Stderr = os.Stderr
	err := realCmd.Run()
	if err != nil {
		return err
	}
	return nil
}
