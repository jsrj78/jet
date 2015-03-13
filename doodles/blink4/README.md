The JET/blink demo in React, with a ticker on the host side.

(adapted from the blink2 demo)

The enabled state is sent to the host as: `{enabled: <bool>}`  
The blink state comes back from the host as: `{blink: <bool>}`  
So the input and output are in the browser, but the ticker logic is on the host.

**Status** - Working as intended.

### Usage

    go run server.go
    open http://localhost:8000

Press refresh to see changes, JSX gets compiled on-the-fly in the browser.
