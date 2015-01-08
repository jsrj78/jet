JeeLabs JET System
==================

[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/jeelabs/jet?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

A dataflow framework and server for multi-node embedded systems.

The "**J**eelabs **E**mbello **T**oolkit" is intended to provide a simple visual way to wire pieces of software together that run on tiny embedded systems and on a back-end server, all connected through a variety of network technologies, in particular simple RF communication.

At this stage JET is just being conceived, so there's nothing "ready-to-use" here yet.

System components
-----------------

JET consists of a number of components. The reason for having multiple processes is two-fold: for one, it's a distributed system consisting of many embedded computers, but the other reason is to support configurations where there is a stable part of the system that may be "in production" and other parts that are in development and being restarted and broken repetitively.

![](web/img/jet-overview.png)

The components are:

- The **JET/Hub** process is a small always-on server and supervisor, acting as inter-connect between various communication interfaces and the rest of the system. The hub process also includes a web server for HTTP(S) and WebSocket connections to drive web-based front ends.
- One or more hardware **JET/Bridges** are used between network technologies to connect to the hub. Typically these bridges are JeeNodes with both an RFM12 or RFM69 RF module and a serial/USB/I2C link, an EtherCard, or a Wifi adapter.
- **JET/Packs** are processes, usually running a set of "circuits". These packs communicate with the rest of the system through the hub, which they can also use to present a browser front end.
- **JET/Web** runs in the browser. It communicates with the central web server over the network and is implemented using HTML5, CSS3, and JavaScript.
- **Circuits** are the building blocks for packs. They define most of the actual logic in the system and can be written / designed using the visual **JET/Flow** dataflow framework.
- **Nodes** are embedded systems, such as JeeNodes, that run a simplified C++ dataflow engine called **JET/Mugs** and communicate with a bridge, typically via RF links.

The overall intention is to have one hub that is always-on and then one or more packs for the rest of the system. Some packs may be always-on and run "production" code, but others may bounce as they are being developed. All packs and their circuits can access all nodes attached to the hub, which makes it easy to deploy new nodes in an existing RF network and separate new nodes with their corresponding circuits from older nodes / circuits that need to continue running undisturbed.

