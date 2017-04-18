package docker

import (
	"testing"

	"github.com/frozzare/go-assert"
)

func TestExecCmd(t *testing.T) {
	err := ExecCmd("ls", false)
	assert.Nil(t, err)

	err = ExecCmd("ls", true)
	assert.Nil(t, err)
}
