restify = require 'restify'
users = require './routes/users'

server = null
###
# INIT
###
exports.init = ->
  server = restify.createServer()
  server.use restify.bodyParser()
  
  exports._initRoutes server
  
exports._initRoutes = (server) ->
  
  server.post '/moin', require './routes/moin'
  # Users
  server.get '/user/:username', users.getUser
  server.post '/user', users.newUser
  server.post '/user/session', users.signIn

###
# RUN
###
exports.start = ->
  port = process.env.PORT || 3000
  
  server.listen port, () ->
    console.log "Server listening on port #{port}."
