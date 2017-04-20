package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	assert "github.com/frozzare/go-assert"
)

func TestConfig(t *testing.T) {
	path, err := os.Getwd()
	assert.Nil(t, err)

	dat := `domain: example.com`

	err = ioutil.WriteFile(filepath.Join(path, "shed.yml"), []byte(dat), 0644)
	assert.Nil(t, err)

	config, err := NewConfig(path)
	assert.Nil(t, err)
	assert.Equal(t, config.Domain, "example.com")

	os.Remove(filepath.Join(path, "shed.yml"))
}

func TestCustomConfigFile(t *testing.T) {
	path, err := os.Getwd()
	assert.Nil(t, err)

	dat := `domain: example.com`

	err = ioutil.WriteFile(filepath.Join(path, "shed-custom.yml"), []byte(dat), 0644)
	assert.Nil(t, err)

	config, err := NewConfig("shed-custom.yml")
	assert.Nil(t, err)
	assert.Equal(t, config.Domain, "example.com")

	os.Remove(filepath.Join(path, "shed-custom.yml"))
}
