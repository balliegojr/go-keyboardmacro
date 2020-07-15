package main

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	evdev "github.com/gvalkov/golang-evdev"
)

type handler interface {
	handle()
}

type execHandler struct {
	command string
}

func (e execHandler) handle() {
	parts := strings.Fields(e.command)
	command := parts[0]

	cmd := exec.Command(command, parts...)
	err := cmd.Start()

	if err != nil {
		fmt.Println(err)
	}
}

const inputAbsolutePath = "/dev/input/by-id/"

// getCompatibleDevices list all kbd and mouse devices
func getCompatibleDevices(nameOnly bool) ([]string, error) {
	devicesNames, _ := filepath.Glob(inputAbsolutePath + "*")
	if devicesNames == nil {
		return []string{}, errors.New("No input devices found")
	}

	devices := []string{}
	for _, d := range devicesNames {
		if strings.HasSuffix(d, "kbd") || strings.HasSuffix(d, "mouse") {

			if nameOnly {
				devices = append(devices, d[len(inputAbsolutePath):])
			} else {
				devices = append(devices, d)
			}
		}
	}

	return devices, nil
}

func getKeyCode(c string) (int, error) {
	code, err := strconv.Atoi(c)
	if err != nil {
		return 0, errors.New("Invalid key code")
	}

	if _, ok := evdev.KEY[code]; ok {
		return code, nil
	}

	return 0, errors.New("Key code not found")

}

// getHook return the hook for the given binding. Only exec hook is available today
func getHandler(b *Binding) handler {
	return execHandler{b.Exec}
}

// Transform given configuration into hooks.
// The result will be a slice of hooks for a key code and device
// ex: hooks := result['device_name'][30]
func hooksFromConfig(c *Config) (map[string]map[int][]handler, error) {
	devicesMap := map[string]map[int][]handler{}

	for _, d := range c.Devices {
		if !evdev.IsInputDevice(inputAbsolutePath + d.Name) {
			return map[string]map[int][]handler{}, errors.New("Input device not found")
		}

		dm, ok := devicesMap[d.Name]
		if !ok {
			dm = map[int][]handler{}
			devicesMap[d.Name] = dm
		}

		for code, m := range d.Bindings {
			c, err := getKeyCode(code)
			if err != nil {
				return map[string]map[int][]handler{}, err
			}

			dm[c] = append(dm[c], getHandler(&m))
		}
	}

	return devicesMap, nil

}

func listenEvents(deviceName string, events chan keyevent, grabDevice bool) {
	device, err := evdev.Open(inputAbsolutePath + deviceName)

	if err != nil {
		fmt.Println(err)
		return
	}

	if grabDevice {
		device.Grab()
	}

	for {
		ev, _ := device.ReadOne()
		if ev.Type != evdev.EV_KEY {
			continue
		}

		e := evdev.NewKeyEvent(ev)
		if e.State == evdev.KeyUp {
			events <- keyevent{deviceName: deviceName, keycode: int(e.Scancode)}
		}
	}
}
