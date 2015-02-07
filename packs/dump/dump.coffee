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

client.on 'message', (topic, payload, packet) ->
  # a leading "R:" indicates that the message is a resent/retained one
  topic = 'R: ' + topic  if packet.retain

  s = JSON.stringify(msgpack.unpack(payload))

  # sanitise the output so binary data doesn't pass through unescaped
  #
  limit = 78 - topic.length
  limit = 20  if limit < 20
  if s.length > limit
    s = s.substr(0, limit-3) + '...'
  else
    s = s.substr(0, limit)
  s = s.replace(/[\x00-\x1F\x7F-\xFF]/g, '?')

  console.log topic, s
