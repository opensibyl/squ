package object

import openapi "github.com/opensibyl/sibyl-go-client"

type FunctionWithState struct {
	*openapi.ObjectFunctionWithSignature

	Reachable bool
	ReachBy   []string
}

type DiffMap = map[string][]int
type DiffFuncMap = map[string][]*FunctionWithState
