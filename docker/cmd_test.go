package docker

import (
	"testing"

	"github.com/frozzare/go-assert"
)

func TestExecCmd(t *testing.T) {
	err := execcmd("ls", false)
	assert.Nil(t, err)

	err = execcmd("ls", true)
	assert.Nil(t, err)
}
