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

func (c *Circuit) Request(args ...Message) {
	c.feed <- incoming{msg: args}
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
func (c *Circuit) Control(cmd []Message) {
	fmt.Println("Circuit control:", cmd)
	switch cmd[0] {

	case AddSym: // Add a new gadget to the circuit.
		typ := cmd[1].(string)
		g := registry[Sym(typ)]()
		c.gadgets = append(c.gadgets, g.install(g, c))

	case ConnectSym: // Add a new wire connection.
		fidx := cmd[1].(int)
		fpin := cmd[2].(int)
		tidx := cmd[3].(int)
		tpin := cmd[4].(int)
		c.gadgets[fidx].Outlet(fpin).connect(c.gadgets[tidx].Inlet(tpin))

	case DisconnectSym: // Remove an existing wire connection.
		fidx := cmd[1].(int)
		fpin := cmd[2].(int)
		tidx := cmd[3].(int)
		tpin := cmd[4].(int)
		c.gadgets[fidx].Outlet(fpin).disconnect(c.gadgets[tidx].Inlet(tpin))

	case SendToPinSym: // Set a pin to a specified value.
		idx := cmd[1].(int)
		pin := cmd[2].(int)
		msg := cmd[3].(Message)
		sendToInlet(c.gadgets[idx].Inlet(pin), msg)

	default:
		c.Gadget.Control(cmd)
	}
}
