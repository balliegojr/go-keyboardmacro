package main

import (
	"io"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

//Binding represents the binding section in the configuration file
type Binding struct {
	Exec string
	Type string `yaml:"omitempty"`
}

// Device represents the device section in the configuration file
type Device struct {
	Name     string
	Bindings map[string]Binding
	Grab     bool `yaml:"grab,omitempty"`
}

// Config represents device configuration and its macros
type Config struct {
	Devices []Device
}

// parseConfig reads the yaml configuration from the Reader
func parseConfig(r io.Reader) Config {
	c := Config{}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return c
}
