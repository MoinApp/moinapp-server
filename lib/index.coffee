{ MoinWebServer } = require './rest/server'

main = ->
  port = process.env.PORT || 3000
  
  server = new MoinWebServer
  
  server.start port

setImmediate ->
  # start a new stack
  main()
