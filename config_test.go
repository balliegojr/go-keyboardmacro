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
    grab: true
    bindings:
      A:
        exec: 'command'
      31: 
        exec: 'another command'

  - device:
    name: 'other device'
    bindings:
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

	if inputkbd.Grab == false {
		t.Fatalf("got device grab %v, expeted true", inputkbd.Grab)
	}

	if len(inputkbd.Bindings) != 2 {
		t.Fatalf("got length %d, expected 2", len(inputkbd.Bindings))
	}

	if inputkbd.Bindings["A"].Exec != "command" {
		t.Fatalf(`got binding["A"] %v, expected "command"`, inputkbd.Bindings["A"].Exec)
	}

	if inputkbd.Bindings["31"].Exec != "another command" {
		t.Fatalf(`got binding["31"] %v, expected "another command"`, inputkbd.Bindings["A"].Exec)
	}

	otherdevice := c.Devices[1]
	if otherdevice.Name != "other device" {
		t.Fatalf("got device name %v, expected 'other device'", otherdevice.Name)
	}

	if len(otherdevice.Bindings) != 1 {
		t.Fatalf("got bindings lenght %d, expected 1", len(otherdevice.Bindings))
	}

	if otherdevice.Bindings["B"].Exec != "B command" {
		t.Fatalf(`got binding["B"] %v, expected "B command"`, otherdevice.Bindings["B"].Exec)
	}

}
