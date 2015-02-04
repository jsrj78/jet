This is a demo "pack" to illustrate basic use of the "hub":

* make sure the hub is running properly: `jet start`
* open a new shell to show all messages: `cd jet/pack/dump && node dump.js`
* in this demo dir you can now launch the demo: `make`

This will start two instances of the demo pack, which connect to the hub and
send a message to the "hub/hello" topic. Then, a second later, the demo pack
will quit and disconnect from the hub. As a result, the hub will send another
message to the "hub/goodbye" topic (using the "last will" mechanism).

Here is a transcript of the message dump output:

    hub/hello [1,"demo",1423039751,10347,"jaycee.local"]
    /test/haha [123,null,"abc"]
    hub/hello [1,"demo",1423039752,10353,"jaycee.local"]
    hub/goodbye [1,"demo",1423039751,10347,"jaycee.local"]
    /test/haha [123,null,"abc"]
    hub/goodbye [1,"demo",1423039752,10353,"jaycee.local"]

The hub/{hello,goodbye} messages contain: [version, pack-name, time, pid, now].

The hub system log itself can be tracked with: `cd jet/run && tail -f hub.log`.
