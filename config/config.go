package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

var (
	ErrNoShedFile = errors.New("No .shed.yml or shed.yml file found")
)

// Docker represents a docker config section.
type Docker struct {
	Machine string      `yaml:"machine"`
	Proxy   DockerProxy `yaml:"proxy"`
	Volumes List        `yaml:"volumes"`
}

// DockerProxy represents a docker proxy config section.
type DockerProxy struct {
	Image     string `yaml:"image"`
	Env       List   `yaml:"env"`
	HTTPPort  string `yaml:"http_port"`
	HTTPSPort string `yaml:"https_port"`
	Recreate  bool   `yaml:"recreate"`
	Volumes   List   `yaml:"volumes"`
}

// Git represents a git config section.
type Git struct {
	Branch string `yaml:"branch"`
	Path   string `yaml:"path"`
}

// Config represents a config file.
type Config struct {
	AfterScript  List              `yaml:"after_script"`
	BeforeScript List              `yaml:"before_script"`
	Branches     map[string]Config `yaml:"branches"`
	Docker       Docker            `yaml:"docker"`
	Git          Git               `yaml:"git"`
	Host         string            `yaml:"host"`
	HTTPS        bool              `yaml:"https"`
	Script       List              `yaml:"script"`
}

// NewConfig creates a new config struct from a yaml file.
func NewConfig(args ...string) (Config, error) {
	var file string
	var path string
	var err error

	if len(args) > 0 && args[0] != "" {
		if _, err := os.Stat(args[0]); err == nil {
			file = args[0]
			path = path = filepath.Dir(file)
		} else {
			path = args[0]
		}
	} else {
		path, err = os.Getwd()
	}

	if err != nil {
		return Config{}, err
	}

	files := []string{".shed.yml", "shed.yml"}
	if len(file) > 0 {
		files = append([]string{file}, files...)
	}

	var dat []byte
	for _, name := range files {
		if len(dat) > 0 {
			break
		}

		file := filepath.Join(path, name)
		dat, err = ioutil.ReadFile(file)
	}

	if err != nil {
		return Config{}, ErrNoShedFile
	}

	var config Config

	if err := yaml.Unmarshal(dat, &config); err != nil {
		return Config{}, err
	}

	if config.HTTPS {
		config.Docker.Proxy.HTTPSPort = "443:443"
	}

	return config, nil
}

// Def returns second parameter if first paramter length is zero.
func Def(a, b string) string {
	if len(a) == 0 {
		return b
	}

	return a
}

// DefList returns second parameter if first paramter length is zero.
func DefList(a, b []string) []string {
	if len(a) == 0 {
		return b
	}

	return a
}
