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

func scanEvents() {
	events := make(chan keyevent)

	devices, err := getCompatibleDevices(true)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	for _, d := range devices {
		go listenEvents(d, events, false)
	}

	fmt.Println("Press ctrl+C to stop")

	for event := range events {
		fmt.Printf("%s: %d\n", event.deviceName, event.keycode)
	}
}

func handleEvents(configPath string) {

	if configPath == "" {
		flag.PrintDefaults()
		os.Exit(-1)
	}

	configFile, err := os.Open(configPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	config := parseConfig(configFile)
	hooks, err := hooksFromConfig(&config)

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	events := make(chan keyevent)
	for _, d := range config.Devices {
		go listenEvents(d.Name, events, d.Grab)
	}

	for event := range events {
		for _, handler := range hooks[event.deviceName][event.keycode] {
			handler.handle()
		}
	}
}

func main() {
	configPtr := flag.String("config", "", "path to config file")
	scanPtr := flag.Bool("scan", false, "output events from available devices")

	flag.Parse()

	if *scanPtr == true {
		scanEvents()
	} else {
		handleEvents(*configPtr)
	}

}
