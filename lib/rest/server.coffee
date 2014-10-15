restify = require 'restify'
routes = require './routes'

###
Server Config
###
defaultThrottle = restify.throttle {
  # requests per seconds
  rate: 1
  burst: 2
  ip: true
}

###
Server Class
###
class MoinWebServer
  constructor: ->
    @server = restify.createServer({
      version: "2.0.0" # REST version
    })
    @configureServer()
    
  configureServer: ->
    @server.use defaultThrottle
    
    # MIDDLEWARE #
    # enable body requests
    @server.use restify.bodyParser()
    # enable query parameters
    @server.use restify.queryParser()
    
    # ROUTES #
    # These routes require not login
    @server.get '/', routes.index.GETindex
    
    # Login methods
    @server.post '/api/auth', routes.session.POSTsignin
    @server.post '/api/signup', routes.session.POSTsignup
    # Authorized methods
    @server.post '/api/moin', routes.session.checkAuthentication, routes.moin.POSTmoin
    
    @server.get '/api/user/:username', routes.session.checkAuthentication, routes.user.GETuser
  
  start: (port) ->
    @server.listen port, ->
      console.log "MoinWebServer started on port #{port}."
  
module.exports.MoinWebServer = MoinWebServer
