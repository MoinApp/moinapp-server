restify = require 'restify'
users = require './routes/users'

server = null
###
# INIT
###
exports.init = ->
  server = restify.createServer()
  
  exports._initRoutes server
  
exports._initRoutes = (server) ->
  
  server.post '/moin', require './routes/moin'
  # Users
  server.get '/user/:name', users.getUser

###
# RUN
###
exports.start = ->
  server.listen 3000, () ->
    console.log "Server listening on port 3000."
