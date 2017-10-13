from __future__ import print_function

import paho.mqtt.client as mqtt
import flow, gadgets
import json

PREFIX = "s/pyf1"

circuits = {}
client = None

class ConnectedCircuit(flow.Circuit):
    def __init__(self, name):
        flow.Circuit.__init__(self)
        self.name = name
        circuits[name] = self

    def control(self, msg):
        print("CONTROL:", self.name, msg)

    def incoming(self, inum, msg):
        print("IN:", self.name, inum, msg)
        self.feed(inum, msg)

    def emit(self, onum, msg):
        print("OUT:", onum, msg)
        topic = "%s/%s/out/%d" % (PREFIX, self.name, onum)
        client.publish(topic, json.dumps(msg))

def on_connect(client, userdata, flags, rc):
    print("Connected: code", rc)
    subs = [(PREFIX, 0)]
    for name in circuits:
        subs += [("%s/%s" % (PREFIX, name), 0),
                 ("%s/%s/in/+" % (PREFIX, name), 0)]
    client.subscribe(subs)

def on_message(client, userdata, msg):
    payload = json.loads(msg.payload)
    if msg.topic == PREFIX:
        print("CMD:", payload)
    else:
        topic = msg.topic[len(PREFIX)+1:]
        parts = topic.split('/')
        cob = circuits[parts[0]]
        if len(parts) == 1:
            cob.control(payload)
        else:
            assert(len(parts) == 3 and parts[1] == 'in')
            cob.incoming(int(parts[2]), payload)

# loop back test circuit: print msgs from inlet 0 and pass them to outlet 0
c = ConnectedCircuit('circ1')
c.add('inlet')
c.add('pass')
c.add('print', 'got:')
c.add('outlet')
c.wire(0, 0, 1, 0)
c.wire(1, 0, 2, 0)
c.wire(1, 0, 3, 0)

client = mqtt.Client()
client.on_connect = on_connect
client.on_message = on_message

client.connect("localhost")
client.loop_forever()
