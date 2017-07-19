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
		g.NumOutlets(1)
		g.AddInlet(func(m Msg) {
			g.Emit(0, m)
		})
		return g
	}

	Registry["circuit"] = func() Gadgetry {
		g := new(Circuit)
		return g
	}
}