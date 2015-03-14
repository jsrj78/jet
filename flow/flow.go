// Package flow implements a Pure-data like dataflow mechanism.
package flow

// Message is a generic data item which can be sent between gadgets.
type Message interface{}

// Vec is a slice of messages (and is also a message).
type Vec []Message

// NewVec constructs a Vec from a list of messages.
func NewVec(v ...Message) Vec {
	return v
}

// Nth returns the n'th item of a vec, or NoSym if it doesn't exist.
func (v *Vec) Nth(index int) Message {
	if index >= len(*v) {
		return nil
	}
	return (*v)[index]
}

// Tag returns the first item of a Vec as SymVal, or NoSym if it can't.
func (v *Vec) Tag() SymVal {
	if len(*v) > 0 {
		switch x := (*v)[0].(type) {
		// TODO maybe also int?
		case SymVal:
			return x
		case string:
			return Sym(x)
		}
	}
	return NoSym
}

// Map is a hash of messages with string keys (and is also a message).
type Map map[string]Message

// Circuitry is the collective name for gadgets and circuits.
type Circuitry interface {
	NumInlets() int
	Inlet(n int) *Inlet
	IsHot(i *Inlet) bool

	NumOutlets() int
	Outlet(n int) *Outlet

	Setup()
	Control(Vec)
	Trigger()
	Cleanup()

	install(self Circuitry, owner *Circuit) *Gadget
}

// the empty string, i.e. SymVal zero, is special
const NoSym = SymVal(0)

var symbols = map[string]SymVal{"": NoSym}
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
