JET heater example mapped to MQTT
=================================

Some definitions:

- Circuit: a circuit is a collection of gadgets or mugs that all run within one execution container.
- Execution container: runs a circuit and can be a Node (embedded system), a Pack (server-side system)
  or a WebApp (browser)
- Gadget or mug: an atomic component in the program with inputs and outputs (circuits can be used as
  gadgets in a herarchical manner)
- Wire: connects the output of one gadget to the input of another gadget irrespective of circuit boundaries
- Message: atomic communication unit that carries data (and often control flow) from the output of one
  gadget to the input of another along a wire

Question: a circuit is both a static entity (piece of code) that can be instantiated many times
and it's a dynamic instantiated, running, executing entity. Do we want two different names? I'm
using is purely in its dynamic meaning here.

Identifiers:

It is assumed here that all entities have identifiers. It's not clear what the nature of the
identifiers is, however. This is because within a Node it's probably best to use small integers, for
example, all the gadgets within a circuit could be numbered starting with 0. But within a
Pack it may be more helpful overall to generate IDs that include the name of the gadget as a string,
such as "time_series_water_temp", possibly with a unique numerical suffix where necessary.

Communication:

It is assumed here that all wires within a circuit are handled internally within the circuit.
Only wires that cross circuit boundaries require communication and involve MQTT.
However, the proposals below assign an MQTT path to all wires, including the internal ones,
with the thought that this may be useful for debugging purposes, i.e., to allow internal
wires to be exposed easily if necessary.
Plus, naming all wires seems like a more consistent way to go.

MQTT paths:

- Each circuit is assigned a unique MQTT path subtree:
- `<type>/<circuit_id>/`
where:
- type is either `node`, `pack` or `web`
- circuit_id for nodes is a 16 or 32-bit number (node_uid), and for packs and web apps is either
  a 32-bit unique id or a string unique id (we need to choose)
- it is assumed for now that nodes have a single circuit (can debate this later...)

Each circuit has a configuration path to which the circuit configurator listens for instructions
  on what to do with the circuit, such as wire it up, shut it down, etc.:
- `<type>/<circuit_id>/conf`

Each gadget in each circuit is assigned a unique MQTT path subtree below it's circuit:
- `<type>/<container_id>/<gadget_id>/`
where:
- gadget_id is a sequential numbering of gadgets, typically optimized to give gadgets with wires
  that cross circuit boundaries a low number (it could also be a string unique id where we don't
  care about the length of the id)

Each output of each gadget in each circuit is also assigned a unique MQTT path and each message
output by the gadget is sent to that path:
- `<type>/<container_id>/<gadget_id>/<output_num>`
where:
- output_num is the sequential numbering of outputs on the gadget from left-to-right

When a circuit is instantiated its inputs must be hooked up to the outputs of other existing
circuits and the inputs of existing circuits that need to be connected to the new circuit's
outputs also need to be hooked up. The hooking up is done by sending messages to the configuration
path of each circuit. The messages are of the form (in pseudo-code):
- Subscribe your <gadget_id>/<input_num> to <mqtt_path_of_other_circuit_output>
where:
- <gadget_id>/<input_num> is a numbering of gadgets and inputs similar to the one done for outputs

Node optimizations:

Because handling long string-form MQTT paths in nodes is awkward, as is the whole publishing and
subscription protocol, the driver in the hub facilitates the whole process. Specifically:
- nodes simply transmit gadget outputs using their RF node_id as source address and using a 1-byte
  message tag consisting of a 5-bit gadget_id and a 3-bit output_id
- drivers expand the received messages to the full MQTT publish operation by mapping the
  source address to a 16/32-bit node_uid
- if more than 31 gadgets have cross-circuit wires or a gadget has more than 8 outputs then a
  0xFF byte escapes to a 16-bit format (10 bits of gadget_id and 6 bits of output_num?)
- drivers subscribe to nodes' configuration path and pre-process messages. Messages about
  input subscriptions are handled by the driver who then forwards the messages received (see
  next bullet). The driver does forward messages for things like shutting down a circuit to the
  node).
- nodes receive gadget inputs as messages with their node_id as destination address and a 1-byte
  message tag consisting of a 5-bit gadget_id and a 3_bit input_id (sequential numbering of
  inputs on the gadget). However, gadget 0 input 7 is reserved for configuration messages

An alternative would be to number cross-circuit outputs sequentially allowing for 256 outputs
regardless of how many gadgets are involved. Similarly we could number the inputs allowing for
255 inputs (one being reserved for the circuit configuration input).

Something to consider: if the overhead is acceptable the code could be structured such that
every output of every gadget can be marked to send a message at runtime for debugging purposes.
Similarly, a configuration message could be send to allow a data message to be sent to any input
of any gadget for debugging purposes. Of course, instead of making this run-time dynamic stuff
the circuit could be recompiled and transparently reloaded.

Circuit loading:

Loading a circuit causes it to "run". In general applications consist of multiple circuits that
are loaded, for example, a circuit in each participating node, one in a pack, and one in every
web app (browser that connects). However, each circuit is loaded on its own, or put differently,
loading a multi-circuit app happens one circuit at a time.

To load a circuit:
- create a unique ID for the circuit (in the case of a node this may be the built-in node ID instead
  of a freshly generated ID)
- send the code for the circuit to its container (pack, node, web browser)
- for all inputs to the circuit that are connected to outputs of existing other circuits send
  a message to the new circuit's configuration path
- for all outputs of the new circuit that are connected to inputs of existing other circuits send
  a message to the configuration input of the existing circuits
- celebrate

To unload a circuit:
- send config messages to the circuit to unsubscribe from other circuit's outputs
- tell the container to drain the circuit (?)
- send config messages to other circuits to unsubscribe from this circuit
- tell the container to stop running the circuit and unload it

To reload a circuit to update its code:
- assume that the circuit keeps gadget numbering the same as much as possible
- keep all subscriptions that stay intact so MQTT sees it as a disconnect and reconnect, for
  that purpose send the circuit a disconnect message so it disconnects from MQTT without
  actually unsubscribing (I believe this notion exists)
- send configuration messages to drop subscriptions where wires actually go away
- tell the container to reload the circuit
- send the circuit config messages to subscribe where new wires are appearing, and to
  reconnect to existing subscriptions
- send other circuits messages to subscribe where new wires appear

