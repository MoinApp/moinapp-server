restify = require 'restify'
moinMiddleware = require './moin-middleware'
logger = require './log-middleware'
apiV200 = require('./routes').v200
apiV300 = require('./routes').v300

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
  constructor: ->
    @server = restify.createServer({
      #version: apiV200.version # REST version
      versions: [ apiV200.version, apiV300.version ] # REST version
    })

    crashOnError = true
    if crashOnError
      @server.on 'uncaughtException', (req, res, route, err) ->
        throw err

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

    @routeV200 moinController
    @routeV300 moinController

  routeV300: (moinController) ->
    # These routes require not login
    @server.get { path: '/', version: apiV300.version }, apiV300.index.GETindex

    # Login methods
    @server.post { path: '/api/auth', version: apiV300.version }, apiV300.session.POSTsignin
    @server.post { path: '/api/signup', version: apiV300.version }, apiV300.session.POSTsignup
    # Authorized methods
    @server.post { path: '/api/moin', version: apiV300.version }, moinThrottle, apiV300.session.checkAuthentication, moinMiddleware(moinController), apiV300.moin.POSTmoin
    @server.get { path: '/api/user/recents', version: apiV300.version }, apiV300.session.checkAuthentication, apiV300.user.GETusersRecents
    @server.get { path: '/api/user/:username', version: apiV300.version }, apiV300.session.checkAuthentication, apiV300.user.GETuser
    @server.get { path: '/api/user', version: apiV300.version}, apiV300.session.checkAuthentication, apiV300.user.GETuserSearch
    @server.post { path: '/api/user/addgcm', version: apiV300.version }, apiV300.session.checkAuthentication, apiV300.user.POSTaddGcm


  routeV200: (moinController) ->
    # These routes require not login
    @server.get { path: '/', version: apiV200.version }, apiV200.index.GETindex

    # Login methods
    @server.post { path: '/api/auth', version: apiV200.version }, apiV200.session.POSTsignin
    @server.post { path: '/api/signup', version: apiV200.version }, apiV200.session.POSTsignup
    # Authorized methods
    @server.post { path: '/api/moin', version: apiV200.version }, moinThrottle, apiV200.session.checkAuthentication, moinMiddleware(moinController), apiV200.moin.POSTmoin
    @server.get { path: '/api/user/:username', version: apiV200.version }, apiV200.session.checkAuthentication, apiV200.user.GETuser
    @server.post { path: '/api/user/addgcm', version: apiV200.version }, apiV200.session.checkAuthentication, apiV200.user.POSTaddGcm

  start: (port) ->
    @server.listen port, ->
      console.log "#{@constructor.name} started on port #{port}."

module.exports.MoinWebServer = MoinWebServer
