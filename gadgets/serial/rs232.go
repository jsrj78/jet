package serial

import (
	"time"

	"github.com/chimera/rs232"
	"github.com/dataence/glog"
	"github.com/jeelabs/jet/flow"
)

func init() {
	flow.Register("serial-rs232", func() flow.Circuitry { return new(rs232G) })
}

type rs232G struct {
	flow.Gadget

	In  flow.Inlet
	Dev flow.Inlet
	Out flow.Outlet

	dev *rs232.Port
}

var (
	openSym  = flow.Sym("open")
	resetSym = flow.Sym("reset")
	bootSym  = flow.Sym("boot")
)

func (g *rs232G) Trigger() {
	switch v := g.In.(type) {
	case flow.Vec:
		switch v.Tag() {

		case openSym:
			if g.dev != nil {
				glog.Errorln("already open:", g.dev)
				return
			}

			port := v[1].(string)
			baud := v[2].(int)

			opt := rs232.Options{
				BitRate:  uint32(baud),
				DataBits: 8,
				StopBits: 1,
			}
			dev, err := rs232.Open(port, opt)
			if err != nil {
				glog.Errorf("port %q: %q", port, err)
				return
			}
			g.dev = dev

			go func() {
				var buf [250]byte
				for {
					n, err := g.dev.Read(buf[:])
					if err != nil {
						glog.Error(err)
						break
					}
					g.Out.Send(buf[:n])
				}
			}()

		case resetSym:
			g.reset(false)

		case bootSym:
			g.reset(true)

		default:
			g.Gadget.Trigger()
		}

	case string:
		g.dev.Write([]byte(v + "\n"))
	case []byte:
		g.dev.Write(v)

	default:
		g.Gadget.Trigger()
	}
}

func (g *rs232G) reset(isp bool) {
	g.dev.SetRTS(isp)  // keep RTS low to enter isp mode
	g.dev.SetDTR(true) // pulse DTR to reset

	// drop any pending input data
	time.Sleep(10 * time.Millisecond)
	if n, err := g.dev.BytesAvailable(); err == nil {
		g.dev.Read(make([]byte, n))
	}

	g.dev.SetDTR(false)
	time.Sleep(time.Millisecond)
	g.dev.SetRTS(false)
}
