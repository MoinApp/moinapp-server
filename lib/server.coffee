restify = require 'restify'

server = null
###
# INIT
###
exports.init = ->
  server = restify.createServer()
  
  exports._initRoutes server
  
exports._initRoutes = (server) ->
  
  server.post '/moin', require './routes/moin'

###
# RUN
###
exports.start = ->
  server.listen 3000, () ->
    console.log "Server listening on port 3000."
