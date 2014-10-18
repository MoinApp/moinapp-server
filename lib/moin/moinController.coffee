db = require '../db/'
{ GCMPush } = require './push/gcmPush'

class MoinController
  constructor: (sender, receipient) ->
    @receipient = receipient
    @sender = sender
    
    @androidPush = new GCMPush process.env.GCM_API_KEY
    
  setUsersFromNames: (senderName, receipientName, callback) ->
    @_resolveUser senderName, (err, sender) =>
      return callback err if !!err
      return callback new Error 'User "' + senderName + '" not found.' if !sender
      
      @sender = sender
      
      @_resolveUser receipientName, (err, receipient) =>
        return callback err if !!err
        return callback new Error 'User "' + receipientName + '" not found.' if !receipient
        
        @receipient = receipient
        
        callback null
    
  _resolveUser: (username, callback) ->
    db.User.find({
      where: {
        username
      }
    }).complete callback
        
  sendMoin: (callback) ->
    warnings = []
    minimumSuccessCount = 2
    
    @androidPush.send @sender, @receipient, (err, results) ->
      # do not crash because we want to send other notifications, too!
      if !!err
        minimumSuccessCount--
        warnings.push {
          type: "android",
          error: err.message
        }
      
      iOSPush = (a, b, cb) ->
        cb? new Error 'Not implemented.'
        
      iOSPush @sender, @receipient, (err, results) ->
        if !!err
          minimumSuccessCount--
          warnings.push {
            type: "iOS",
            error: err.message
          }
        
        err = null
        if minimumSuccessCount <= 0
          err = new Error 'No Push succeeded.'
          console.log "Invalid push:", warnings
        callback? err, warnings
    # TODO: add iOS

module.exports.MoinController = MoinController
