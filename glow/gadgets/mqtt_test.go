package gadgets

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/jeelabs/jet/glow"
)

const SKIP = true // skip this test by default, it depends on my local setup

func TestMqttConnect(t *testing.T) {
	if SKIP || testing.Short() || os.Getenv("USER") != "jcw" {
		t.SkipNow()
	} else {
		tmp := glow.Debug
		defer func() { glow.Debug = tmp }()
		b := &bytes.Buffer{}
		glow.Debug = b

		c := glow.NewCircuit()
		c.Add(glow.LookupGadget("mqtt", "hub/1hz", "tcp://mohse:1883"))
		c.Add(glow.LookupGadget("print"))
		c.AddWire(0, 0, 1, 0)

		time.Sleep(3 * time.Second)

		if b.String() == "" {
			t.Error("no messages received")
		}
	}
}
