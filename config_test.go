package main

import (
	"strings"
	"testing"
)

func TestParseConfig(t *testing.T) {
	sampleConfig := `version: 1.2
devices: 
  - device:
    name: 'input_kbd'
    macros:
      A:
        exec: 'command'
      31: 
        exec: 'another command'

  - device:
    name: 'other device'
    macros:
      B: 
        exec: 'B command'

`

	r := strings.NewReader(sampleConfig)
	c := parseConfig(r)
	t.Logf("%v", c)

	if len(c.Devices) != 2 {
		t.Fatalf("got devices length %d, expected 2", len(c.Devices))
	}

	inputkbd := c.Devices[0]

	if inputkbd.Name != "input_kbd" {
		t.Fatalf("got device name %v, expected input_kbd", inputkbd.Name)
	}

	if len(inputkbd.Macros) != 2 {
		t.Fatalf("got length %d, expected 2", len(inputkbd.Macros))
	}

	if inputkbd.Macros["A"].Exec != "command" {
		t.Fatalf(`got macro["A"] %v, expected "command"`, inputkbd.Macros["A"].Exec)
	}

	if inputkbd.Macros["31"].Exec != "another command" {
		t.Fatalf(`got macro["31"] %v, expected "another command"`, inputkbd.Macros["A"].Exec)
	}

	otherdevice := c.Devices[1]
	if otherdevice.Name != "other device" {
		t.Fatalf("got device name %v, expected 'other device'", otherdevice.Name)
	}

	if len(otherdevice.Macros) != 1 {
		t.Fatalf("got macros lenght %d, expected 1", len(otherdevice.Macros))
	}

	if otherdevice.Macros["B"].Exec != "B command" {
		t.Fatalf(`got macro["B"] %v, expected "B command"`, otherdevice.Macros["B"].Exec)
	}

}
