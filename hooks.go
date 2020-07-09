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

type hook interface {
	execute()
}

type execHook struct {
	command string
}

func (e execHook) execute() {
	err := exec.Command("sh", "-c", e.command).Start()
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

// getHook return the hook for the given macro. Only exec hook is available today
func getHook(m *Macro) hook {
	return execHook{m.Exec}
}

// Transform given configuration into hooks.
// The result will be a slice of hooks for a key code and device
// ex: hooks := result['device_name'][30]
func hooksFromConfig(c *Config) (map[string]map[int][]hook, error) {
	devicesMap := map[string]map[int][]hook{}

	for _, d := range c.Devices {
		if !evdev.IsInputDevice(inputAbsolutePath + d.Name) {
			return map[string]map[int][]hook{}, errors.New("Input device not found")
		}

		dm, ok := devicesMap[d.Name]
		if !ok {
			dm = map[int][]hook{}
			devicesMap[d.Name] = dm
		}

		for code, m := range d.Macros {
			c, err := getKeyCode(code)
			if err != nil {
				return map[string]map[int][]hook{}, err
			}

			dm[c] = append(dm[c], getHook(&m))
		}
	}

	return devicesMap, nil

}

func openDevice(deviceName string, events chan keyevent) {
	device, _ := evdev.Open(inputAbsolutePath + deviceName)

	for {
		ev, _ := device.ReadOne()
		e := evdev.NewKeyEvent(ev)

		if e.State == evdev.KeyUp {
			events <- keyevent{deviceName: deviceName, keycode: int(e.Scancode)}
		}
	}
}
