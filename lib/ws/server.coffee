socketio = require 'socket.io'
{ SessionHandler } = require './../auth/sessionHandler'

class MoinWebSocketServer
  constructor: (webServer) ->
    @server = socketio.listen webServer
    
  configureServer: (@moinController) ->
    @server.sockets.on 'connection', (socket) =>
      @onNewConnection socket
      
  start:  ->
    console.log "#{@constructor.name} running."
    
  onNewConnection: (socket) ->
    socket.on 'login', (session, cb) =>
      SessionHandler.getInstance().getUserForSessionToken session, (err, user) =>
        @moinController.on 'moin', (sender, receipient) ->
          if receipient.username == user.username
            socket.emit 'moin', { sender: sender.username }
            
        socket.on 'moin', (receipientName, fn) =>
          @moinController.sendMoin user.username, receipientName, fn
        
        cb? err, !!user
    
module.exports.MoinWebSocketServer = MoinWebSocketServer
