package main

func init() {
	Registry["print"] = func() Gadgetry {
		g := new(Gadget)
		g.AddInlet(func(m Msg) {
			Debug.Write([]byte(m.AsString()))
		})
		return g
	}
}
