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

func (b *BaseRunner) GetRunCommand(_ []*openapi.ObjectFunctionWithSignature) []string {
	// TODO implement me
	panic("implement me")
}

func (b *BaseRunner) Run(cases []*openapi.ObjectFunctionWithSignature, ctx context.Context) error {
	cmd := b.GetRunCommand(cases)
	realCmd := exec.CommandContext(ctx, cmd[0], cmd[1:]...)
	realCmd.Dir = b.config.SrcDir
	realCmd.Stdout = os.Stdout
	realCmd.Stderr = os.Stderr
	err := realCmd.Run()
	if err != nil {
		return err
	}
	return nil
}
