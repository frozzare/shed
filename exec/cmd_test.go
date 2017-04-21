package exec

import (
	"testing"

	"github.com/frozzare/go-assert"
)

func TestCmd(t *testing.T) {
	err := ExecCmd("ls", false)
	assert.Nil(t, err)

	err = ExecCmd("ls", true)
	assert.Nil(t, err)
}
