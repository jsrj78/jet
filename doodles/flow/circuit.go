package flow

import "fmt"

// Circuit is a collection of gadgets.
type Circuit struct {
	Gadget
	gadgets []*Gadget
}

// NewCircuit creates a new empty circuit.
func NewCircuit() *Circuit {
	c := new(Circuit)
	c.install(c, nil)
	return c
}

func (c *Circuit) request(v ...Message) {
	c.feed <- incoming{msg: v}
}

// Add a new gadget to the circuit.
func (c *Circuit) Add(typ string) {
	c.request(AddSym, typ)
}

// Add a new wire connection to a circuit.
func (c *Circuit) Connect(fidx, fpin, tidx, tpin int) {
	c.request(ConnectSym, fidx, fpin, tidx, tpin)
}

// Remove an existing wire connection from a circuit.
func (c *Circuit) Disconnect(fidx, fpin, tidx, tpin int) {
	c.request(DisconnectSym, fidx, fpin, tidx, tpin)
}

// Set a pin to a specified value.
func (c *Circuit) SendToPin(idx, pin int, m Message) {
	c.request(SendToPinSym, idx, pin, m)
}

// Terminate all the gadgets in the circuit.
func (c *Circuit) Terminate() {
	for _, g := range c.gadgets {
		g.Terminate()
	}
	//c.Gadget.Terminate()
}

var (
	// TODO these could be pre-defined as enum, if collected centrally
	AddSym        = Sym("add")
	ConnectSym    = Sym("connect")
	DisconnectSym = Sym("disconnect")
	SendToPinSym  = Sym("sendtopin") // temporary
)

// Control gets called with messages sent to the special nil inlet.
func (c *Circuit) Control(m Message) {
	if v, ok := m.([]Message); ok {
		fmt.Println("Circuit control:", v)
		switch v[0] {
		case AddSym:
			typ := v[1].(string)
			g := registry[Sym(typ)]()
			c.gadgets = append(c.gadgets, g.install(g, c))
		case ConnectSym:
			fidx := v[1].(int)
			fpin := v[2].(int)
			tidx := v[3].(int)
			tpin := v[4].(int)
			c.gadgets[fidx].Outlet(fpin).Connect(c.gadgets[tidx].Inlet(tpin))
		case DisconnectSym:
			fidx := v[1].(int)
			fpin := v[2].(int)
			tidx := v[3].(int)
			tpin := v[4].(int)
			c.gadgets[fidx].Outlet(fpin).Disconnect(c.gadgets[tidx].Inlet(tpin))
		case SendToPinSym:
			idx := v[1].(int)
			pin := v[2].(int)
			msg := v[3].(Message)
			sendToInlet(c.gadgets[idx].Inlet(pin), msg)
		default:
			c.Gadget.Control(m)
		}
	}
}
