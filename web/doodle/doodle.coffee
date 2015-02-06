app = angular.module 'DoodleApp', ['websocket']

app.value 'wsPort', 1111 # websocket may not be on same server as static files

app.controller 'DoodleCtrl', ($scope, $timeout, $websocket, wsPort) ->
  wsProto = if "https:" is document.location.protocol then "wss" else "ws"
  ws = $websocket.connect "#{wsProto}://#{location.hostname}:#{wsPort}/ws"

  ws.register '', (topic, body) ->
    console.log 'mqtt:', topic, body

  ws.emit '/doodle/dah', [7, 8, 9]

  $timeout ->
    ws.emit '/doodle', [1, 2, 3]
  , 1000
  $timeout ->
    ws.emit '/doodledah', [4, 5, 6]
  , 2000
