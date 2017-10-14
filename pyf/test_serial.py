from __future__ import print_function

import paho.mqtt.client as mqtt
import json, sys, time

def on_connect(client, userdata, flags, rc):
    client.subscribe('s/pyf-demo/test/out/0')

def on_message(client, userdata, msg):
    print('reply:', msg.topic, str(msg.payload))

client = mqtt.Client()

def send(topic, msg):
    client.publish(topic, json.dumps(msg))

client.on_connect = on_connect
client.on_message = on_message
client.connect('localhost')
client.loop_start()

port = '/dev/cu.usbmodem34208131'  # jeelabs mac default, but can be overriden
if len(sys.argv) > 1:
    port = sys.argv[1]

send('s/pyf-demo', ['create', 'test'])
send('s/pyf-demo/test', [['inlet'],
                         ['serial', port],
                         ['outlet'],
                         [0, 0, 1, 0],
                         [1, 0, 2, 0]])

send('s/pyf-demo/test/in/0', '1 2 + .')
time.sleep(1)
send('s/pyf-demo/test/in/0', '11 22 + .')
time.sleep(1)
