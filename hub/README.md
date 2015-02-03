JET/Hub - central hub & router
==============================

JET/Hub is the central always-on server which supervises and interconnects all
other parts of a JET system.

Godoc documentation: <http://godoc.org/github.com/jeelabs/jet/hub>  
MIT license, see the LICENSE file: <https://github.com/jeelabs/jet>

Running JET/Hub
---------------

* do `cd jet/hub && make` to build the executable (jet/build/jethub)
* put a symlink to the `jet/jet` script somewhere in your executable PATH
* to start the hub as background process: `jet start`
* to see the last line of the log: `jet status`
* stop the hub again with: `jet stop`
* unless configured otherwise, logs are appended to `jet/run/hub.log`
* wanna know more? `jet help`

Functions
---------

- Communicate with embedded nodes either via directly connected ports, such as serial or I2C, or
  via bridges reachable over UDP
- Communicate with JET/Pack servers via http/websockets
- Route messages between embedded nodes and JET/Packs
- Store messages short term to allow JET/Packs to replay messages they missed if they restart,
  or for other purposes such as reprocessing them
- Launch and supervisor JET/Pack processes

APIs
----

The hub has at least 4 APIs:
- Packs connect to the hub and communicate with nodes and also with browser apps (JET/Web), they need
  to send & receive messages and serve requests made by the browser
- Browsers connect to the hub and get served web pages with apps that then make requests to
  server-side applications running in Packs
- Drivers for node connectivity use an internal API to feed RF messages into the system and pull
  messages to be transmitted to the nodes
- A UDP driver extends the internal node connectivity API over UDP to allow for distant drivers
The APIs for packs and browsers are really one with different sets of requests for the two quite 
distinct functions.

Pack API
--------

The pack API allows server-side apps to connect to the hub, discover the nodes available, subscribe
messages from these nodes, send messages to nodes, and modify the programs that nodes run (i.e.,
what code they bootstrap).

The pack API supports the following requests (all this can be formulated as restful http
requests but it's stated generically here):
- Retrieve the list of nodes with info about each node, including unique ID, network
  address, what code it's running, when it was last seen, when it last booted
- Update a node allowing to set the name, the network address, and the program
- Reboot a node (sends it a message to reboot)
- Retrieve the list of programs with info about each program
- Add/update a program (i.e. upload HEX code)
- Send a message to a node (or broadcast)
- Subscribe to a feed from a node or a set of nodes starting at time T (this will replay messages
  if T is in the past)
- Subscribe to requests sent by browsers to a specific URL subtree (and ability to respond to
  those requests)
- Upload a static web page to the hub for serving to browsers (this could be phrased in various
  ways, including mounting a filesystem directory to a URL root path)
- ??

The browser API supports the following requests:
- GET a static page (that was prepared by a pack)
- POST to a URL to which a pack has subscribed and upgrade to websockets to "talk to the pack"

The internal API for drivers consists of:
- Create/delete network (a way for the driver to register a new network)
- Regular message received from node
- Regular message send to node
- Pairing message received from node
- Boot request from node
- Boot data to node
The boot protocol could be mapped into regular messages for example by using a special node ID, but
in general, a pack that subscribes to messages from node 55 shouldn't receive the boot messages
that node 55 sends out, instead, a boot server that subscribes to boot messages should get those.
