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
  server.listen 3000, () ->
    console.log "Server listening on port 3000."
