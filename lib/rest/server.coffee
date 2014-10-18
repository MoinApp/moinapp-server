restify = require 'restify'
moinMiddleware = require './moin-middleware'
routes = require './routes'

###
Server Config
###
defaultThrottle = restify.throttle {
  # requests per second
  rate: 2
  burst: 5
  ip: true
}
moinThrottle = restify.throttle {
  # requests per second
  # should be 1 / 5 but does not seem to allow any requests at all then
  rate: 1,
  burst: 1,
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
    # MIDDLEWARE #
    # only accept requests we can respond to
    @server.use restify.acceptParser(@server.acceptable)
    # sync date
    @server.use restify.dateParser()
    # enable body requests
    @server.use restify.bodyParser()
    # enable query parameters
    @server.use restify.queryParser()
    # sanitize paths
    @server.use restify.pre.sanitizePath()
    # should enable logging. Does not do anything?
    @server.use restify.requestLogger()
    # enable GZIP responses
    @server.use restify.gzipResponse()
    # Add all headers to the response
    @server.use restify.fullResponse()
    
    # ROUTES #
    @server.use defaultThrottle
    
    # These routes require not login
    @server.get '/', routes.index.GETindex
    
    # Login methods
    @server.post '/api/auth', routes.session.POSTsignin
    @server.post '/api/signup', routes.session.POSTsignup
    # Authorized methods
    @server.post '/api/moin', moinThrottle, routes.session.checkAuthentication, moinMiddleware(moinController), routes.moin.POSTmoin
    @server.get '/api/user/:username', routes.session.checkAuthentication, routes.user.GETuser
    @server.post '/api/user/addgcm', routes.session.checkAuthentication, routes.user.POSTaddGcm
  
  start: (port) ->
    @server.listen port, ->
      console.log "MoinWebServer started on port #{port}."
  
module.exports.MoinWebServer = MoinWebServer
