package gadgets

import (
	// these packages register themselves as side-effect
	_ "github.com/jeelabs/jet/gadgets/net/mqtt"
	_ "github.com/jeelabs/jet/gadgets/serial"
)

// no code here, this package only exists to load other gadget collections
