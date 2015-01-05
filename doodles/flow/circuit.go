package flow

// Circuit is a collection of gadgets.
type Circuit struct {
	Gadget
	gadgets []*Gadget
}

// NewCircuit creates a new empty circuit.
func NewCircuit() *Circuit {
	return &Circuit{}
}

// Add a new gadget to the circuit.
func (c *Circuit) Add(typ string) int {
	g := registry[Sym(typ)]()
	i := len(c.gadgets)
	c.gadgets = append(c.gadgets, g.install(g, c))
	return i
}

// Add a new wire connection to a circuit.
func (c *Circuit) Connect(fidx, fpin, tidx, tpin int) {
	c.gadgets[fidx].Outlet(fpin).Connect(c.gadgets[tidx].Inlet(tpin))
}

// Remove an existing wire connection from a circuit.
func (c *Circuit) Disconnect(fidx, fpin, tidx, tpin int) {
	c.gadgets[fidx].Outlet(fpin).Disconnect(c.gadgets[tidx].Inlet(tpin))
}

// Set a pin to a specified value.
func (c *Circuit) SendToPin(idx, pin int, m Message) {
	sendToInlet(c.gadgets[idx].Inlet(pin), m)
}

// Terminate all the gadgets in the circuit.
func (c *Circuit) Terminate() {
	for _, g := range c.gadgets {
		g.Terminate()
	}
	//c.Gadget.Terminate()
}
