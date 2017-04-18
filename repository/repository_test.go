package repository

import (
	"os"
	"path/filepath"
	"testing"

	git "gopkg.in/src-d/go-git.v4"

	assert "github.com/frozzare/go-assert"
	"github.com/frozzare/shed/config"
)

func TestRepository(t *testing.T) {
	path, err := os.Getwd()
	assert.Nil(t, err)

	_, err = git.PlainInit(path, false)
	assert.Nil(t, err)

	config, err := NewRepository(config.Git{
		Path: path,
	})
	assert.Nil(t, err)
	assert.Equal(t, config.Branch, "master")

	os.RemoveAll(filepath.Join(path, ".git"))
}
