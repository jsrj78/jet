package flow

// Circuit is a collection of gadgets.
type Circuit struct {
	Gadget
	gadgets map[string]*Gadget
}

// Add a new gadget to the circuit.
func (c *Circuit) Add(name, typ string) {
	g := registry[typ]()
	c.gadgets[name] = g.install(g, name, c)
}

// Add a new wire connection to a circuit.
func (c *Circuit) Connect(fname string, fpin int, tname string, tpin int) {
	fg := c.gadgets[fname]
	tg := c.gadgets[tname]
	fg.Outlet(fpin).Connect(tg.Inlet(tpin))
}

// Remove an existing wire connection from a circuit.
func (c *Circuit) Disconnect(fname string, fpin int, tname string, tpin int) {
	fg := c.gadgets[fname]
	tg := c.gadgets[tname]
	if fg != nil && tg != nil {
		fg.Outlet(fpin).Disconnect(tg.Inlet(tpin))
	}
}

// Set a pin to a specified value.
func (c *Circuit) SendToPin(name string, pin int, m Message) {
	g := c.gadgets[name]
	sendToInlet(g.Inlet(pin), m)
}

// Terminate all the gadgets in the circuit.
func (c *Circuit) Terminate() {
	for _, g := range c.gadgets {
		g.Terminate()
	}
	//close(c.feed)
	//<-c.done
}

// NewCircuit creates a new empty circuit.
func NewCircuit() *Circuit {
	return &Circuit{
		gadgets: make(map[string]*Gadget),
	}
}
