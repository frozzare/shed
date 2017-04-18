package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

// Docker represents a docker config section.
type Docker struct {
	Endpoint string `yaml:"endpoint"`
}

// Git represents a git config section.
type Git struct {
	Branch string `yaml:"branch"`
	Path   string `yaml:"path"`
}

// Config represents a config file.
type Config struct {
	Branches map[string]Config `yaml:"branches"`
	Docker   Docker            `yaml:"docker"`
	Domain   string            `yaml:"domain"`
	Git      Git               `yaml:"git"`
}

// NewConfig creates a new config struct from a yaml file.
func NewConfig(args ...string) (Config, error) {
	var path string
	var err error

	if len(args) > 0 && args[0] != "" {
		path = args[0]
	} else {
		path, err = os.Getwd()
	}

	if err != nil {
		return Config{}, err
	}

	file := filepath.Join(path, "shed.yml")

	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return Config{}, err
	}

	var config Config

	if err := yaml.Unmarshal(dat, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}
