JET/Hub
=======

This is the central command centre for Jet. It must always be running.  
It requires [node.js][N] (plus all dependencies installed with `npm install`).

### Development

Dev mode depends on [Leiningen][L] (and Java). Start them as follows:

    lein figwheel

Keep this running, then launch the main Node.js process in a separate session:

    node figwheel.js

### Production

To generate a production-ready app, run `lein cljsbuild once prod`.  
Production requires neither Leiningen nor Java, just launch the hub as:

    node server.js

The hub is the central server, it needs to stay up and running at all times.

[N]: https://nodejs.org
[L]: https://github.com/technomancy/leiningen
