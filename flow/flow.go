// Flow implements a Pure-data like dataflow mechanism.
package flow

// Message is a generic data item which can be sent between gadgets.
type Message interface{}

// List constructs a slice from a list of messages.
func List(v ...Message) []Message {
	return v
}

// Circuitry is the collective name for gadgets and circuits.
type Circuitry interface {
	NumInlets() int
	Inlet(n int) *Inlet
	IsHot(i *Inlet) bool

	NumOutlets() int
	Outlet(n int) *Outlet

	Setup()
	Control(cmd []Message)
	Trigger()
	Cleanup()

	install(self Circuitry, owner *Circuit) *Gadget
}

// the empty string, i.e. SymVal zero, is special
var symbols = map[string]SymVal{"": SymVal(0)}
var symList = []string{""}

// SymVal is a unique small number representing a symbol.
type SymVal int

// String representation of a SymVal is its original string.
func (y SymVal) String() string {
	return symList[y]
}

// Sym returns a SymVal for the specified string, creating a new one if needed.
func Sym(s string) SymVal {
	if v, ok := symbols[s]; ok {
		return v
	}
	i := SymVal(len(symList))
	symbols[s] = i
	symList = append(symList, s)
	return i
}

// the registry is used to map some symbols for use as gadget constructors
var registry = map[SymVal]func() Circuitry{}

// Register a constructor for a named gadget type.
func Register(name string, f func() Circuitry) {
	registry[Sym(name)] = f
}
