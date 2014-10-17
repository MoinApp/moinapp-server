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
    @androidPush.send @sender, @receipient, callback
    # TODO: add iOS

module.exports.MoinController = MoinController
