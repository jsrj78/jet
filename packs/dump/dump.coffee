#/usr/bin/env nodejs

console.log 'Dump of all hub messages:'

mqtt = require 'mqtt'
msgpack = require 'msgpack'

client = mqtt.connect 'localhost', { keepalive: 0 }
client.subscribe '#'

client.on 'message', (topic, payload) ->
  console.log topic, JSON.stringify(msgpack.unpack(payload))
