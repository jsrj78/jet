package gadgets

import (
	"github.com/jeelabs/jet/glow"
	"github.com/eclipse/paho.mqtt.golang"
)

func init() {
	glow.Registry["mqtt"] = func(args glow.Message) glow.Gadgetry {
		pattern := args.At(0).AsString()
		if pattern == "" {
			pattern = "#"
		}
		broker := args.At(1).AsString()
		if broker == "" {
			broker = "tcp://localhost:1883"
		}

		g := new(glow.Gadget)
		g.AddOutlets(1)

		opts := mqtt.NewClientOptions()
		opts.AddBroker(broker)

		c := mqtt.NewClient(opts)
		if t := c.Connect(); t.Wait() && t.Error() != nil {
			panic(t.Error())
		}

		var f mqtt.MessageHandler = func(c mqtt.Client, m mqtt.Message) {
			g.Emit(0, glow.Message{string(m.Topic()), string(m.Payload())})
		}
		if t := c.Subscribe(pattern, 0, f); t.Wait() && t.Error() != nil {
			panic(t.Error())
		}

		return g
	}
}
