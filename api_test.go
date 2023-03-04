package UnitSqueezer

import (
	"fmt"
	"testing"

	"github.com/opensibyl/UnitSqueezor/object"
)

func TestApi(t *testing.T) {
	// dead loop
	t.Skip()
	MainFlow(object.DefaultConfig())
}

func TestA(t *testing.T) {
	fmt.Println("aa")
}

func TestB(t *testing.T) {
	fmt.Println("bb")
}
