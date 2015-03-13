The JET/blink demo using a host-side only setup in Go, based on MQTT.

The input is an MQTT topic called "enabled", accepting a boolean as JSON.
The output is an MQTT topic called "blink", sending out booleans as JSON.
The ticker logic is included in the server and is controlled via pubsub.

**Status** - Working as intended.

### Usage

Launch the server in a separate terminal window:

    go run server.go

Now launch a client which displays all data going through MQTT:

    go run dumper.go

Finally, in a third terminal window, send enable/disable events:

    go run enabler.go true

or

    go run enabler.go false

The results will be visible in the dumper output.
