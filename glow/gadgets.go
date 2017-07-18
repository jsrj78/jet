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
		g.AddOutlets(1)
		g.AddInlet(func(m Msg) {
			g.Emit(0, m)
		})
		return g
	}
}
