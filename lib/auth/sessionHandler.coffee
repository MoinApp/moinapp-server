restify = require 'restify'
db = require './../db/'

class SessionHandler
  @instance = null
  @getInstance = ->
    if !SessionHandler.instance
      SessionHandler.instance = new SessionHandler
    SessionHandler.instance
  
  constructor: ->
    
  getUserForSessionToken: (token, callback) ->
    db.Session.find({
      where: {
        uid: token
      }
    }).complete (err, session) ->
      return callback err if !!err
      
      if !session
        callback new restify.NotAuthorizedError 'Invalid session token.'
      else
        session.getUser().complete callback
        
module.exports.SessionHandler = SessionHandler
