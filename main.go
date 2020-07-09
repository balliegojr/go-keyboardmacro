package main

import (
	"flag"
	"fmt"
	"os"
)

type keyevent struct {
	deviceName string
	keycode    int
}

func main() {

	configPtr := flag.String("config", "", "path to config file")

	flag.Parse()

	if *configPtr == "" {
		flag.PrintDefaults()
		os.Exit(-1)
	}

	configFile, err := os.Open(*configPtr)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	events := make(chan keyevent)

	config := parseConfig(configFile)
	hooks, err := hooksFromConfig(&config)

	for _, d := range config.Devices {
		go openDevice(d.Name, events)
	}

	for event := range events {
		for _, hook := range hooks[event.deviceName][event.keycode] {
			hook.execute()
		}
	}
}
