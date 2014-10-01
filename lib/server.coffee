restify = require 'restify'
moin = require './routes/moin'
users = require './routes/users'
session = require './routes/session'

server = null
###
# INIT
###
exports.init = ->
  server = restify.createServer()
  # enable throttling
  server.use restify.throttle {
    # requests per second
    rate: 1,
    burst: 2,
    ip: true
  }
  
  # make body requests possible
  server.use restify.bodyParser()
  # make query parameters available
  server.use restify.queryParser()
  
  exports._initRoutes server
  
exports._initRoutes = (server) ->
  
  # Routes without login
  
  server.get '/', (req, res, next) ->
    pkg = require '../package'
    
    res.setHeader 'Location', pkg.homepage # redirect here
    res.send 302
  
  server.post '/user', users.newUser
  server.post '/user/session', users.signIn
  server.post '/user/session/validate', users.validateSession
  
  # Routes that require a login
  
  server.use session.requireLogin
  
  server.post '/moin', restify.throttle {
    rate: 1/5,
    burst: 1/5,
    ip: true
  }, moin.moin
  # Users
  server.get '/user/:username', users.getUser
  server.post '/user/gcm', users.addGCMId

###
# RUN
###
exports.start = ->
  port = process.env.PORT || 3000
  
  server.listen port, () ->
    console.log "Server listening on port #{port}."
