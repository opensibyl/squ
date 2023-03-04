package UnitSqueezer

import (
	"fmt"
	"testing"

	"github.com/opensibyl/UnitSqueezor/object"
)

func TestApi(t *testing.T) {
	conf := object.DefaultConfig()
	conf.Dry = true
	MainFlow(conf)
}

func TestA(t *testing.T) {
	fmt.Println("aa")
}

func TestB(t *testing.T) {
	fmt.Println("bb")
}
