package main

import (
	"io"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Config represents device configuration and its macros
type Config struct {
	Devices []struct {
		Name   string
		Macros map[string]string
	}
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
