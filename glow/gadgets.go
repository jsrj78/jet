package main

func init() {
	Registry["print"] = func() Gadgetry {
		g := new(Gadget)
		g.AddInlet(func(m Msg) {
			Debug.Write([]byte(m.AsString()))
		})
		return g
	}

	Registry["pass"] = func() Gadgetry {
		g := new(Gadget)
		g.AddOutlet()
		g.AddInlet(func(m Msg) {
			g.Emit(0, m)
		})
		return g
	}

	Registry["circuit"] = func() Gadgetry {
		g := new(Circuit)
		return g
	}

	Registry["inlet~"] = func() Gadgetry {
		g := new(Gadget)
		g.AddOutlet()
		g.onAdded = func(c *Circuit) {
			c.AddInlet(func(m Msg) {
				g.Emit(0, m)
			})
		}
		return g
	}

	Registry["outlet~"] = func() Gadgetry {
		g := new(Gadget)
		g.onAdded = func(c *Circuit) {
			i := c.AddOutlet()
			g.AddInlet(func(m Msg) {
				c.Emit(i, m)
			})
		}
		return g
	}
}
