package service_config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	ConfigPathEnvVar = "CONFIG_PATH"
	ConfigEnvVar     = "CONFIG"
)

var NoConfigError = errors.New("No Config or Config Path Specified. Please supply one of the following: -config, -configPath, CONFIG, or CONFIG_PATH")

type ServiceConfig struct {
	configFlag     string
	configPathFlag string
}

func New() *ServiceConfig {
	return &ServiceConfig{}
}

func (c *ServiceConfig) AddFlags(flagSet *flag.FlagSet) {
	flagSet.StringVar(&c.configFlag, "config", "", "json encoded configuration string")
	flagSet.StringVar(&c.configPathFlag, "configPath", "", "path to configuration file with json encoded content")
}

func (c ServiceConfig) ConfigBytes() ([]byte, error) {
	if c.configFlag != "" {
		return []byte(c.configFlag), nil
	}

	if c.configPathFlag != "" {
		absolutePath, err := filepath.Abs(c.configPathFlag)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Making config file path absolute: %s", err.Error()))
		}

		bytes, err := ioutil.ReadFile(absolutePath)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Reading config file: %s", err.Error()))
		}

		return bytes, nil
	}

	config := os.Getenv(ConfigEnvVar)
	if config != "" {
		return []byte(config), nil
	}

	configPath := os.Getenv(ConfigPathEnvVar)
	if configPath != "" {
		absolutePath, err := filepath.Abs(configPath)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Making config file path absolute: %s", err.Error()))
		}

		bytes, err := ioutil.ReadFile(absolutePath)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Reading config file: %s", err.Error()))
		}

		return bytes, nil
	}

	return nil, NoConfigError
}

func (c ServiceConfig) ConfigPath() string {
	return c.configPathFlag
}

func (c ServiceConfig) Read(model interface{}) error {
	bytes, err := c.ConfigBytes()
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, model)
	if err != nil {
		return errors.New(fmt.Sprintf("Unmarshaling config: %s", err.Error()))
	}

	return nil
}
