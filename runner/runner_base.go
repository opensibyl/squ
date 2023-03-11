package runner

import (
	"github.com/opensibyl/UnitSqueezor/object"
	openapi "github.com/opensibyl/sibyl-go-client"
)

type BaseRunner struct {
	config    *object.SharedConfig
	apiClient *openapi.APIClient
}
