angular-websocket-service is a [WebSocket][] service for
[AngularJS][].  Use it to send and receive messages consisting of a
topic string (without spaces) followed by a space and a [JSON][]
payload.  For example:

    /state/put {"name": "Alice", "friends": ["Bob", "Charlie"]}

The user-facing API is:

* `$websocket.connect(endpoint)` to return a new wrapped WebSocket
  object connected to the WebSocket server at `endpoint`.

Wrapped WebSocket instances have the following user-facing API:

* `emit(topic, body)` to send outgoing messages.  Messages send before
  the WebSocket connects are queued until the connection completes.
* `register(topic, callback[, options])` to add a callback for
  incoming messages.  The callback fingerprint should be
  `callback(topic, body)`.  The callback is triggered for messages
  who's topic starts with the registered topic.  For example,
  callbacks registered for `/state` and `/state/get` will both be
  invoked for a `/state/get` message.  You can optionally pass in an
  `options` object to configure the registration.  For example, if you
  *do* want an exact topic match (instead of a prefix match), use
  `{exact: true}` (a callback registered for `/state` will `exact` set
  to true will not be invoked for a `/state/get` message).

Similar packages
================

* [angular-websocket][] sends and receives opaque messages using a
  `send` method and events with a configurable prefix (MIT).  The
  event-generating code is not clear to me, so I'm not actually sure
  how the prefix bit works.
* [angular-websocket-provider][] sends and receives JSON messages
  using a `request` method and `websocket.<METHOD>` and
  `websocket.message` events.  `<METHOD>` is set using the payload's
  `method` value (GPLv2).
* [angular-reconnecting-websocket][] wraps websockets to automatically
  reconnect after accidental disconnections (MIT).

[WebSocket]: http://www.w3.org/TR/websockets/
[AngularJS]: https://angularjs.org/
[JSON]: http://json.org/
[angular-websocket]: https://github.com/gdi2290/angular-websocket
[angular-websocket-provider]: https://github.com/instabledesign/angular-websocket
[angular-reconnecting-websocket]: https://github.com/adieu/angular-reconnecting-websocket
