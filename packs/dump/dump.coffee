#/usr/bin/env nodejs

console.log 'Dump of all hub messages:'

mqtt = require 'mqtt'
msgpack = require 'msgpack'

client = mqtt.connect 'localhost', { keepalive: 300 }

client.on 'error', (e) -> console.log 'error:', e

#client.on 'offline', -> console.log 'offline'
#client.on 'close', -> console.log 'close'

client.on 'connect', ->
  console.log 'connect:', Date()
  client.subscribe '#'

client.on 'message', (topic, payload) ->
  console.log topic, JSON.stringify(msgpack.unpack(payload))
