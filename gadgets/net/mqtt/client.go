package mqtt

import (
	"github.com/dataence/glog"
	"github.com/jeelabs/jet/flow"
	"github.com/surge/surgemq/message"
	"github.com/surge/surgemq/service"
)

func init() {
	flow.Register("mqtt-client", func() flow.Circuitry { return new(clientG) })
}

type clientG struct {
	flow.Gadget
	Cmd   flow.Inlet
	Topic flow.Inlet
	Out   flow.Outlet
	clt   *service.Client
}

func (g *clientG) Setup() {
	g.clt = &service.Client{}

	msg := message.NewConnectMessage()
	msg.SetVersion(4)
	msg.SetKeepAlive(300)

	if err := g.clt.Connect("tcp://127.0.0.1:1883", msg); err != nil {
		glog.Fatal(err)
	}
}

func (g *clientG) Trigger() {
}
