package gadgets

import (
	"fmt"

	"github.com/jeelabs/jet/glow"
)

func init() {
	glow.Registry["print"] = func(args glow.Message) glow.Gadgetry {
		g := glow.NewGadget()
		g.AddInlet(func(m glow.Message) {
			if args.IsBang() {
				fmt.Fprintln(glow.Debug, m)
			} else {
				fmt.Fprintln(glow.Debug, args, m)
			}
		})
		return g
	}

	glow.Registry["pass"] = func(args glow.Message) glow.Gadgetry {
		g := glow.NewGadget()
		g.AddOutlets(1)
		g.AddInlet(func(m glow.Message) {
			g.Emit(0, m)
		})
		return g
	}

	glow.Registry["inlet"] = func(args glow.Message) glow.Gadgetry {
		g := glow.NewGadget()
		g.AddOutlets(1)
		g.OnAdded = func(c *glow.Circuit) {
			c.AddInlet(func(m glow.Message) {
				g.Emit(0, m)
			})
		}
		return g
	}

	glow.Registry["outlet"] = func(args glow.Message) glow.Gadgetry {
		g := glow.NewGadget()
		g.OnAdded = func(c *glow.Circuit) {
			o := c.AddOutlets(1)
			g.AddInlet(func(m glow.Message) {
				c.Emit(o, m)
			})
		}
		return g
	}

	glow.Registry["swap"] = func(args glow.Message) glow.Gadgetry {
		g := glow.NewGadget()
		g.AddOutlets(2)
		g.AddInlet(func(m glow.Message) {
			g.Emit(1, m)
			g.Emit(0, args)
		})
		g.AddInlet(func(m glow.Message) {
			args = m
		})
		return g
	}

	glow.Registry["send"] = func(args glow.Message) glow.Gadgetry {
		var parent *glow.Circuit
		g := glow.NewGadget()
		g.OnAdded = func(c *glow.Circuit) {
			parent = c
		}
		g.AddInlet(func(m glow.Message) {
			parent.Notify("msg:"+args.String(), m...)
		})
		return g
	}
	glow.Registry["s"] = glow.Registry["send"]

	glow.Registry["receive"] = func(args glow.Message) glow.Gadgetry {
		g := glow.NewGadget()
		g.AddOutlets(1)
		g.OnAdded = func(c *glow.Circuit) {
			c.On("msg:"+args.String(), func(m glow.Message) {
				g.Emit(0, m)
			})
		}
		return g
	}
	glow.Registry["r"] = glow.Registry["receive"]

	glow.Registry["metro"] = func(args glow.Message) glow.Gadgetry {
		// TODO start on hot inlet, add 2nd inlet for changing period
		g := glow.NewGadget()
		g.AddOutlets(1)
		g.OnAdded = func(c *glow.Circuit) {
			glow.SetPeriodic(args.AsInt(), func() {
				g.Emit(0, nil)
			})
		}
		return g
	}

	glow.Registry["smooth"] = func(args glow.Message) glow.Gadgetry {
		state, order := 0, 0
		g := glow.NewGadget()
		g.AddOutlets(1)
		g.AddInlet(func(m glow.Message) {
			state = (order*state + m.AsInt()) / (order + 1)
			order = args.AsInt()
			g.Emit(0, glow.Message{state})
		})
		g.AddInlet(func(m glow.Message) {
			args = m
		})
		return g
	}

	glow.Registry["change"] = func(args glow.Message) glow.Gadgetry {
		var last glow.Message
		g := glow.NewGadget()
		g.AddOutlets(1)
		g.AddInlet(func(m glow.Message) {
			if last == nil || m.AsInt() != last.AsInt() {
				last = m
				g.Emit(0, last)
			}
		})
		return g
	}

	glow.Registry["moses"] = func(args glow.Message) glow.Gadgetry {
		g := glow.NewGadget()
		g.AddOutlets(2)
		g.AddInlet(func(m glow.Message) {
			if m.AsInt() < args.AsInt() {
				g.Emit(0, m)
			} else {
				g.Emit(1, m)
			}
		})
		g.AddInlet(func(m glow.Message) {
			args = m
		})
		return g
	}
}
