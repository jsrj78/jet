**The `blink<N>` doodles are explorations around a common theme & setup:**

* an input to toggle between "enabled" and "disabled" states
* an output which indicates the current "blink" state
* a task which periodically toggles the blink state whenever enabled is true

In other words: _think of a blinking light, controlled by a switch._

But the key point here is that each of these "gadgets" can exist in different
contexts: in the browser, in a host server, or in an embedded microcontroller.

## Blink 1

First attempt to implement an all-in-the-browser version, based on React.

## Blink 2

Simple all-in-the-browser version, based on React. In JavaScript ES5, not ES6.

## Blink 3

Early attempt to implement the same again using ClosureScript. Unfinished.

## Blink 4

This is _Blink 2_, but with the periodic task running on the host, implemented
in Go and connected to the browser via a websocket. Uses JSON as protocol.

## Blink 5

An all-on-the-host implementation, using Go and MQTT (in Go). The "input" and
"output" are now MQTT topics.
