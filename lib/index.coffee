{ MoinWebServer } = require './rest/server'
{ MoinWebSocketServer } = require './ws/server'
{ MoinController } = require './moin/moinController'

main = ->
  port = process.env.PORT || 3000
  
  moin = new MoinController
  
  server = new MoinWebServer
  ws = new MoinWebSocketServer server.server
  
  server.configureServer moin
  ws.configureServer moin
  
  server.start port
  ws.start port

setImmediate ->
  # start a new stack
  main()
