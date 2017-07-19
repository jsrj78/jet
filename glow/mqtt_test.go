package glow

import (
	"bytes"
	"os"
	"testing"
	"time"
)

const SKIP = true // skip this test by default, it depends on my local setup

func TestMqttConnect(t *testing.T) {
	if SKIP || testing.Short() || os.Getenv("USER") != "jcw" {
		t.SkipNow()
	} else {
		tmp := Debug
		defer func() { Debug = tmp }()
		b := &bytes.Buffer{}
		Debug = b

		c := new(Circuit)
		c.Add(NewGadget("mqtt", "hub/1hz", "tcp://mohse:1883"))
		c.Add(NewGadget("print"))
		c.AddWire(0, 0, 1, 0)

		time.Sleep(3 * time.Second)

		if b.String() == "" {
			t.Error("no messages received")
		}
	}
}
