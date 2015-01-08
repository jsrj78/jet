JeeLabs Jet System
==================

[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/jeelabs/jet?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Dataflow framework and server for JeeLabs embedded systems

The Jet system is intended to provide a simple visual way to wire pieces of software together that run on tiny embedded systems and on a back-end server, all connected through a variety of network technologies, in particular simple RF communication. At this stage Jet is just being conceived, so there's not much "ready-to-use" here yet.

System components:
-----------------
Jet consists of a number of components. The reason for having multiple parts is two-fold: for one, it's a distributed system consitsing of many tiny embedded computers, but the other reason is that we want to support configurations where there is a stable part of the system that may be "in production" and other parts that are in development and being restarted and broken repetitively. The parts are:
- The Jet Bridge (or Hub?) is a simple always-up server process that bridges communication between all the embedded systems (typically JeeNodes reachable over RFM12 or RFM69 RF networks) and the various application servers.
- The Jet Server runs applications that are written/designed using the visual Jet Flow dataflow framework. The application servers communicate with embedded systems through the Jet Bridge and often have a Jet Web browser application (javascript world)
- The Jet Web applications run in the browser and communicate with Jet Servers over HTTP and websocket
- The Jet Nodes are embedded systems, such as JeeNodes, that run a simplified C++ version of Jet Flow and communicate with Jet Servers through the Jet Bridge, typically via RF links
- Optionally, a number of Jet Gateways bridge [uhh, should these be Jet Bridges?] between network technologies to connect up to a Jet Bridge, typically these gateways are JeeNodes with both an RFM12 or RFM69 RF module and a serial link, a EtherCard, or a Wifi adapter.

The overall intention is to have one Jee Bridge that is always-on and then one or sevaral Jee Servers that run applications. Some Jee Servers may also be always-on and run "production" applications, but others may bounce as their apps are being developed. All the Jee Servers and their applications can access all nodes attached to the bridge, which makes it easy to deploy new nodes in an existing RF network and separate the new nodes and the corresponding new apps from older nodes and apps that need to continue running undisturbed.
