package main

import (
	"testing"
)

func TestGetKeyCode(t *testing.T) {
	c, _ := getKeyCode("30")
	if c != 30 {
		t.Fatalf("got key code %v, expected 30", c)
	}

	c, err := getKeyCode("A")
	if c != 0 {
		t.Fatalf("got key code %v, expected 0", c)
	}

	if err.Error() != "Invalid key code" {
		t.Fatalf("got err %v, expected 'Ivalid key code'", err)
	}

	c, err = getKeyCode("-1")
	if c != 0 {
		t.Fatalf("got key code %v, expected 0", c)
	}

	if err.Error() != "Key code not found" {
		t.Fatalf("got err %v, expected 'Key code not found'", err)
	}

}

func TestHooksFromConfig(t *testing.T) {
	h, err := hooksFromConfig(&Config{})
	if len(h) != 0 {
		t.Fatalf("got hooks length %v, expected 0", len(h))
	}

	if err != nil {
		t.Fatalf("got err %v, expected nil", err)
	}

	c := Config{
		Devices: []Device{
			Device{"invalid device", map[string]Macro{}},
		},
	}

	h, err = hooksFromConfig(&c)
	if err.Error() != "Input device not found" {
		t.Fatalf("got err %v, expected 'Input device not found'", err)
	}

	devices, err := getCompatibleDevices(true)
	if err != nil {
		t.Fatal(err)
	}

	c = Config{
		Devices: []Device{
			Device{
				Name: devices[0],
				Macros: map[string]Macro{
					"30": Macro{"command X", ""},
					"31": Macro{"command y", ""},
				},
			},
		},
	}

	h, err = hooksFromConfig(&c)
	if len(h) != 1 {
		t.Fatalf("got length devices %v, expected 1", len(h))
	}

	if len(h[devices[0]]) != 2 {
		t.Fatalf(`got hooks[%v] length %v, expected 2`, devices[0], len(h[devices[0]]))
	}

	if len(h[devices[0]][30]) != 1 {
		t.Fatalf(`got hooks[%v][30] length %v, expected 1`, devices[0], len(h[devices[0]][30]))
	}

	if len(h[devices[0]][31]) != 1 {
		t.Fatalf(`got hooks[%v][31] length %v, expected 1`, devices[0], len(h[devices[0]][31]))
	}

}
