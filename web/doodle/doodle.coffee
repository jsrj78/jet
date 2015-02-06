ng = angular.module 'DoodleApp', ['websocket']

ng.controller 'DoodleCtrl', ($scope, $timeout, $websocket) ->
  wsProto = if "https:" is document.location.protocol then "wss" else "ws"
  ws = $websocket.connect "#{wsProto}://#{location.hostname}:1111/ws", ["jet"]

  ws.register '', (topic, body) ->
    console.log 'mqtt:', topic, body

  $timeout ->
    ws.emit '/doodle', [1, 2, 3]
  , 1000
  $timeout ->
    ws.emit '/doodledah', [4, 5, 6]
  , 2000
  $timeout ->
    ws.emit '/doodle/dah', [7, 8, 9]
  , 3000
