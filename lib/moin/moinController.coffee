{ EventEmitter } = require 'events' 
db = require '../db/'
{ GCMPush } = require './push/gcmPush'

class MoinController extends EventEmitter
  constructor: (sender, receipient) ->
    @androidPush = new GCMPush process.env.GCM_API_KEY
    
  _getUsersFromNames: (senderName, receipientName, callback) ->
    @_resolveUser senderName, (err, sender) =>
      return callback err if !!err
      return callback new Error 'User "' + senderName + '" not found.' if !sender
      
      @_resolveUser receipientName, (err, receipient) =>
        return callback err if !!err
        return callback new Error 'User "' + receipientName + '" not found.' if !receipient
        
        callback null, sender, receipient
    
  _resolveUser: (username, callback) ->
    db.User.find({
      where: {
        username
      }
    }).complete callback
   
  _sendMoinViaPush: (sender, receipient, callback) ->
    # send push notifications
    warnings = []
    minimumSuccessCount = 2
  
    @androidPush.send sender, receipient, (err, results) ->
      # do not crash because we want to send other notifications, too!
      if !!err
        minimumSuccessCount--
        warnings.push {
          type: "android",
          error: err.message
        }
    
      iOSPush = (a, b, cb) ->
        # TODO: add iOS
        cb? new Error 'Not implemented.'
      
      iOSPush sender, receipient, (err, results) ->
        if !!err
          minimumSuccessCount--
          warnings.push {
            type: "iOS",
            error: err.message
          }
      
        err = null
        if minimumSuccessCount <= 0
          console.log "Unsuccessful push:", warnings
          
          warningsString = ""
          warnings.forEach (warning) ->
            if warningsString != ""
              warningsString += ", "
            warningsString += JSON.stringify warning
          err = new Error 'No Push succeeded: ' + warningsString
        callback? err, warnings
  _sendMoinEvent: (sender, receipient) ->
    @emit 'moin', sender, receipient
  sendMoin: (senderName, receipientName, callback) ->
    @_getUsersFromNames senderName, receipientName, (err, sender, receipient) =>
      return callback err if !!err
      
      # users are validated at this point
      # send event
      @_sendMoinEvent sender, receipient
      # now onto the push task...
      
      @_sendMoinViaPush sender, receipient, callback

module.exports.MoinController = MoinController
