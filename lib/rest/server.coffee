restify = require 'restify'
moinMiddleware = require './moin-middleware'
logger = require './log-middleware'
apiV200 = require('./routes').v200

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
      version: apiV200.version # REST version
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
    # enable GZIP responses
    @server.use restify.gzipResponse()
    # Add all headers to the response
    @server.use restify.fullResponse()
    
    @server.use logger
    
    # ROUTES #
    @server.use defaultThrottle
    
    # These routes require not login
    @server.get '/', apiV200.index.GETindex
    
    # Login methods
    @server.post '/api/auth', apiV200.session.POSTsignin
    @server.post '/api/signup', apiV200.session.POSTsignup
    # Authorized methods
    @server.post '/api/moin', moinThrottle, apiV200.session.checkAuthentication, moinMiddleware(moinController), apiV200.moin.POSTmoin
    @server.get '/api/user/:username', apiV200.session.checkAuthentication, apiV200.user.GETuser
    @server.post '/api/user/addgcm', apiV200.session.checkAuthentication, apiV200.user.POSTaddGcm
  
  start: (port) ->
    @server.listen port, ->
      console.log "MoinWebServer started on port #{port}."
  
module.exports.MoinWebServer = MoinWebServer
