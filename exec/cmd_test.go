package exec

import (
	"testing"

	"github.com/frozzare/go-assert"
)

func TestCmd(t *testing.T) {
	err := Cmd("ls", false)
	assert.Nil(t, err)

	err = Cmd("ls", true)
	assert.Nil(t, err)
}
