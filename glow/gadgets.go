package glow

import "fmt"

func init() {
	Registry["print"] = func(args Msg) Gadgetry {
		g := new(Gadget)
		g.AddInlet(func(m Msg) {
			fmt.Fprint(Debug, m)
		})
		return g
	}

	Registry["pass"] = func(args Msg) Gadgetry {
		g := new(Gadget)
		g.AddOutlets(1)
		g.AddInlet(func(m Msg) {
			g.Emit(0, m)
		})
		return g
	}
	/*
		Registry["circuit"] = func(args Msg) Gadgetry {
			g := new(Circuit)
			return g
		}
	*/
	Registry["inlet~"] = func(args Msg) Gadgetry {
		g := new(Gadget)
		g.AddOutlets(1)
		g.onAdded = func(c *Circuit) {
			c.AddInlet(func(m Msg) {
				g.Emit(0, m)
			})
		}
		return g
	}

	Registry["outlet~"] = func(args Msg) Gadgetry {
		g := new(Gadget)
		g.onAdded = func(c *Circuit) {
			o := c.AddOutlets(1)
			g.AddInlet(func(m Msg) {
				c.Emit(o, m)
			})
		}
		return g
	}
}
