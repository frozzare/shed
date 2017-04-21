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

	for file, arg := range map[string]string{"shed.yml": path, "shed-custom.yml": "shed-custom.yml"} {
		dat := `host: example.com`

		err = ioutil.WriteFile(filepath.Join(path, file), []byte(dat), 0644)
		assert.Nil(t, err)

		config, err := NewConfig(arg)
		assert.Nil(t, err)
		assert.Equal(t, config.Host, "example.com")

		os.Remove(filepath.Join(path, file))
	}
}
