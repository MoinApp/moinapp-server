{ EventEmitter } = require 'events'
socketio = require 'socket.io'
{ SessionHandler } = require './../auth/sessionHandler'
user = require './user'

class MoinWebSocketServer
  constructor: (webServer) ->
    @server = socketio.listen webServer
    
  configureServer: (@moinController) ->
    @server.sockets.on 'connection', (socket) =>
      @onNewConnection socket
      
  start:  ->
    console.log "#{@constructor.name} running."
    
  onNewConnection: (socket) ->
    new MoinWebSocketConnection socket, @moinController
        
class MoinWebSocketConnection extends EventEmitter
  constructor: (@socket, @moinController) ->
      # client sends login message
      @socket.on 'auth', (session, callback) =>
        console.log "validate login", session
        # we validate the session token, just like /api/auth would do
        @_validateSessionToken session, (err, user) =>
          return callback? err if !!err
        
          # we are authenticated!
          @user = user
        
          @_setupEvents()
        
          callback? err, !!user
      
  _validateSessionToken: (session, callback) ->
    SessionHandler.getInstance().getUserForSessionToken session, callback
  _setupEvents: ->
    # moin by others
    @moinController.on 'moin', moinControllerListener = (sender, receipient) =>
      # if this user should receive the moin
      if receipient.username == @user.username
        @socket.emit 'moin', {
          sender: sender.username
        }
        
    @socket.on 'disconnect', =>
      @moinController.removeListener 'moin', moinControllerListener
    # moin by the client
    @socket.on 'moin', (receipientName, callback) =>
      @moinController.sendMoin @user.username, receipientName, callback
      
    @socket.on 'getUser', user.getUser
    
module.exports.MoinWebSocketServer = MoinWebSocketServer
module.exports.MoinWebSocketConnection = MoinWebSocketConnection
