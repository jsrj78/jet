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
}
