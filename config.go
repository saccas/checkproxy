package main

import (
	"fmt"

	"github.com/unprofession-al/objectstore"
	"gopkg.in/yaml.v2"
)

// Config holds the whole configuration.
type Config struct {
	PersistanceBase string `yaml:"presistance_base"`
	PathPrefix      string `yaml:"path_prefix"`
	Auth            struct {
		RWTokens []string `yaml:"rw_tokens"`
		WTokens  []string `yaml:"w_tokens"`
		RTokens  []string `yaml:"r_tokens"`
	} `yaml:"auth"`
}

// NewConfig reads the file at the path provided, unmarshalls its content into a config
// struct and returns the result.
func NewConfig(path string) (*Config, error) {
	c := &Config{}

	o, err := objectstore.New(path)
	if err != nil {
		errOut := fmt.Errorf("error while parsing config file name '%s': %s", path, err)
		return c, errOut
	}

	data, err := o.Read()
	if err != nil {
		errOut := fmt.Errorf("error while reading config file '%s': %s", path, err)
		return c, errOut
	}

	err = yaml.Unmarshal(data, c)
	if err != nil {
		errOut := fmt.Errorf("error while unmarshalling config file '%s': %s", path, err)
		return c, errOut
	}

	return c, nil
}
