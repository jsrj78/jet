// Flow implements a Pure-data like dataflow mechanism.
package flow

// Message is a generic data item which can be sent between gadgets.
type Message interface{}

// Circuitry is the collective name for gadgets and circuits.
type Circuitry interface {
	NumInlets() int
	Inlet(n int) *Inlet

	NumOutlets() int
	Outlet(n int) *Outlet

	Setup()
	Loop()
	Trigger()
	Cleanup()

	install(self Circuitry, name string, owner *Circuit) *Gadget
}

var registry = map[string]func() Circuitry{}

// Register a constructor for a named gadget type.
func Register(name string, f func() Circuitry) {
	registry[name] = f
}
