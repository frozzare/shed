package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	assert "github.com/frozzare/go-assert"
)

func TestList(t *testing.T) {
	path, err := os.Getwd()
	assert.Nil(t, err)

	for _, dat := range []string{"before_script: test", "before_script:\n - test"} {
		err = ioutil.WriteFile(filepath.Join(path, "shed.yml"), []byte(dat), 0644)
		assert.Nil(t, err)

		config, err := NewConfig(path)
		assert.Nil(t, err)
		assert.Equal(t, config.BeforeScript.Values, []string{"test"})

		os.Remove(filepath.Join(path, "shed.yml"))
	}
}
