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

func (c *Circuit) Request(req string, args ...Message) {
	cmd := append([]Message{Sym(req)}, args...)
	c.feed <- incoming{msg: cmd}
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
	addSym        = Sym("add")
	connectSym    = Sym("connect")
	disconnectSym = Sym("disconnect")
	send2pinSym   = Sym("send2pin") // temporary
)

// Control gets called with messages sent to the special nil inlet.
func (c *Circuit) Control(cmd []Message) {
	fmt.Println("Circuit control:", cmd)
	switch cmd[0] {

	case addSym: // Add a new gadget to the circuit.
		typ := cmd[1].(string)
		g := registry[Sym(typ)]()
		c.gadgets = append(c.gadgets, g.install(g, c))

	case connectSym: // Add a new wire connection.
		fidx := cmd[1].(int)
		fpin := cmd[2].(int)
		tidx := cmd[3].(int)
		tpin := cmd[4].(int)
		c.gadgets[fidx].Outlet(fpin).connect(c.gadgets[tidx].Inlet(tpin))

	case disconnectSym: // Remove an existing wire connection.
		fidx := cmd[1].(int)
		fpin := cmd[2].(int)
		tidx := cmd[3].(int)
		tpin := cmd[4].(int)
		c.gadgets[fidx].Outlet(fpin).disconnect(c.gadgets[tidx].Inlet(tpin))

	case send2pinSym: // Set a pin to a specified value.
		idx := cmd[1].(int)
		pin := cmd[2].(int)
		msg := cmd[3].(Message)
		sendToInlet(c.gadgets[idx].Inlet(pin), msg)

	default:
		c.Gadget.Control(cmd)
	}
}
