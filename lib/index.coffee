{ MoinWebServer } = require './rest/server'
{ MoinController } = require './moin/moinController'

main = ->
  port = process.env.PORT || 3000
  
  moin = new MoinController
  server = new MoinWebServer moin
  
  server.start port

setImmediate ->
  # start a new stack
  main()
