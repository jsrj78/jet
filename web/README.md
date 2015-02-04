JET/Web is the web browser stuff. Or rather, a placeholder for things to come.

### Development note

The `config.codekit` file is for a Mac-specific dev tool called CodeKit,
which can do file serving, live reload, pre-processing, and minification.
This shouldn't prevent using Node.js during development on other platforms.

To try this out, you need to start the hub and the "front" pack:

* do `jet start` and then `cd ../packs/front && go run front.go`
* leave that running, then open a browser on http://localhost:1111
* open the JavaScript console to see what's going on
