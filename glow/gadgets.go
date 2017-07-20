package glow

import "fmt"

func init() {
	Registry["print"] = func(args Message) Gadgetry {
		g := new(Gadget)
		g.AddInlet(func(m Message) {
			if args.IsBang() {
				fmt.Fprintln(Debug, m)
			} else {
				fmt.Fprintln(Debug, args, m)
			}
		})
		return g
	}

	Registry["pass"] = func(args Message) Gadgetry {
		g := new(Gadget)
		g.AddOutlets(1)
		g.AddInlet(func(m Message) {
			g.Emit(0, m)
		})
		return g
	}

	Registry["inlet"] = func(args Message) Gadgetry {
		g := new(Gadget)
		g.AddOutlets(1)
		g.onAdded = func(c *Circuit) {
			c.AddInlet(func(m Message) {
				g.Emit(0, m)
			})
		}
		return g
	}

	Registry["outlet"] = func(args Message) Gadgetry {
		g := new(Gadget)
		g.onAdded = func(c *Circuit) {
			o := c.AddOutlets(1)
			g.AddInlet(func(m Message) {
				c.Emit(o, m)
			})
		}
		return g
	}

	Registry["swap"] = func(args Message) Gadgetry {
		g := new(Gadget)
		g.AddOutlets(2)
		g.AddInlet(func(m Message) {
			g.Emit(1, m)
			g.Emit(0, args)
		})
		g.AddInlet(func(m Message) {
			args = m
		})
		return g
	}
}
