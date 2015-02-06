app = angular.module 'MyApp', ['websocket']

app.run ($rootScope) ->
  $rootScope.dummyMsg = 'JET says hello...'

app.controller 'MainCtrl', ($scope, $timeout, $websocket) ->
  wsProto = if "https:" is document.location.protocol then "wss" else "ws"
  ws = $websocket.connect "#{wsProto}://#{location.hostname}:1111/ws", ["jet"]

  ws.register '/haha', (topic, body) ->
    console.log 'ws1 got:', topic, body
  ws.register '/haha', (topic, body) ->
    console.log 'ws2 got:', topic, body
  , { exact: true }

  $timeout ->
    ws.emit '/haha', [1, 2, 3]
  , 1000
  $timeout ->
    ws.emit '/hahaha', [4, 5, 6]
  , 2000
  $timeout ->
    ws.emit '/haha/ha', [7, 8, 9]
  , 3000
