package mqtt

import (
	"github.com/golang/glog"
	"github.com/jeelabs/jet/flow"
	"github.com/surge/surgemq/service"
)

func init() {
	flow.Register("mqtt-server", func() flow.Circuitry { return new(serverG) })
}

type serverG struct {
	flow.Gadget
	Cmd   flow.Inlet
	Topic flow.Inlet
	srv   *service.Server
}

func (g *serverG) Setup() {
	g.srv = &service.Server{}

	if err := g.srv.ListenAndServe("tcp://:1883"); err != nil {
		glog.Fatal(err)
	}
}

func (g *serverG) Trigger() {
}
