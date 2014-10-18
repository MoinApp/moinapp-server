restify = require 'restify'
moinMiddleware = require './moin-middleware'
routes = require './routes'

###
Server Config
###
defaultThrottle = restify.throttle {
  # requests per seconds
  rate: 1
  burst: 3
  ip: true
}

###
Server Class
###
class MoinWebServer
  constructor: (moinController) ->
    @server = restify.createServer({
      version: "2.0.0" # REST version
    })
    @configureServer moinController
    
  configureServer: (moinController) ->
    @server.use defaultThrottle
    
    # MIDDLEWARE #
    # only accept requests we can respond to
    @server.use restify.acceptParser(@server.acceptable)
    # enable GZIP responses
    @server.use restify.gzipResponse()
    # Add all headers to the response
    @server.use restify.fullResponse()
    # enable query parameters
    @server.use restify.queryParser()
    # enable body requests
    @server.use restify.bodyParser()
    # should enable logging. Does not do anything?
    @server.use restify.requestLogger()
    
    # ROUTES #
    # These routes require not login
    @server.get '/', routes.index.GETindex
    
    # Login methods
    @server.post '/api/auth', routes.session.POSTsignin
    @server.post '/api/signup', routes.session.POSTsignup
    # Authorized methods
    @server.post '/api/moin', routes.session.checkAuthentication, moinMiddleware(moinController), routes.moin.POSTmoin
    @server.post '/api/user/addgcm', routes.session.checkAuthentication, routes.user.POSTaddGcm
    @server.get '/api/user/:username', routes.session.checkAuthentication, routes.user.GETuser
  
  start: (port) ->
    @server.listen port, ->
      console.log "MoinWebServer started on port #{port}."
  
module.exports.MoinWebServer = MoinWebServer
