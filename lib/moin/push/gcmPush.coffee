gcm = require 'node-gcm'
uuid = require 'node-uuid'

class GCMPush
  constructor: (gcmApiKey) ->
    @gcmSender = new gcm.Sender gcmApiKey
    
  send: (sender, receipient, callback) ->
    if !sender?.getPublicModel?
      return callback new Error 'Must provide database object for "sender".'
    if !receipient?.getPublicModel?
      return callback new Error 'Must provide database object for "receipient".'
    
    user.getGcmIDs().complete (err, gcmIDs) ->
      return callback err if !!err
      return callback new Error 'No device registered for this user.' if !gcmIDs
        
      # convert into public model
      gcmIDs ( gcmID.getPublicModel() for gcmID in gcmIDs )
    
      # now that we have the devices' push tokens
      message = new gcm.Message()
      message.addDataWithObject sender.getPublicModel()
      message.addDataWithKeyValue 'moin-uuid', uuid.v1()
      
      numberOfRetries = 1
      gcmSender?.send message, gcmIDs, numberOfRetries, callback

module.exports.GCMPush = GCMPush
