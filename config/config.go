package config

import (
	"encoding/json"
	"io/ioutil"
)

type ConfigStruct struct {
	IP          string
	Port        string
	DatabaseURL string
	Secret      string
}

var config *ConfigStruct

func Get() *ConfigStruct {
	return config
}

func Export(configPath string) error {
	var c *ConfigStruct

	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(configFile, &c)
	if err != nil {
		return err
	}

	config = c
	return nil
}
